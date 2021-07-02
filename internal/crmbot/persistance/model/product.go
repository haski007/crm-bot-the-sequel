package model

import "errors"

type Unit string

func (u Unit) String() string {
	return string(u)
}

const (
	PieceUnit Unit = "шт"
	LiterUnit Unit = "л"
	GramUnit  Unit = "гр"

	UnknownUnit Unit = ""
)

func NewUnit(input string) (Unit, error) {
	switch input {
	case PieceUnit.String():
		return PieceUnit, nil
	case LiterUnit.String():
		return LiterUnit, nil
	case GramUnit.String():
		return GramUnit, nil
	default:
		return UnknownUnit, errors.New("NO_SUCH_UNIT")
	}
}

type Product struct {
	ID              string  `json:"_id" bson:"_id"`
	Title           string  `json:"title" bson:"title"`
	PurchasingPrice Money   `json:"pur_price" bson:"pur_price"`
	BidPrice        Money   `json:"bid_price" bson:"bid_price"`
	Quantity        float64 `json:"quantity" bson:"quantity"`
	Unit            Unit    `json:"unit" bson:"unit"`
	CategoryID      string  `json:"category_id" bson:"category_id"`
	SupplierID      string  `json:"supplier_id" bson:"supplier_id"`
	Description     string  `json:"description" bson:"description"`
}

const (
	ProductEditTitle       ProductEditField = "Название"
	ProductEditPurPrice    ProductEditField = "Цена закупки"
	ProductEditBidPrice    ProductEditField = "Цена продажи"
	ProductEditUnit        ProductEditField = "Единица измерения"
	ProductEditSupplier    ProductEditField = "Поставщик"
	ProductEditCategory    ProductEditField = "Категория"
	ProductEditDescription ProductEditField = "Описание"
)

type ProductEditField string

func (f ProductEditField) String() string {
	return string(f)
}

func (f ProductEditField) BsonField() string {
	switch f {
	case ProductEditTitle:
		return "title"
	case ProductEditPurPrice:
		return "pur_price"
	case ProductEditBidPrice:
		return "bid_price"
	case ProductEditUnit:
		return "unit"
	case ProductEditSupplier:
		return "supplier_id"
	case ProductEditCategory:
		return "category_id"
	case ProductEditDescription:
		return "description"
	}
	return ""
}

func NewProductEditField(input string) (ProductEditField, error) {
	switch input {
	case ProductEditTitle.String():
		return ProductEditTitle, nil
	case ProductEditPurPrice.String():
		return ProductEditPurPrice, nil
	case ProductEditBidPrice.String():
		return ProductEditBidPrice, nil
	case ProductEditUnit.String():
		return ProductEditUnit, nil
	case ProductEditSupplier.String():
		return ProductEditSupplier, nil
	case ProductEditCategory.String():
		return ProductEditCategory, nil
	case ProductEditDescription.String():
		return ProductEditDescription, nil
	}
	return ProductEditField(""), ErrNoSuchProductEditField
}

var (
	ErrNoSuchProductEditField = errors.New("NO_SUCH_EDIT_FIELD")
)

type ProductEdit struct {
	ID    string
	Field ProductEditField
}
