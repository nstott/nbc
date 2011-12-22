package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	// "launchpad.net/gobson/bson"
	"os"
	"strings"
)

var train *bool 		= flag.Bool("train", true, "training mode")
var class *string 		= flag.String("class", "true", "The class associated with this training set")
var filename *string 	= flag.String("filename", "./nbc.go", "the filename to read from in training mode")
var forget 				= flag.Bool("nuke", false, "forget the learned data")
var collection 			= flag.String("collection", "data", "The db collection to use")
var laplaceConstant 	= flag.Float64("k", 1, "The laplacian smoothing constant to use")
var nGramSize			= flag.Int("n", 3, "The size of the ngrams")
var verbose				= flag.Bool("v", false, "Be verbose")

func main() {
	flag.Parse()

	mongoConnect()
    defer mongoDisconnect()

    if *forget {
    	fmt.Printf("Forgetting learned data in %s.%s\n",mongoDB,mongoCollection)
    	forgetData()
    	os.Exit(0)
    }

	data, err := readFile(*filename)
	if err != nil {
		panic(err)
	}

	d1 := tokenize(string(data))
	ngrams := AggregateNGrams(GenerateNGrams(d1, *nGramSize, *class), *class)

	if *train {

		for _, v := range ngrams {
			fmt.Printf("%d -> %s\n", v.Count[*class], v.Hash )
		}
		dumpNGramsToMongo(ngrams, *class)
		updateClass(*class, 1)
		counts, total := getClassTotals()
		probabilities := classProbabilities(counts, total)
		fmt.Printf("%v\n",probabilities)

	} else {

		classCount := getClassCount()

		counts, total := getClassTotals()
		cb := classProbabilities(counts, total)
		for class, v := range cb {
			totalngrams := getTotalNGrams(class)
			probabilities := make([]float64, len(ngrams))
			idx := 0
			for _, v := range ngrams {
				instanceCount := getInstanceCount(v.Hash, class)
				probabilities[idx] = laplaceSmoothing(instanceCount, totalngrams, classCount)

				if *verbose {
					fmt.Printf("%s: %d + 1 / %d + %d = %f\n", 
						v.Hash, instanceCount, totalngrams, classCount, probabilities[idx] )
				}
				idx += 1
			}
			p := totalProbability(probabilities, v)
			fmt.Printf("Processing Class: %s = %f\n", class, p)
		}

	}
}


func readFile(str string) (string, os.Error) {
	data, err := ioutil.ReadFile(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func tokenize(str string) []string {
	return strings.Fields(str)
}

func laplaceSmoothing(n int, N int, classCount int) float64 {
	return ( float64(n) + *laplaceConstant ) / ( float64(N) + float64(classCount) )
}








