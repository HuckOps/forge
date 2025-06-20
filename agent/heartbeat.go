package agent

import (
	"context"
	"fmt"
	"github.com/HuckOps/forge/config"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	heartbeatCtx    context.Context
	heartbeatCancel context.CancelFunc
	heartbeatWg     sync.WaitGroup
)

func init() {
	heartbeatCtx, heartbeatCancel = context.WithCancel(context.Background())
}

// StartHeartBeat 启动心跳服务
func StartHeartBeat() {
	uuid := GetOrGenUUID()
	log.Printf("启动心跳服务，UUID: %s", uuid)

	heartbeatWg.Add(1)
	go func() {
		defer heartbeatWg.Done()
		heartbeatLoop(uuid)
	}()
}

// StopHeartBeat 停止心跳服务
func StopHeartBeat() {
	log.Println("停止心跳服务...")
	heartbeatCancel()
	heartbeatWg.Wait()
	log.Println("心跳服务已停止")
}

// 心跳主循环
func heartbeatLoop(uuid string) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// 初始发送一次心跳
	if err := sendHeartBeat(uuid); err != nil {
		log.Printf("初始心跳失败: %v", err)
	}

	for {
		select {
		case <-heartbeatCtx.Done():
			log.Println("收到停止信号，退出心跳循环")
			return
		case <-ticker.C:
			if err := sendHeartBeat(uuid); err != nil {
				log.Printf("心跳发送失败: %v", err)

				// 失败后重试
				go func() {
					time.Sleep(2 * time.Second)
					if err := sendHeartBeat(uuid); err != nil {
						log.Printf("心跳重试失败: %v", err)
					} else {
						log.Println("心跳重试成功")
					}
				}()
			}
		}
	}
}

// 发送心跳请求
func sendHeartBeat(uuid string) error {
	url := fmt.Sprintf("%s/api/heartbeat?uuid=%s", config.AgentConfig.RegistryCenter, uuid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求发送失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("非预期状态码: %d", resp.StatusCode)
	}

	log.Println("心跳发送成功")
	return nil
}
