package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tumbnail/internal/handlers"
	"tumbnail/internal/repository"
	"tumbnail/internal/services"
	"tumbnail/server"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewRedisClient(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Password: os.Getenv("DB_PASSWORD"),
		DB:       viper.GetInt("db.db"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	service := services.NewService(repos)
	handlers := handlers.NewHandler(service)

	srv := server.NewServer(handlers)
	go func() {
		if err := srv.Run(viper.GetString("port")); err != nil {
			log.Fatalf("error occurred while running gRPC server: %s", err.Error())
		}
	}()

	log.Print("Thumbnail-server Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Thumbnail-server Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occurred on server shutting down: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
