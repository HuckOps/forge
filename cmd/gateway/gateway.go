package main

import (
	"context"
	"fmt"
	"github.com/HuckOps/forge/config"
	"github.com/HuckOps/forge/gateway"
	"github.com/HuckOps/forge/mq"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config.InitGatewayConfig("./config/gateway.yaml")
	mq.InitAMQP(ctx, config.GatewayConfig.AMQP)
	r := gin.Default()
	gateway.RegistryRouter(r)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.GatewayConfig.Host, config.GatewayConfig.Port),
		Handler: r,
		// 关键优化参数
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,

		// 启用 Keep-Alive
		ConnState: func(conn net.Conn, state http.ConnState) {
			if state == http.StateIdle {
				// 设置空闲连接超时
				conn.SetDeadline(time.Now().Add(120 * time.Second))
			}
		},
	}
	server.SetKeepAlivesEnabled(true)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
