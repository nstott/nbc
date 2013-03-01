package main

import (
	// "fmt"
	"io/ioutil"
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
	return strings.Join(in, " ")
}

/* exists performs a query against mongo to detirmine if a specific ngram exists in the database or not.
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

// Document holds the text that we are training with or classifying
type Document struct {
	filename string
	tokens []string
	// ngrams []nGram
	totalNgrams int
	class *ClassData
	ngrams map[string]nGram

}

// NewDocument creates a new Document 
func NewDocument() *Document {
	d := &Document{}
	d.class = &ClassData{}
	return d
}

// TokenizeFile reads a file from disk and tokenizes it by splitting on spaces
func (d *Document) TokenizeFile(fn string) {
	d.filename = fn 
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	d.tokens = strings.Fields(string(data))
}

// GenerateNGrams organizes the already tokenized text into ngrams of a specified size and class
func (d *Document) GenerateNGrams(n int, class string) {
	d.class.Name = class
	out := make([]nGram, 0)
	for i := 0; i <= len(d.tokens) - n; i += 1 {
		out = append(out, NewNGram(n, d.tokens[i:i+n], class))
	}
	d.totalNgrams = len(out)

	d.ngrams = make(map[string]nGram)	 

	for _, v := range out {
		_, ok := d.ngrams[v.Hash]
		if ok {
			d.ngrams[v.Hash].Count[class]++
		} else {
			v.Count[class] = 1
			d.ngrams[v.Hash] = v
		}
	}
}

// DumpToMongo commits the Document to mongo
func (d *Document)DumpToMongo() {
	collection := getCollection()	
	field := "cound." + d.class.Name
		
	for _, ngram := range d.ngrams {
		if ngram.exists() {
			q := bson.M{"hash": ngram.Hash}
			err := collection.Update(q, bson.M{"$inc": bson.M{field: 1}})
			if err != nil {
				panic(err)
			}
		} else {
			// straight up insert
			collection.Insert(ngram)
		}
	}
	d.class.Update()
}

