package nbc

import (
	"strings"
	"io/ioutil"
)

// Document holds the text that we are training with or classifying
type Document struct {
	filename string
	tokens []string
	totalNgrams int
	class *Classification
	ngrams map[string]*nGram

}

func NewDocument(class Classification) Document {
	d := Document{}
	c := NewClassification(class.Name)
	d.class = &c
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

func (d *Document) TokenizeString(s string) {
	d.tokens = strings.Fields(s)
}

// GenerateNGrams organizes the already tokenized text into ngrams of a specified size and class
func (d *Document) GenerateNGrams(n int) {
	out := make([]*nGram, 0)
	for i := 0; i <= len(d.tokens) - n; i += 1 {
		out = append(out, NewNGram(d.tokens[i:i+n]))
	}
	d.totalNgrams = len(out)
	for _, v := range out {
		_, ok := d.class.ngrams[v.Hash]
		if ok {
			d.class.ngrams[v.Hash].Count++
		} else {
			d.class.ngrams[v.Hash] = v
		}
	}
}
