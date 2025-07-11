package pushgateway

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/HuckOps/forge/agent"
	"github.com/HuckOps/forge/config"
	gw_config "github.com/HuckOps/forge/gateway/handler/config"
	"github.com/HuckOps/forge/internal/logger"
	"github.com/HuckOps/forge/server/common/restful"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"
	"time"
)

type PushGateway struct {
	Path    string
	Port    int
	Version string
	Pid     int
	Ctx     context.Context
	Cancel  context.CancelFunc
	Cmd     *exec.Cmd
}

func NewPushGateway(ctx context.Context, path string, port int, version string) *PushGateway {
	ctx, cancel := context.WithCancel(ctx)
	return &PushGateway{
		Path:    path,
		Port:    port,
		Version: version,
		Ctx:     ctx,
		Cancel:  cancel,
	}
}

func (p *PushGateway) Deploy() error {
	//	判断文件是否存在，不存在则下载
	if _, err := os.Stat(p.Path); os.IsNotExist(err) {
		if err := os.MkdirAll(p.Path, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// 判断是否已经部署过

	cmd := exec.Command(fmt.Sprintf("%s/pushgateway-%s.linux-amd64/pushgateway", p.Path, p.Version), "--version")
	output, err := cmd.Output()
	if err == nil {
		if version, err := extractVersion(string(output)); err == nil {
			if version == p.Version {
				logger.Logger.Info(fmt.Sprintf("Version verify success, current version: %v", p.Version))
				return nil
			} else {
				logger.Logger.Warn(fmt.Sprintf("Version mismatch: expected %s, got %s", p.Version, version))
			}
		} else {
			logger.Logger.Warn("Failed to extract version from output")
		}
	} else {
		logger.Logger.Warn("Failed to execute version command")
	}

	url := fmt.Sprintf("https://gh-proxy.net/github.com/prometheus/pushgateway/releases/download/v%s/pushgateway-%s.linux-amd64.tar.gz", p.Version, p.Version)

	tmpFile := fmt.Sprintf("/tmp/pushgateway-%s.linux-amd64.tar.gz", p.Version)
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer out.Close()

	logger.Logger.Info(
		fmt.Sprintf("Downloading %s to %s, version: %s", url, tmpFile, p.Version))
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	file, err := os.Open(tmpFile)
	if err != nil {
		return fmt.Errorf("open tmpfile failed: %v", err)
	}
	defer file.Close()

	logger.Logger.Info(
		"Downloaded file success")

	logger.Logger.Info(
		fmt.Sprintf("Unzipping %s to %s", tmpFile, p.Path))
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(p.Path, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, p.Path); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown typeflag %v", header.Typeflag)

		}
	}

	logger.Logger.Info(fmt.Sprintf("Unzipping %s to %s success", tmpFile, p.Path))

	logger.Logger.Info(fmt.Sprintf("Verify downloaded file version"))

	cmd = exec.Command(fmt.Sprintf("%s/pushgateway-%s.linux-amd64/pushgateway", p.Path, p.Version), "--version")
	output, err = cmd.Output()
	if err != nil {
		return err
	}
	version, err := extractVersion(string(output))
	if err != nil {
		return err
	}

	if version != p.Version {
		logger.Logger.Error(fmt.Sprintf("Pushgateway version does not match %s != %s", p.Version, version))
		return fmt.Errorf("pushgateway version mismatch. Expected %s, got %s", p.Version, version)
	}

	logger.Logger.Info(fmt.Sprintf("Pushgateway version is %s, verify success", version))

	return nil
}

