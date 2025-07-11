package config

import "encoding/json"

type PushGatewayConfig struct {
	Version string `json:"version" redis:"version"`
	Port    int    `json:"port" redis:"port"`
}

// 实现 encoding.BinaryMarshaler 接口
func (p *PushGatewayConfig) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

// 实现 encoding.BinaryUnmarshaler 接口
func (p *PushGatewayConfig) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
