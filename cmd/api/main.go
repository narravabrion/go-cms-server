package main

import (
	"expvar"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/narravabrion/go-cms-server/internal/auth"
	"github.com/narravabrion/go-cms-server/internal/db"
	"github.com/narravabrion/go-cms-server/internal/env"
	"github.com/narravabrion/go-cms-server/internal/mailer"
	"github.com/narravabrion/go-cms-server/internal/ratelimiter"
	"github.com/narravabrion/go-cms-server/internal/store"
	"github.com/narravabrion/go-cms-server/internal/store/cache"
	"go.uber.org/zap"
)

const version = "0.0.1"
//	@title			go cms server
//	@version		1.0
//	@description	Thi is a simple blog cms.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	err := godotenv.Load("../../.env")
	if err != nil {
		logger.Fatal(err)
	}
	config := config{
		addr:   env.GetStringEnv("ADDR", ":8081"),
		apiURL: env.GetStringEnv("EXTERNAL_URL", "localhost:8081"),
		frontEndURL: env.GetStringEnv("FRONTEND_URL", "http://127.0.0.1:3000"),
		db: dbConfig{
			connString:   env.GetStringEnv("CONN_STRING", "postgres://postgres:password@localhost/go_cms?sslmode=disable"),
			maxOpenConns: env.GetIntEnv("DB_MAX_OPEN_CONNS", 20),
			maxIdleConns: env.GetIntEnv("DB_MAX_IDLE_CONNS", 10),
			maxIdleTIme:  env.GetTimeEnv("DB_MAX_IDLE_TIME", 15*time.Minute),
		},
		redisConfig: redisConfig{
			addr: env.GetStringEnv("REDIS_ADDR","127.0.0.1:6379"),
			password: env.GetStringEnv("REDIS_PASSWORD",""),
			db: env.GetIntEnv("REDIS_DB",0),
			enabled: false,
		},
		env: env.GetStringEnv("ENV", "development"),
		mail: mailConfig{
			exp: 3 * 3 * time.Hour,
			fromEmail: env.GetStringEnv("SENDGRID_FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey:    env.GetStringEnv("SENDGRID_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetStringEnv("AUTH_BASIC_USER", "admin"),
				password: env.GetStringEnv("AUTH_BASIC_PASSWORD", "admin"),
			},
			token: tokenConfig{
				secret: env.GetStringEnv("AUTH_TOKEN_SECRET", "somesecret"),
				exp: time.Hour *24 * 3,
				iss: "go-cms",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestPerTimeFrame: env.GetIntEnv("RATELIMITER_REQUESTS_COUNT", 20),
			TimeFrame: time.Second*5,
			Enabled: true,
		},
	}

	db, err := db.New(
		config.db.connString,
		config.db.maxOpenConns,
		config.db.maxIdleConns,
		config.db.maxIdleTIme,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("connected to Db!")

	var redisDB *redis.Client
	if config.redisConfig.enabled {
		redisDB = cache.NewRedisClient(config.redisConfig.addr, config.redisConfig.password, config.redisConfig.db)
		logger.Info("redis connection established")
	}

	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		config.rateLimiter.RequestPerTimeFrame,
		config.rateLimiter.TimeFrame,
	)

	store := store.NewStrorage(db)
	cacheStorage := cache.NewRedisStorage(redisDB)

	mailer := mailer.NewSendGRid(config.mail.sendGrid.apiKey, config.mail.fromEmail)

	jwtAuthenticator := auth.NewJWTAuthenticator(config.auth.token.secret, config.auth.token.iss, config.auth.token.iss)
	api := &api{
		config: config,
		store:  store,
		logger: logger,
		mailer: mailer,
		authenticator: jwtAuthenticator,
		cacheStorage: cacheStorage,
		rateLimiter: rateLimiter,
	}

	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	logger.Fatal(api.run(api.muxHandler()))
}
