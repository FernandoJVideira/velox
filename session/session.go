package session

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

type Session struct {
	CookieLifetime string
	CookiePersist  string
	CookieName     string
	CookieDomain   string
	CookieSecure   string
	SessionType    string
	DBPool         *sql.DB
	RedisPool      *redis.Pool
}

func (v *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	//How long the session will last
	minutes, err := strconv.Atoi(v.CookieLifetime)
	if err != nil {
		minutes = 60
	}

	//Should the session persist after the browser is closed
	if strings.ToLower(v.CookiePersist) == "true" {
		persist = true
	}

	//Must Cookie be secure
	if strings.ToLower(v.CookieSecure) == "true" {
		secure = true
	}

	//Create the session
	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Name = v.CookieName
	session.Cookie.Secure = secure
	session.Cookie.Domain = v.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	//Which session type to use
	switch strings.ToLower(v.SessionType) {
	case "redis":
		session.Store = redisstore.New(v.RedisPool)
	case "mysql", "mariadb":
		session.Store = mysqlstore.New(v.DBPool)
	case "postgresql", "postgres":
		session.Store = postgresstore.New(v.DBPool)
	default:
		//Cookie
	}

	return session
}
