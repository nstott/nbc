package main

import (
	"fmt"
	"launchpad.net/gobson/bson")

type DataClass struct {
	ClassCounts bool
	Class string
	Count int
}

type Person struct {
	Name string
	Phone string
}


func updateClass(class string, count int) {
	collection := getCollection()
	err := collection.Update(bson.M{"classcounts": true, "class": class}, bson.M{"$inc": bson.M{"count": 1}})
	if (err != nil) {
		err = collection.Insert(&DataClass{true, class, 1})
	}
}

func getClassTotals() (map[string]int, int) {
	collection := getCollection()
	var result DataClass

	counts := make(map[string]int)
	var total int

	iter := collection.Find(bson.M{"classcounts": true}).Limit(100).Iter()
	for iter.Next(&result) {
		total += result.Count
		counts[result.Class] = result.Count
    	fmt.Printf("%s: %d\n", result.Class, result.Count)
    }	

    return counts, total
}