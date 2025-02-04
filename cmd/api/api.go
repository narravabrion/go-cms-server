package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type api struct {
	config config
}

type config struct {
	addr string
}

func (api *api) muxHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", api.healthCheckHandler)
	})

	return r
}

func (api *api) run(mux http.Handler) error {

	server := &http.Server{
		Addr:         api.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  3 * time.Minute,
	}
	log.Printf("server started on port %s", api.config.addr)
	return server.ListenAndServe()
}
