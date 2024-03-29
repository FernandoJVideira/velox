package velox

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (v *Velox) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	if v.Debug {
		mux.Use(middleware.Logger)
	}
	mux.Use(middleware.Recoverer)
	mux.Use(v.SessionLoad)
	mux.Use(v.NoSurf)
	mux.Use(v.CheckForMaintenanceMode)

	return mux
}

// Routes are velox specific routes, which are mounted in the routes file
// in Velox application
func Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/test-v", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Velox!"))
	})

	return r
}
