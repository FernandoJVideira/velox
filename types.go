package velox

import "database/sql"

// initPaths is a struct that holds the root path and the folder names
type initPaths struct {
	RootPath    string
	FolderNames []string
}

// config is a struct that holds the cookie configuration
type cookieConfig struct {
	name     string
	lifetime string
	presist  string
	secure   string
	domain   string
}

type dbConfig struct {
	dsn      string
	database string
}

type Database struct {
	DbType string
	Pool   *sql.DB
}

type redisConfig struct {
	host     string
	password string
	prefix   string
}
