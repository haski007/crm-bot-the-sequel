package mongodb

import (
	"time"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type CashRepository struct {
	Coll *mgo.Collection
}

func (r *CashRepository) InitConn() {
	r.Coll = session.DB(cfg.DBName).C("cash")
}

func (r *CashRepository) ChangeAmount(diff model.Money) error {
	queryUpdate := bson.M{
		"$inc": bson.M{
			"amount": diff,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	return r.Coll.UpdateId(config.MainCashID, queryUpdate)
}
