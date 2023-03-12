package main

import (
	"flag"
	"notification-service/internal/app"
)

var (
	configPath = flag.String("config", "./config/config.yaml", "path to config file. default: config.yaml")
)

func main() {
	flag.Parse()

	application := app.NewApp(*configPath)
	application.Run()
}
