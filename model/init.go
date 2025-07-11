package model

import (
	"errors"
	"fmt"
	"github.com/HuckOps/forge/db"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/net/context"
	"reflect"
	"strings"
)

type Indexes interface {
	Indexes() []mongo.IndexModel
}

type TableName interface {
	TableName() string
}

func InitCollections(ctx context.Context) error {
	collectionList := []interface{}{
		&Node{}, &Label{}, &Exporter{}, &PushGateway{}, &NodeLabel{}, &Federation{},
	}
	for _, collection := range collectionList {
		collectionName := getCollectionName(ctx, collection)
		names, err := db.MongoDB.ListCollectionNames(ctx, bson.M{"name": collectionName})
		if err != nil {
			return err
		}
		collectionExists := false
		for _, n := range names {
			if n == collectionName {
				collectionExists = true
				break
			}
		}
		if !collectionExists {
			if err := db.MongoDB.CreateCollection(ctx, collectionName); err != nil {
				if !isNamespaceExistsError(err) {
					return err
				}
			}
		} else {
			continue
		}
		if m, ok := collection.(Indexes); ok {
			if _, err := db.MongoDB.Collection(collectionName).Indexes().CreateMany(ctx, m.Indexes()); err != nil {
				if !isIndexConflictError(err) {
					return fmt.Errorf("failed to create indexes: %v", err)
				}
			}
		}
	}
	return nil
}

func getCollectionName(ctx context.Context, collection interface{}) string {
	if m, ok := collection.(TableName); ok {
		return m.TableName()
	}
	t := reflect.TypeOf(collection)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return strings.ToLower(t.Name())
}

func isNamespaceExistsError(err error) bool {
	var cmdErr mongo.CommandError
	if errors.As(err, &cmdErr) {
		return cmdErr.Code == 48 // MongoDB 错误码 48: NamespaceExists
	}
	return false
}

func isIndexConflictError(err error) bool {
	return strings.Contains(err.Error(), "IndexOptionsConflict")
}
