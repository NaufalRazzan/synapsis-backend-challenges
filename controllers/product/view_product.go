package product

import (
	"context"
	"synapsis-backend-test/models"
	"synapsis-backend-test/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func ViewProductsListByCategory(category string) ([]models.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("products")
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	var results []models.Products

	filter := bson.M{
		"category": category,
	}

	cursor, err := collectionName.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
