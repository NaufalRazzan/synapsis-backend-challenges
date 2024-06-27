package models

type Users struct {
	User_id   string `bson:"user_id" json:"user_id"`
	Full_name string `bson:"full_name" json:"full_name" validate:"required"`
	Email     string `bson:"email" json:"email" validate:"required"`
	Password  string `bson:"password" json:"password" validate:"required"`
}

type UsersLogin struct{
	Email     string `bson:"email" json:"email" validate:"required"`
	Password  string `bson:"password" json:"password" validate:"required"`
}