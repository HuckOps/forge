package consumer

import (
	"context"
	"encoding/json"
	"github.com/HuckOps/forge/db"
	"github.com/HuckOps/forge/model"
	"github.com/HuckOps/forge/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

func Registry(d amqp.Delivery) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msg := &mq.RegisteyMessage{}
	if err := json.Unmarshal(d.Body, msg); err != nil {
		log.Println("json unmarshal err:", err)
		return
	}

	node := model.Node{}
	opts := options.UpdateOne().SetUpsert(true)

	update := bson.M{
		"$set": bson.M{
			"hostname":         msg.Hostname,
			"ip":               msg.IP,
			"heartbeat":        time.Now(),
			"heartbeat_status": true,
		},
		"$setOnInsert": bson.M{
			"uuid":       msg.UUID,
			"created_at": time.Now(),
			"updated_at": time.Now(),
		},
	}

	_, err := db.MongoDB.Collection(node.TableName()).UpdateOne(ctx, bson.M{"uuid": msg.UUID}, update, opts)
	if err != nil {
		log.Fatalln("update failed:", err)
	}
}
