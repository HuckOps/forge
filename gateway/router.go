package gateway

import (
	"github.com/HuckOps/forge/gateway/handler/config"
	"github.com/HuckOps/forge/gateway/handler/list"
	"github.com/HuckOps/forge/gateway/handler/registry_center"
	"github.com/gin-gonic/gin"
)

func RegistryRouter(e *gin.Engine) {
	RegistryCenter(e)
	Config(e)
}

func RegistryCenter(e *gin.Engine) {
	r := e.Group("/api")
	{
		r.POST("/register", registry_center.Registry)
		r.GET("/heartbeat", registry_center.Heartbeat)
	}
}

func Config(e *gin.Engine) {
	r := e.Group("/config")
	{
		r.GET("/pushgateway", list.GetPushGatewayList)
		r.GET("/pushgateway/:uuid", config.GetPushGatewayConfig)
	}
}
