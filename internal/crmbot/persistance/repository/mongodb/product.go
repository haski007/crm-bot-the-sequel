package mongodb

import (
	"reflect"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type ProductRepository struct {
	Coll *mgo.Collection
}

func (r *ProductRepository) InitConn(session *mgo.Session, dbName string) {
	r.Coll = session.DB(dbName).C("products")
}

func (r *ProductRepository) Add(category model.Product) error {
	if r.IsProductExists(category.Title) {
		return repository.ErrDocAlreadyExists
	}

	return r.Coll.Insert(category)
}

func (r *ProductRepository) FindAll(products *[]*model.Product) error {
	return r.Coll.Find(nil).All(products)
}

func (r *ProductRepository) FindByID(productID string, product *model.Product) error {
	if !r.isIDProductExists(productID) {
		return repository.ErrDocDoesNotExist
	}

	return r.Coll.FindId(productID).One(product)
}

func (r *ProductRepository) RemoveByID(productID string) error {
	if !r.isIDProductExists(productID) {
		return repository.ErrDocDoesNotExist
	}

	return r.Coll.RemoveId(productID)
}

func (r *ProductRepository) RemoveAllByCategoryID(categoryID string) error {
	query := bson.M{
		"category_id": categoryID,
	}
	_, err := r.Coll.RemoveAll(query)
	return err
}

func (r *ProductRepository) RemoveAllBySupplierID(supplierID string) error {
	query := bson.M{
		"supplier_id": supplierID,
	}
	_, err := r.Coll.RemoveAll(query)
	return err
}

func (r *ProductRepository) FindAllByCategoryID(categoryID string, products *[]*model.Product) error {
	query := bson.M{
		"category_id": categoryID,
	}
	return r.Coll.Find(query).All(products)
}

func (r *ProductRepository) UpdateField(productID, field string, input interface{}) error {
	if !r.isIDProductExists(productID) {
		return repository.ErrDocDoesNotExist
	}

	switch reflect.ValueOf(input).Kind() {
	case reflect.Float64:
		input = input.(float64)
	default:
		input = input.(string)
	}

	query := bson.M{
		"$set": bson.M{
			field: input,
		},
	}

	return r.Coll.UpdateId(productID, query)
}

func (r *ProductRepository) UpdateFieldByTitle(productTitle, field string, input interface{}) error {
	queryFind := bson.M{
		"title": productTitle,
	}

	queryUpdate := bson.M{
		"$set": bson.M{
			field: input,
		},
	}

	return r.Coll.Update(queryFind, queryUpdate)
}

func (r *ProductRepository) AddQuantity(productTitle string, input interface{}) error {
	queryFind := bson.M{
		"title": productTitle,
	}

	queryUpdate := bson.M{
		"$inc": bson.M{
			"quantity": input,
		},
	}

	return r.Coll.Update(queryFind, queryUpdate)
}

func (r *ProductRepository) FindProductByTitle(title string, product *model.Product) error {
	if !r.IsProductExists(title) {
		return repository.ErrDocDoesNotExist
	}

	query := bson.M{
		"title": title,
	}

	return r.Coll.Find(query).One(product)
}

func (r *ProductRepository) GetFieldSum(field string) (float64, error) {

	var quantities []float64

	if err := r.Coll.Find(nil).Distinct(field, &quantities); err != nil {
		return 0, err
	}

	var total float64
	for _, v := range quantities {
		total += v
	}

	return total, nil
}

func (r *ProductRepository) FindTitlesByCategoryID(categoryID string, products *[]string) error {
	if !r.isProductExistsInCategory(categoryID) {
		return repository.ErrDocDoesNotExist
	}

	query := bson.M{
		"category_id": categoryID,
	}

	return r.Coll.Find(query).Distinct("title", products)
}

// Utils
func (r *ProductRepository) IsProductExists(title string) bool {
	query := bson.M{
		"title": title,
	}

	if n, _ := r.Coll.Find(query).Count(); n > 0 {
		return true
	}
	return false
}

func (r *ProductRepository) isProductExistsInCategory(categoryID string) bool {
	query := bson.M{
		"category_id": categoryID,
	}

	if n, _ := r.Coll.Find(query).Count(); n > 0 {
		return true
	}
	return false
}

func (r *ProductRepository) isIDProductExists(productID string) bool {
	if n, _ := r.Coll.FindId(productID).Count(); n > 0 {
		return true
	}
	return false
}
