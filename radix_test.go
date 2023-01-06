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

	tree := New()
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("expected no panic: %v\n", err)
			}
		}()
		tree.Insert(c.key, c.val)
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

	tree := New()
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("expected no panic: %v\n", err)
			}
		}()
		tree.Insert(c.key, c.val)
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

	tree := New()
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("expected no panic: %v\n", err)
			}
		}()
		tree.Insert(c.key, c.val)
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

	tree := New()
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("expected no panic: %v\n", err)
			}
		}()
		tree.Insert(c.key, c.val)
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

	tree := New()
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("expected no panic: %v\n", err)
			}
		}()
		tree.Insert(c.key, c.val)
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
		// {
		// 	name:      "param-2",
		// 	key:       "/f/:f",
		// 	val:       "param-2",
		// 	getKey:    "/f/1",
		// 	expVal:    "param-2",
		// 	expParams: map[string]string{":f": "1"},
		// },
	}

	tree := New()
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("expected no panic: %v\n", err)
			}
		}()
		tree.Insert(c.key, c.val)
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
		{
			name:      "param-2",
			key:       "/f/:f",
			val:       "param-2",
			getKey:    "/f/1",
			expVal:    "param-2",
			expParams: map[string]string{":f": "1"},
		},
	}

	tree := New()
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("expected no panic: %v\n", err)
			}
		}()
		tree.Insert(c.key, c.val)
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

	tree := New()
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("expected no panic: %v\n", err)
			}
		}()
		tree.Insert(c.key, c.val)
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

	tree := New()
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("expected no panic: %v\n", err)
			}
		}()
		tree.Insert(c.key, c.val)
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

	tree := New()
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("expected no panic: %v\n", err)
			}
		}()
		tree.Insert(c.key, c.val)
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
