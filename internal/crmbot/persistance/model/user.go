package model

type User struct {
	ID        string `json:"_id" bson:"_id"`
	TgID      int    `json:"tgid" bson:"tgid"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Username  string `json:"username" bson:"username"`

	Role string `json:"role" bson:"role"`
}
