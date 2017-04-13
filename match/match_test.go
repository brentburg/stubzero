package match

import (
	"regexp"
	"testing"
)

var matchOne = Custom(func(v interface{}) bool {
	return v.(int) == 1
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
	re := regexp.MustCompile("true$")

	t.Run("with bytes", func(t *testing.T) {
		if !Regexp(re)([]byte("should be true")) {
			t.Error("matcher expected to return true for matching bytes")
		}

		if Regexp(re)([]byte("should be false")) {
			t.Error("matcher expected to return false for non-matching bytes")
		}
	})

	t.Run("with string", func(t *testing.T) {
		if !Regexp(re)("should be true") {
			t.Error("matcher expected to return true for matching string")
		}

		if Regexp(re)("should be false") {
			t.Error("matcher expected to return false for non-matching string")
		}
	})

	t.Run("with other types", func(t *testing.T) {
		if Regexp(re)(1) {
			t.Error("matcher expected to return false for other types")
		}
	})
}

func TestKeys(t *testing.T) {
	t.Run("with values", func(t *testing.T) {
		cases := []struct {
			m1 map[int]int
			m2 map[int]int
			r  bool
		}{
			{map[int]int{1: 2}, map[int]int{1: 2}, true},
			{map[int]int{1: 2}, map[int]int{1: 2, 3: 4}, true},
			{map[int]int{}, map[int]int{1: 2}, true},
			{map[int]int{1: 2, 3: 4}, map[int]int{1: 2}, false},
			{map[int]int{1: 2}, map[int]int{}, false},
		}

		for _, c := range cases {
			if Keys(c.m1)(c.m2) != c.r {
				t.Errorf("matcher expected to return %t for %v and %v", c.r, c.m1, c.m2)
			}
		}
	})

	t.Run("with sub-matchers", func(t *testing.T) {
		cases := []struct {
			m1 map[int]Matcher
			m2 map[int]int
			r  bool
		}{
			{map[int]Matcher{1: matchOne}, map[int]int{1: 1}, true},
			{map[int]Matcher{}, map[int]int{1: 2}, true},
			{map[int]Matcher{1: matchOne}, map[int]int{1: 2}, false},
			{map[int]Matcher{1: matchOne, 3: matchOne}, map[int]int{1: 1}, false},
			{map[int]Matcher{1: matchOne}, map[int]int{}, false},
		}

		for _, c := range cases {
			if Keys(c.m1)(c.m2) != c.r {
				t.Errorf("matcher expected to return %t for %v", c.r, c.m2)
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

func TestFields(t *testing.T) {
	t.Run("with values", func(t *testing.T) {
		cases := []struct {
			s1 interface{}
			s2 interface{}
			r  bool
		}{
			{struct{ A int }{1}, struct {
				A int
				B string
			}{1, "hello"}, true},
			{struct{ A int }{1}, struct{ A int }{1}, true},
			{struct{ B string }{"hi"}, struct {
				A int
				B string
			}{1, "hi"}, true},
			{struct{ A int }{1}, struct{ B string }{"hello"}, false},
			{struct{ A int }{1}, struct{ B int }{2}, false},
			{struct{ A int }{1}, struct{}{}, false},
			{struct {
				A int
				B string
			}{1, "hello"}, struct{ A int }{1}, false},
		}

		for _, c := range cases {
			if Fields(c.s1)(c.s2) != c.r {
				t.Errorf("matcher expected to be %t for %+v and %+v", c.r, c.s1, c.s2)
			}
		}
	})

	t.Run("with matchers", func(t *testing.T) {
		cases := []struct {
			s1 interface{}
			s2 interface{}
			r  bool
		}{
			{struct{ A Matcher }{matchOne}, struct {
				A int
				B string
			}{1, "hello"}, true},
			{struct{ A Matcher }{matchOne}, struct{ A int }{1}, true},
			{struct{ A int }{1}, struct{ B string }{"hello"}, false},
			{struct{ A int }{1}, struct{ A int }{2}, false},
			{struct{ A Matcher }{matchOne}, struct{}{}, false},
			{struct{ B Matcher }{}, struct{ A int }{1}, false},
		}

		for _, c := range cases {
			if Fields(c.s1)(c.s2) != c.r {
				t.Errorf("matcher expected to be %t for matchOne and %+v", c.r, c.s2)
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
