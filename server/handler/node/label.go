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

func GetLabels(ctx *gin.Context) {
	skip := ctx.DefaultQuery("skip", "0")
	limit := ctx.DefaultQuery("limit", "10")
	filterKeys := []string{"name", "code"}

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

	result, err := node.GetNodeLabelsByPagination(ctx, skipInt, limitInt, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, restful.Restful[map[string]interface{}]{
			Code: restful.SearchError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[pagination.PaginationResult[model.Label]]{
		Data: result,
		Code: restful.Success,
	})

}

func CreateLabel(ctx *gin.Context) {
	req := &model.Label{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}
	label := model.Label{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}
	err := node.CreateLabel(ctx, label)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, restful.Restful[map[string]interface{}]{
			Code: restful.InsertError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[map[string]interface{}]{
		Code: restful.Success,
	})

}

func GetLabelDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := node.GetLabelById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, restful.Restful[map[string]interface{}]{
			Code: restful.SearchError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[*model.Label]{
		Code: restful.Success,
		Data: result,
	})
}

func GetLabelNodes(ctx *gin.Context) {
	id := ctx.Param("id")
	skip := ctx.DefaultQuery("skip", "0")
	limit := ctx.DefaultQuery("limit", "10")
	skipInt, err1 := strconv.Atoi(skip)
	limitInt, err2 := strconv.Atoi(limit)
	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  "skip or limit unsupported type",
		})
		return
	}

	result, err := node.GetLabelNodesByIDWithPagination(ctx.Request.Context(), id, skipInt, limitInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, restful.Restful[map[string]interface{}]{
			Code: restful.SearchError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[*pagination.PaginationResult[model.Node]]{
		Code: restful.Success,
		Data: result,
	})
}
