package federation

import (
	"github.com/HuckOps/forge/model"
	"github.com/HuckOps/forge/server/repository/pagination"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/net/context"
	"time"
)

func CreateFederation(ctx context.Context, federation model.Federation, uuid string) error {
	node, err := (&model.Node{}).FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	federation.CreatedAt = time.Now()
	federation.UpdatedAt = time.Now()
	federation.NodeID = node.ID

	_, err = (&model.Federation{}).Repository().Create(ctx, federation)
	return err
}

type Federation struct {
	model.Federation `json:",inline" bson:",inline"`
	UUID             string `json:"uuid" bson:"uuid"`
	IP               string `json:"ip" bson:"ip"`
}

func GetFederationByPagination(ctx context.Context, skip, limit int) (pagination.PaginationResult[Federation], error) {
	pipline := bson.A{
		bson.M{
			"$match": bson.M{
				"deleted_at": bson.M{"$exists": false},
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "node",
				"localField":   "node_id",
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
	result := pagination.PaginationResult[Federation]{}
	cursor, err := (&model.Federation{}).Repository().Collection().Aggregate(ctx, pipline)
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
