package model

import (
	"gopkg.in/mgo.v2"
)

type MongoStore struct {
	Session *mgo.Session
}