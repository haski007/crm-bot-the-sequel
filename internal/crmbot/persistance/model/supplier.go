package model

import (
	"errors"
)

type Supplier struct {
	ID          string `json:"_id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Phone       string `json:"phone" bson:"phone"`
	Description string `json:"description" bson:"description"`
}

type SupplierEditField string

func (f SupplierEditField) String() string {
	return string(f)
}

func (f SupplierEditField) BsonField() string {
	switch f {
	case SupplierEditName:
		return "name"
	case SupplierEditPhone:
		return "phone"
	case SupplierEditDescription:
		return "description"
	}
	return ""
}

func NewSupplierEditField(input string) (SupplierEditField, error) {
	switch input {
	case SupplierEditName.String():
		return SupplierEditName, nil
	case SupplierEditPhone.String():
		return SupplierEditPhone, nil
	case SupplierEditDescription.String():
		return SupplierEditDescription, nil
	}
	return SupplierEditField(""), ErrNoSuchSupplierEditField
}

var (
	ErrNoSuchSupplierEditField = errors.New("NO_SUCH_EDIT_FIELD")
)

const (
	SupplierEditName        SupplierEditField = "Ф.И.О."
	SupplierEditPhone       SupplierEditField = "Номер телефона"
	SupplierEditDescription SupplierEditField = "Описание"
)

type SupplierEdit struct {
	ID    string
	Field SupplierEditField
}
