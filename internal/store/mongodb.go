package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient() *mongo.Client {
	ctx := context.Background()
	opts := options.Client().ApplyURI("mongodb://localhost:27017").SetConnectTimeout(5 * time.Second)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	return client
}
