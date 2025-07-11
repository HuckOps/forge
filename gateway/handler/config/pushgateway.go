package config

import (
	"fmt"
	"github.com/HuckOps/forge/db"
	"github.com/HuckOps/forge/model"
	"github.com/HuckOps/forge/server/common/restful"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/net/context"
	"log"
	"time"
)

func GetPushGatewayConfig(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(400, restful.Restful[[]PushGatewayConfig]{
			Code: restful.SearchError,
			Msg:  "uuid parameter is required",
		})
		return
	}

	// 查询缓存是否命中
	cacheKey := fmt.Sprintf("gateway_%s", uuid)
	data, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedResult []PushGatewayConfig
		if err := msgpack.Unmarshal([]byte(data), &cachedResult); err == nil {
			ctx.JSON(200, restful.Restful[[]PushGatewayConfig]{
				Code: restful.Success,
				Data: cachedResult,
			})
			return
		}
		log.Printf("Failed to unmarshal cached data for uuid %s: %v", uuid, err)
	} else if err != redis.Nil { // 忽略键不存在的错误
		log.Printf("Redis error for uuid %s: %v", uuid, err)
	}

	// 未命中则从数据库加载
	node, err := (&model.Node{}).FindByUUID(ctx, uuid)
	if err != nil {
		log.Printf("Failed to find node with uuid %s: %v", uuid, err)
		ctx.JSON(400, restful.Restful[[]PushGatewayConfig]{
			Code: restful.SearchError,
			Msg:  "failed to find node",
		})
		return
	}

	var result []PushGatewayConfig
	pushgws, err := (&model.PushGateway{}).Repository().FindByFilter(
		ctx.Request.Context(), bson.M{"node": node.ID})
	if err != nil {
		log.Printf("Failed to query push gateways for node %s: %v", uuid, err)
		ctx.JSON(400, restful.Restful[[]PushGatewayConfig]{
			Code: restful.SearchError,
			Msg:  "failed to query database",
		})
		return
	}

	for _, pushgw := range pushgws {
		result = append(result, PushGatewayConfig{
			Port:    pushgw.Port,
			Version: pushgw.Version,
		})
	}

	// 异步更新缓存
	go func() {
		data, err := msgpack.Marshal(result)
		if err != nil {
			log.Printf("Failed to marshal data for caching (uuid %s): %v", uuid, err)
			return
		}
		if err := db.RedisClient.Set(
			context.Background(), // 使用新的context，因为原请求可能已结束
			cacheKey,
			data,
			15*time.Minute,
		).Err(); err != nil {
			log.Printf("Failed to cache data for uuid %s: %v", uuid, err)
		}
	}()

	ctx.JSON(200, restful.Restful[[]PushGatewayConfig]{
		Code: restful.Success,
		Data: result,
	})
}
