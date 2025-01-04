package webserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vs0uz4/rate-limit/config"
)

func Start(cfg *config.Config) error {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	log.Printf("Starting server on port %s...", cfg.WebServerPort)
	return http.ListenAndServe(":"+cfg.WebServerPort, r)
}
