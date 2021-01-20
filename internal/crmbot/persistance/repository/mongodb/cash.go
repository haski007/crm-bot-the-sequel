package mongodb

import "github.com/globalsign/mgo"

type CashRepository struct {
	Coll *mgo.Collection
}

func (r *CashRepository) InitConn() {
	r.Coll = session.DB(cfg.DBName).C("cash")
}
