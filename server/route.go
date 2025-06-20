package server

import (
	"github.com/HuckOps/forge/server/handler/node"
	"github.com/gin-gonic/gin"
)

func RegistryRouter(e *gin.Engine) {
	NodeManagementRouter(e)
}

func NodeManagementRouter(e *gin.Engine) {
	routerGroup := e.Group("/api/v1")
	{
		routerGroup.GET("/nodes", node.GetNodeList)
	}
}
