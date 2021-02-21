package mongodb

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type UserRepository struct {
	Coll *mgo.Collection
}

func (r *UserRepository) InitConn() {
	r.Coll = session.DB(cfg.DBName).C("users")
}

func (r *UserRepository) AddUser(user model.User) error {
	if r.isUserExistsByTgID(user.TgID) {
		return repository.ErrDocAlreadyExists
	}

	return r.Coll.Insert(user)
}

// ---> Utils

func (r *UserRepository) isUserExistsByTgID(tgID int) bool {
	query := bson.M{
		"tg_id": tgID,
	}

	if n, _ := r.Coll.Find(query).Count(); n > 0 {
		return true
	}
	return false
}
