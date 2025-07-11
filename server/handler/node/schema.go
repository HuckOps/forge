package node

type CreateLabelRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}

type NodeLabelRequest struct {
	Nodes  []string `json:"nodes" binding:"required"`
	Labels []string `json:"labels" binding:"required"`
}
