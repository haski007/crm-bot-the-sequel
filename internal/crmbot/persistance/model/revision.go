package model

import "time"

type RevisionStatusType string

const (
	RevisionInitiated  RevisionStatusType = "Initiated"
	RevisionProcessing RevisionStatusType = "Processing"
	RevisionCompleted  RevisionStatusType = "Completed"
	RevisionFailed     RevisionStatusType = "Failed"
	RevisionRejected   RevisionStatusType = "Rejected"
)

func (t RevisionStatusType) String() string {
	return string(t)
}

type Revision struct {
	ID           string             `json:"_id" bson:"_id"`
	Result       string             `json:"result" bson:"result"`
	Status       RevisionStatusType `json:"status" bson:"status"`
	Author       string             `json:"author" bson:"author"`
	ProductsCost Money              `json:"products_cost" bson:"products_cost"`
	Cash         Money              `json:"cash" bson:"cash"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}
