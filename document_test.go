package main

import (
	"testing"
)

func Test_TokenizeFile(t *testing.T) {
	var d = []struct{
		in string
		want []string
	}{
		{"corpus/Test_TokenizeFile.txt", []string{"play", "sports", "today"}},
	}

	for _, v := range d {
		doc := NewDocument()
		doc.TokenizeFile(v.in)
		if !equalStringSlice(doc.tokens, v.want) {
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
		doc := NewDocument()
		doc.TokenizeString(v.in)
		if !equalStringSlice(doc.tokens, v.want) {
			t.Errorf("TokenizeFile(%s) != %v, got %v", v.in, v.want, doc.tokens)
		}
	}
}

func Test_GenerateNGrams(t *testing.T) {
	class := "class1"
	num := 2

	var d = []struct{
		in string
		want map[string]nGram
	}{
		{"play sports today", map[string]nGram{
			"8e332df73afd1944b529f1ee94eb0d7d": nGram{Length: num, Tokens: []string{"play", "sports"}, Hash: "8e332df73afd1944b529f1ee94eb0d7d", Count: map[string]int{class: 1}},
			"d3364f66e254f86cfef25c00cb30fe59": nGram{Length: num, Tokens: []string{"sports", "today"}, Hash: "d3364f66e254f86cfef25c00cb30fe59", Count: map[string]int{class: 1}},
			},
		}, 
		{"play play play sports today", map[string]nGram{
			"8e332df73afd1944b529f1ee94eb0d7d": nGram{Length: num, Tokens: []string{"play", "sports"}, Hash: "8e332df73afd1944b529f1ee94eb0d7d", Count: map[string]int{class: 1}},
			"d3364f66e254f86cfef25c00cb30fe59": nGram{Length: num, Tokens: []string{"sports", "today"}, Hash: "d3364f66e254f86cfef25c00cb30fe59", Count: map[string]int{class: 1}},
			"ec7841687efd9cf97ac07f0c80c48e8e": nGram{Length: num, Tokens: []string{"play", "play"}, Hash: "ec7841687efd9cf97ac07f0c80c48e8e", Count: map[string]int{class: 2}},
			}, 
		},
	}

	for _, v := range d {
		doc := NewDocument()
		doc.TokenizeString(v.in)
		doc.GenerateNGrams(num, class)

		if !equalNgramMap(doc.ngrams, v.want) {
			t.Errorf("TokenizeFile(%s) \n\t%v, got \n\t%v", v.in, v.want, doc.ngrams)
		}
	}
}


// TODO use reflect.DeepEqual instead

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalNgramMap(a, b map[string]nGram) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if !v.equal(b[k]) {
			return false
		}
	}
	return true
}

func (a *nGram) equal(b nGram) bool {
	if a.Length != b.Length || a.Hash != b.Hash || a.Length != b.Length {
		return false
	}

	if len(a.Count) != len(b.Count) {
		return false
	}
	for k, v := range a.Count {
		if b.Count[k] != v {
			return false
		}
	}

	if len(a.Tokens) != len(b.Tokens) {
		return false
	}

	for i := range a.Tokens {
		if a.Tokens[i] != b.Tokens[i] {
			return false
		}
	}

	return true
}


