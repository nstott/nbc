package nbc

import (
	"crypto/md5"
	"fmt"
	"io"
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
