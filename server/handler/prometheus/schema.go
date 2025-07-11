package prometheus

type CreatePushGatewayRequest struct {
	Version string `json:"version"`
	Port    int    `json:"port"`
	UUID    string `json:"uuid"`
}

type UpdatePushGatewayRequest struct {
	Port    int    `json:"port"`
	Version string `json:"version"`
}
