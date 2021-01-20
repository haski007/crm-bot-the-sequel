package model

type Supplier struct {
	ID    string `json:"_id" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
}
