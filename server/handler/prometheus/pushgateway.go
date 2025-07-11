package prometheus

import (
	"github.com/HuckOps/forge/model"
	"github.com/HuckOps/forge/server/common/restful"
	"github.com/HuckOps/forge/server/logic/prometheus"
	"github.com/HuckOps/forge/server/repository/pagination"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
	"strconv"
)

func GetPushGateways(ctx *gin.Context) {
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
	result, err := prometheus.GetPushGatewayByPagination(ctx.Request.Context(), skipInt, limitInt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.SearchError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[pagination.PaginationResult[prometheus.Pushgateway]]{
		Code: restful.Success,
		Data: result,
	})
}

func CreatePushGateway(ctx *gin.Context) {
	req := CreatePushGatewayRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}
	gw := model.PushGateway{
		Port:    req.Port,
		Version: req.Version,
	}
	err := prometheus.CreatePushGateway(ctx.Request.Context(), gw, req.UUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.InsertError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, restful.Restful[map[string]interface{}]{
		Code: restful.Success,
	})

}

func DeletePushGateway(ctx *gin.Context) {
	id := ctx.Param("id")
	idHex, err := bson.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}
	err = prometheus.DeletePushGateway(ctx.Request.Context(), idHex)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[map[string]interface{}]{
		Code: restful.Success,
		Data: nil,
	})
}

func SearchPushGateway(ctx *gin.Context) {
	id := ctx.Param("id")
	idHex, err := bson.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}
	gw, err := prometheus.GetPushGatewayByID(ctx.Request.Context(), idHex)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[*prometheus.Pushgateway]{
		Code: restful.Success,
		Data: gw,
	})
}

func UpdatePushGateway(ctx *gin.Context) {
	id := ctx.Param("id")
	idHex, err := bson.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}
	req := UpdatePushGatewayRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}

	gw := model.PushGateway{
		Port:    req.Port,
		Version: req.Version,
	}

	err = prometheus.UpdatePushGateway(ctx.Request.Context(), idHex, gw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[map[string]interface{}]{
		Code: restful.Success,
		Data: nil,
	})
}
