package product

import (
	"context"
	"fmt"
	"synapsis-backend-test/models"
	"synapsis-backend-test/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertProductToShoppingCart(product_id string, user_id string, amount uint64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
		"$inc": -amount,
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

	trxPayload := models.Transactions{
		Product_id:  product_id,
		User_id:     user_id,
		Amount:      amount,
		Total_price: amount * resProduct.Price,
		Has_bought:  false,
	}

	res, err := trxCollection.InsertOne(ctx, trxPayload)
	if err != nil {
		return "", err
	}

	resultID, _ := res.InsertedID.(primitive.ObjectID)

	return resultID.String(), nil
}

func ViewShoppingCartLists(user_id string) ([]models.ViewShoppingCartResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, trxCollection, err := db.ConnectDB("transactions")
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"$and": []bson.M{
					{"user_id": user_id,},
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
				"_id":         0,
				"transaction_id": "$transaction_id",
				"amount":      "$amount",
				"total_price": "$total_price",
				"products": bson.A{
					bson.M{
						"product_id":   "$joinedProducts.product_id",
						"product_name": "$joinedProducts.product_name",
						"price":        "$joinedProducts.price",
						"stock":        "$joinedProducts.stock",
						"category":     "$joinedProducts.category",
					},
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, trxCollection, err := db.ConnectDB("transactions")
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	filter := bson.M{
		"transaction_id": trx_id,
	}

	result, err := trxCollection.DeleteOne(ctx, filter)
	if err != nil{
		return err
	}
	if result.DeletedCount == 0{
		return fmt.Errorf("no data exists to be deleted")
	}

	return nil
}

func CheckoutToPayment(trx_id string) error{
	var trxData models.Transactions

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trxClient, trxCollection, err := db.ConnectDB("transactions")
	if err != nil{
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

	
	if err := trxCollection.FindOneAndUpdate(ctx, trxFilter, trxUpdate).Decode(&trxData); err != nil{
		return err
	}

	// update product stock
	productClient, productCollection, err := db.ConnectDB("products")
	if err != nil{
		return err
	}
	defer productClient.Disconnect(ctx)

	productFilter := bson.M{
		"product_id": trxData.Product_id,
	}
	productUpdate := bson.M{
		"$inc": trxData.Amount,
	}

	if _, err := productCollection.UpdateOne(ctx, productFilter, productUpdate); err != nil{
		return err
	}

	return nil
}
