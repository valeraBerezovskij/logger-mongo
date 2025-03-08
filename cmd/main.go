package main

import (
	"context"
	"fmt"
	"github.com/valeraBerezovskij/logger-mongo/internal/config"
	"github.com/valeraBerezovskij/logger-mongo/internal/repository"
	"github.com/valeraBerezovskij/logger-mongo/internal/server"
	"github.com/valeraBerezovskij/logger-mongo/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Programm started", time.Now())

	//Config init
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	//Context init
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Mongo connection
	opts := options.Client()
	opts.SetAuth(options.Credential{
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
	})
	opts.ApplyURI(cfg.DB.URI)

	dbClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	if err := dbClient.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mongo started", time.Now())

	//Database init
	db := dbClient.Database(cfg.DB.Database)

	//repo, service, server init
	auditRepo := repository.NewAudit(db)
	auditService := service.NewAudit(auditRepo)
	auditSrv := server.NewAuditServer(auditService)
	srv, err := server.NewServer("amqp://guest:guest@localhost:5672/", "logs", auditSrv)
	if err != nil {
		log.Fatal("failed to initialize server: ", err)
	}
	defer srv.Close()

	fmt.Println("Server started", time.Now())

	srv.ConsumeMessages(ctx)
}
