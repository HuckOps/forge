package main

import (
	"context"
	"fmt"
	"github.com/HuckOps/forge/config"
	"github.com/HuckOps/forge/db"
	"github.com/HuckOps/forge/server"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config.InitServerConfig("./config/server.yaml")
	db.InitMongo(ctx, config.ServerConfig.MongoURL)
	r := gin.Default()
	server.RegistryRouter(r)
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", config.ServerConfig.Port),

		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
