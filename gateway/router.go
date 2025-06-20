package gateway

import (
	"github.com/HuckOps/forge/gateway/handler/registry_center"
	"github.com/gin-gonic/gin"
)

func RegistryRouter(e *gin.Engine) {
	RegistryCenter(e)
}

func RegistryCenter(e *gin.Engine) {
	r := e.Group("/api")
	{
		r.POST("/register", registry_center.Registry)
		r.GET("/heartbeat", registry_center.Heartbeat)
	}
}
