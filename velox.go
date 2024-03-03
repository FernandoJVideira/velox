package velox

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FernandoJVideira/velox/filesystems/miniofilesystem"
	"github.com/FernandoJVideira/velox/filesystems/s3filesystem"
	"github.com/FernandoJVideira/velox/filesystems/sftpfilesystem"
	"github.com/FernandoJVideira/velox/filesystems/webdavfilesystem"
	"github.com/FernandoJVideira/velox/mailer"

	"github.com/dgraph-io/badger/v3"
	"github.com/robfig/cron/v3"

	"github.com/CloudyKit/jet/v6"
	"github.com/FernandoJVideira/velox/cache"
	"github.com/FernandoJVideira/velox/render"
	"github.com/FernandoJVideira/velox/session"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

// Velox Version
const version = "1.0.0"

var redisCache *cache.RedisCache
var badgerCache *cache.BadgerCache
var redisPool *redis.Pool
var badgerConn *badger.DB

var maintenanceMode bool

// Velox is the overall struct for the framework. Members are exported so they can be used by the user.
type Velox struct {
	AppName       string
	Debug         bool
	Version       string
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	RootPath      string
	Routes        *chi.Mux
	Render        *render.Render
	Session       *scs.SessionManager
	DB            Database
	JetViews      *jet.Set
	config        config
	EncryptionKey string
	Cache         cache.Cache
	Scheduler     *cron.Cron
	Mail          mailer.Mail
	Server        Server
	FileSystems   map[string]interface{}
	S3            s3filesystem.S3
	SFTP          sftpfilesystem.SFTP
	WebDAV        webdavfilesystem.WebDAV
	Minio         miniofilesystem.Minio
}

type Server struct {
	ServerName string
	Port       string
	Secure     bool
	URL        string
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    dbConfig
	redis       redisConfig
	uploads     uploadConfig
}

type uploadConfig struct {
	allowedMimeTypes []string
	maxUploadSize    int64
}

