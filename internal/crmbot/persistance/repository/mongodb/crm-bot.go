package mongodb

import (
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/caarlos0/env"
	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
)

const (
	collName = "chats"
)

var session *mgo.Session
var cfg config.MongoCfg

type ChatRepository struct {
	coll *mgo.Collection
}

func (r *ChatRepository) InitDBConn() {
	if err := env.Parse(&cfg); err != nil {
		logrus.Fatalf("[env Parse] MongoCfg err: %s", err)
	}

	cfg.Addr = "mongodb://" + cfg.Username + ":" + cfg.Password + "@" + cfg.HostName + ":" + cfg.Port
	session, err := mgo.Dial(cfg.Addr)
	if err != nil {
		logrus.Fatalf("[mgo Dial] addr: %s | err: %s", cfg.Addr, err)
	}

	if err = session.Ping(); err != nil {
		logrus.Fatalf("[mgo Ping] addr: %s | err: %s", cfg.Addr, err)
	}

	r.coll = session.DB(cfg.DBName).C(collName)
}
