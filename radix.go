// It refers to the implementation of go-radix
// See: https://github.com/armon/go-radix
package radix

import (
	"sort"
	"strings"
)

// leafNode is used to represent a value
type leafNode struct {
	key string
	val string
}

// child is used to represent an child node
type child struct {
	label byte
	node  *node
}

type node struct {
	// leaf is used to store possible leaf
	leaf *leafNode

	// prefix is the common prefix we ignore
	prefix string

	// children should be stored in-order for iteration.
	// We avoid a fully materialized slice to save memory,
	// since in most cases we expect to be sparse
	children children
}

func (n *node) addChild(e child) {
	num := len(n.children)
	idx := sort.Search(num, func(i int) bool {
		return n.children[i].label >= e.label
	})

	n.children = append(n.children, child{})
	copy(n.children[idx+1:], n.children[idx:])
	n.children[idx] = e
}

func (n *node) updateChild(label byte, node *node) {
	// TODO: Read later
	num := len(n.children)
	// NOTE: Refactoring. this codes also used in getChild.
	idx := sort.Search(num, func(i int) bool {
		return n.children[i].label >= label
	})
	if idx < num && n.children[idx].label == label {
		n.children[idx].node = node
		return
	}
	// TODO: change messages or return error.
	panic("replacing missing child")
}

func (n *node) getChild(label byte) *node {
	// TODO: Read later
	num := len(n.children)
	// NOTE: Refactoring. this codes also used in updateChild.
	idx := sort.Search(num, func(i int) bool {
		return n.children[i].label >= label
	})
	if idx < num && n.children[idx].label == label {
		return n.children[idx].node
	}
	return nil
}

type children []child

func (e children) Len() int {
	return len(e)
}

func (e children) Less(i, j int) bool {
	return e[i].label < e[j].label
}

func (e children) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e children) Sort() {
	sort.Sort(e)
}

// Tree implements a radix tree. This can be treated as a
// Dictionary abstract data type. The main advantage over
// a standard hash map is prefix-based lookups and
// ordered iteration,
type Tree struct {
	root *node
}

// New returns an empty Tree
func New() *Tree {
	return &Tree{root: &node{}}
}

// NOTE: This could probably be written differently. Try writing some test code.
// longestPrefix finds the length of the shared prefix
// of two strings
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

// Insert is used to add a newentry or update
// an existing entry. Returns true if an existing record is updated.
func (t *Tree) Insert(s, v string) {
	var parent *node
	n := t.root
	search := s
	for {
		// NOTE: If the string to be inserted is empty, it can be an error if it is a router (return an error instead of cleanpath), so there is no need for a conditional branch here
		// Handle key exhaution
		if len(search) == 0 {
			if n.leaf != nil {
				n.leaf.val = v
				return
			}

			n.leaf = &leafNode{
				key: s,
				val: v,
			}
			return
		}

		// Look for the child
		parent = n
		n = n.getChild(search[0])

		// No child, create one
		if n == nil {
			parent.addChild(child{
				label: search[0],
				node: &node{
					leaf: &leafNode{
						key: s,
						val: v,
					},
					prefix: search,
				},
			})
			return
		}

		// Determine longest prefix of the search key on match
		commonPrefix := longestPrefix(search, n.prefix)
		if commonPrefix == len(n.prefix) {
			search = search[commonPrefix:]
			continue
		}

		// Split the node
		spln := &node{
			prefix: search[:commonPrefix],
		}
		parent.updateChild(search[0], spln)

		// Restore the existing node
		spln.addChild(child{
			label: n.prefix[commonPrefix],
			node:  n,
		})
		n.prefix = n.prefix[commonPrefix:]

		// Create a new leaf node
		leaf := &leafNode{
			key: s,
			val: v,
		}

		// If the new key is a subset, add to this node
		search = search[commonPrefix:]
		if len(search) == 0 {
			spln.leaf = leaf
			return
		}

		// Create a new child for the node
		spln.addChild(child{
			label: search[0],
			node: &node{
				leaf:   leaf,
				prefix: search,
			},
		})
		return
	}
}

// Get is used to lookup a specific key, returning
// the value and if it was found
func (t *Tree) Get(s string) string {
	n := t.root
	search := s
	for {
		// Check for key exhaution
		if len(search) == 0 {
			if n.leaf != nil {
				return n.leaf.val
			}
			break
		}

		// Look for an child
		n = n.getChild(search[0])
		if n == nil {
			break
		}

		// Consume the search prefix
		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			break
		}
	}
	return ""
}
