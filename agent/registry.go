package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/HuckOps/forge/config"
	"github.com/HuckOps/forge/internal/logger"
	uuid2 "github.com/google/uuid"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"runtime"
)

func GetOrGenUUID() string {
	uuidPath := ""
	switch runtime.GOOS {
	case "darwin":
		uuidPath = "/etc/forge/uuid"
	case "linux":
		uuidPath = "/etc/forge/uuid"
	case "windows":
		uuidPath = "c:\\forge\\uuid"
	default:
		panic("Unsupported os platform")
	}

	data, err := os.ReadFile(uuidPath)
	uuid := ""
	if os.IsNotExist(err) {
		u, _ := uuid2.NewUUID()
		uuid = u.String()
		err := os.WriteFile(uuidPath, []byte(u.String()), 0666)
		if err != nil {
			panic(err)
		}
	} else {
		uuid = string(data)
	}

	_, err = uuid2.Parse(uuid)
	if err != nil {
		panic(err)
	}
	return uuid
}

func getPrimaryIP() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		// 跳过回环接口和非活动接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP
			// 选择全局单播IPv4地址
			if ip.IsGlobalUnicast() && ip.To4() != nil {
				return ip, nil
			}
		}
	}
	return nil, fmt.Errorf("no valid IPv4 address found")
}
func Registry() {
	ip, err := getPrimaryIP()
	if err != nil {
		panic(err)
	}
	hostname, _ := os.Hostname()
	uuid := GetOrGenUUID()
	// 注册agent
	client := http.Client{}
	body := map[string]string{
		"uuid":     uuid,
		"hostname": hostname,
		"ip":       ip.String(),
	}
	logger.Logger.Info("Registry starting...",
		zap.String("uuid", uuid),
		zap.String("hostname", hostname),
		zap.String("ip", ip.String()))
	jsonData, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/register", config.AgentConfig.RegistryCenter),
		bytes.NewBuffer(jsonData))
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		logger.Logger.Error(fmt.Sprintf("Registry failed to register with registry status code: %v",
			resp), zap.Error(err))
		panic("Failed to register forge agent")
	}
	defer resp.Body.Close()

}
