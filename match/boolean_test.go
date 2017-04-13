package match

import (
	"testing"
)

var trueMatcher = Custom(func(_ interface{}) bool {
	return true
})

var falseMatcher = Custom(func(_ interface{}) bool {
	return false
})

func TestAnd(t *testing.T) {
	if !And(trueMatcher, trueMatcher, trueMatcher)(struct{}{}) {
		t.Error("matcher expected to return true if all matchers are true")
	}
	if And(trueMatcher, trueMatcher, falseMatcher)(struct{}{}) {
		t.Error("matcher expected to return false if any matchers are false")
	}
	if And(falseMatcher, falseMatcher, falseMatcher)(struct{}{}) {
		t.Error("matcher expected to return false if all matchers are false")
	}
}

func TestOr(t *testing.T) {
	if !Or(trueMatcher, trueMatcher, trueMatcher)(struct{}{}) {
		t.Error("matcher expected to return true if all matchers are true")
	}
	if !Or(falseMatcher, trueMatcher, falseMatcher)(struct{}{}) {
		t.Error("matcher expected to return true if any matchers are true")
	}
	if Or(falseMatcher, falseMatcher, falseMatcher)(struct{}{}) {
		t.Error("Or matcher expected to return false if no matchers are true")
	}
}

func TestXor(t *testing.T) {
	if !Xor(falseMatcher, trueMatcher)(struct{}{}) {
		t.Error("matcher expected to return true if only one matcher is true")
	}
	if Xor(falseMatcher, falseMatcher)(struct{}{}) {
		t.Error("matcher expected to return false if both matchers are false")
	}
	if Xor(trueMatcher, trueMatcher)(struct{}{}) {
		t.Error("matcher expected to return false if both matchers are true")
	}
}
