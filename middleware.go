package velox

import (
	"github.com/justinas/nosurf"
	"net/http"
	"strconv"
)

func (v *Velox) SessionLoad(next http.Handler) http.Handler {
	return v.Session.LoadAndSave(next)
}

func (v *Velox) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(v.config.cookie.secure)

	csrfHandler.ExemptGlob("/api/*")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   v.config.cookie.domain,
	})

	return csrfHandler
}
