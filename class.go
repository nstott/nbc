package main

import (
	"fmt"
	"launchpad.net/gobson/bson")

type class struct {
	class string
	count int
}

func updateClass(class string, count int) {
	collection := getCollection()
	err := collection.Update(bson.M{"class": class}, bson.M{"$inc": bson.M{"count": 1}})
	if (err != nil) {
		collection.Insert(bson.M{"class": class, "count": 1})
		fmt.Printf("Err inserting class: %v", err)
	}
}