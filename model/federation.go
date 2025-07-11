package model

import (
	"github.com/HuckOps/forge/server/repository/generic"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Federation struct {
	BaseModel `bson:",inline"`
	Version   string        `bson:"version" json:"version"`
	Port      int           `bson:"port" json:"port"`
	NodeID    bson.ObjectID `bson:"node_id" json:"node_id"`
}

func (model *Federation) Repository() *generic.Repository[Federation] {
	return generic.NewRepository[Federation]("federation",
		generic.WithSoftDelete[Federation](),
	)
}
func (model *Federation) Indexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "node_id", Value: 1},
				bson.E{Key: "port", Value: 1},
				bson.E{Key: "deleted_at", Value: -1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
}
