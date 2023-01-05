// It refers to the implementation of go-radix
// See: https://github.com/armon/go-radix
package radix

import (
	"fmt"
	"sort"
	"strings"
)

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
	key string
	val string
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
		if path == "" {
			panic(fmt.Sprintf("duplicate path registration. path: %v prefix: %v", k, n.prefix))
		}
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
			path = path[commonPrefix:]
			continue
		}

		if strings.Count(path, "/") == strings.Count(n.prefix, "/") {
			if strings.Contains(path, ":") || strings.Contains(n.prefix, ":") {
				panic(fmt.Sprintf("conflicts path parameter. path: %v prefix: %v", path, n.prefix))
			}
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

// TODO: ここは後でsync.pool
var parameters = map[string]string{}

// Get gets a value from tree by key.
func (t *Tree) Get(k string) string {
	n := t.root
	path := k
	for {
		if len(path) == 0 {
			if n.leaf != nil {
				fmt.Printf("%#v\n", parameters)
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
			commonPrefix := longestPrefix(path, n.prefix)
			param := n.prefix[commonPrefix:]
			isParamPrefix := strings.Contains(param, ":")
			if isParamPrefix {
				paramPath := path[commonPrefix:]
				parameters[param] = paramPath
				path = path[commonPrefix+len(paramPath):]
			} else {
				break
			}
		}
	}
	return ""
}
