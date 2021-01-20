package mongodb

import "github.com/globalsign/mgo"

type ProductRepository struct {
	Coll *mgo.Collection
}

func (r *ProductRepository) InitConn() {
	r.Coll = session.DB(cfg.DBName).C("products")
}
