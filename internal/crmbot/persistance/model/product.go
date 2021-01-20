package model

type Unit string

const (
	PieceUnit Unit = "шт"
	LiterUnit Unit = "л"
	Gram      Unit = "гр"
)

type Product struct {
	ID              string   `json:"_id" bson:"_id"`
	Title           string   `json:"title" bson:"title"`
	Description     string   `json:"description" bson:"description"`
	PurchasingPrice Money    `json:"purchasing_price" bson:"purchasing_price"`
	BidPrice        Money    `json:"bid_price" bson:"bid_price"`
	Quantity        int      `json:"quantity" bson:"quantity"`
	Category        Category `json:"category" bson:"category"`
	Supplier        Supplier `json:"supplier" bson:"supplier"`
	Unit            Unit     `json:"unit" bson:"unit"`
}
