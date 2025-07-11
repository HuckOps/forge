package model

import (
	"github.com/HuckOps/forge/server/repository/generic"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PushGateway struct {
	BaseModel `bson:",inline"`

	Port    int           `bson:"port" json:"port"`
	Version string        `bson:"version" json:"version"`
	Node    bson.ObjectID `bson:"node" json:"node"`
}

func (model *PushGateway) Repository() *generic.Repository[PushGateway] {
	return generic.NewRepository[PushGateway]("pushgateway",
		generic.WithSoftDelete[PushGateway](),
	)
}

func (model *PushGateway) Indexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "node", Value: 1},
				bson.E{Key: "port", Value: 1},
				bson.E{Key: "deleted_at", Value: -1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
}
