package stubzero

import (
	"testing"

	"github.com/brentburg/stubzero/match"
)

func TestCallCalledWith(t *testing.T) {
	t.Run("with equal values", func(t *testing.T) {
		c := newCall(1, 2, 3)
		if !c.CalledWith(1, 2) {
			t.Error("expected true when calling with equal values")
		}
		if c.CalledWith(1, 3) {
			t.Error("expected false when calling with non-equal values")
		}
	})

	t.Run("with deeply equal values", func(t *testing.T) {
		c := newCall([]int{1, 2}, 3)
		if !c.CalledWith([]int{1, 2}) {
			t.Error("expected true when calling with equal values")
		}
		if c.CalledWith([]int{2, 3}) {
			t.Error("expected false when calling with non-equal values")
		}
	})

	t.Run("with matchers", func(t *testing.T) {
		c := newCall([]int{1, 2}, 3)
		if !c.CalledWith(match.Contains(1)) {
			t.Error("expected true when calling with valid matcher")
		}
		if c.CalledWith(match.Contains(3)) {
			t.Error("expected false when calling with invalid matcher")
		}
	})
}

func TestCallCalledWithExactly(t *testing.T) {
	c := newCall(1, 2, 3)
	if !c.CalledWithExactly(1, 2, 3) {
		t.Error("expected true when calling exact number of arguments")
	}
	if c.CalledWithExactly(1, 2) {
		t.Error("expected false when calling with incorrect number of arguments")
	}
}

func TestCallCalledBefore(t *testing.T) {
	first := newCall()
	second := newCall()
	if !first.CalledBefore(second) {
		t.Error("expected first to be called before second")
	}
	if second.CalledBefore(first) {
		t.Error("expected second to not be called before first")
	}
}

func TestCallCalledAfter(t *testing.T) {
	first := newCall()
	second := newCall()
	if !second.CalledAfter(first) {
		t.Error("expected second to be called after first")
	}
	if first.CalledAfter(second) {
		t.Error("expected first to not be called after second")
	}
}
