package db

import (
	"context"
	"synapsis-backend-test/configs"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectDB(collectionName string) (*mongo.Client, *mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(configs.GetConfig().MongoDBURI))
	if err != nil {
		return nil, nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, nil, err
	}

	collection := client.Database(configs.GetConfig().MONGODBName).Collection(collectionName)

	return client, collection, nil
}
