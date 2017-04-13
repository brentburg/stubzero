package match

import (
	"fmt"
	"regexp"
	"testing"
)

var matchOne = Custom(func(v interface{}) bool {
	i, ok := v.(int)
	if !ok {
		return false
	}
	return i == 1
})

func TestMatch(t *testing.T) {
	t.Run("with matcher and value", func(t *testing.T) {
		if !Match(matchOne, 1) {
			t.Error("expected to return true if matcher returns true for value")
		}
		if Match(matchOne, 2) {
			t.Error("expected to return false if matcher returns false for value")
		}
	})

	t.Run("with deeply equal values", func(t *testing.T) {
		cases := [][2]interface{}{
			{struct{ a int }{1}, struct{ a int }{1}},
			{[]int{1, 2, 3}, []int{1, 2, 3}},
			{map[int]int{1: 2}, map[int]int{1: 2}},
			{"hello", "hello"},
			{1, 1},
			{nil, nil},
			{true, true},
			{false, false},
		}

		for _, c := range cases {
			if !Match(c[0], c[1]) {
				t.Errorf("expected to be true for %+v and %+v", c[0], c[1])
			}
		}
	})

	t.Run("with values not deeply equal", func(t *testing.T) {
		cases := [][2]interface{}{
			{struct{ a int }{1}, struct{ a int }{2}},
			{struct{ a int }{1}, struct{ v int }{1}},
			{[]int{1, 2, 3}, []int{1, 2, 4}},
			{map[int]int{1: 2}, map[int]int{1: 3}},
			{"hello", "hola"},
			{1, 2},
			{true, false},
		}

		for _, c := range cases {
			if Match(c[0], c[1]) {
				t.Errorf("expected to be false for %+v and %+v", c[0], c[1])
			}
		}
	})
}

func TestAny(t *testing.T) {
	cases := []interface{}{
		struct{ Name string }{"name"},
		[2]int{1, 2},
		"string",
		1,
		true,
		false,
		nil,
	}
	for _, c := range cases {
		if !Any(c) {
			t.Errorf("matcher expected to not return false for %+v", c)
		}
	}
}

func TestRegexp(t *testing.T) {
	exps := []interface{}{"true$", regexp.MustCompile("true$")}

	for _, exp := range exps {
		t.Run(fmt.Sprintf("%T with []byte", exp), func(t *testing.T) {
			if !Regexp(exp)([]byte("should be true")) {
				t.Error("matcher expected to return true for matching bytes")
			}

			if Regexp(exp)([]byte("should be false")) {
				t.Error("matcher expected to return false for non-matching bytes")
			}
		})

		t.Run(fmt.Sprintf("%T with string", exp), func(t *testing.T) {
			if !Regexp(exp)("should be true") {
				t.Error("matcher expected to return true for matching string")
			}

			if Regexp(exp)("should be false") {
				t.Error("matcher expected to return false for non-matching string")
			}
		})
	}
	t.Run("with other types", func(t *testing.T) {
		if Regexp("1")(1) {
			t.Error("matcher expected to return false for other types")
		}
	})
}

func TestKey(t *testing.T) {
	t.Run("with values", func(t *testing.T) {
		cases := []struct {
			k int
			v int
			m map[int]int
			r bool
		}{
			{1, 2, map[int]int{1: 2}, true},
			{1, 2, map[int]int{1: 2, 3: 4}, true},
			{1, 2, map[int]int{}, false},
			{1, 2, map[int]int{1: 3}, false},
			{1, 2, map[int]int{2: 2}, false},
		}

		for _, c := range cases {
			if Key(c.k, c.v)(c.m) != c.r {
				t.Errorf(
					"matcher expected to be %t for %v: %v and %v",
					c.r, c.k, c.v, c.m,
				)
			}
		}
	})

	t.Run("with sub-matchers", func(t *testing.T) {
		cases := []struct {
			k int
			v Matcher
			m map[int]int
			r bool
		}{
			{1, matchOne, map[int]int{1: 1}, true},
			{1, matchOne, map[int]int{1: 1, 2: 3}, true},
			{1, matchOne, map[int]int{1: 2}, false},
			{1, matchOne, map[int]int{2: 1}, false},
			{1, matchOne, map[int]int{}, false},
		}

		for _, c := range cases {
			if Key(c.k, c.v)(c.m) != c.r {
				t.Errorf("matcher expected to be %t for %v of %v", c.r, c.k, c.m)
			}
		}
	})
}

func TestContains(t *testing.T) {
	t.Run("with a value", func(t *testing.T) {
		if !Contains("hello")([]string{"hello", "goodbye"}) {
			t.Error("matcher expected to match slice containing value")
		}

		if Contains("hello")([]string{"goodbye"}) {
			t.Error("matcher expected to not match slice without value")
		}
	})

	t.Run("with a matcher", func(t *testing.T) {
		if !Contains(matchOne)([]int{1, 2}) {
			t.Error("matcher expected to match slice containing matching value")
		}

		if Contains(matchOne)([]int{2}) {
			t.Error("matcher expected to not match slice without matching value")
		}
	})
}

func TestField(t *testing.T) {
	t.Run("with values", func(t *testing.T) {
		cases := []struct {
			n string
			v interface{}
			s interface{}
			r bool
		}{
			{"A", 1, struct {
				A int
				B string
			}{1, "hello"}, true},
			{"A", 1, struct{ A int }{1}, true},
			{"B", "hello", struct {
				A int
				B string
			}{1, "hello"}, true},
			{"A", 1, struct{ B string }{"hello"}, false},
			{"A", 1, struct{ A int }{2}, false},
			{"A", 1, struct{ A string }{"hello"}, false},
			{"A", 1, struct{}{}, false},
		}

		for _, c := range cases {
			if Field(c.n, c.v)(c.s) != c.r {
				t.Errorf(
					"matcher expected to be %t for %s: %v and %+v",
					c.r, c.n, c.v, c.s,
				)
			}
		}
	})

	t.Run("with matchers", func(t *testing.T) {
		cases := []struct {
			n string
			v Matcher
			s interface{}
			r bool
		}{
			{"A", matchOne, struct {
				A int
				B string
			}{1, "hello"}, true},
			{"A", matchOne, struct{ A int }{1}, true},
			{"A", matchOne, struct{ B string }{"hello"}, false},
			{"A", matchOne, struct{ A string }{"hello"}, false},
			{"A", matchOne, struct{ A int }{2}, false},
			{"A", matchOne, struct{}{}, false},
		}

		for _, c := range cases {
			if Field(c.n, c.v)(c.s) != c.r {
				t.Errorf(
					"matcher expected to be %t for %s: matchOne and %+v",
					c.r, c.n, c.s,
				)
			}
		}
	})
}

func TestCustom(t *testing.T) {
	cm := Custom(func(v interface{}) bool {
		return v.(string) == "custom"
	})
	if !Match(cm, "custom") {
		t.Error("matcher expected to match value 'custom'")
	}
	if Match(cm, "other") {
		t.Error("matcher expected to not match value 'other'")
	}
}
