package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/HuckOps/forge/config"
	"github.com/HuckOps/forge/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"
	"os"
	"sync"
	"time"
)

var (
	gatewayConfigCtx    context.Context
	gatewayConfigCancel context.CancelFunc
	gatewayConfigWg     sync.WaitGroup
)

func init() {
	gatewayConfigCtx, gatewayConfigCancel = context.WithCancel(context.Background())
}

func StartGenGatewayConfig() {
	gatewayConfigWg.Add(1)
	go func() {
		defer gatewayConfigWg.Done()

		genGatewayConfigLoop(gatewayConfigCtx)
	}()

}

func genGatewayConfigLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			genGatewayConfig(ctx)
		}
	}
}

func genGatewayConfig(ctx context.Context) {
	gateways, err := (&model.PushGateway{}).Repository().FindByFilter(ctx, bson.M{})
	if err != nil {
		log.Printf("get gateways err: %v", err)
		return
	}

	log.Printf("get gateways total: %v", len(gateways))

	type GatewayConfig struct {
		Version string `json:"version"`
		Port    int    `json:"port"`
	}

	for _, gateway := range gateways {
		c := GatewayConfig{
			Version: gateway.Version,
			Port:    gateway.Port,
		}

		configBytes, err := json.Marshal(c)
		if err != nil {
			log.Printf("genGatewayConfigLoop err: %v", err)
			return
		}
		fmt.Println("%s/pushgateway/%s", config.WorkerConfig.ConfigDir, gateway.Node.Hex())
		err = os.WriteFile(
			fmt.Sprintf("%s/pushgateway/%s", config.WorkerConfig.ConfigDir, gateway.Node.Hex()),
			configBytes,
			0644,
		)
		fmt.Printf("genGatewayConfigLoop err: %v", err)
	}
}
