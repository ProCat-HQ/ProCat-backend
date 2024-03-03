package main

import (
	"github.com/procat-hq/procat-backend/internal/app/handler"
	"github.com/procat-hq/procat-backend/internal/app/server"
	"log"
)

func main() {
	srv := new(server.Server)

	handlers := new(handler.Handler)

	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("Error while running server %s", err.Error())
	}
}
