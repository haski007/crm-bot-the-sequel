package model

import "time"

type RevisionStatusType string

const (
	Initiated RevisionStatusType = "Initiated"
	Completed RevisionStatusType = "Completed"
	Failed    RevisionStatusType = "Failed"
	Rejected  RevisionStatusType = "Rejected"
)

type Revision struct {
	ID        string             `json:"_id" bson:"_id"`
	Result    string             `json:"result" bson:"result"`
	Status    RevisionStatusType `json:"status" bson:"status"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
