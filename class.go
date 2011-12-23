package main

import (
	"launchpad.net/gobson/bson"
)



type ClassData struct {
	Name string
	Count int
}

func (c *ClassData) Update() {
	collection := getClassCollection()
	err := collection.Update(bson.M{"class": c.Name}, bson.M{"$inc": bson.M{"count": 1}})
	if (err != nil) {
		err = collection.Insert(&ClassData{c.Name, 1})
	}
}

// func GetClassTotals() (map[string]int, int) {
// 	collection := getClassCollection()
// 	var result ClassData

// 	counts := make(map[string]int)
// 	var total int

// 	iter := collection.Find(bson.M{}).Limit(100).Iter()
// 	for iter.Next(&result) {
// 		total += result.Count
// 		counts[result.Name] = result.Count
//     }	

//     return counts, total
// }

// func classProbabilities(counts map[string]int, total int) map[string]float64 {
// 	var classCount = len(counts)
// 	probabilities := make(map[string]float64)
// 	for k, v := range counts {
// 		probabilities[k] = laplaceSmoothing(v, total, classCount)
// 	} 
// 	return probabilities
// }

func GetClassProbabilities() map[string]float64 {
	collection := getClassCollection()
	var result ClassData

	counts := make(map[string]int)
	var total int

	iter := collection.Find(bson.M{}).Limit(100).Iter()
	for iter.Next(&result) {
		total += result.Count
		counts[result.Name] = result.Count
    }	

	classCount := len(counts)
	probabilities := make(map[string]float64)

	for k, v := range counts {
		probabilities[k] = laplaceSmoothing(v, total, classCount)
	} 
	return probabilities
}