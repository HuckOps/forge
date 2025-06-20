package consumer

import (
	"context"
	"encoding/json"
	"github.com/HuckOps/forge/db"
	"github.com/HuckOps/forge/model"
	"github.com/HuckOps/forge/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"
	"time"
)

func Heartbeat(d amqp.Delivery) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	msg := &mq.HeartBeatMessage{}
	if err := json.Unmarshal(d.Body, msg); err != nil {
		log.Println("json unmarshal err:", err)
	}
	node := model.Node{}
	_, err := db.MongoDB.Collection(node.TableName()).UpdateOne(ctx, bson.M{"uuid": msg.UUID}, bson.M{"$set": bson.M{
		"heartbeat":        time.Now(),
		"heartbeat_status": true,
	}})
	if err != nil {
		log.Println("update failed:", err)
	}
}
