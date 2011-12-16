package main

import (
	"fmt"
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

type MongoNGram struct {
	count map[string]int // class : key, count : int
	probability map[string]float64 // class : key, prob: float
	hash string
}

func mongoConnect() *mgo.Session {
	var err os.Error
	session, err = mgo.Mongo(mongoHost)
    if err != nil {
            panic(err)
    }
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

// func addNgram(ngram NGram, class string) {
// 	item := &MongoNGram{1, 0, ngram.genhash}
// 	fmt.Println(item)	
// }