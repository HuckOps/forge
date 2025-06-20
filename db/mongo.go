package db

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoDB *mongo.Database

func InitMongo(ctx context.Context, mongoURL string) {
	client, err := mongo.Connect(options.Client().ApplyURI(
		mongoURL))
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	MongoDB = client.Database("forge")
}
