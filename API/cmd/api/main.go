package main

import (
	"log"

	"github.com/spf13/viper"
	"ryanlawton.art/photospace-api/config"
	"ryanlawton.art/photospace-api/server"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	app := server.NewApp()

	if err := app.Run(viper.GetString("port")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
