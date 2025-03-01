package main

import (
	"context"
	"fmt"
	"github.com/valeraBerezovskij/logger/internal/config"
	"github.com/valeraBerezovskij/logger/internal/repository"
	"github.com/valeraBerezovskij/logger/internal/server"
	"github.com/valeraBerezovskij/logger/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	//Config init
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

	//Database init
	db := dbClient.Database(cfg.DB.Database)

	//repo, service, server init
	auditRepo := repository.NewAudit(db)
	auditService := service.NewAudit(auditRepo)
	auditSrv := server.NewAuditServer(auditService)
	srv := server.New(auditSrv)

	fmt.Println("Server started", time.Now())

	if err := srv.ListenAndServe(cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
