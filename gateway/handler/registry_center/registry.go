package registry_center

import (
	"encoding/json"
	"github.com/HuckOps/forge/mq"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegistryRequest struct {
	Hostname string `json:"hostname" binding:"required"`
	IP       string `json:"ip" binding:"required,ip"`
	UUID     string `json:"uuid" binding:"required,uuid"`
}

func Registry(ctx *gin.Context) {
	request := RegistryRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msg, _ := json.Marshal(mq.RegisteyMessage{Hostname: request.Hostname, UUID: request.UUID,
		IP: request.IP})
	err := mq.RabbitMQClient.PublishMessage(
		mq.REGISTRYCHANNEL,
		msg,
		false,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
