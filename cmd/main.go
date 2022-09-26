package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	server "homeworkdeliverysystem"
	"homeworkdeliverysystem/pkg/handler"
	"homeworkdeliverysystem/pkg/repository"
	"homeworkdeliverysystem/pkg/service"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}

	ds, err := repository.InitDS()
	if err != nil {
		log.Fatalf("Failed to initialize dbs: %s", err.Error())
	}

	newRepository := repository.NewRepository(ds)
	newService := service.NewService(newRepository)
	handlers := handler.NewHandler(newService)

	srv := new(server.Server)
	if err := srv.Run(os.Getenv("SERVER_PORT"), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}
}
