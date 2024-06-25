package models

type Users struct{
	User_id string `bson:"user_id" json:"user_id"`
	Full_name string `bson:"full_name" json:"full_name"`
	Email string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}