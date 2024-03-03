package velox

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/justinas/nosurf"
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

func (v *Velox) CheckForMaintenanceMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlEnvVar := os.Getenv("ALLOWED_URLS")
		allowedURLs := strings.Split(urlEnvVar, ",")

		if maintenanceMode {
			if !strings.Contains(r.URL.Path, "/public/maintenance.html") {
				if !sliceItemStartsWith(r.URL.Path, allowedURLs) {
					w.WriteHeader(http.StatusServiceUnavailable)
					w.Header().Set("Retry-After", "300")
					w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
					http.ServeFile(w, r, fmt.Sprintf("%s/public/maintenance.html", v.RootPath))
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func sliceItemStartsWith(s string, slice []string) bool {
	for _, item := range slice {
		if strings.HasPrefix(s, item) {
			return true
		}
	}
	return false
}
