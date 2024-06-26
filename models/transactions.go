package models

type Transactions struct {
	Transaction_id string `bson:"transaction_id" json:"transaction_id"`
	User_id        string `bson:"user_id" json:"user_id"`
	Product_id     string `bson:"product_id" json:"product_id"`
	Amount         uint64 `bson:"amount" json:"amount"`
	Total_price    uint64 `bson:"total_price" json:"total_price"`
	Has_bought     bool   `bson:"has_bought" json:"has_bought"`
}

type ViewShoppingCartResponse struct {
	Product     Products `json:"products"`
	Amount      uint64   `json:"amount"`
	Total_price uint64   `json:"total_price"`
}
