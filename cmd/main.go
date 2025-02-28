package main

import (
	"log"
	"onepercentdev_server/internal/app"
)

func main() {
	app, err := app.NewApp()

	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

}
