package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/narravabrion/go-cms-server/internal/env"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("ADDR")
	fmt.Printf("the port is: %s", port)
	app := &api{
		config: config{
			addr: env.GetEnv("ADDR", ":8081"),
		},
	}

	log.Fatal(app.run(app.muxHandler()))
}
