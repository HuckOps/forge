package prometheus

import (
	"context"
	"errors"
	"github.com/HuckOps/forge/model"
	"github.com/HuckOps/forge/server/repository/pagination"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type Pushgateway struct {
	model.PushGateway `bson:",inline" json:",inline"`
	UUID              string `json:"uuid" bson:"uuid"`
	IP                string `json:"ip" bson:"ip"`
}

func GetPushGatewayByPagination(ctx context.Context, skip, limit int) (pagination.PaginationResult[Pushgateway], error) {
	pipline := bson.A{
		bson.M{
			"$match": bson.M{
				"deleted_at": bson.M{"$eq": nil},
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "node",
				"localField":   "node",
				"foreignField": "_id",
				"as":           "node",
			},
		},
		bson.M{
			"$unwind": "$node",
		},
		bson.M{
			"$addFields": bson.M{
				"uuid": "$node.uuid",
				"ip":   "$node.ip",
				"node": "$node._id",
			},
		},
		bson.M{
			"$facet": bson.M{
				"data": bson.A{
					bson.M{"$skip": skip},
					bson.M{"$limit": limit},
				},
				"totalCount": bson.A{
					bson.M{"$count": "count"},
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"total": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$totalCount.count", 0}},
						0,
					},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"data":  1,
				"total": 1,
				"_id":   0, // 可选，移除MongoDB默认的_id字段
			},
		},
	}
	result := pagination.PaginationResult[Pushgateway]{}
	cursor, err := (&model.PushGateway{}).Repository().Collection().Aggregate(ctx, pipline)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)
	for !cursor.Next(ctx) {
		return result, nil
	}

	if err := cursor.Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}

func CreatePushGateway(ctx context.Context, gw model.PushGateway, uuid string) error {
	result, err := (&model.Node{}).FindByUUID(ctx, uuid)
	if err != nil {
		return errors.New("node not found")
	}
	gw.Node = result.ID
	gw.CreatedAt = time.Now()
	gw.UpdatedAt = time.Now()
	_, err = (&model.PushGateway{}).Repository().Create(ctx, gw)
	return err
}

func DeletePushGateway(ctx context.Context, id bson.ObjectID) error {
	return (&model.PushGateway{}).Repository().Delete(ctx, id)
}

func GetPushGatewayByID(ctx context.Context, id bson.ObjectID) (*Pushgateway, error) {
	gw, err := (&model.PushGateway{}).Repository().FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	node, err := (&model.Node{}).Repository().FindByID(ctx, gw.Node)
	if err != nil {
		return nil, err
	}
	return &Pushgateway{
		PushGateway: *gw,
		UUID:        node.UUID,
		IP:          node.IP,
	}, nil
}

func UpdatePushGateway(ctx context.Context, id bson.ObjectID, gw model.PushGateway) error {
	_, err := (&model.PushGateway{}).Repository().Update(ctx, id, bson.M{
		"port":    gw.Port,
		"version": gw.Version,
	})
	return err
}
