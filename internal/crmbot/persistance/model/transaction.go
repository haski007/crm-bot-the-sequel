package model

import "time"

type Money float64

func NewMoney(input float64) Money {
	return Money(input)
}

func (m Money) Add(input float64) Money {
	return NewMoney(m.Float64() + input)
}

func (m Money) Float64() float64 {
	return float64(m)
}

type TxType string

const (
	TxAddCash  TxType = "TxAddCash"
	TxGetCash  TxType = "TxGetCash"
	TxAddGoods TxType = "TxAddGoods"
	TxGetGoods TxType = "TxGetGoods"
)

const (
	TxPurchaseComment = "Purchase"
)

type Transaction struct {
	ID        string    `json:"_id" bson:"_id"`
	Author    string    `json:"author" bson:"author"`
	Amount    Money     `json:"Amount" bson:"Amount"`
	Type      TxType    `json:"type" bson:"type"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Comment   string    `json:"comment" bson:"comment"`
}

type TransactionStock struct {
	ID           string    `json:"_id" bson:"_id"`
	Author       string    `json:"author" bson:"author"`
	ProductTitle string    `json:"product_title" bson:"product_title"`
	ProductPrice Money     `json:"product_price" bson:"product_price"`
	Amount       float64   `json:"Amount" bson:"Amount"`
	Type         TxType    `json:"type" bson:"type"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}
