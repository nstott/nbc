package main

import (
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
	"os"
)

const (
	mongoHost = "localhost"
	mongoDB = "nbc"
	mongoCollection = "data"
)

var session *mgo.Session

func mongoConnect() *mgo.Session {
	var err os.Error
	session, err = mgo.Mongo(mongoHost)
    if err != nil {
        panic(err)
    }
    session.SetSafe(&mgo.Safe{})
    indexCollection()
    return session
}

func mongoDisconnect() {
	session.Close()
}

func getCollection() mgo.Collection {
	return session.DB(mongoDB).C(mongoCollection)	
}

func indexCollection() {
	col := getCollection()
	index := mgo.Index{
		Key: []string{"hash"},
		Unique: true,
		DropDups: true,
		Background: true,
		Sparse: true,
	}
	_ = col.EnsureIndex(index)
}

func forgetData() {
	c := getCollection()
	c.RemoveAll(bson.M{"name": 1})
}
