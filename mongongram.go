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


func exists(ngram nGram) bool {
	collection := getCollection()

	if ngram.hash == "" {
		panic("NGram has unitialized Hash") 
	} 

	c, err := collection.Find(bson.M{"hash": ngram.hash}).Count()
	if err != nil || c == 0 {
		return false
	}
	return true
}

func addNgram(ngram nGram, class string) {
	if exists(ngram) {
		// the hash already exists, update the counts
	} else {
		// insert this ngram into the ddb
	}
	// fmt.Printf("does ngram (%v) exist? %v\n", ngram.tokens, exists(ngram));
}

func dumpNGramsToMongo(ngrams map[string]nGram) {
	for _, v := range ngrams {
		
	}
}

