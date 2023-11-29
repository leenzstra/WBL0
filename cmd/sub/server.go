package main

import (
	"os"

	"github.com/leenzstra/WBL0/config"
	"github.com/leenzstra/WBL0/internal/app"

)

func main() {
	config := config.New(os.Getenv("POSTGRES_URL"), os.Getenv("NATS_URL"))
	app.Run(&config)
}

