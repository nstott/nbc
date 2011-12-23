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

// mongoConnect sets up a global *mgo.Session object, and ensures that the main collection is indexed
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

// mongoDisconnect disconnects from the underlying data source
func mongoDisconnect() {
	session.Close()
}

// getCollection returns the mongo collection that is used to store the ngram data
func getCollection() mgo.Collection {
	return session.DB(mongoDB).C(*collection)	
}

//getClassCollection returns the mongo collection used to store information about 
// the different classes that have been learned so far
func getClassCollection() mgo.Collection {
	return session.DB(mongoDB).C(*collection + "_classes")
}

// indexCollection ensures that the mongo collection for the data is properly indexed
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


// forgetData clears the data out of the mongo collection
func forgetData() {
	c := getCollection()
	c.RemoveAll(bson.M{})

	c = getClassCollection()
	c.RemoveAll(bson.M{})
}
