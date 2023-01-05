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

func TestUnderstandGetChild(t *testing.T) {
	tree := New()
	// tree.Insert("/foo/bar", "foobar")
	tree.Insert("/foo/:bar/baz", "param-bar")
	// tree.Insert("/foo/:bar/:baz", "param-bar")
	// tree.Insert("/foo/:bar", "param-bar")
	// tree.Insert("/foo/baz", "foobaz")
	// tree.Insert("/foo/:ba", "param-bar")
	// tree.Insert("/foo/bazz", "foobaz")
	pb := tree.Get("/foo/param-bar")
	// b := tree.Get("/foo/bar")
	fmt.Printf("%#v\n", pb)
}

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
