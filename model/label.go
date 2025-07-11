package model

import (
	"github.com/HuckOps/forge/server/repository/generic"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Label struct {
	BaseModel   `bson:",inline"`
	Name        string `json:"name" bson:"name"`
	Code        string `json:"code" bson:"code"`
	Description string `json:"description" bson:"description"`
}

func (model *Label) Repository() *generic.Repository[Label] {
	return generic.NewRepository[Label]("label",
		generic.WithSoftDelete[Label](),
	)
}

func (model *Label) Indexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "code", Value: 1},
				bson.E{Key: "deleted_at", Value: -1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
}
