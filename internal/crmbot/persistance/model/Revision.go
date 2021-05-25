package model

import "time"

type Revision struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
