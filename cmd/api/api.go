package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/narravabrion/go-cms-server/docs"
	"github.com/narravabrion/go-cms-server/internal/auth"
	"github.com/narravabrion/go-cms-server/internal/mailer"
	"github.com/narravabrion/go-cms-server/internal/ratelimiter"
	"github.com/narravabrion/go-cms-server/internal/store"
	"github.com/narravabrion/go-cms-server/internal/store/cache"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type api struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	cacheStorage  cache.Storage
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontEndURL string
	auth        authConfig
	redisConfig redisConfig
	rateLimiter ratelimiter.Config
}

type redisConfig struct {
	addr     string
	password string
	db       int
	enabled  bool
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	user     string
	password string
}

type mailConfig struct {
	sendGrid  sendGridConfig
	exp       time.Duration
	fromEmail string
}

type dbConfig struct {
	connString   string
	maxOpenConns int
	maxIdleConns int
	maxIdleTIme  time.Duration
}

type sendGridConfig struct {
	apiKey string
}

func (api *api) muxHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)


	if api.config.rateLimiter.Enabled {
		r.Use(api.RateLimiterMiddleware)
	}

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.With(api.BasicAuthMiddleware()).Get("/health", api.healthCheckHandler)
		docsURL := fmt.Sprintf("%s/swagger/doc.json", api.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
		r.Route("/posts", func(r chi.Router) {
			r.Use(api.AuthTokenMiddleware)
			r.Post("/", api.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(api.postContextMiddleware)
				r.Get("/", api.getPostHandler)
				r.Delete("/", api.checkPostOwnership("admin", api.deletePostHandler))
				r.Patch("/", api.checkPostOwnership("moderator", api.updatePostHandler))
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", api.activateUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(api.AuthTokenMiddleware)
				// r.Use(api.userContextMiddleware)
				r.Get("/", api.getUserHandler)
				r.Delete("/", api.deleteUserHandler)
				r.Patch("/", api.updateUserHandler)
				r.Put("/follow", api.followUserHandler)
				r.Put("/unfollow", api.unfollowUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(api.AuthTokenMiddleware)
				r.Get("/feed", api.getUserFeedHandler)
			})
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/user", api.registerUserHandler)
			r.Post("/token", api.createTokenHandler)
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

	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		api.logger.Infow("signal caught", "signal", s.String())
		shutdown <- server.Shutdown(ctx)
	}()
	api.logger.Infow("server started", "addr", api.config.addr, "env", api.config.env)
	err := server.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdown
	if err != nil {
		return err
	}
	api.logger.Infow("server stopped", "addr", api.config.addr, "env", api.config.env)
	return nil
}
