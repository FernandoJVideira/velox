package velox

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
}

func (v *Velox) New(rootPath string) error {
	//Create folder structure if it doesn't exist
	pathConfig := initPaths{
		RootPath:    rootPath,
		FolderNames: []string{"handlers", "migrations", "views", "mail", "data", "public", "tmp", "logs", "middleware"},
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
		redisCache = v.CreateClientRedisCache()
		v.Cache = redisCache
		redisPool = redisCache.Conn
	}

	if os.Getenv("CACHE") == "badger" {
		badgerCache = v.CreateClientBadgerCache()
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

	// Set config
	v.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			presist:  os.Getenv("COOKIE_PERSISTS"),
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

	// Create session
	sess := session.Session{
		CookieLifetime: v.config.cookie.lifetime,
		CookiePersist:  v.config.cookie.presist,
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

// ListenAndServe starts the web server
func (v *Velox) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     v.ErrorLog,
		Handler:      v.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	if v.DB.Pool != nil {
		defer v.DB.Pool.Close()
	}

	if redisPool != nil {
		defer redisPool.Close()
	}

	if badgerConn != nil {
		defer badgerConn.Close()
	}

	v.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	v.ErrorLog.Fatal(err)
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

func (v *Velox) CreateClientBadgerCache() *cache.BadgerCache {
	badgerCache := cache.BadgerCache{
		Conn: v.createBadgerConn(),
	}
	return &badgerCache
}

func (v *Velox) CreateClientRedisCache() *cache.RedisCache {
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
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}
	default:
	}

	return dsn
}
