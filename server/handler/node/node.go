package node

import (
	"github.com/HuckOps/forge/server/common/restful"
	"github.com/HuckOps/forge/server/logic/node"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
	"strconv"
)

func GetNodeList(ctx *gin.Context) {
	skip := ctx.DefaultQuery("skip", "0")
	limit := ctx.DefaultQuery("limit", "10")
	filterKeys := []string{"uuid", "hostname", "ip", "heartbeat_status"}

	filter := bson.M{}
	for _, key := range filterKeys {
		value := ctx.Query(key)
		if value != "" {
			filter[key] = value
		}
	}

	skipInt, err1 := strconv.Atoi(skip)
	limitInt, err2 := strconv.Atoi(limit)
	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  "skip or limit unsupported type",
		})
		return
	}
	nodes, total, err := node.GetNodesByPagination(ctx.Request.Context(), skipInt, limitInt, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, restful.Restful[map[string]interface{}]{
			Code: restful.SearchError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[map[string]interface{}]{
		Data: gin.H{
			"nodes": nodes,
			"total": total,
		},
		Code: restful.Success,
	})
}
