package main

import (
	"fmt"
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
	"math"
	"strings"
	)

/* an ngram */
type nGram struct {
	Length int
	Tokens []string
	Hash string
	Count map[string]int
}

func NewNGram(n int, tokens []string, class string) nGram  {
	return nGram{n, tokens, genhash(tokens), map[string]int{class: 1}}
}

func genhash(in []string) string {
	return strings.Join(in, " ")
}

func GenerateNGrams(in []string, n int, class string) []nGram {
	out := make([]nGram, 0)
	for i := 0; i <= len(in) - n; i += 1 {
		out = append(out, NewNGram(n, in[i:i+n], class))
	}
	return out
}


func AggregateNGrams(ngrams []nGram, class string) map[string]nGram {
	ret := make(map[string]nGram)
	var mng nGram // declare these here
	var ok bool	 

	for _, v := range ngrams {
		mng, ok = ret[v.Hash]
		if ok {
			mng.Count[class]++
		} else {
			v.Count[class] = 1
			ret[v.Hash] = v
		}
	}
	return ret
}

// does an ngram exist
func exists(ngram nGram) bool {
	collection := getCollection()

	if ngram.Hash == "" {
		panic("NGram has unitialized Hash") 
	} 

	c, err := collection.Find(bson.M{"hash": ngram.Hash}).Count()
	if err != nil || c == 0 {
		return false
	}
	return true
}


func dumpNGramsToMongo(ngrams map[string]nGram, class string) {
	collection := getCollection()		
	for _, ngram := range ngrams {
		if exists(ngram) {
			fmt.Println("This should be an upsert, but we're not that smart yet")
			q := bson.M{"hash": ngram.Hash}
			err := collection.Update(q, bson.M{"$inc": bson.M{"count."+class: 1}})
			if err != nil {
				fmt.Printf("err: %s\n", err)
			}
		} else {
			// straight up insert
			collection.Insert(ngram)
		}
	}
}


func getTotalNGrams(class string) int {

	collection := getCollection()
	var field = "count." + class
	job := mgo.MapReduce{
        Map:      "function() { emit(\"total\", this.count."+class+")}",
        Reduce:   "function(key, values) { var t = 0; values.forEach(function (i) {t += i});return t; }",
	}
	var result []struct { Id string "_id"; Value int }
	q := collection.Find(bson.M{field: bson.M{"$gt": 0}})
	_, err := q.MapReduce(job, &result)
	if err != nil  {
	    panic(err)
	}
	if len(result) > 0 {
		return result[0].Value
	}
	return 0;
}

func getClassCount() int {
	collection := getCollection()
	count, err := collection.Find(bson.M{}).Count()
	if err != nil {
		panic(err)
	}
	return count
}

func getInstanceCount(hash, class string) int {
	collection := getCollection()
	var ngram nGram
	err := collection.Find(bson.M{"hash": hash}).One(&ngram)
	if err != nil {
		return 0
	}
	return ngram.Count[class]
}

func totalProbability(probabilities []float64, classProbability float64) float64 {
	ret := classProbability
	for _, v := range probabilities {
		ret += math.Log(v)
	}	
	return ret
}