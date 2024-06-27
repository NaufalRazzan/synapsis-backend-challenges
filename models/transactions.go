package models

type Transactions struct {
	Transaction_id string `bson:"transaction_id" json:"transaction_id"`
	User_id        string `bson:"user_id" json:"user_id"`
	Product_id     string `bson:"product_id" json:"product_id"`
	Amount         int64  `bson:"amount" json:"amount"`
	Total_price    int64  `bson:"total_price" json:"total_price"`
	Has_bought     bool   `bson:"has_bought" json:"has_bought"`
}

type ProductResponse struct {
	Product_name string `bson:"product_name" json:"product_name"`
	Price        int64  `bson:"price" json:"price"`
	Stock        int64  `bson:"stock" json:"stock"`
	Category     string `bson:"category" json:"category"`
}

type ViewShoppingCartResponse struct {
	Transaction_id string          `bson:"transaction_id" json:"transaction_id"`
	Product        ProductResponse `bson:"products" json:"products"`
	Amount         int64           `bson:"amount" json:"amount"`
	Total_price    int64           `bson:"total_price" json:"total_price"`
}

type InsertTrxpayload struct {
	Product_id string `bson:"product_id" json:"product_id" validate:"required"`
	User_name  string `bson:"user_name" json:"user_name" validate:"required"`
	Amount     int64  `bson:"amount" json:"amount" validate:"required"`
}
