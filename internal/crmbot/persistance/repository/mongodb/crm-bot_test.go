package mongodb

import (
	"testing"

	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
)

const (
	mongoAddr = "mongodb://demian:password@172.20.0.2:27017"
)

func getCollection(collName string) *mgo.Collection {
	session, err := mgo.Dial(mongoAddr)
	if err != nil {
		logrus.Fatalf("[mgo Dial] addr: %s | err: %s", cfg.Addr, err)
	}

	if err = session.Ping(); err != nil {
		logrus.Fatalf("[mgo Ping] addr: %s | err: %s", cfg.Addr, err)
	}

	return session.DB(cfg.DBName).C(collName)
}

func TestChatRepository_PushNewPublusher(t *testing.T) {
	t.Skip()
}
