package node

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"

	//"github.com/HuckOps/forge/db"
	"github.com/HuckOps/forge/model"
	"github.com/HuckOps/forge/server/repository/pagination"
)

func GetNodeLabelsByPagination(ctx context.Context, skip, limit int, filter bson.M) (pagination.PaginationResult[model.Label], error) {
	query := pagination.PaginationQuery[model.Label]{
		Repository: (&model.Label{}).Repository(),
		Skip:       skip,
		Limit:      limit,
		Filter:     filter,
		Sort:       bson.D{{"created_at", -1}},
	}
	return pagination.GetByPagination[model.Label](ctx, query)
}

func CreateLabel(ctx context.Context, label model.Label) error {
	label.CreatedAt = time.Now()
	label.UpdatedAt = time.Now()
	_, err := (&model.Label{}).Repository().Create(ctx, label)
	return err
}

func GetLabelById(ctx context.Context, id string) (*model.Label, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return (&model.Label{}).Repository().FindByID(ctx, objID)
}

func GetLabelNodesByIDWithPagination(
	ctx context.Context,
	id string,
	skip, limit int,
) (*pagination.PaginationResult[model.Node], error) {
	// 转换ID
	idHex, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	// 构建聚合管道
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"label_id":   idHex,
				"deleted_at": bson.M{"$eq": nil},
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
			"$replaceRoot": bson.M{
				"newRoot": "$node",
			},
		},
		bson.M{
			"$facet": bson.M{
				"data": bson.A{
					bson.M{"$skip": skip},
					bson.M{"$limit": limit},
				},
				"total": bson.A{
					bson.M{"$count": "count"},
				},
			},
		},
	}

	// 定义正确的接收结构
	var result struct {
		Data  []model.Node `bson:"data"`
		Total []struct {
			Count int64 `bson:"count"`
		} `bson:"total"`
	}

	// 执行聚合查询
	cursor, err := (&model.NodeLabel{}).Repository().Collection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation failed: %w", err)
	}
	defer cursor.Close(ctx)

	// 解码结果
	if !cursor.Next(ctx) {
		return &pagination.PaginationResult[model.Node]{
			Data:  []model.Node{},
			Total: 0,
		}, nil
	}

	if err := cursor.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode result: %w", err)
	}

	// 处理可能没有结果的情况
	totalCount := int64(0)
	if len(result.Total) > 0 {
		totalCount = result.Total[0].Count
	}

	return &pagination.PaginationResult[model.Node]{
		Data:  result.Data,
		Total: totalCount,
	}, nil
}
