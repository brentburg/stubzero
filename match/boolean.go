package match

func And(matchers ...Matcher) Matcher {
	return func(v interface{}) bool {
		for _, matcher := range matchers {
			if !matcher(v) {
				return false
			}
		}
		return true
	}
}

func Or(matchers ...Matcher) Matcher {
	return func(v interface{}) bool {
		for _, matcher := range matchers {
			if matcher(v) {
				return true
			}
		}
		return false
	}
}

func Xor(m1, m2 Matcher) Matcher {
	return func(v interface{}) bool {
		return m1(v) != m2(v)
	}
}
