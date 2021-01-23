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
