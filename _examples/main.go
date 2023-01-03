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

	g := r.Get("/foobar")
	fmt.Printf("%#v\n", g)
}
