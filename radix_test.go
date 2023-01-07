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

func TestStatic(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       string
		getKey    string
		expVal    string
		expParams map[string]string
	}{
		{
			name:      "root",
			key:       "/",
			val:       "root",
			getKey:    "/",
			expVal:    "root",
			expParams: map[string]string{},
		},
		{
			name:      "static-1",
			key:       "/foo",
			val:       "static-1",
			getKey:    "/foo",
			expVal:    "static-1",
			expParams: map[string]string{},
		},
		{
			name:      "static-2",
			key:       "/foo/bar",
			val:       "static-2",
			getKey:    "/foo/bar",
			expVal:    "static-2",
			expParams: map[string]string{},
		},
		{
			name:      "static-3",
			key:       "/foo/bar/baz",
			val:       "static-3",
			getKey:    "/foo/bar/baz",
			expVal:    "static-3",
			expParams: map[string]string{},
		},
		{
			name:      "static-split-node-short",
			key:       "/fo/ba/ba",
			val:       "static-split-node-short",
			getKey:    "/fo/ba/ba",
			expVal:    "static-split-node-short",
			expParams: map[string]string{},
		},
		{
			name:      "static-split-node-long",
			key:       "/fooo/barr/bazz",
			val:       "static-split-node-long",
			getKey:    "/fooo/barr/bazz",
			expVal:    "static-split-node-long",
			expParams: map[string]string{},
		},
	}

	testAssertGet(t, cases)
}

func TestRootWithOneParam(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       string
		getKey    string
		expVal    string
		expParams map[string]string
	}{
		{
			name:      "root",
			key:       "/",
			val:       "root",
			getKey:    "/",
			expVal:    "root",
			expParams: map[string]string{},
		},
		{
			name:      "param-1",
			key:       "/foo/:bar",
			val:       "param-1",
			getKey:    "/foo/bar",
			expVal:    "param-1",
			expParams: map[string]string{":bar": "bar"},
		},
		// this conflics with /foo/:bar, can't define like this.
		// {
		// 	name:      "param-2",
		// 	key:       "/foo2/:bar2/:baz2",
		// 	val:       "param-2",
		// 	getKey:    "/foo/bar",
		// 	expVal:    "param-2",
		// 	expParams: map[string]string{},
		// },
	}

	testAssertGet(t, cases)
}

func TestRootWithTwoParam(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       string
		getKey    string
		expVal    string
		expParams map[string]string
	}{
		{
			name:      "root",
			key:       "/",
			val:       "root",
			getKey:    "/",
			expVal:    "root",
			expParams: map[string]string{},
		},
		{
			name:      "param-2",
			key:       "/foo/:bar/:baz",
			val:       "param-2",
			getKey:    "/foo/bar/baz",
			expVal:    "param-2",
			expParams: map[string]string{":bar": "bar", ":baz": "baz"},
		},
	}

	testAssertGet(t, cases)
}

func TestOnlyRoot(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       string
		getKey    string
		expVal    string
		expParams map[string]string
	}{
		{
			name:      "only-root",
			key:       "/",
			val:       "only-root",
			getKey:    "/",
			expVal:    "only-root",
			expParams: map[string]string{},
		},
	}

	testAssertGet(t, cases)
}

func TestWithoutRootOnlyStatic(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       string
		getKey    string
		expVal    string
		expParams map[string]string
	}{
		{
			name:      "static-1",
			key:       "/foo",
			val:       "static-1",
			getKey:    "/foo",
			expVal:    "static-1",
			expParams: map[string]string{},
		},
		{
			name:      "static-2",
			key:       "/foo/bar",
			val:       "static-2",
			getKey:    "/foo/bar",
			expVal:    "static-2",
			expParams: map[string]string{},
		},
	}

	testAssertGet(t, cases)
}

func TestWithoutRootOnlyOneParam(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       string
		getKey    string
		expVal    string
		expParams map[string]string
	}{
		{
			name:      "param-1",
			key:       "/foo/:foo",
			val:       "param-1",
			getKey:    "/foo/1",
			expVal:    "param-1",
			expParams: map[string]string{":foo": "1"},
		},
	}

	testAssertGet(t, cases)
}

func TestWithoutRootOnlyTwoParam(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       string
		getKey    string
		expVal    string
		expParams map[string]string
	}{
		{
			name:      "param-1",
			key:       "/foo/:foo",
			val:       "param-1",
			getKey:    "/foo/1",
			expVal:    "param-1",
			expParams: map[string]string{":foo": "1"},
		},
		// TODO: conflicts error
		{
			name:      "param-2",
			key:       "/f/:f",
			val:       "param-2",
			getKey:    "/f/1",
			expVal:    "param-2",
			expParams: map[string]string{":f": "1"},
		},
	}

	testAssertGet(t, cases)
}

