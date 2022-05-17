package infrastructure

import (
  "context"

  "go.mongodb.org/mongo-driver/mongo"
)

func insertOne(ctx context.Context, collection *mongo.Collection, obj interface{}) error {
  _, err := collection.InsertOne(ctx, obj)
  if err != nil {
    return err
  }
  return nil
}
