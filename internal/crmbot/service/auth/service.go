package auth

import (
	"crypto/tls"
	"net"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"github.com/caarlos0/env"
	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	UsersColl *mgo.Collection
}

var mgoCFG config.MongoCfg

func NewAuthService(cfg Config) (*AuthService, error) {
	if err := env.Parse(&cfg); err != nil {
		logrus.Fatalf("[env Parse] MongoCfg err: %s", err)
	}

	tlsConfig := &tls.Config{}

	dialInfo := &mgo.DialInfo{
		Addrs: []string{"cluster0-shard-00-00.k2lrx.mongodb.net:27017",
			"cluster0-shard-00-01.k2lrx.mongodb.net:27017",
			"cluster0-shard-00-02.k2lrx.mongodb.net:27017"},
		Database: "admin",
		Username: mgoCFG.Username,
		Password: mgoCFG.Password,
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		logrus.Fatalf("[NewAuthService] mgo Dial | err: %s", err)
	}
	session.SetMode(mgo.PrimaryPreferred, false)

	if err = session.Ping(); err != nil {
		return nil, err
	}

	return &AuthService{
		UsersColl: session.DB(mgoCFG.DBName).C(cfg.UsersCollName),
	}, nil
}
