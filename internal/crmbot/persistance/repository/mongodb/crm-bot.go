package mongodb

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
)

func InitDBConnection(cfg config.Mongo) *mgo.Session {

	tlsConfig := &tls.Config{}

	logrus.Println("Connecting to mongoDB...")
	dialInfo := &mgo.DialInfo{
		Timeout:  time.Second * 10,
		Addrs:    cfg.Addrs,
		Database: "admin",
		Username: cfg.Username,
		Password: cfg.Password,
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		logrus.Fatalf("[mgo Dial] err: %s", err)
	}

	session.SetMode(mgo.PrimaryPreferred, false)

	if err = session.Ping(); err != nil {
		logrus.Fatalf("[mgo Ping] err: %s", err)
	}

	return session
}
