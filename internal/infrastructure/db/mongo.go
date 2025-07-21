package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Test authentication by actually pinging the database
	err = client.Ping(ctx, nil)
	if err != nil {
		client.Disconnect(ctx) // Clean up on failure
		return nil, err
	}

	// Test actual database operations to catch auth errors early
	db := client.Database("ledger")
	err = db.RunCommand(ctx, map[string]interface{}{"ping": 1}).Err()
	if err != nil {
		client.Disconnect(ctx) // Clean up on failure
		return nil, err
	}

	return client, nil
}