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
		{
			// TODO: 正規表現がある場合は正規表現を考慮する必要がある
			key: `/foo/:bar[^\D+$]`,
			val: "5-reg",
		},
		{
			// TODO: routerの仕様としてこれをどう扱うか？ 先勝ちにするか、そもそも登録できないしようにするか。TestInsertAndGetForHTTPRouter
			// 登録できないよう仕様にするほうが安全な気がする
			key: "/foo/:ba",
			val: "9",
		},
		{
			key: "/foo/bar/:baz",
			val: "6",
		},
		{
			key: "/foo/:bar/baz",
			val: "7",
		},
		{
			// NOTE: 上書きされるだけ?
			key: "/foo/:bar/baz",
			val: "7-x",
		},
		{
			key: "/bar",
			val: "8",
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
			key:    "/",
			expVal: "1",
		},
		{
			key:    "/foo",
			expVal: "2",
		},
		{
			key:    "/foo/bar",
			expVal: "3",
		},
		{
			key:    "/foo/bar/baz",
			expVal: "4",
		},
		{
			// /foo/1
			key:    "/foo/:bar",
			expVal: "5",
		},
		{
			// /foo/one
			key:    `/foo/:bar[^\D+$]`,
			expVal: "5-reg",
		},
		{
			// /foo/1
			key:    "/foo/:ba",
			expVal: "9",
		},
		{
			// /foo/bar/1
			key:    "/foo/bar/:baz",
			expVal: "6",
		},
		{
			// /foo/1/baz
			key:    "/foo/:bar/baz",
			expVal: "7",
		},
		{
			// /foo/1/baz
			key:    "/foo/:bar/baz",
			expVal: "7-x",
		},
		{
			// /bar
			key:    "/bar",
			expVal: "8",
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
