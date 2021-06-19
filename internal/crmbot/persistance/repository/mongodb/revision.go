package mongodb

import (
	"reflect"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type RevisionRepository struct {
	Coll *mgo.Collection
}

func (r *RevisionRepository) InitConn(session *mgo.Session, dbName string) {
	r.Coll = session.DB(dbName).C("revisions")
}

func (r *RevisionRepository) Add(revision model.Revision) error {
	return r.Coll.Insert(revision)
}

func (r *RevisionRepository) FindAll(products *[]*model.Product) error {
	return r.Coll.Find(nil).All(products)
}

func (r *RevisionRepository) FindByID(id string, revision *model.Revision) error {
	err := r.Coll.FindId(id).One(revision)
	if err == mgo.ErrNotFound {
		return repository.ErrDocDoesNotExist
	}

	return err
}

func (r *RevisionRepository) RemoveByID(id string) error {
	err := r.Coll.RemoveId(id)
	if err == mgo.ErrNotFound {
		return repository.ErrDocDoesNotExist
	}

	return err
}

func (r *RevisionRepository) UpdateField(revisionID, field string, input interface{}) error {

	switch reflect.ValueOf(input).Kind() {
	case reflect.Float64:
		input = input.(float64)
	case reflect.String:
		input = input.(string)
	}

	query := bson.M{
		"$set": bson.M{
			field: input,
		},
	}

	err := r.Coll.UpdateId(revisionID, query)
	if err == mgo.ErrNotFound {
		return repository.ErrDocDoesNotExist
	}
	return err
}

// Utils
func (r *RevisionRepository) isIDRevisionExists(productID string) bool {
	if n, _ := r.Coll.FindId(productID).Count(); n > 0 {
		return true
	}
	return false
}
