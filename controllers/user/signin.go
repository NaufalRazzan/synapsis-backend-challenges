package user

import (
	"context"
	"fmt"
	"synapsis-backend-test/models"
	"synapsis-backend-test/pkg/db"
	"synapsis-backend-test/pkg/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(email string, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("users")
	if err != nil{
		return "", err
	}
	defer client.Disconnect(ctx)

	var user models.Users

	filter := bson.M{
		"$and": []bson.M{
			{"email": email},
			{"password": password},
		},
	}

	if err := collectionName.FindOne(ctx, filter).Decode(&user); err != nil{
		if err == mongo.ErrNoDocuments{
			return "", fmt.Errorf("invalid email or password")
		}
		return "", err
	}

	// verify password
	if err := utils.ComparePasswords(password, user.Password); err != nil{
		if err == bcrypt.ErrMismatchedHashAndPassword{
			return "", fmt.Errorf("invalid email or password")
		}
		return "", err
	}

	// assign jwt token
	acc_token, err := utils.GenerateJWT(user)
	if err != nil{
		return "", err
	}

	return acc_token, nil
}