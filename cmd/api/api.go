package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/narravabrion/go-cms-server/docs"
	"github.com/narravabrion/go-cms-server/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type api struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
	mail   mailConfig
}

type mailConfig struct {
	exp time.Duration
}

type dbConfig struct {
	connString   string
	maxOpenConns int
	maxIdleConns int
	maxIdleTIme  time.Duration
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
		docsURL := fmt.Sprintf("%s/swagger/doc.json", api.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", api.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(api.postContextMiddleware)
				r.Get("/", api.getPostHandler)
				r.Delete("/", api.deletePostHandler)
				r.Patch("/", api.updatePostHandler)
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}",api.activateUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(api.userContextMiddleware)
				r.Get("/", api.getUserHandler)
				r.Delete("/", api.deleteUserHandler)
				r.Patch("/", api.updateUserHandler)
				r.Put("/follow", api.followUserHandler)
				r.Put("/unfollow", api.unfollowUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Get("/feed", api.getUserFeedHandler)
			})
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/user", api.registerUserHandler)
		})
	})

	return r
}

func (api *api) run(mux http.Handler) error {

	// Docs
	docs.SwaggerInfo.Version = "1.0.1"
	docs.SwaggerInfo.Host = api.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	server := &http.Server{
		Addr:         api.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  3 * time.Minute,
	}
	api.logger.Infow("server started", "addr", api.config.addr, "env", api.config.env)
	return server.ListenAndServe()
}
