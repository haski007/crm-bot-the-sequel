package model

type Unit string

func (u Unit) String() string {
	return string(u)
}

const (
	PieceUnit Unit = "шт"
	LiterUnit Unit = "л"
	GramUnit  Unit = "гр"
)

type Product struct {
	ID              string `json:"_id" bson:"_id"`
	Title           string `json:"title" bson:"title"`
	Description     string `json:"description" bson:"description"`
	PurchasingPrice Money  `json:"purchasing_price" bson:"purchasing_price"`
	BidPrice        Money  `json:"bid_price" bson:"bid_price"`
	Quantity        int    `json:"quantity" bson:"quantity"`
	CategoryID      string `json:"category_id" bson:"category_id"`
	SupplierID      string `json:"supplier_id" bson:"supplier_id"`
	Unit            Unit   `json:"unit" bson:"unit"`
}
