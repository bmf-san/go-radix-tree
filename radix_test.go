package radix

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddChild(t *testing.T) {
	tree := New()
	var b byte
	c := child{
		label: b,
		node:  &node{},
	}
	tree.root.addChild(c)

	exp := c
	act := tree.root.children[0]
	if !reflect.DeepEqual(exp, act) {
		t.Fatalf("expected: %v actual: %v", exp, act)
	}
}

func TestUpdateChild(t *testing.T) {
	tree := New()
	var b byte
	c := child{
		label: b,
		node: &node{
			prefix: "before",
		},
	}
	tree.root.addChild(c)
	var b2 byte
	p := "after"
	tree.root.updateChild(b2, &node{
		prefix: p,
	})

	exp := p
	act := tree.root.children[0].node.prefix
	if !reflect.DeepEqual(exp, act) {
		t.Fatalf("expected: %v actual: %v", exp, act)
	}
}

func TestGetChild(t *testing.T) {
	tree := New()
	var b byte
	c := child{
		label: b,
		node:  &node{},
	}
	tree.root.addChild(c)

	exp := tree.root.children[0].node
	act := tree.root.getChild(b)
	if !reflect.DeepEqual(exp, act) {
		t.Fatalf("expected: %v actual: %v", exp, act)
	}
}

func TestNew(t *testing.T) {
	exp := New()
	act := &Tree{root: &node{}}
	if !reflect.DeepEqual(exp, act) {
		t.Fatalf("expected: %v actual: %v", exp, act)
	}
}

type insertItem struct {
	key string
	val string
}

func TestStaticAndParams(t *testing.T) {
	cases := []struct {
		name   string
		key    string
		val    string
		getKey string
		expVal string
		// TODO: expValsを正常系として、異常系も追加
	}{
		// static routes
		// {
		// 	name:   "root",
		// 	key:    "/",
		// 	val:    "root",
		// 	getKey: "/",
		// 	expVal: "root",
		// },
		// {
		// 	name:   "static-1",
		// 	key:    "/foo",
		// 	val:    "static-1",
		// 	getKey: "/foo",
		// 	expVal: "static-1",
		// },
		{
			name:   "static-2",
			key:    "/foo/bar",
			val:    "static-2",
			getKey: "/foo/bar",
			expVal: "static-2",
		},
		// {
		// 	name:   "static-3",
		// 	key:    "/foo/bar/baz",
		// 	val:    "static-3",
		// 	getKey: "/foo/bar/baz",
		// 	expVal: "static-3",
		// },
		// {
		// 	name:   "static-split-node-short",
		// 	key:    "/fo/ba/ba",
		// 	val:    "static-split-node-short",
		// 	getKey: "/fo/ba/ba",
		// 	expVal: "static-split-node-short",
		// },
		// {
		// 	name:   "static-split-node-long",
		// 	key:    "/fooo/barr/bazz",
		// 	val:    "static-split-node-long",
		// 	getKey: "/fooo/barr/bazz",
		// 	expVal: "static-split-node-long",
		// },
		// param routes
		{
			name:   "param-1",
			key:    "/foo/:bar",
			val:    "param-1",
			getKey: "/foo/1",
			expVal: "param-1",
		},
	}

	tree := New()
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				// if !c.hasPanic {
				t.Errorf("expected no panic: %v\n", err)
				// }
			}
		}()
		tree.Insert(c.key, c.val)
	}
	for _, c := range cases {
		actVal := tree.Get(c.getKey)
		fmt.Printf("%#v\n", parameters)
		if c.expVal != actVal {
			t.Fatalf("expected: %v actual: %v", c.expVal, actVal)
		}
	}
}

