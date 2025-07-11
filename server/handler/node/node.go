package node

import (
	"github.com/HuckOps/forge/model"
	"github.com/HuckOps/forge/server/common/restful"
	"github.com/HuckOps/forge/server/logic/node"
	"github.com/HuckOps/forge/server/repository/pagination"
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
	result, err := node.GetNodesByPagination(ctx.Request.Context(), skipInt, limitInt, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, restful.Restful[map[string]interface{}]{
			Code: restful.SearchError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[pagination.PaginationResult[model.Node]]{
		Data: result,
		Code: restful.Success,
	})
}

func SetNodeLabel(ctx *gin.Context) {
	req := NodeLabelRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}
	nodeIDHexs := []bson.ObjectID{}
	for _, id := range req.Nodes {
		nodeIDHex, err := bson.ObjectIDFromHex(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
				Code: restful.RequestError,
				Msg:  err.Error(),
			})
			return
		}
		nodeIDHexs = append(nodeIDHexs, nodeIDHex)
	}
	labelIDHexs := []bson.ObjectID{}
	for _, labelID := range req.Labels {
		labelIDHex, err := bson.ObjectIDFromHex(labelID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
				Code: restful.RequestError,
				Msg:  err.Error(),
			})
			return
		}
		labelIDHexs = append(labelIDHexs, labelIDHex)
	}
	err := node.SetNodeLabel(ctx.Request.Context(), nodeIDHexs, labelIDHexs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, restful.Restful[map[string]interface{}]{
			Code: restful.InsertError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[map[string]interface{}]{
		Code: restful.Success,
		Data: nil,
	})
}