// New reads the .env file, creates our application config, populates the Velox type with settings
// based on .env values, and creates necessary folders and files if they don't exist
func (v *Velox) New(rootPath string) error {
	//Create folder structure if it doesn't exist
	pathConfig := initPaths{
		RootPath:    rootPath,
		FolderNames: []string{"handlers", "migrations", "views", "mail", "data", "public", "tmp", "logs", "middleware", "screenshots"},
	}

	err := v.Init(pathConfig)
	if err != nil {
		return err
	}
	// Create .env file if it doesn't exist
	err = v.checkDotEnv(rootPath)
	if err != nil {
		return err
	}
	//Read .env file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}
	//Create loggers
	infoLog, errorLog := v.startLoggers()
	//Connect to database
	if os.Getenv("DATABASE_TYPE") != "" {
		db, err := v.OpenDb(os.Getenv("DATABASE_TYPE"), v.BuildDSN())
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		v.DB = Database{
			DbType: os.Getenv("DATABASE_TYPE"),
			Pool:   db,
		}
	}

	scheduler := cron.New()
	v.Scheduler = scheduler

	if os.Getenv("CACHE") == "redis" || os.Getenv("SESSION_TYPE") == "redis" {
		redisCache = v.createClientRedisCache()
		v.Cache = redisCache
		redisPool = redisCache.Conn
	}

	if os.Getenv("CACHE") == "badger" {
		badgerCache = v.createClientBadgerCache()
		v.Cache = badgerCache
		badgerConn = badgerCache.Conn

		_, err := v.Scheduler.AddFunc("@daily", func() {
			_ = badgerCache.Conn.RunValueLogGC(0.7)
		})
		if err != nil {
			return err
		}
	}

	//Populate Velox struct
	v.InfoLog = infoLog
	v.ErrorLog = errorLog
	v.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	v.Version = version
	v.RootPath = rootPath
	v.Mail = v.createMailer()
	v.Routes = v.routes().(*chi.Mux)

	//File Uploads
	exploded := strings.Split(os.Getenv("ALLOWED_FILETYPES"), ",")
	var mimeTypes []string

	for _, m := range exploded {
		mimeTypes = append(mimeTypes, m)
	}

	var maxUploadSize int64

	if max, err := strconv.Atoi(os.Getenv("MAX_UPLOAD_SIZE")); err != nil {
		maxUploadSize = 10 << 20
	} else {
		maxUploadSize = int64(max)
	}

	// Set config
	v.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSISTS"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: dbConfig{
			dsn:      v.BuildDSN(),
			database: os.Getenv("DATABASE_NAME"),
		},
		redis: redisConfig{
			host:     os.Getenv("REDIS_HOST"),
			password: os.Getenv("REDIS_PASS"),
			prefix:   os.Getenv("REDIS_PREFIX"),
		},
		uploads: uploadConfig{
			maxUploadSize:    maxUploadSize,
			allowedMimeTypes: mimeTypes,
		},
	}

	// Server config
	secure := true
	if strings.ToLower(os.Getenv("SECURE")) == "false" {
		secure = false
	}

	v.Server = Server{
		ServerName: os.Getenv("SERVER_NAME"),
		Port:       os.Getenv("PORT"),
		Secure:     secure,
		URL:        os.Getenv("APP_URL"),
	}

	// create session

	sess := session.Session{
		CookieLifetime: v.config.cookie.lifetime,
		CookiePersist:  v.config.cookie.persist,
		CookieName:     v.config.cookie.name,
		SessionType:    v.config.sessionType,
		CookieDomain:   v.config.cookie.domain,
	}

	switch v.config.sessionType {
	case "redis":
		sess.RedisPool = redisCache.Conn
	case "mysql", "postgres", "postgresql", "mariadb":
		sess.DBPool = v.DB.Pool
	default:
		// Idk
	}

	v.Session = sess.InitSession()
	v.EncryptionKey = os.Getenv("KEY")

	// Jet views
	if v.Debug {
		var views = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
			jet.InDevelopmentMode(),
		)
		v.JetViews = views
	} else {
		var views = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		)
		v.JetViews = views
	}

	// Create renderer
	v.CreateRenderer()
	v.FileSystems = v.createFileSystems()
	go v.Mail.ListenForMail()

	return nil
}

// Init creates the folder structure
func (v *Velox) Init(p initPaths) error {
	root := p.RootPath
	for _, path := range p.FolderNames {
		// Create folder if it doesn't exist
		err := v.CreateFolderIfNotExists(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

// checkDotEnv checks if the .env file exists, if not it creates it
func (v *Velox) checkDotEnv(path string) error {
	err := v.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

// startLoggers creates the loggers
func (v *Velox) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

// CreateRenderer creates the renderer
func (v *Velox) CreateRenderer() {
	rend := render.Render{
		Renderer: v.config.renderer,
		RootPath: v.RootPath,
		Port:     v.config.port,
		JetViews: v.JetViews,
		Session:  v.Session,
	}
	v.Render = &rend
}

func (v *Velox) createMailer() mailer.Mail {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	m := mailer.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Templates:   v.RootPath + "/mail",
		Host:        os.Getenv("SMTP_HOST"),
		Port:        port,
		Username:    os.Getenv("SMTP_USERNAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		Encryption:  os.Getenv("SMTP_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
		Jobs:        make(chan mailer.Message, 20),
		Results:     make(chan mailer.Result, 20),
		API:         os.Getenv("MAILER_API"),
		APIKey:      os.Getenv("MAILER_KEY"),
		APIUrl:      os.Getenv("MAILER_URL"),
	}
	return m
}

func (v *Velox) createClientBadgerCache() *cache.BadgerCache {
	badgerCache := cache.BadgerCache{
		Conn: v.createBadgerConn(),
	}
	return &badgerCache
}

func (v *Velox) createClientRedisCache() *cache.RedisCache {
	cacheClient := cache.RedisCache{
		Conn:   v.createRedisPool(),
		Prefix: v.config.redis.prefix,
	}
	return &cacheClient
}

func (v *Velox) createRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				v.config.redis.host,
				redis.DialPassword(v.config.redis.password))
		},

		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}

func (v *Velox) createBadgerConn() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions(v.RootPath + "/tmp/badger"))
	if err != nil {
		return nil
	}
	return db
}

