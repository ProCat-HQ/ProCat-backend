package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/procat-hq/procat-backend/internal/app/handler"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"github.com/procat-hq/procat-backend/internal/app/server"
	"github.com/procat-hq/procat-backend/internal/app/service"
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {

}

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Error while reading configs %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error while loading .env file: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.name"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		log.Fatalf("Error occured while init DB: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	bindAddr := viper.GetString("bind_addr")

	if err := srv.Run(bindAddr, handlers.InitRoutes()); err != nil {
		log.Fatalf("Error while running server %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
