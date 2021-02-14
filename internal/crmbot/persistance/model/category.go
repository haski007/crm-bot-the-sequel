package model

import "errors"

type Category struct {
	ID          string `json:"_id" bson:"_id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

type CategoryEditField string

func (f CategoryEditField) String() string {
	return string(f)
}

func (f CategoryEditField) BsonField() string {
	switch f {
	case CategoryEditTitle:
		return "title"
	case CategoryEditDescription:
		return "description"
	}
	return ""
}

func NewCategoryEditField(input string) (CategoryEditField, error) {
	switch input {
	case CategoryEditTitle.String():
		return CategoryEditTitle, nil
	case CategoryEditDescription.String():
		return CategoryEditDescription, nil
	}
	return CategoryEditField(""), ErrNoSuchCategoryEditField
}

var (
	ErrNoSuchCategoryEditField = errors.New("NO_SUCH_EDIT_FIELD")
)

const (
	CategoryEditTitle       CategoryEditField = "Название"
	CategoryEditDescription CategoryEditField = "Описание"
)

type CategoryEdit struct {
	ID    string
	Field CategoryEditField
}
