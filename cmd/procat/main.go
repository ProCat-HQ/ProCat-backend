package main

import (
	"github.com/procat-hq/procat-backend/internal/app/handler"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"github.com/procat-hq/procat-backend/internal/app/server"
	"github.com/procat-hq/procat-backend/internal/app/service"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Error while reading configs %s", err.Error())
	}
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(server.Server)

	if err := srv.Run(viper.GetString("bind_addr"), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error while running server %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
