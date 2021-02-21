package auth

import (
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
)

func (s *AuthService) IsUser(tgID int) bool {
	query := bson.M{
		"tg_id": tgID,
	}

	if count, err := s.UsersColl.Find(query).Count(); err != nil {
		logrus.Fatalf("[authService] IsUser | err: %s", err)
	} else if count > 0 {
		return true
	}
	return false
}

func (s *AuthService) IsAdmin(tgID int) bool {
	query := bson.M{
		"tg_id":  tgID,
		"access": "admin",
	}

	if count, err := s.UsersColl.Find(query).Count(); err != nil {
		logrus.Fatalf("[authService] IsUser | err: %s", err)
	} else if count > 0 {
		return true
	}
	return false
}
