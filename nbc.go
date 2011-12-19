package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var train *bool = flag.Bool("train", true, "training mode")
var trainingClass *string = flag.String("class", "true", "The class associated with this training set")
var trainingFilename *string = flag.String("filename", "./nbc.go", "the filename to read from in training mode")
var forget = flag.Bool("nuke", false, "forget the learned data")



func main() {
	flag.Parse()

	mongoConnect()
    defer mongoDisconnect()

    if *forget {
    	fmt.Printf("Forgetting learned data in %s.%s\n",mongoDB,mongoCollection)
    	forgetData()
    }

	if *train {
		data, err := readFile(*trainingFilename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// universe := make(
		d1 := tokenize(string(data))
		ngrams := AggregateNGrams(GenerateNGrams(d1, 3, *trainingClass), *trainingClass)
		for _, v := range ngrams {
			fmt.Printf("%d -> %s\n", v.count[*trainingClass], v.hash )
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

func parseRootDir(root string) ([]string, os.Error) {
	fi, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	// we return strings and not file info's
	var ret = make([]string, len(fi))

	for i := 0; i < len(fi); i++ {
		ret[i] = fi[i].Name
	}

	return ret, nil
}









