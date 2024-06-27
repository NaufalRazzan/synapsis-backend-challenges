package product

import (
	"context"
	"fmt"
	"synapsis-backend-test/models"
	"synapsis-backend-test/pkg/db"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUserId(user_name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("users")
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)

	filter := bson.M{
		"full_name": user_name,
	}

	var userData models.Users

	if err := collectionName.FindOne(ctx, filter).Decode(&userData); err != nil {
		fmt.Println("error here")
		return "", err
	}

	return userData.User_id, nil

}

func InsertProductToShoppingCart(product_id string, user_name string, amount int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), (30*time.Second)*2)
	defer cancel()

	// get product data
	productClient, productCollection, err := db.ConnectDB("products")
	if err != nil {
		return "", err
	}
	defer productClient.Disconnect(ctx)

	var resProduct models.Products

	filter := bson.M{
		"$and": []bson.M{
			{
				"product_id": product_id,
			},
			{
				"stock": bson.M{
					"$gt": 0,
				},
			},
		},
	}

	// decrement product stock
	update := bson.M{
		"$inc": bson.M{
			"stock": -amount,
		},
	}

	if err := productCollection.FindOneAndUpdate(ctx, filter, update).Decode(&resProduct); err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("products does not exists")
		}
		return "", err
	}

	if resProduct.Stock == 0 {
		return fmt.Sprintf("stock is empty for %s", resProduct.Product_name), nil
	}
	if amount > resProduct.Stock {
		return "bought amount is out of stock", nil
	}

	// insert to transaction
	trxClient, trxCollection, err := db.ConnectDB("transactions")
	if err != nil {
		return "", err
	}
	defer trxClient.Disconnect(ctx)

	// get user id
	user_id, err := getUserId(user_name)
	if err != nil {
		return "", err
	}

	trxPayload := models.Transactions{
		Transaction_id: uuid.NewString(),
		Product_id:     product_id,
		User_id:        user_id,
		Amount:         amount,
		Total_price:    amount * resProduct.Price,
		Has_bought:     false,
	}

	res, err := trxCollection.InsertOne(ctx, trxPayload)
	if err != nil {
		return "", err
	}

	resultID, _ := res.InsertedID.(primitive.ObjectID)

	return resultID.String(), nil
}

func ViewShoppingCartLists(user_name string) ([]models.ViewShoppingCartResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, trxCollection, err := db.ConnectDB("transactions")
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	user_id, err := getUserId(user_name)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"$and": []bson.M{
					{"user_id": user_id},
					{"has_bought": false},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "products",
				"localField":   "product_id",
				"foreignField": "product_id",
				"as":           "joinedProducts",
			},
		},
		{
			"$unwind": "$joinedProducts",
		},
		{
			"$project": bson.M{
				"_id":            0,
				"transaction_id": "$transaction_id",
				"amount":         "$amount",
				"total_price":    "$total_price",
				"products": bson.M{
					"product_id":   "$joinedProducts.product_id",
					"product_name": "$joinedProducts.product_name",
					"price":        "$joinedProducts.price",
					"stock":        "$joinedProducts.stock",
					"category":     "$joinedProducts.category",
				},
			},
		},
	}

	cursor, err := trxCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var result []models.ViewShoppingCartResponse
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func DeleteProductFromShoppingCart(trx_id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, trxCollection, err := db.ConnectDB("transactions")
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	filter := bson.M{
		"transaction_id": trx_id,
	}

	var deletedData models.Transactions

	if err := trxCollection.FindOneAndDelete(ctx, filter).Decode(&deletedData); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("no data exists to be deleted")
		}
		return err
	}

	// update product stock
	productClient, productCollection, err := db.ConnectDB("products")
	if err != nil {
		return err
	}
	defer productClient.Disconnect(ctx)

	productFilter := bson.M{
		"product_id": deletedData.Product_id,
	}
	productUpdate := bson.M{
		"$inc": bson.M{
			"stock": deletedData.Amount,
		},
	}

	if _, err := productCollection.UpdateOne(ctx, productFilter, productUpdate); err != nil {
		return err
	}

	return nil
}

func CheckoutToPayment(trx_id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	trxClient, trxCollection, err := db.ConnectDB("transactions")
	if err != nil {
		return err
	}
	defer trxClient.Disconnect(ctx)

	trxFilter := bson.M{
		"transaction_id": trx_id,
	}

	trxUpdate := bson.M{
		"$set": bson.M{
			"has_bought": true,
		},
	}

	res, err := trxCollection.UpdateOne(ctx, trxFilter, trxUpdate)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("no document modified")
	}

	return nil
}
