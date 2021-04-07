package auth

import (
	"github.com/globalsign/mgo"
)

type AuthService struct {
	UsersColl *mgo.Collection
}

func NewAuthService(cfg Config) (*AuthService, error) {

	return &AuthService{
		UsersColl: cfg.MgoSession.DB(cfg.Mongo.DBName).C(cfg.UsersCollName),
	}, nil
}
