package list

import "github.com/HuckOps/forge/model"

type PushGatewayResponse struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

type NodePushGatewayResult struct {
	model.Node   `bson:",inline"`
	PushGateways []model.PushGateway `bson:"pushgateway"`
}
