package model

import (
	"github.com/HuckOps/forge/server/repository/generic"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/net/context"
	"time"
)

type Node struct {
	BaseModel `bson:",inline"`

	UUID string `json:"uuid" bson:"uuid"`

	HostName string `json:"hostname" bson:"hostname"`
	IP       string `json:"ip" bson:"ip"`

	HeartBeat       time.Time `json:"heartbeat" bson:"heartbeat"`
	HeartBeatStatus bool      `json:"heartbeat_status" bson:"heartbeat_status"`
}

func (model *Node) TableName() string {
	return "node"
}

func (model *Node) Repository() *generic.Repository[Node] {
	return generic.NewRepository[Node]("node",
		generic.WithSoftDelete[Node](),
	)
}

func (model *Node) Indexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "uuid", Value: 1},
				//bson.E{Key: "hostname", Value: 1},
				bson.E{Key: "ip", Value: 1},
				bson.E{Key: "deleted_at", Value: -1},
			},
			Options: options.Index().SetName("idx_uuid_ip_deleteted_at_unique").SetUnique(true),
		},
	}
}

func (model *Node) FindByUUID(ctx context.Context, uuid string) (*Node, error) {
	err := model.Repository().Collection().FindOne(ctx, bson.M{"uuid": uuid}).Decode(model)
	return model, err
}

type NodeLabel struct {
	NodeID  bson.ObjectID `bson:"node_id"`
	LabelID bson.ObjectID `bson:"label_id"`
}

func (model *NodeLabel) Repository() *generic.Repository[NodeLabel] {
	return generic.NewRepository[NodeLabel]("node_label")
}
