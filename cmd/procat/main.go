package main

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/procat-hq/procat-backend/internal/app/handler"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"github.com/procat-hq/procat-backend/internal/app/server"
	"github.com/procat-hq/procat-backend/internal/app/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error while reading configs %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error while loading .env file: %s", err.Error())
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
		logrus.Fatalf("Error occured while init DB: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	bindAddr := viper.GetString("bind_addr")

	go func() {
		if err := srv.Run(bindAddr, handlers.InitRoutes()); !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Error while running server %s", err.Error())
		}
	}()

	logrus.Printf("Server started on port %s", bindAddr)

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-quitSignal

	logrus.Printf("Server shuts down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Can't terminate server: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("Can't close DB connection: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("cmd/procat/config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
