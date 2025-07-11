package federation

type CreateFederationReq struct {
	Version string `json:"version"`
	Port    int    `json:"port"`
	UUID    string `json:"uuid"`
}
