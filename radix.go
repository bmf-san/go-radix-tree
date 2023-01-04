// It refers to the implementation of go-radix
// See: https://github.com/armon/go-radix
package radix

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// paramsPool is a pool for parameters
var paramsPool sync.Pool

// Tree is a Radix tree.
type Tree struct {
	root *node
}

// node is a node in tree.
type node struct {
	leaf     *leafNode
	prefix   string
	children children
}

// leafNode is the node that doesn't have a node.
type leafNode struct {
	key       string
	val       string
	hasParams bool
}

// child is the nodes that have a node.
type child struct {
	label byte
	node  *node
}

// children are the nodes that have a node.
type children []child

// addChild adds child to a node.
func (n *node) addChild(c child) {
	num := len(n.children)
	idx := sort.Search(num, func(i int) bool {
		return n.children[i].label >= c.label
	})

	n.children = append(n.children, child{})
	copy(n.children[idx+1:], n.children[idx:])
	n.children[idx] = c
}

// updateChild updates a child in a node by label.
func (n *node) updateChild(label byte, node *node) {
	num := len(n.children)
	idx := sort.Search(num, func(i int) bool {
		return n.children[i].label >= label
	})
	if idx < num && n.children[idx].label == label {
		n.children[idx].node = node
		return
	}
	panic("replacing missing child")
}

// getChild gets a child from a node by label.
func (n *node) getChild(label byte) *node {
	num := len(n.children)
	idx := sort.Search(num, func(i int) bool {
		return n.children[i].label >= label
	})
	if idx < num && n.children[idx].label == label {
		return n.children[idx].node
	}
	return nil
}

// New creates a Tree.
func New() *Tree {
	return &Tree{root: &node{}}
}

// longestPrefix returns prefix length of two strings.
func longestPrefix(k1, k2 string) int {
	max := len(k1)
	if l := len(k2); l < max {
		max = l
	}
	var i int
	for i = 0; i < max; i++ {
		if k1[i] != k2[i] {
			break
		}
	}
	return i
}

// Insert inserts a key and value to tree.
func (t *Tree) Insert(k, v string) {
	var parent *node
	n := t.root
	path := k
	for {
		parent = n
		n = n.getChild(path[0])

		if n == nil {
			parent.addChild(child{
				label: path[0],
				node: &node{
					leaf: &leafNode{
						key: k,
						val: v,
					},
					prefix: path,
				},
			})
			return
		}

		commonPrefix := longestPrefix(path, n.prefix)
		if commonPrefix == len(n.prefix) {
			if path == n.prefix {
				panic(fmt.Sprintf("duplicate path registration. path: %v prefix: %v", path, n.prefix))
			}
			path = path[commonPrefix:]
			continue
		}

		if strings.Count(path, "/") == strings.Count(n.prefix, "/") && strings.Contains(n.prefix, ":") {
			panic(fmt.Sprintf("conflicts path parameter. path: %v prefix: %v", path, n.prefix))

		}

		spln := &node{
			prefix: path[:commonPrefix],
		}
		parent.updateChild(path[0], spln)

		spln.addChild(child{
			label: n.prefix[commonPrefix],
			node:  n,
		})
		n.prefix = n.prefix[commonPrefix:]

		leaf := &leafNode{
			key: k,
			val: v,
		}

		path = path[commonPrefix:]
		if len(path) == 0 {
			spln.leaf = leaf
			return
		}

		spln.addChild(child{
			label: path[0],
			node: &node{
				leaf:   leaf,
				prefix: path,
			},
		})
		return
	}
}

// Get gets a value from tree by key.
func (t *Tree) Get(k string) string {
	n := t.root
	path := k
	for {
		if len(path) == 0 {
			if n.leaf != nil {
				return n.leaf.val
			}
			break
		}

		n = n.getChild(path[0])
		if n == nil {
			break
		}

		if strings.HasPrefix(path, n.prefix) {
			path = path[len(n.prefix):]
		} else {
			break
		}
	}
	return ""
}

func getParams(k string) []string {
	// TODO: これ不要かも
	// /foo/:bar → []string{"bar"} みたいなparamsを抜き出す処理を書く
	return []string{}
}