func extractVersion(output string) (string, error) {
	// 定义匹配版本号的正则表达式
	re := regexp.MustCompile(`version (\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", fmt.Errorf("无法从输出中提取版本号")
	}
	return matches[1], nil
}

func (p *PushGateway) ExecutePushGateway() error {
	cmd := exec.CommandContext(p.Ctx,
		fmt.Sprintf("%s/pushgateway-%s.linux-amd64/pushgateway", p.Path, p.Version),
		fmt.Sprintf("--web.listen-address=:%d", p.Port),
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("get stdout pipline error: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("get stderr pipline error: %w", err)
	}

	logger.Logger.Info(fmt.Sprintf("Start pushgateway process, port: %d", p.Port))
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start command error: %w", err)
	}

	p.Pid = cmd.Process.Pid
	p.Cmd = cmd

	go logger.ScanOutput(logger.Logger, stdout, "stdout",
		zap.Int("pid", p.Pid),
		zap.String("svc", "pushgateway"),
	)
	go logger.ScanOutput(logger.Logger, stderr, "stderr",
		zap.Int("pid", p.Pid),
		zap.String("svc", "pushgateway"),
	)

	// 🛑 不再调用 cmd.Wait()，交由 StopPushGateway 来处理
	return nil
}

var (
	pushGatewayWg   sync.WaitGroup
	PushGatewayList []*PushGateway
	listMutex       sync.Mutex
)

func StartPushGatewayCron(ctx context.Context) {
	uuid := agent.GetOrGenUUID()
	pushGatewayWg.Add(1)
	go func() {
		defer pushGatewayWg.Done()
		pushgatewayLoop(ctx, uuid)
	}()

}

func pushgatewayLoop(parentCtx context.Context, uuid string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-parentCtx.Done():
			logger.Logger.Info("PushGateway同步循环收到停止信号")
			return
		case <-ticker.C:
			err := syncPushGateway(parentCtx, uuid)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func StopPushGateway() {
	listMutex.Lock()
	defer listMutex.Unlock()

	if len(PushGatewayList) == 0 {
		logger.Logger.Info("No PushGateway instances to stop")
		return
	}

	logger.Logger.Info("Stopping all PushGateway instances",
		zap.Int("count", len(PushGatewayList)),
	)

	var wg sync.WaitGroup
	for _, pg := range PushGatewayList {
		wg.Add(1)
		go func(p *PushGateway) {
			defer wg.Done()

			// 1. 首先取消Context
			p.Cancel()

			// 2. 发送SIGTERM信号
			if p.Cmd != nil && p.Cmd.Process != nil {
				logger.Logger.Info("Sending SIGTERM to PushGateway",
					zap.Int("pid", p.Pid),
					zap.Int("port", p.Port),
				)

				if err := p.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
					logger.Logger.Error("Failed to send SIGTERM",
						zap.Int("pid", p.Pid),
						zap.Error(err),
					)
				}

				// 3. 等待进程退出
				_, err := p.Cmd.Process.Wait()
				logger.Logger.Info("PushGateway process stopped",
					zap.Int("pid", p.Pid),
					zap.Int("port", p.Port),
					zap.Error(err),
				)
			}
		}(pg)
	}

	// 等待所有进程停止或超时
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Logger.Info("All PushGateway instances stopped successfully")
	case <-time.After(15 * time.Second):
		logger.Logger.Warn("Timeout while waiting for some PushGateway instances to stop")
	}

	// 清空列表
	PushGatewayList = nil
}

func syncPushGateway(ctx context.Context, uuid string) error {
	url := fmt.Sprintf("%s/config/pushgateway/%s", config.AgentConfig.RegistryCenter, uuid)
	logger.Logger.Info(fmt.Sprintf("Start sync push gateway from %s, uuid: %s", url, uuid))

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	data := restful.Restful[[]gw_config.PushGatewayConfig]{}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	logger.Logger.Info(fmt.Sprintf("Sync push gateway from %s, result: %v", url, data.Data))

	listMutex.Lock()
	defer listMutex.Unlock()

	newList := []*PushGateway{}

	for _, gwConfig := range data.Data {
		found := false

		// 查找是否已存在对应的网关
		for _, existing := range PushGatewayList {
			if existing.Port == gwConfig.Port && existing.Version == gwConfig.Version {
				newList = append(newList, existing)
				found = true
				break
			}
		}

		if !found {
			pushGateway := NewPushGateway(ctx, "/home/huck/pushgateway", gwConfig.Port, gwConfig.Version)
			logger.Logger.Info(fmt.Sprintf("Start new gateway from %s, version: %s, port: %d",
				pushGateway.Path, gwConfig.Version, gwConfig.Port))
			if err := pushGateway.Deploy(); err == nil {
				newList = append(newList, pushGateway)
				go pushGateway.ExecutePushGateway()
			}

		}
	}

	// 停止未在新配置中的网关
	for _, oldGateway := range PushGatewayList {
		keep := false
		for _, newGateway := range newList {
			if oldGateway.Port == newGateway.Port && oldGateway.Version == newGateway.Version {
				keep = true
				break
			}
		}
		if !keep {
			logger.Logger.Info(fmt.Sprintf("Stopping obsolete gateway: version: %s, port: %d", oldGateway.Version, oldGateway.Port))
			//go oldGateway.StopPushGateway() // 假设你有这个方法
			err := oldGateway.Cmd.Process.Signal(syscall.SIGTERM)
			if err != nil {
				logger.Logger.Error("Failed to send SIGTERM",
					zap.Int("pid", oldGateway.Pid),
					zap.Error(err))
			}
			go func() {
				_, err = oldGateway.Cmd.Process.Wait()
				if err != nil {
					logger.Logger.Error("Failed to wait for SIGTERM",
						zap.Int("pid", oldGateway.Pid),
						zap.Error(err))
				}
			}()

		}
	}

	PushGatewayList = newList
	return nil
}
