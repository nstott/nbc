package main

import (
	"launchpad.net/gobson/bson")

type DataClass struct {
	Class string
	Count int
}



func updateClass(class string, count int) {
	collection := getClassCollection()
	err := collection.Update(bson.M{"class": class}, bson.M{"$inc": bson.M{"count": 1}})
	if (err != nil) {
		err = collection.Insert(&DataClass{class, 1})
	}
}

func getClassTotals() (map[string]int, int) {
	collection := getClassCollection()
	var result DataClass

	counts := make(map[string]int)
	var total int

	iter := collection.Find(bson.M{}).Limit(100).Iter()
	for iter.Next(&result) {
		total += result.Count
		counts[result.Class] = result.Count
    }	

    return counts, total
}

func classProbabilities(counts map[string]int, total int) map[string]float64 {
	var classCount = len(counts)
	probabilities := make(map[string]float64)
	for k, v := range counts {
		probabilities[k] = laplaceSmoothing(v, total, classCount)
	} 
	return probabilities
}