package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName           = "Ecommerce"
	UsersCollectionName    = "Users"
	ProductsCollectionName = "Products"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err = client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	return client, nil
}

func UsersCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(DatabaseName).Collection(UsersCollectionName)
}

func ProductsCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(DatabaseName).Collection(ProductsCollectionName)
}

// UserData and ProductData remain as aliases for backward compatibility.
func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database(DatabaseName).Collection(collectionName)
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database(DatabaseName).Collection(collectionName)
}
