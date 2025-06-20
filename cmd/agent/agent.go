package main

import (
	"github.com/HuckOps/forge/agent"
	"github.com/HuckOps/forge/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化配置
	config.InitAgentConfig("./config/agent.yaml")

	// 注册到注册中心
	agent.Registry()

	// 启动心跳服务
	agent.StartHeartBeat()

	// 设置优雅退出
	setupGracefulShutdown()
}

// 设置优雅退出机制
func setupGracefulShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 阻塞等待退出信号
	sig := <-sigCh
	log.Printf("收到信号: %v，开始关闭程序...", sig)

	// 停止心跳服务
	agent.StopHeartBeat()

	log.Println("程序已优雅退出")
	os.Exit(0)
}
