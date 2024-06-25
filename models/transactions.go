package models

type Transactions struct{
	User Users `bson:"user" json:"user"`
	Product_name string `bson:"product_name" json:"product_name"`
	Amount uint32 `bson:"amount" json:"amount"`
	Total_price uint64 `bson:"total_price" json:"total_price"`
	Has_bought bool `bson:"has_bought" json:"has_bought"`
}