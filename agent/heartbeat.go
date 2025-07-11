package agent

import (
	"context"
	"fmt"
	"github.com/HuckOps/forge/config"
	"github.com/HuckOps/forge/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

var (
	heartbeatWg sync.WaitGroup
)

// StartHeartBeat starts the heartbeat loop with the provided context.
func StartHeartBeat(ctx context.Context) {
	uuid := GetOrGenUUID()
	logger.Logger.Info("Starting heartbeat service") // Optional structured field
	// zap.String("uuid", uuid),

	heartbeatWg.Add(1)
	go func() {
		defer heartbeatWg.Done()
		heartbeatLoop(ctx, uuid)
	}()
}

// StopHeartBeat waits for the heartbeat loop to exit.
func StopHeartBeat() {
	logger.Logger.Info("Waiting for heartbeat service to stop...")
	heartbeatWg.Wait()
	logger.Logger.Info("Heartbeat service stopped")
}

// heartbeatLoop is the main heartbeat routine.
func heartbeatLoop(ctx context.Context, uuid string) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Send initial heartbeat
	if err := sendHeartBeat(uuid); err != nil {
		logger.Logger.Warn("Initial heartbeat failed",
			// zap.String("uuid", uuid),
			zap.Error(err),
		)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Logger.Info("Heartbeat loop received shutdown signal")
			return
		case <-ticker.C:
			if err := sendHeartBeat(uuid); err != nil {
				logger.Logger.Error("Failed to send heartbeat",
					zap.Error(err),
				)

				// Retry once asynchronously
				go func() {
					time.Sleep(2 * time.Second)
					if err := sendHeartBeat(uuid); err != nil {
						logger.Logger.Error("Heartbeat retry failed",
							zap.Error(err),
						)
					} else {
						logger.Logger.Info("Heartbeat retry succeeded")
					}
				}()
			}
		}
	}
}

// sendHeartBeat sends a single heartbeat request.
func sendHeartBeat(uuid string) error {
	url := fmt.Sprintf("%s/api/heartbeat?uuid=%s", config.AgentConfig.RegistryCenter, uuid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Logger.Error("Failed to create heartbeat request",
			zap.Error(err),
		)
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Error("Heartbeat request failed",
			zap.Error(err),
		)
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Error("Heartbeat request returned non-200 status",
			// zap.String("url", url),
			// zap.String("uuid", uuid),
			zap.Int("status", resp.StatusCode),
		)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	logger.Logger.Info("Heartbeat sent successfully")
	return nil
}
