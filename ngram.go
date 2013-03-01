package main

import (
	// "fmt"
	"crypto/md5"
	"fmt"
	"io"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"strings"
	)

/* nGram is a collection of n tokens.  Count refers to the 
 * number of times this ngram has been seen in a specific document/class.  
 * ngrams are considered unique based upon their hash
 */  
type nGram struct {
	Length int
	Tokens []string
	Hash string
	Count map[string]int
}

// NewNGram returns a new ngram
func NewNGram(n int, tokens []string, class string) nGram  {
	return nGram{n, tokens, genhash(tokens), map[string]int{class: 1}}
}

/*
 * genhash is a hashing function that returns a unique representation of the tokens.  the hash also serves as the primary key in the mongo collection.
 * currently this is just the tokens joined together,
 */ 
func genhash(in []string) string {
	h := md5.New()
	io.WriteString(h, strings.Join(in, ""))
	return fmt.Sprintf("%x", h.Sum([]byte{}))
}

/* exists performs a query against mongo to determine if a specific ngram exists in the database or not.
 * this is part of an 'insert on duplicate key update' type situation
 */
func (n *nGram) exists() bool {
	collection := getCollection()

	if n.Hash == "" {
		panic("NGram has unitialized Hash") 
	} 

	c, err := collection.Find(bson.M{"hash": n.Hash}).Count()
	if err != nil || c == 0 {
		return false
	}
	return true
}

// GetInstanceCount returns the number of times an ngram has been seen in a class
func (n *nGram) GetInstanceCount(class string) int {
	collection := getCollection()
	var ngram nGram
	err := collection.Find(bson.M{"hash": n.Hash}).One(&ngram)
	if err != nil {
		return 0
	}
	return ngram.Count[class]
}

// getTotalNGrams returns the total number of times a specific ngram has been seen in a class across all documents
func GetTotalNGrams(class string) int {
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

/* 
 * CountDistinctNGrams returns the size of our vocabulary.  
 * i.e. the number of distinct ngrams across all documents and classes
 */ 
func CountDistinctNGrams() int {
	collection := getCollection()
	count, err := collection.Find(bson.M{}).Count()
	if err != nil {
		panic(err)
	}
	return count
}

