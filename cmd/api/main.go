package main

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/narravabrion/go-cms-server/internal/db"
	"github.com/narravabrion/go-cms-server/internal/env"
	"github.com/narravabrion/go-cms-server/internal/store"
	"go.uber.org/zap"
)

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
		addr: env.GetStringEnv("ADDR", ":8081"),
		apiURL: env.GetStringEnv("EXTERNAL_URL", "localhost:8081"),
		db: dbConfig{
			connString:   env.GetStringEnv("CONN_STRING", "postgres://postgres:password@localhost/go_cms?sslmode=disable"),
			maxOpenConns: env.GetIntEnv("DB_MAX_OPEN_CONNS", 20),
			maxIdleConns: env.GetIntEnv("DB_MAX_IDLE_CONNS", 10),
			maxIdleTIme:  env.GetTimeEnv("DB_MAX_IDLE_TIME", 15*time.Minute),
		},
	}

	db, err := db.New(
		config.db.connString,
		config.db.maxOpenConns,
		config.db.maxIdleConns,
		config.db.maxIdleTIme,
	)
	defer db.Close()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("connected to Db!")
	store := store.NewStrorage(db)
	api := &api{
		config: config,
		store:  store,
		logger: logger,
	}

	logger.Fatal(api.run(api.muxHandler()))
}
