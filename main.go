package main

import (
	"log"

	"github.com/septalfauzan/saku-api/app/server"
	"github.com/septalfauzan/saku-api/app/server/config"
	"github.com/septalfauzan/saku-api/app/server/datasources"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Println("Warning: .env file not found", err)
	}

	ds := datasources.NewDatasources(
		conf.GeminiAPIKey,
		conf.GeminiAPIUrl,
		conf.GeminiOCRPrompt,
		conf.GeminiOCRModel,
	)

	port := conf.Port
	if port == "" {
		port = "3000"
	}

	app := server.NewServer(ds)

	log.Fatal(app.Listen(":" + port))
}
