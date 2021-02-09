package mongodb

import (
	"crypto/tls"
	"net"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/caarlos0/env"
	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
)

var session *mgo.Session
var cfg config.MongoCfg

func init() {
	if err := env.Parse(&cfg); err != nil {
		logrus.Fatalf("[env Parse] MongoCfg err: %s", err)
	}

	tlsConfig := &tls.Config{}

	dialInfo := &mgo.DialInfo{
		Addrs: []string{"cluster0-shard-00-00.k2lrx.mongodb.net:27017",
			"cluster0-shard-00-01.k2lrx.mongodb.net:27017",
			"cluster0-shard-00-02.k2lrx.mongodb.net:27017"},
		Database: "admin",
		Username: cfg.Username,
		Password: cfg.Password,
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	var err error
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		logrus.Fatalf("[mgo Dial] addr: %s | err: %s", cfg.Addr, err)
	}
	session.SetMode(mgo.PrimaryPreferred, false)

	if err = session.Ping(); err != nil {
		logrus.Fatalf("[mgo Ping] addr: %s | err: %s", cfg.Addr, err)
	}
}
