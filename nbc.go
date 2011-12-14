package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	)

func main() {

	data, err := ioutil.ReadFile("nbc.go")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("data: \n%s\n",data)

	d1 := tokenize(string(data))
	for i := 0; i < len(d1); i++ {
		fmt.Printf("thing: %s\n", strings.TrimSpace(d1[i]))
	}


	d2, err := parseRootDir("./")

	for i := 0; i < len(d2); i++ {
		fmt.Printf("%s\n", d2[i])
	}

}


func tokenize(str string) []string {
	fmt.Printf("In Tokenize: %s\n", str)
	return strings.Fields(str)
}

func parseRootDir(root string) ([]string, error) {
	fi, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	// we return strings and not file info's
	var ret = make([]string, len(fi))

	for i := 0; i < len(fi); i++ {
		ret[i] = fi[i].Name()
	}

	return ret, nil
}