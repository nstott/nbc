package main

import (
	"flag"
	"fmt"
	"math"
	"os"
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

    doc := NewDocument()
    doc.TokenizeFile(*filename)	
	doc.GenerateNGrams(*nGramSize, *class)

	if *train {

		if *verbose { // dump out the ngrams we've discovered
			for _, v := range doc.ngrams {
				fmt.Printf("%d -> %s\n", v.Count[*class], v.Hash )
			}
		}
		doc.DumpToMongo()

	} else {

		classCount := CountDistinctNGrams()
		cb := GetClassProbabilities()

		if *verbose {
			for k, v := range cb {
				fmt.Printf("P(%s) = %f\n", k, v)
			}
		}

		for class, v := range cb {
			totalngrams := GetTotalNGrams(class)
			probabilities := make([]float64, doc.totalNgrams)
			idx := 0
			for _, v := range doc.ngrams {
				instanceCount := v.GetInstanceCount(class)
				probabilities[idx] = laplaceSmoothing(instanceCount, totalngrams, classCount)

				if *verbose {
					fmt.Printf("P(%s|%s) = (%d+1)/(%d+%d) = %f\n", 
						class, v.Hash, instanceCount, totalngrams, classCount, probabilities[idx] )
				}
				idx += 1
			}
			p := totalProbability(probabilities, v)
			fmt.Printf("P(%s|Message) = %f\n", class, p)
		}

	}
}

func totalProbability(probabilities []float64, classProbability float64) float64 {
	ret := classProbability
	for _, v := range probabilities {
		ret += math.Log(v)
	}	
	return ret
}

func laplaceSmoothing(n int, N int, classCount int) float64 {
	return ( float64(n) + *laplaceConstant ) / ( float64(N) + float64(classCount) )
}
