package nbc

import (
	"reflect"
	"testing"
)

var dummyClass = NewClassification("ugh")

func Test_TokenizeFile(t *testing.T) {
	var d = []struct{
		in string
		want []string
	}{
		{"corpus/Test_TokenizeFile.txt", []string{"play", "sports", "today"}},
	}

	for _, v := range d {
		doc := NewDocument(dummyClass)
		doc.TokenizeFile(v.in)
		if !reflect.DeepEqual(doc.tokens, v.want) {
			t.Errorf("TokenizeFile(%s) != %v, got %v", v.in, v.want, doc.tokens)
		}
	}
}

func Test_TokenizeString(t *testing.T) {
	var d = []struct{
		in string
		want []string
	}{
		{"play sports today", []string{"play", "sports", "today"}},
	}

	for _, v := range d {
		doc := NewDocument(dummyClass)
		doc.TokenizeString(v.in)
		if !reflect.DeepEqual(doc.tokens, v.want) {
			t.Errorf("TokenizeFile(%s) != %v, got %v", v.in, v.want, doc.tokens)
		}
	}
}

func Test_GenerateNGrams(t *testing.T) {

	num := 2

	var d = []struct{
		in string
		want map[string]nGram
	}{
		{"play sports today", map[string]nGram{
			"8e332df73afd1944b529f1ee94eb0d7d": nGram{Length: num, Tokens: []string{"play", "sports"}, Hash: "8e332df73afd1944b529f1ee94eb0d7d", Count: map[string]int{dummyClass.Name: 1}},
			"d3364f66e254f86cfef25c00cb30fe59": nGram{Length: num, Tokens: []string{"sports", "today"}, Hash: "d3364f66e254f86cfef25c00cb30fe59", Count: map[string]int{dummyClass.Name: 1}},
			},
		}, 
		{"play play play sports today", map[string]nGram{
			"8e332df73afd1944b529f1ee94eb0d7d": nGram{Length: num, Tokens: []string{"play", "sports"}, Hash: "8e332df73afd1944b529f1ee94eb0d7d", Count: map[string]int{dummyClass.Name: 1}},
			"d3364f66e254f86cfef25c00cb30fe59": nGram{Length: num, Tokens: []string{"sports", "today"}, Hash: "d3364f66e254f86cfef25c00cb30fe59", Count: map[string]int{dummyClass.Name: 1}},
			"ec7841687efd9cf97ac07f0c80c48e8e": nGram{Length: num, Tokens: []string{"play", "play"}, Hash: "ec7841687efd9cf97ac07f0c80c48e8e", Count: map[string]int{dummyClass.Name: 2}},
			}, 
		},
	}

	for _, v := range d {
		doc := NewDocument(dummyClass)
		doc.TokenizeString(v.in)
		doc.GenerateNGrams(num)

		if !reflect.DeepEqual(doc.ngrams, v.want) {
			t.Errorf("TokenizeFile(%s) \n\t%v, got \n\t%v", v.in, v.want, doc.ngrams)
		}
	}
}
