package main

import (
	// "fmt"
	"io/ioutil"
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
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


// does an ngram exist
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


func (n *nGram) GetInstanceCount(class string) int {
	collection := getCollection()
	var ngram nGram
	err := collection.Find(bson.M{"hash": n.Hash}).One(&ngram)
	if err != nil {
		return 0
	}
	return ngram.Count[class]
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

func CountDistinctNGrams() int {
	collection := getCollection()
	count, err := collection.Find(bson.M{}).Count()
	if err != nil {
		panic(err)
	}
	return count
}




type Document struct {
	filename string
	tokens []string
	// ngrams []nGram
	totalNgrams int
	class *ClassData
	ngrams map[string]nGram

}

func NewDocument() *Document {
	d := &Document{}
	d.class = &ClassData{}
	return d
}

func (d *Document) TokenizeFile(fn string) {
	d.filename = fn 
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	d.tokens = strings.Fields(string(data))
}

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

