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
	Tokens []string
	Hash string
	Count int
}

func NewNGram(tokens []string) *nGram  {
	return &nGram{tokens, genhash(tokens), 1}
}

/*
 * genhash is a hashing function that returns a unique representation of the tokens.  
 * The hash also serves as the primary key in the mongo collection.
 * currently this is just the tokens joined together,
 */ 
func genhash(in []string) string {
	h := md5.New()
	io.WriteString(h, strings.Join(in, ""))
	return fmt.Sprintf("%x", h.Sum([]byte{}))
}
