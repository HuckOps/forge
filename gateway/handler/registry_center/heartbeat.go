package registry_center

import (
	"encoding/json"
	"github.com/HuckOps/forge/mq"
	"github.com/gin-gonic/gin"
)

func Heartbeat(ctx *gin.Context) {
	uuid := ctx.Query("uuid")
	if uuid == "" {
		ctx.JSON(400, gin.H{})
	}
	msg := mq.HeartBeatMessage{
		UUID: uuid,
	}
	m, _ := json.Marshal(msg)
	mq.RabbitMQClient.PublishMessage(mq.HEARTBEAT, m, false)
	ctx.JSON(200, gin.H{})
}