// BuildDSN builds the datasource name for our database, and returns it as a string
func (v *Velox) BuildDSN() string {
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"))

		// we check to see if a database password has been supplied, since including "password=" with nothing
		// after it sometimes causes postgres to fail to allow a connection.
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}

	case "mysql", "mariadb":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?collation=utf8_unicode_ci&timeout=5s&parseTime=true&tls=%s&readTimeout=5s",
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_PASS"),
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"))

	default:

	}

	return dsn
}

func (v *Velox) createFileSystems() map[string]interface{} {
	fileSystems := make(map[string]interface{})

	if os.Getenv("MINIO_SECRET") != "" {
		useSSL := false
		if strings.ToLower(os.Getenv("MINIO_USESSL")) == "true" {
			useSSL = true
		}

		minio := miniofilesystem.Minio{
			Endpoint: os.Getenv("MINIO_ENDPOINT"),
			Key:      os.Getenv("MINIO_KEY"),
			Secret:   os.Getenv("MINIO_SECRET"),
			UseSSL:   useSSL,
			Region:   os.Getenv("MINIO_REGION"),
			Bucket:   os.Getenv("MINIO_BUCKET"),
		}

		fileSystems["MINIO"] = minio
		v.Minio = minio
	}

	if os.Getenv("SFTP_HOST") != "" {
		sftp := sftpfilesystem.SFTP{
			Host: os.Getenv("SFTP_HOST"),
			User: os.Getenv("SFTP_USER"),
			Pass: os.Getenv("SFTP_PASS"),
			Port: os.Getenv("SFTP_PORT"),
		}
		fileSystems["SFTP"] = sftp
		v.SFTP = sftp
	}

	if os.Getenv("WEBDAV_HOST") != "" {
		webDav := webdavfilesystem.WebDAV{
			Host: os.Getenv("WEBDAV_HOST"),
			User: os.Getenv("WEBDAV_USER"),
			Pass: os.Getenv("WEBDAV_PASS"),
		}
		fileSystems["WEBDAV"] = webDav
		v.WebDAV = webDav
	}

	if os.Getenv("S3_KEY") != "" {
		s3 := s3filesystem.S3{
			Key:      os.Getenv("S3_KEY"),
			Secret:   os.Getenv("S3_SECRET"),
			Region:   os.Getenv("S3_REGION"),
			Endpoint: os.Getenv("S3_ENDPOINT"),
			Bucket:   os.Getenv("S3_BUCKET"),
		}
		fileSystems["S3"] = s3
		v.S3 = s3
	}

	return fileSystems
}

type RPCServer struct{}

func (r *RPCServer) MaintenanceMode(inMaintenanceMode bool, resp *string) error {
	if inMaintenanceMode {
		maintenanceMode = true
		*resp = "Server is in maintenance mode"
	} else {
		maintenanceMode = false
		*resp = "Server Live!"
	}
	return nil
}

func (v *Velox) listenRPC() {
	//if nothing specified for RPC port, don't start the server
	if os.Getenv("RPC_PORT") != "" {
		v.InfoLog.Println("Starting RPC server on port", os.Getenv("RPC_PORT"))
		err := rpc.Register(new(RPCServer))
		if err != nil {
			v.ErrorLog.Println(err)
			return
		}
		listen, err := net.Listen("tcp", "127.0.0.1:"+os.Getenv("RPC_PORT"))
		if err != nil {
			v.ErrorLog.Println(err)
			return
		}

		for {
			rpcConn, err := listen.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(rpcConn)
		}
	}
}