func TestOnlyRootAndOneStatic(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       string
		getKey    string
		expVal    string
		expParams map[string]string
	}{
		{
			name:      "root",
			key:       "/",
			val:       "root",
			getKey:    "/",
			expVal:    "root",
			expParams: map[string]string{},
		},
		{
			name:      "/foo",
			key:       "/foo",
			val:       "foo",
			getKey:    "/foo",
			expVal:    "foo",
			expParams: map[string]string{},
		},
	}

	testAssertGet(t, cases)
}

func TestOnlyRootAndTwoStatic(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       string
		getKey    string
		expVal    string
		expParams map[string]string
	}{
		{
			name:      "root",
			key:       "/",
			val:       "root",
			getKey:    "/",
			expVal:    "root",
			expParams: map[string]string{},
		},
		{
			name:      "/foo",
			key:       "/foo",
			val:       "foo",
			getKey:    "/foo",
			expVal:    "foo",
			expParams: map[string]string{},
		},
		{
			name:      "/foo/bar",
			key:       "/foo/bar",
			val:       "foobar",
			getKey:    "/foo/bar",
			expVal:    "foobar",
			expParams: map[string]string{},
		},
	}

	testAssertGet(t, cases)
}

func TestParamIfNotMatchStatic(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       string
		getKey    string
		expVal    string
		expParams map[string]string
	}{
		{
			name:      "param-2",
			key:       "/foo/bar",
			val:       "static-2",
			getKey:    "/foo/bar",
			expVal:    "static-2", // not /foo/:bar but /foo/bar
			expParams: map[string]string{},
		},
		{
			name:      "param-1",
			key:       "/foo/:bar",
			val:       "param-bar",
			getKey:    "/foo/1",
			expVal:    "param-bar", // not /foo/bar　but /foo/:bar
			expParams: map[string]string{":bar": "1"},
		},
	}

	testAssertGet(t, cases)
}

func TestExample(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       string
		getKey    string
		expVal    string
		expParams map[string]string
	}{
		// {
		// 	name:      "healthcheck",
		// 	key:       "/healthcheck",
		// 	val:       "healthcheck",
		// 	getKey:    "/healthcheck",
		// 	expVal:    "healthcheck",
		// 	expParams: map[string]string{},
		// },
		// {
		// 	name:      "posts",
		// 	key:       "/posts",
		// 	val:       "posts",
		// 	getKey:    "/posts",
		// 	expVal:    "posts",
		// 	expParams: map[string]string{},
		// },
		// {
		// 	name:      "posts/categories/:name",
		// 	key:       "/posts/categories/:name",
		// 	val:       "posts/categories/:name",
		// 	getKey:    "/posts/categories/name",
		// 	expVal:    "posts/categories/:name",
		// 	expParams: map[string]string{":name": "name"},
		// },
		// {
		// 	name:      "posts/tags/:name",
		// 	key:       "/posts/tags/:name",
		// 	val:       "posts/tags/:name",
		// 	getKey:    "/posts/tags/name",
		// 	expVal:    "posts/tags/:name",
		// 	expParams: map[string]string{":name": "name"},
		// },
		// {
		// 	name:      "posts/:title",
		// 	key:       "/posts/:title",
		// 	val:       "posts/:title",
		// 	getKey:    "/posts/title",
		// 	expVal:    "posts/:title",
		// 	expParams: map[string]string{":title": "title"},
		// },
		// {
		// 	name:      "posts/:title/comments",
		// 	key:       "/posts/:title/comments",
		// 	val:       "posts/:title/comments",
		// 	getKey:    "/posts/title/comments",
		// 	expVal:    "posts/:title/comments",
		// 	expParams: map[string]string{":title": "title"},
		// },
		{
			name:      "categories/:name",
			key:       "/categories/:name",
			val:       "categories/:name",
			getKey:    "/categories/cate-name",
			expVal:    "categories/:name",
			expParams: map[string]string{":name": "cate-name"},
		},
		// {
		// 	name:      "tags",
		// 	key:       "/tags",
		// 	val:       "tags",
		// 	getKey:    "/tags",
		// 	expVal:    "tags",
		// 	expParams: map[string]string{},
		// },
		{
			name:      "tags/:name",
			key:       "/tags/:name",
			val:       "tags/:name",
			getKey:    "/tags/tag-name",
			expVal:    "tags/:name",
			expParams: map[string]string{":name": "tag-name"},
		},
		// TODO: categories/:nameに一致してしまう
		{
			name:      "signin",
			key:       "/signin",
			val:       "signin",
			getKey:    "/signin",
			expVal:    "signin",
			expParams: map[string]string{},
		},
		// TODO: add private api routes
	}

	testAssertGet(t, cases)
}

func TestGoblin(t *testing.T) {
	// TODO: goblinのテストケースも
}

