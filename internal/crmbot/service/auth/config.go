package auth

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/globalsign/mgo"
)

type Config struct {
	MgoSession    *mgo.Session `json:"mgo_session"`
	Mongo         config.Mongo `json:"mongo"`
	UsersCollName string       `json:"users_coll_name"`
}
