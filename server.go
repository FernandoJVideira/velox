package velox

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func (v *Velox) ListenAndServe() error {
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

	go v.listenRPC()

	v.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	return srv.ListenAndServe()
}
