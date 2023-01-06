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

		if strings.Contains(path, ":") && strings.Contains(n.prefix, ":") {
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

// TODO: ここは後でsync.pool
var parameters = map[string]string{}

// Get gets a value from tree by key.
func (t *Tree) Get(k string) string {
	n := t.root
	path := k
	if path == "/" {
		n = n.getChild(path[0])
		if n == nil {
			return ""
		}
		if n.leaf != nil {
			return n.leaf.val
		}
		return ""
	}
	var tmppx string // for parammatch
	for {
		var tmpn *node
		if len(path) == 0 {
			if n.leaf != nil {
				// fmt.Printf("%#v\n", parameters)
				return n.leaf.val
			}
			break
		}
		if n.getChild(path[0]) != nil {
			n = n.getChild(path[0])
			if n.prefix == "/" {
				path = path[len(n.prefix):]
			} else {
				ppcp := longestPrefix(path, n.prefix)
				path = path[ppcp:]
			}

			continue
		}
		for i := 0; i < len(n.children); i++ {
			// prefix match
			ncp := n.children[i].node.prefix
			if strings.HasPrefix(path, ncp) {
				path = path[len(ncp):]
				tmpn = n.children[i].node
				tmppx = ncp
			}
			// param match
			if strings.Contains(ncp, ":") {
				cp := longestPrefix(path, ncp)

				tmppk := ncp[cp:]
				pki := strings.Index(tmppk, "/")
				pk := tmppk
				if pki > 0 {
					pk = tmppk[:pki]
				}

				tmppv := path[cp:]
				pvi := strings.Index(tmppv, "/")
				pv := tmppv
				if pvi > 0 {
					pv = tmppv[:pvi]
				}

				// TODO: /foo/:id/:id みたいにするとオーバーライドしてしまう
				parameters[pk] = pv // ex. key :one val one

				tmppx = ncp[:cp+len(pk)]

				tmpn = n.children[i].node
				if len(n.children[i].node.children) == 0 {
					path = path[cp+len(pv):]
					if path == "" {
						break
					}
					if n.children[i].node.leaf != nil {
						epk := explodePath(n.children[i].node.leaf.key[len(tmppx):])
						epp := explodePath(path)
						if len(epk) != len(epp) {
							return ""
						}
						for j := 0; j < len(epk); j++ {
							if strings.Contains(epk[j], ":") {
								// TODO: /foo/:id/:id みたいにするとオーバーライドしてしまう
								parameters[epk[j]] = epp[j]
							}
							if epk[j] == epp[j] {
								continue
							}
						}
					}
					path = ""
					break
				} else {
					path = path[cp+len(pv)+1:] // 1 is for /
				}
			}
		}

		n = tmpn

		if n == nil {
			break
		}

		// TODO: ここはforの冒頭だけで良い条件かも
		// if len(path) == 0 {
		// 	if n.leaf != nil {
		// 		// fmt.Printf("%#v\n", parameters)
		// 		return n.leaf.val
		// 	}
		// 	break
		// }
	}
	return ""
}

func explodePath(path string) []string {
	splitFn := func(c rune) bool {
		return c == '/'
	}
	return strings.FieldsFunc(path, splitFn)
}
