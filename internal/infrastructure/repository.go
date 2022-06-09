package infrastructure

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func insertOne(ctx context.Context, collection *mongo.Collection, obj interface{}) error {
	_, err := collection.InsertOne(ctx, obj)
	if err != nil {
		return err
	}
	return nil
}

func findOne[T any](ctx context.Context, collection *mongo.Collection, searchKey bson.D) (*T, error) {
	var obj T
	err := collection.FindOne(ctx, searchKey).Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
