package hardware

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCPUInfo(t *testing.T) {
	info, err := GetCPUInfo()
	assert.Nil(t, err)
	assert.NotEmpty(t, info)
}

func TestGetMemoryInfo(t *testing.T) {
	info, err := GetMemoryInfo()
	assert.Nil(t, err)
	assert.NotNil(t, info)
	fmt.Println(info)
}
