package model

type Category struct {
	ID          string `json:"_id" bson:"_id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}
