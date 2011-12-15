package main

type MongoNGram struct {
	count map[string]int // class : key, count : int
	probability map[string]float64 // class : key, prob: float
	ngram NGram
}


