package user

import (
	"context"
	"fmt"
	"synapsis-backend-test/models"
	"synapsis-backend-test/pkg/db"
	"synapsis-backend-test/pkg/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/google/uuid"
)

func SignUp(user models.Users) (string, error){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("users")
	if err != nil{
		return "", err
	}
	defer client.Disconnect(ctx)

	uniquefilter := bson.M{
		"full_name": user.Full_name,
	}

	// check unique
	countUser, err := collectionName.CountDocuments(ctx, uniquefilter)
	if err != nil{
		return "", err
	}
	if countUser > 0{
		return "", fmt.Errorf("user name already exists")
	}

	// hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil{
		return "", err
	}
	user.Password = hashedPassword

	// generate uuid
	id := uuid.New()
	user.User_id = id.String()
	
	// insert to mongo
	result, err := collectionName.InsertOne(ctx, user)
	if err != nil{
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).String(), nil
}