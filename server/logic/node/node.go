package node

import (
	"context"
	"github.com/HuckOps/forge/db"
	"github.com/HuckOps/forge/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func GetNodesByPagination(ctx context.Context, skip, limit int, filter bson.M) ([]model.Node, int64, error) {
	n := model.Node{}
	collection := db.MongoDB.Collection(n.TableName())

	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	nodes := make([]model.Node, 0)
	if err = cursor.All(ctx, &nodes); err != nil {
		return nil, 0, err
	}

	if nodes == nil {
		nodes = make([]model.Node, 0)
	}

	return nodes, total, nil
}
