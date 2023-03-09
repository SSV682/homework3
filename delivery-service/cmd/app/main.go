package main

import (
	"delivery-service/internal/app"
	"flag"
)

var (
	configPath = flag.String("config", "./config/config.yaml", "path to config file. default: config.yaml")
)

func main() {
	flag.Parse()

	application := app.NewApp(*configPath)
	application.Run()
}
