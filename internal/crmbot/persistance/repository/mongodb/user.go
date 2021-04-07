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

func (r *UserRepository) InitConn(session *mgo.Session, dbName string) {
	r.Coll = session.DB(dbName).C("users")
}

func (r *UserRepository) AddUser(user model.User) error {
	if r.isUserExistsByTgID(user.TgID) {
		return repository.ErrDocAlreadyExists
	}

	return r.Coll.Insert(user)
}

func (r *UserRepository) GetAll(users *[]model.User) error {
	if count, err := r.Coll.Find(nil).Count(); count == 0 {
		return repository.ErrDocDoesNotExist
	} else if err != nil {
		return err
	}

	return r.Coll.Find(nil).All(users)
}

func (r *UserRepository) RemoveByID(userID string) error {
	if !r.isUserExistsByID(userID) {
		return repository.ErrDocDoesNotExist
	}

	return r.Coll.RemoveId(userID)
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

func (r *UserRepository) isUserExistsByID(userID string) bool {
	query := bson.M{
		"_id": userID,
	}

	if n, _ := r.Coll.Find(query).Count(); n > 0 {
		return true
	}
	return false
}
