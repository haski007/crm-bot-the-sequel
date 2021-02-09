package mongodb

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type SupplierRepository struct {
	Coll *mgo.Collection
}

func (r *SupplierRepository) InitConn() {
	r.Coll = session.DB(cfg.DBName).C("suppliers")
}

func (r *SupplierRepository) Add(supplier model.Supplier) error {
	if r.isSupplierExists(supplier.Name) {
		return repository.ErrDocAlreadyExists
	}

	return r.Coll.Insert(supplier)
}

func (r *SupplierRepository) FindAll(suppliers *[]*model.Supplier) error {
	return r.Coll.Find(nil).All(suppliers)
}

func (r *SupplierRepository) FindByName(name string, supplier *model.Supplier) error {
	if !r.isSupplierExists(name) {
		return repository.ErrDocDoesNotExist
	}

	query := bson.M{
		"name": name,
	}

	return r.Coll.Find(query).One(supplier)
}

func (r *SupplierRepository) DistinctNames(suppliers *[]string) error {
	return r.Coll.Find(nil).Distinct("name", suppliers)
}

func (r *SupplierRepository) FindByID(supplierID string, supplier *model.Supplier) error {
	if !r.isIDSupplierExists(supplierID) {
		return repository.ErrDocDoesNotExist
	}

	return r.Coll.FindId(supplierID).One(supplier)
}

func (r *SupplierRepository) RemoveByID(suppplierID string) error {
	if !r.isIDSupplierExists(suppplierID) {
		return repository.ErrDocDoesNotExist
	}

	return r.Coll.RemoveId(suppplierID)
}

func (r *SupplierRepository) UpdateField(supplierID, field string, value interface{}) error {
	if !r.isIDSupplierExists(supplierID) {
		return repository.ErrDocDoesNotExist
	}

	query := bson.M{
		"$set": bson.M{
			field: value,
		},
	}

	return r.Coll.UpdateId(supplierID, query)
}

// Utils
func (r *SupplierRepository) isSupplierExists(name string) bool {
	query := bson.M{
		"name": name,
	}

	if n, _ := r.Coll.Find(query).Count(); n > 0 {
		return true
	}
	return false
}

func (r *SupplierRepository) isIDSupplierExists(supplierID string) bool {
	if n, _ := r.Coll.FindId(supplierID).Count(); n > 0 {
		return true
	}
	return false
}
