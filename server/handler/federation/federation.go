package federation

import (
	"github.com/HuckOps/forge/model"
	"github.com/HuckOps/forge/server/common/restful"
	"github.com/HuckOps/forge/server/logic/federation"
	"github.com/HuckOps/forge/server/repository/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateFederation(ctx *gin.Context) {
	req := CreateFederationReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}
	fede := model.Federation{
		Port:    req.Port,
		Version: req.Version,
	}
	err := federation.CreateFederation(ctx, fede, req.UUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[map[string]interface{}]{
		Code: restful.Success,
	})
}

func GetFederationList(ctx *gin.Context) {
	skip := ctx.DefaultQuery("skip", "0")
	limit := ctx.DefaultQuery("limit", "10")
	skipInt, err1 := strconv.Atoi(skip)
	limitInt, err2 := strconv.Atoi(limit)
	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusBadRequest, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  "skip or limit invalid",
		})
		return
	}
	result, err := federation.GetFederationByPagination(ctx.Request.Context(), skipInt, limitInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, restful.Restful[map[string]interface{}]{
			Code: restful.RequestError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, restful.Restful[pagination.PaginationResult[federation.Federation]]{
		Code: restful.Success,
		Data: result,
	})

}
