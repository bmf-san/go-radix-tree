package main

import (
	"fmt"

	"github.com/bmf-san/go-radix-tree"
)

func main() {
	// Create a tree
	r := radix.New()
	// r.Insert("", "v-/")
	r.Insert("/", "v-/")
	r.Insert("/foo", "v-foo")
	r.Insert("/fo", "v-fo")
	r.Insert("/bar", "v-bar")
	r.Insert("/foo/bar", "v-foo/bar")
	r.Insert("/foobar", "v-foobar")

	// Find the longest prefix match
	// m, _, _ := r.LongestPrefix("/f")
	// fmt.Printf("%#v\n", m)

	g, _ := r.Get("/foobar")
	fmt.Printf("%#v\n", g)
}