func TestPriority(t *testing.T) {
	cases := []struct {
		name   string
		key    string
		val    string
		getKey string
		expVal string
		// TODO: expValsを正常系として、異常系も追加
	}{
		// TODO:　優先度の仕様を分かるように
		{
			name:   "static-2",
			key:    "/foo/bar",
			val:    "static-2",
			getKey: "/foo/bar",
			expVal: "static-2",
		},
		{
			name:   "param-1",
			key:    "/foo/:bar",
			val:    "param-1",
			getKey: "/foo/bar",
			expVal: "static-2",
		},
	}

	tree := New()
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				// if !c.hasPanic {
				t.Errorf("expected no panic: %v\n", err)
				// }
			}
		}()
		tree.Insert(c.key, c.val)
	}
	for _, c := range cases {
		actVal := tree.Get(c.getKey)
		fmt.Printf("%#v\n", parameters)
		if c.expVal != actVal {
			t.Fatalf("expected: %v actual: %v", c.expVal, actVal)
		}
	}
}

func TestHTTPRouters(t *testing.T) {
	cases := []struct {
		name     string
		items    []insertItem
		hasPanic bool
		getKeys  []string
		expVals  []string
		// TODO: expValsを正常系として、異常系も追加
	}{
		// param test cases
		// {
		// 	name: "param-1",
		// 	items: []insertItem{
		// 		{
		// 			key: "/foo/:one",
		// 			val: "param-1",
		// 		},
		// 	},
		// 	hasPanic: false,
		// 	getKeys:  []string{"/foo/one"},
		// 	expVals:  []string{"param-1"},
		// },
		// {
		// 	name: "param-2",
		// 	items: []insertItem{
		// 		{
		// 			key: "/foo/:one/:two",
		// 			val: "param-2",
		// 		},
		// 	},
		// 	hasPanic: false,
		// 	getKeys:  []string{"/foo/one/two"},
		// 	expVals:  []string{"param-2"},
		// },

		// TODO:
		// /foo/:one/two/:three

		// TODO:
		// /foo/bar/baz
		// /foo/:bar/baz

		// {
		// 	name: "param-3",
		// 	items: []insertItem{
		// 		{
		// 			key: "/foo/:one/:two/:three",
		// 			val: "param-3",
		// 		},
		// 	},
		// 	hasPanic: false,
		// 	getKeys:  []string{"/foo/one/two/three"},
		// 	expVals:  []string{"param-3"},
		// },
		// {
		// 	name: "param-1-middle",
		// 	items: []insertItem{
		// 		{
		// 			key: "/foo/:one/bar",
		// 			val: "param-1-middle",
		// 		},
		// 	},
		// 	hasPanic: false,
		// 	getKeys:  []string{"/foo/one/bar"},
		// 	expVals:  []string{"param-1-middle"},
		// },
		{
			name: "param-2-middle",
			items: []insertItem{
				{
					key: "/foo/:one/bar/:two",
					val: "param-2-middle",
				},
			},
			hasPanic: false,
			getKeys:  []string{"/foo/one/bar/two"},
			expVals:  []string{"param-2-middle"},
		},
		// TODO: param all case
		// priority test cases
		// panic test cases
		// general test cases
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				err := recover()
				if err != nil {
					if !c.hasPanic {
						t.Errorf("expected no panic: %v\n", err)
					}
				}
			}()
			tree := New()
			for _, i := range c.items {
				tree.Insert(i.key, i.val)
			}
			var actVals []string
			for _, k := range c.getKeys {
				actVals = append(actVals, tree.Get(k))
			}
			if !reflect.DeepEqual(c.expVals, actVals) {
				t.Fatalf("expected: %v actual: %v", c.expVals, actVals)
			}
		})
	}
}

func TestExample(t *testing.T) {
	// TODO: bmf-techのrouting
}

func TestLongestPrefix(t *testing.T) {
	cases := []struct {
		k1  string
		k2  string
		exp int
	}{
		{
			k1:  "f",
			k2:  "foo",
			exp: 1,
		},
		{
			k1:  "foo",
			k2:  "foo",
			exp: 3,
		},
		{
			k1:  "foo",
			k2:  "bar",
			exp: 0,
		},
		{
			k1:  "foo",
			k2:  "fooo",
			exp: 3,
		},
	}
	for _, c := range cases {
		exp := c.exp
		act := longestPrefix(c.k1, c.k2)
		if exp != act {
			t.Fatalf("expected: %v actual: %v", exp, act)
		}
	}
}

func BenchmarkInsert(b *testing.B) {
	r := New()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.Insert("key", "value")
	}
}
