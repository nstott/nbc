package nbc

import (
	"fmt"
	"math"
)

type Classifier struct {
	engine StorageEngine
	laPlaceConstant int
	nGramSize int
}

func newClassifier(engine StorageEngine, laplaceConstant, nGramSize int) Classifier {
	engine.Setup()
	return Classifier{engine, laplaceConstant, nGramSize}
}

type StorageEngine interface {
	Setup() 
	TearDown()
	DumpDocument(d *Document)
	GetInstanceCount(n *nGram, class string) int
	CountDistinctNGrams() int
	GetTotalNGrams(class string) int
	GetClassProbabilities() map[string]float64
}

func (classifier *Classifier) TrainFile(filename string, class Classification) {
    doc := classifier.processFile(filename, class)
	classifier.engine.DumpDocument(&doc)
}

func (classifier *Classifier) ClassifyFile(filename string, class Classification) {
	doc := classifier.processFile(filename, class)
	classCount := classifier.engine.CountDistinctNGrams()
	cb := classifier.engine.GetClassProbabilities()

	for class, v := range cb {
		totalngrams := classifier.engine.GetTotalNGrams(class)
		probabilities := make([]float64, doc.totalNgrams)
		idx := 0
		for _, v := range doc.ngrams {
			instanceCount := classifier.engine.GetInstanceCount(&v, class)
			probabilities[idx] = laplaceSmoothing(instanceCount, totalngrams, classCount)

			fmt.Printf("P(%s|%s) = (%d+1)/(%d+%d) = %f\n", 
				class, v.Hash, instanceCount, totalngrams, classCount, probabilities[idx] )
			idx += 1
		}
		p := totalProbability(probabilities, v)
		fmt.Printf("P(%s|Message) = %f\n", class, p)
	}
}

func (classifier *Classifier) processFile(filename string,  class Classification) Document {
	doc := NewDocument(class)
    doc.TokenizeFile(filename)	
	doc.GenerateNGrams(classifier.nGramSize)
	return doc
}

func totalProbability(probabilities []float64, classProbability float64) float64 {
	ret := classProbability
	for _, v := range probabilities {
		ret += math.Log(v)
	}	
	return ret
}

func laplaceSmoothing(n, N int, classCount int) float64 {
	K := 1.0
	return ( float64(n) + K) / ( float64(N) + float64(classCount) )
}
