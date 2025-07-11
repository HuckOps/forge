package generic

import (
	"context"
	"fmt"
	"github.com/HuckOps/forge/db"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type Option[T any] func(repository *Repository[T])

func WithSoftDelete[T any]() Option[T] {
	return func(repository *Repository[T]) {
		repository.SoftDelete = true
	}
}

func NewRepository[T any](collection string, opt ...Option[T]) *Repository[T] {
	c := db.MongoDB.Collection(collection)
	repo := &Repository[T]{
		collection: c,
	}
	for _, opt := range opt {
		opt(repo)
	}
	return repo
}

type Repository[T any] struct {
	collection *mongo.Collection
	// 软删除开关
	SoftDelete bool
}

func (r *Repository[T]) Collection() *mongo.Collection {
	return r.collection
}

func (r *Repository[T]) FindByFilter(
	ctx context.Context,
	filter bson.M,
	opts ...options.Lister[options.FindOptions],
) ([]T, error) {
	if r.SoftDelete {
		filter["deleted_at"] = bson.M{"$exists": false}
	}
	cursor, err := r.Collection().Find(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("find failed: %w", err)
	}
	defer func() {
		if cursor != nil {
			cursor.Close(ctx)
		}
	}()

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("decoding results failed: %w", err)
	}

	return results, nil
}

func (r *Repository[T]) FindByID(ctx context.Context, id bson.ObjectID, opts ...options.Lister[options.FindOneOptions]) (*T, error) {
	data := new(T)
	filter := bson.M{"_id": id}
	if r.SoftDelete {
		filter["deleted_at"] = bson.M{"$exists": false}
	}
	if err := r.collection.FindOne(ctx, filter, opts...).Decode(data); err != nil {
		return data, err
	}
	return data, nil
}

func (r *Repository[T]) FindByIDs(ctx context.Context, ids []bson.ObjectID) ([]T, error) {
	data := make([]T, 0)
	cursor, err := r.Collection().Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return data, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &data); err != nil {
		return data, err
	}
	return data, nil
}

func (r *Repository[T]) Create(ctx context.Context, data T) (interface{}, error) {
	result, err := r.collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *Repository[T]) CreateMany(ctx context.Context, data []T) (interface{}, error) {
	result, err := r.collection.InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil
}

func (r *Repository[T]) Update(ctx context.Context, id bson.ObjectID, data interface{}) (interface{}, error) {
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})
	if err != nil {
		return nil, err
	}
	return result.UpsertedID, nil
}

func (r *Repository[T]) UpdateByFilter(ctx context.Context, filter bson.M, data interface{}) (interface{}, error) {
	result, err := r.collection.UpdateMany(ctx, filter, bson.M{"$set": data})
	if err != nil {
		return nil, err
	}
	return result.UpsertedID, nil
}

func (r *Repository[T]) Delete(ctx context.Context, id bson.ObjectID) error {
	if r.SoftDelete {
		_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
		if err != nil {
			return err
		}
	} else {
		_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			return err
		}
	}
	return nil
}
