package list

import (
	"fmt"
	"github.com/HuckOps/forge/internal/logger"
	"github.com/HuckOps/forge/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
	"net/http"
)

func GetPushGatewayList(ctx *gin.Context) {
	//(&model.PushGateway{}).Repository().Collection().Aggregate()
	resp := []PushGatewayResponse{}
	pipeline := bson.A{
		// 主表筛选未软删除的记录
		bson.D{
			{"$match", bson.D{
				{"deleted_at", bson.Null{}}, // 主表的 deleted_at 为 null
			}},
		},
		// 关联子表 pushgateway，并筛选未软删除的子表记录
		bson.D{
			{"$lookup", bson.D{
				{"from", "pushgateway"},
				{"let", bson.D{{"nodeId", "$_id"}}},
				{"pipeline", bson.A{
					bson.D{
						{"$match", bson.D{
							{"$expr", bson.D{
								{"$and", bson.A{
									bson.D{{"$eq", bson.A{"$node", "$$nodeId"}}},
									bson.D{{"$eq", bson.A{"$deleted_at", bson.TypeNull}}}, // deleted_at == null
								}},
							}},
						}},
					},
				}},
				{"as", "pushgateway"},
			}},
		},
	}

	cursor, err := (&model.Node{}).Repository().Collection().Aggregate(ctx.Request.Context(), pipeline)
	if err != nil {
		return
	}
	defer cursor.Close(ctx.Request.Context())
	for cursor.Next(ctx.Request.Context()) {
		var item NodePushGatewayResult
		if err := cursor.Decode(&item); err != nil {
			logger.Logger.Error("decode pushgateway result", zap.Error(err))
			return
		}
		if len(item.PushGateways) == 0 {
			continue
		}
		r := PushGatewayResponse{
			Labels: map[string]string{
				"pushgateway":      item.IP,
				"pushgateway_uuid": item.UUID,
			},
		}
		for _, pushgateway := range item.PushGateways {
			r.Targets = append(r.Targets, fmt.Sprintf("%s:%d", item.IP, pushgateway.Port))
		}
		resp = append(resp, r)
	}
	if err := cursor.Err(); err != nil {
		logger.Logger.Error("cursor error", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
