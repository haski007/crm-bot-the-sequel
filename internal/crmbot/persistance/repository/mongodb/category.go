package mongodb

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type CategoryRepository struct {
	Coll *mgo.Collection
}

func (r *CategoryRepository) InitConn() {
	r.Coll = session.DB(cfg.DBName).C("categories")
}

func (r *CategoryRepository) Add(category model.Category) error {
	if r.isCategoryExists(category.Title) {
		return repository.ErrDocAlreadyExists
	}

	return r.Coll.Insert(category)
}

func (r *CategoryRepository) FindAll(categories *[]*model.Category) error {
	return r.Coll.Find(nil).All(categories)
}

func (r *CategoryRepository) RemoveByID(categoryID string) error {
	if !r.isCategoryExists(categoryID) {
		return repository.ErrDocDoesNotExist
	}

	return r.Coll.RemoveId(categoryID)
}

func (r *CategoryRepository) isCategoryExists(title string) bool {
	query := bson.M{
		"title": title,
	}

	if n, _ := r.Coll.Find(query).Count(); n > 0 {
		return true
	}
	return false
}
