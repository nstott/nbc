package main

import (
	"strings"
	"io/ioutil"
	"launchpad.net/mgo/bson"
)

// Document holds the text that we are training with or classifying
type Document struct {
	filename string
	tokens []string
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