func TestInsert(t *testing.T) {
	cases := []struct {
		name     string
		key      string
		val      string
		hasPanic bool
	}{
		// Those with haspanic true are not registered in the tree.
		{
			name:     "root",
			key:      "/",
			val:      "root",
			hasPanic: false,
		},
		{
			name:     "root",
			key:      "/",
			val:      "root", // duplicate path registration with /
			hasPanic: true,
		},
		{
			name:     "1",
			key:      "/foo",
			val:      "1",
			hasPanic: false,
		},
		{
			name:     "2",
			key:      "/foo",
			val:      "2",
			hasPanic: true, // dupulicate path registration with /foo
		},
		{
			name:     "3",
			key:      "/foo/bar",
			val:      "3",
			hasPanic: false,
		},
		{
			name:     "4",
			key:      "/foo/bar", // duplicate path registration with /foo/bar
			val:      "4",
			hasPanic: true,
		},
		{
			name:     "5",
			key:      "/foo/baz",
			val:      "5",
			hasPanic: false,
		},
		{
			name:     "6",
			key:      "/foo/baz", // duplicate path registration with /foo/baz
			val:      "6",
			hasPanic: true,
		},
		{
			name:     "7",
			key:      "/foo/bar/bar",
			val:      "7",
			hasPanic: false,
		},
		{
			name:     "8",
			key:      "/foo/bar/bar", // dupulicate path registration with /foo/bar/bar
			val:      "8",
			hasPanic: true,
		},
		{
			name:     "9",
			key:      "/foo/bar/baz",
			val:      "9",
			hasPanic: false,
		},
		{
			name:     "10",
			key:      "/foo/:bar",
			val:      "10",
			hasPanic: false,
		},
		{
			name:     "11",
			key:      "/foo/:bar",
			val:      "11",
			hasPanic: true, // dupulicate path registration with /foo/:bar
		},
		{
			name:     "12",
			key:      "/foo/:baz",
			val:      "12",
			hasPanic: true, // conflicts with /foo/:bar
		},
		{
			name:     "13",
			key:      "/foo/:bar/:bar",
			val:      "13",
			hasPanic: false,
		},
		{
			name:     "14",
			key:      "/foo/:bar/:bar",
			val:      "14",
			hasPanic: true, // duplicate path registration with /foo/:bar/:bar
		},
		{
			name:     "15",
			key:      "/foo/:bar/:baz",
			val:      "15",
			hasPanic: true, // conflics with /foo/:bar/:bar
		},
		{
			name:     "16",
			key:      "/foo/:bar/qux",
			val:      "16",
			hasPanic: false,
		},
		{
			name:     "17",
			key:      "/foo/:bar/qux", // duplicates with /foo/:bar/qux
			val:      "17",
			hasPanic: true,
		},
		{
			name:     "18",
			key:      "/foo/:bar/:qux", // duplicates with /foo/:bar/:bar
			val:      "18",
			hasPanic: true,
		},
		{
			name:     "19",
			key:      "/foo/:bar/quux",
			val:      "19",
			hasPanic: false,
		},
		{
			name:     "20",
			key:      "/foo/:bar/qux/quux",
			val:      "20",
			hasPanic: false,
		},
		{
			name:     "21",
			key:      "/foo/:bar/qux/quux", // duplicates with /foo/:bar/qux/quux
			val:      "21",
			hasPanic: true,
		},
		{
			name:     "22",
			key:      "/foo/:bar/qux/:quux",
			val:      "22",
			hasPanic: false,
		},
	}

	testAssertInsert(t, cases)
}

func TestInsertExample(t *testing.T) {
	// TODO: ここから
	cases := []struct {
		name     string
		key      string
		val      string
		hasPanic bool
	}{
		// Those with haspanic true are not registered in the tree.
		{
			name:     "root",
			key:      "/foo/:bar/:baz",
			val:      "root",
			hasPanic: false,
		},
		{
			name:     "root",
			key:      "/foo/:bar/:bam",
			val:      "root",
			hasPanic: false,
		},
	}

	testAssertInsert(t, cases)
}

func testAssertInsert(t *testing.T, cases []struct {
	name     string
	key      string
	val      string
	hasPanic bool
}) {
	tree := New()
	for _, c := range cases {
		// use anonymous functions to continue test cases.
		func() {
			defer func() {
				err := recover()
				if !c.hasPanic {
					if err != nil {
						t.Fatalf("%#v\n", err)
					}
				}
			}()
			tree.Insert(c.key, c.val)
			if c.hasPanic {
				t.Fatal("expected panic, but not")
			}
		}()
	}
}

func testAssertGet(t *testing.T, cases []struct {
	name      string
	key       string
	val       string
	getKey    string
	expVal    string
	expParams map[string]string
}) {
	tree := New()
	for _, c := range cases {
		// use anonymous functions to continue test cases.
		func() {
			defer func() {
				err := recover()
				if err != nil {
					t.Fatalf("expected no panic: %v\n", err)
				}
			}()
			tree.Insert(c.key, c.val)
		}()
	}
	for _, c := range cases {
		actVal := tree.Get(c.getKey)
		if c.expVal != actVal {
			t.Fatalf("expected: %v actual: %v", c.expVal, actVal)
		}
		if !reflect.DeepEqual(c.expParams, parameters) {
			t.Fatalf("expected: %v actual: %v", c.expParams, parameters)
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
