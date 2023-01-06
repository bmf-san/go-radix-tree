package radix

import (
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

func TestInsertAndGet(t *testing.T) {
	items := []struct {
		key string
		val string
	}{
		{
			key: "root",
			val: "1",
		},
		{
			key: "slow",
			val: "2",
		},
		{
			key: "slower",
			val: "3",
		},
		{
			key: "waste",
			val: "4",
		},
		{
			key: "water",
			val: "5",
		},
		{
			key: "watch",
			val: "6",
		},
		{
			key: "watcher",
			val: "7",
		},
	}
	tree := New()
	for _, i := range items {
		tree.Insert(i.key, i.val)
	}
	cases := []struct {
		key    string
		expVal string
	}{
		{
			key:    "root",
			expVal: "1",
		},
		{
			key:    "slow",
			expVal: "2",
		},
		{
			key:    "slower",
			expVal: "3",
		},
		{
			key:    "waste",
			expVal: "4",
		},
		{
			key:    "water",
			expVal: "5",
		},
		{
			key:    "watch",
			expVal: "6",
		},
		{
			key:    "watcher",
			expVal: "7",
		},
	}
	for _, c := range cases {
		exp := c.expVal
		act := tree.Get(c.key)
		if exp != act {
			t.Fatalf("expected: %v actual: %v", exp, act)
		}
	}
}

type insertItem struct {
	key string
	val string
}

func TestHTTPRouter(t *testing.T) {
	cases := []struct {
		name     string
		items    []insertItem
		hasPanic bool
		getKeys  []string
		expVals  []string
		// TODO: expValsを正常系として、異常系も追加
	}{
		// static test cases
		// {
		// 	name: "only-root",
		// 	items: []insertItem{
		// 		{
		// 			key: "/",
		// 			val: "only-root",
		// 		},
		// 	},
		// 	hasPanic: false,
		// 	getKeys:  []string{"/"},
		// 	expVals:  []string{"only-root"},
		// },
		// {
		// 	name: "static-1",
		// 	items: []insertItem{
		// 		{
		// 			key: "/foo",
		// 			val: "static-1",
		// 		},
		// 	},
		// 	hasPanic: false,
		// 	getKeys:  []string{"/foo"},
		// 	expVals:  []string{"static-1"},
		// },
		// {
		// 	name: "static-2",
		// 	items: []insertItem{
		// 		{
		// 			key: "/foo/bar",
		// 			val: "static-2",
		// 		},
		// 	},
		// 	hasPanic: false,
		// 	getKeys:  []string{"/foo/bar"},
		// 	expVals:  []string{"static-2"},
		// },
		// {
		// 	name: "static-3",
		// 	items: []insertItem{
		// 		{
		// 			key: "/foo/bar/baz",
		// 			val: "static-3",
		// 		},
		// 	},
		// 	hasPanic: false,
		// 	getKeys:  []string{"/foo/bar/baz"},
		// 	expVals:  []string{"static-3"},
		// },
		// {
		// 	name: "root-and-static-all",
		// 	items: []insertItem{
		// 		{
		// 			key: "/",
		// 			val: "root",
		// 		},
		// 		{
		// 			key: "/foo",
		// 			val: "static-1",
		// 		},
		// 		{
		// 			key: "/foo/bar",
		// 			val: "static-2",
		// 		},
		// 		{
		// 			key: "/foo/bar/baz",
		// 			val: "static-3",
		// 		},
		// 	},
		// 	hasPanic: false,
		// 	getKeys: []string{
		// 		"/",
		// 		"/foo",
		// 		"/foo/bar",
		// 		"/foo/bar/baz",
		// 	},
		// 	expVals: []string{
		// 		"root",
		// 		"static-1",
		// 		"static-2",
		// 		"static-3",
		// 	},
		// },
		// {
		// 	name: "root-and-static-split-node",
		// 	items: []insertItem{
		// 		{
		// 			key: "/",
		// 			val: "root",
		// 		},
		// 		{
		// 			key: "/foo",
		// 			val: "foo",
		// 		},
		// 		{
		// 			key: "/fo",
		// 			val: "fo",
		// 		},
		// 		{
		// 			key: "/foz",
		// 			val: "foz",
		// 		},
		// 		{
		// 			key: "/fooo",
		// 			val: "fooo",
		// 		},
		// 		{
		// 			key: "/foo/bar",
		// 			val: "foobar",
		// 		},
		// 		{
		// 			key: "/foo/ba",
		// 			val: "fooba",
		// 		},
		// 		{
		// 			key: "/foo/baz",
		// 			val: "foobaz",
		// 		},
		// 		{
		// 			key: "/foo/barr",
		// 			val: "foobarr",
		// 		},
		// 	},
		// 	hasPanic: false,
		// 	getKeys: []string{
		// 		"/",
		// 		"/foo",
		// 		"/fo",
		// 		"/foz",
		// 		"/fooo",
		// 		"/foo/bar",
		// 		"/foo/ba",
		// 		"/foo/baz",
		// 		"/foo/barr",
		// 	},
		// 	expVals: []string{
		// 		"root",
		// 		"foo",
		// 		"fo",
		// 		"foz",
		// 		"fooo",
		// 		"foobar",
		// 		"fooba",
		// 		"foobaz",
		// 		"foobarr",
		// 	},
		// },
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
		{
			name: "param-2",
			items: []insertItem{
				{
					key: "/foo/:one/:two",
					val: "param-2",
				},
			},
			hasPanic: false,
			getKeys:  []string{"/foo/one/two"},
			expVals:  []string{"param-2"},
		},
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
		// {
		// 	name: "param-2-middle",
		// 	items: []insertItem{
		// 		{
		// 			key: "/foo/:one/bar/:two",
		// 			val: "param-2-middle",
		// 		},
		// 	},
		// 	hasPanic: false,
		// 	getKeys:  []string{"/foo/one/bar/two"},
		// 	expVals:  []string{"param-2-middle"},
		// },
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

// TODO: 不要、後で消す
func TestInsertAndGetForHTTPRouter(t *testing.T) {
	items := []struct {
		key string
		val string
	}{
		{
			key: "/",
			val: "1",
		},
		{
			key: "/foo",
			val: "2",
		},
		{
			key: "/foo/bar",
			val: "3",
		},
		{
			key: "/foo/bar/baz",
			val: "4",
		},
		{
			key: "/foo/:bar",
			val: "5",
		},
		// {
		// 	// TODO: routerの仕様としてこれをどう扱うか？ 先勝ちにするか、そもそも登録できないしようにするか。TestInsertAndGetForHTTPRouter
		// 	// 登録できないよう仕様にするほうが安全な気がする
		// 	key: "/foo/:ba",
		// 	val: "9",
		// },
		// TODO: path paramが対応できてから対応する
		// {
		// 	// TODO: 正規表現がある場合は正規表現を考慮する必要があるj
		// 	key: `/foo/:bar[^\D+$]`,
		// 	val: "5-reg",
		// },
		{
			key: "/foo/bar/:baz",
			val: "6",
		},
		{
			key: "/foo/:bar/baz",
			val: "7",
		},
		// {
		// 	// NOTE: 上書きしないでpanicにする。これは無駄な登録なので、そのほうが安全そう
		// 	// NOTE: →上書きできないようにした
		// 	key: "/foo/:bar/baz",
		// 	val: "7-x",
		// },
		{
			key: "/bar",
			val: "8",
		},
		{
			key: "/baz/caz",
			val: "10",
		},
		{
			key: "/baz/cat",
			val: "11",
		},
		{
			key: "/a/b/c/:d/:e/:f",
			val: "12",
		},
		{
			key: "/a/b/c/:d/:e",
			val: "13",
		},
	}
	tree := New()
	for _, i := range items {
		tree.Insert(i.key, i.val)
	}
	cases := []struct {
		key string
		val string
	}{
		// {
		// 	key: "/",
		// 	val: "1",
		// },
		// {
		// 	key: "/foo",
		// 	val: "2",
		// },
		// {
		// 	key: "/foo/bar",
		// 	val: "3",
		// },
		// {
		// 	key: "/foo/bar/baz",
		// 	val: "4",
		// },
		{
			key: "/foo/path-param-bar",
			val: "5",
		},
		// {
		// 	// TODO: routerの仕様としてこれをどう扱うか？ 先勝ちにするか、そもそも登録できないしようにするか。TestInsertAndGetForHTTPRouter
		// 	// 登録できないよう仕様にするほうが安全な気がする
		// 	key: "/foo/:ba",
		// 	val: "9",
		// },
		// TODO: path paramが対応できてから対応する
		// {
		// 	// TODO: 正規表現がある場合は正規表現を考慮する必要があるj
		// 	key: `/foo/:bar[^\D+$]`,
		// 	val: "5-reg",
		// },
		// {
		// 	key: "/foo/bar/path-param-baz",
		// 	val: "6",
		// },
		// {
		// 	key: "/foo/:bar/baz",
		// 	val: "7",
		// },
		// {
		// 	// NOTE: 上書きしないでpanicにする。これは無駄な登録なので、そのほうが安全そう
		// 	// NOTE: →上書きできないようにした
		// 	key: "/foo/:bar/baz",
		// 	val: "7-x",
		// },
		// {
		// 	key: "/bar",
		// 	val: "8",
		// },
		// {
		// 	key: "/baz/caz",
		// 	val: "10",
		// },
		// {
		// 	key: "/baz/cat",
		// 	val: "11",
		// },
		// {
		// 	key: "/a/b/c/path-param-d/path-param-e/path-param-f",
		// 	val: "12",
		// },
		// {
		// 	key: "/a/b/c/path-param-d/path-param-e",
		// 	val: "13",
		// },
	}
	for _, c := range cases {
		exp := c.val
		act := tree.Get(c.key)
		if exp != act {
			t.Fatalf("expected: %v actual: %v", exp, act)
		}
	}
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
