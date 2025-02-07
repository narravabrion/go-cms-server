package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/narravabrion/go-cms-server/internal/db"
	"github.com/narravabrion/go-cms-server/internal/env"
	"github.com/narravabrion/go-cms-server/internal/store"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}
	config := config{
		addr: env.GetStringEnv("ADDR", ":8081"),
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
		log.Panic(err)
	}
	log.Println("connected to Db!")
	store := store.NewStrorage(db)
	api := &api{
		config: config,
		store:  store,
	}

	log.Fatal(api.run(api.muxHandler()))
}
