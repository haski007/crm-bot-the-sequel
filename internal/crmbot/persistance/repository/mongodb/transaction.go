package mongodb

import (
	"time"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/jinzhu/now"
)

type TransactionRepository struct {
	Coll *mgo.Collection
}

func (r *TransactionRepository) InitConn() {
	r.Coll = session.DB(cfg.DBName).C("transactions")
}

func (r *TransactionRepository) Add(transaction model.Transaction) error {
	return r.Coll.Insert(transaction)
}

func (r *TransactionRepository) GetAll(transactions *[]model.Transaction) error {
	if count, err := r.Coll.Find(nil).Count(); count == 0 {
		return repository.ErrDocDoesNotExist
	} else if err != nil {
		return err
	}

	return r.Coll.Find(nil).All(transactions)
}

func (r *TransactionRepository) GetForLastDays(days int, transactions *[]model.Transaction) error {
	date := time.Now().Add(time.Hour * 24 * time.Duration(-days))
	query := bson.M{
		"created_at": bson.M{
			"$gt": date,
		},
	}

	if count, err := r.Coll.Find(query).Count(); count == 0 {
		return repository.ErrDocDoesNotExist
	} else if err != nil {
		return err
	}

	return r.Coll.Find(query).All(transactions)
}

func (r *TransactionRepository) RemoveByID(txID string) error {
	if !r.isTxExistsByID(txID) {
		return repository.ErrDocDoesNotExist
	}

	return r.Coll.RemoveId(txID)
}

func (r *TransactionRepository) GetTxByID(txID string, tx *model.Transaction) error {
	if !r.isTxExistsByID(txID) {
		return repository.ErrDocDoesNotExist
	}

	return r.Coll.FindId(txID).One(tx)
}

func (r *TransactionRepository) GetCurrentMonth(transactions *[]*model.Transaction) error {
	query := bson.M{
		"created_at": bson.M{
			"$gt": now.BeginningOfMonth(),
		},
	}

	return r.Coll.Find(query).All(transactions)
}

// ---> Utils

func (r *TransactionRepository) isTxExistsByID(txID string) bool {
	query := bson.M{
		"_id": txID,
	}

	if n, _ := r.Coll.Find(query).Count(); n > 0 {
		return true
	}
	return false
}
