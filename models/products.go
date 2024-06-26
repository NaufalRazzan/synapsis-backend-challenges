package models

type Products struct {
	Product_id   string `bson:"product_id" json:"product_id"`
	Product_name string `bson:"product_name" json:"product_name"`
	Price        uint64 `bson:"price" json:"price"`
	Stock        uint64 `bson:"stock" json:"stock"`
	Category     string `bson:"category" json:"category"`
}
