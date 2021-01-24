package mongodb

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type ProductRepository struct {
	Coll *mgo.Collection
}

func (r *ProductRepository) InitConn() {
	r.Coll = session.DB(cfg.DBName).C("products")
}

func (r *ProductRepository) Add(category model.Product) error {
	if r.isProductExists(category.Title) {
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

// Utils
func (r *ProductRepository) isProductExists(title string) bool {
	query := bson.M{
		"title": title,
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
