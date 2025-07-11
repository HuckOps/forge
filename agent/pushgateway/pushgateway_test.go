package pushgateway

import (
	"fmt"
	"github.com/HuckOps/forge/internal/logger"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"syscall"
	"testing"
	"time"
)

func init() {
	logger.InitLogger()
}

func TestPushGateway_Deploy(t *testing.T) {
	gw := PushGateway{
		Path:    "/home/huck/pushgateway",
		Version: "1.11.1",
	}
	err := gw.Deploy()
	if err != nil {
		t.Errorf("PushGateway.Deploy() error = %v", err)
	}
	assert.Nil(t, err)
}

func TestPushGateway_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gw := PushGateway{
		Path:    "/home/huck/pushgateway",
		Version: "1.11.1",
		Port:    9999,
		Ctx:     ctx,
	}
	go func() {
		err := gw.ExecutePushGateway()
		if err != nil {
			t.Errorf("PushGateway.ExecutePushGateway() error = %v", err)
		}
		assert.Nil(t, err)
	}()
	time.Sleep(10 * time.Second)
	gw.Cmd.Process.Signal(syscall.SIGTERM)
	fmt.Println(gw.Cmd.Process.Wait())

}
