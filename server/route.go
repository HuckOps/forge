package server

import (
	"github.com/HuckOps/forge/server/handler/federation"
	"github.com/HuckOps/forge/server/handler/node"
	"github.com/HuckOps/forge/server/handler/prometheus"
	"github.com/gin-gonic/gin"
)

func RegistryRouter(e *gin.Engine) {
	NodeManagementRouter(e)
	LabelManagementRouter(e)
	PrometheusRouter(e)
	FederationRouter(e)
}

func NodeManagementRouter(e *gin.Engine) {
	routerGroup := e.Group("/api/v1/nodes")
	{
		routerGroup.GET("", node.GetNodeList)
		routerGroup.POST("/labels", node.SetNodeLabel)
	}
}

func LabelManagementRouter(e *gin.Engine) {
	labelGroup := e.Group("/api/v1/labels")
	{
		labelGroup.GET("", node.GetLabels)
		labelGroup.POST("", node.CreateLabel)

		labelGroup.GET("/:id", node.GetLabelDetail)
		labelGroup.GET("/:id/nodes", node.GetLabelNodes)
	}
}

func PrometheusRouter(e *gin.Engine) {
	prometheusGroup := e.Group("/api/v1/prometheus")
	{
		prometheusGroup.GET("/pushgateway", prometheus.GetPushGateways)
		prometheusGroup.POST("/pushgateway", prometheus.CreatePushGateway)
		prometheusGroup.DELETE("/pushgateway/:id", prometheus.DeletePushGateway)
		prometheusGroup.GET("/pushgateway/:id", prometheus.SearchPushGateway)
		prometheusGroup.PUT("/pushgateway/:id", prometheus.UpdatePushGateway)
	}
}

func FederationRouter(e *gin.Engine) {
	routerGroup := e.Group("/api/v1/prometheus")
	{
		routerGroup.POST("/federation", federation.CreateFederation)
		routerGroup.GET("/federation", federation.GetFederationList)
	}
}
