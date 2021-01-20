package mongodb

import "github.com/globalsign/mgo"

type TransactionRepository struct {
	Coll *mgo.Collection
}

func (r *TransactionRepository) InitConn() {
	r.Coll = session.DB(cfg.DBName).C("transactions")
}
