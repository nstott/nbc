package main

import (
	"fmt"
	"crypto/md5"
	)
/* an ngram */
type NGram struct {
	length int
	tokens []string
	hash []byte
}

func genhash(in []string) []byte {
	md5 := md5.New()
	md5.Write([]byte(fmt.Sprintf("%x", in)))
	return md5.Sum()
}

func GenerateNGrams(in []string, n int) []NGram {
	out := make([]NGram, 0)
	for i := 0; i <= len(in) - n; i += 1 {
		out = append(out, NGram{n, in[i:i+n], genhash(in[i:i+n])})
	}

	for i := 0; i < len(out); i++ {
		fmt.Println(out[i])
	}
	return out
}
