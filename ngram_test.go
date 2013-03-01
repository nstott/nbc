package main

import (
	"log"
	"testing"
)

func Test_genhash(t *testing.T) {
	var d = []struct{
		in []string
		want string
	}{
		{[]string{"one", "two", "three"}, "7b0391feb2e0cd271f1cf39aafb4376f"},
		{[]string{"four"}, "8cbad96aced40b3838dd9f07f6ef5772"},
		{[]string{""}, "d41d8cd98f00b204e9800998ecf8427e"},
		{[]string{"nick nick"}, "c0d30f94487234173e100e7861c57b5d"},
	}

	for _, v := range d {
		out := genhash(v.in)
		if out != v.want {
			t.Errorf("genHash(%s) != %s, got %s instead", v.in, v.want, out)
		}
		log.Println(out)
	}
}