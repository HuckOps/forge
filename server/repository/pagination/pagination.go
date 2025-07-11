package pagination

import (
	"context"
	"github.com/HuckOps/forge/server/repository/generic"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PaginationQuery[T any] struct {
	Repository *generic.Repository[T]
	Filter     bson.M
	Skip       int
	Limit      int
	Sort       bson.D
}

type PaginationResult[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
}

func GetByPagination[T any](ctx context.Context, query PaginationQuery[T]) (PaginationResult[T], error) {
	result := PaginationResult[T]{}

	if query.Repository.SoftDelete {
		query.Filter["deleted_at"] = bson.M{"$exists": false}
	}
	// 获取总数
	total, err := query.Repository.Collection().CountDocuments(ctx, query.Filter)
	if err != nil {
		return result, err
	}
	result.Total = total

	// 设置查询选项
	findOptions := options.Find()
	if query.Limit != -1 {
		findOptions.SetSkip(int64(query.Skip))
		findOptions.SetLimit(int64(query.Limit))
	}
	findOptions.SetSort(query.Sort)

	// 执行查询
	cursor, err := query.Repository.Collection().Find(ctx, query.Filter, findOptions)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	// 解码结果
	var data []T
	if err := cursor.All(ctx, &data); err != nil {
		return result, err
	}

	if data == nil {
		data = make([]T, 0)
	}

	result.Data = data
	return result, nil
}
