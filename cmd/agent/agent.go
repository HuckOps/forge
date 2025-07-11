package main

import (
	"github.com/HuckOps/forge/agent"
	"github.com/HuckOps/forge/agent/pushgateway"
	"github.com/HuckOps/forge/config"
	"github.com/HuckOps/forge/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize logger
	logger.InitLogger()
	logger.Logger.Info("Application starting...")

	// Load configuration
	config.InitAgentConfig("./config/agent.yaml")

	// Register to service registry
	agent.Registry()

	// Start heartbeat service
	agent.StartHeartBeat(ctx)

	// Start PushGateway cron job
	pushgateway.StartPushGatewayCron(ctx)

	// Setup graceful shutdown
	setupGracefulShutdown(ctx, cancel)
}

// Setup graceful shutdown mechanism
func setupGracefulShutdown(ctx context.Context, cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Block until we receive shutdown signal
	sig := <-sigCh
	logger.Logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	logger.Logger.Info("Initiating graceful shutdown...")

	// Cancel the main context to notify all components
	cancel()

	// Stop heartbeat service
	agent.StopHeartBeat()

	// Stop all PushGateway instances
	pushgateway.StopPushGateway()

	// Allow some time for cleanup
	select {
	case <-time.After(5 * time.Second):
		logger.Logger.Info("Cleanup timeout reached, exiting...")
	case <-ctx.Done():
		logger.Logger.Info("All components shutdown completed")
	}

	logger.Logger.Info("Application exited gracefully")
}
