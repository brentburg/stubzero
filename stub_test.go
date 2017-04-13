package stubzero

import (
	"testing"
)

func TestNew(t *testing.T) {
	s := New()
	if s.calls == nil || s.calls.Len() > 0 {
		t.Error("calls list is not initialized as an empty list")
	}
	if s.returns == nil || s.calls.Len() > 0 {
		t.Error("returns list is not initializes as an empty list")
	}
	if s.defaultReturn == nil || len(s.defaultReturn) > 0 {
		t.Error("default return is not initialized as an empty slice")
	}
}

func TestStubReset(t *testing.T) {
	s := New()
	s.ReturnsOnce(1)
	s.Returns(2)
	s.Call(1, 2)
	s.Reset()
	if s.calls == nil || s.calls.Len() > 0 {
		t.Error("calls list is not reset to an empty list")
	}
	if s.returns == nil || s.calls.Len() > 0 {
		t.Error("returns list is not reset to an empty list")
	}
	if s.defaultReturn == nil || len(s.defaultReturn) > 0 {
		t.Error("default return is not reset to an empty slice")
	}
}

func TestStubCall(t *testing.T) {
	s := New()
	s.Call()
	s.Call()
	if s.calls.Len() != 2 {
		t.Error("stub expected to have recorded 2 calls")
	}
}

func TestStubReturnsOnce(t *testing.T) {
	s := New()
	s.Returns(1)
	s.ReturnsOnce(5, 4)
	s.ReturnsOnce(3, 2)
	ret := s.Call()
	if ret[0].(int) != 5 || ret[1].(int) != 4 {
		t.Error("expected to return one time return values in order (first)")
	}
	ret = s.Call()
	if ret[0].(int) != 3 || ret[1].(int) != 2 {
		t.Error("expected to return one time return values in order (second)")
	}
	ret = s.Call()
	if ret[0].(int) != 1 {
		t.Error("expected to return default after all one time returns")
	}
}

func TestStubReturns(t *testing.T) {
	s := New()
	s.Returns(2)
	s.Returns(1, 2)
	ret := s.Call()
	if ret[0].(int) != 1 || ret[1].(int) != 2 {
		t.Error("expected to return values from most recent Returns call")
	}
}

func TestStubCallCount(t *testing.T) {
	s := New()
	if s.CallCount() != 0 {
		t.Error("expected call count to be 0")
	}
	s.Call()
	if s.CallCount() != 1 {
		t.Error("expected call count to be 1")
	}
	s.Call()
	s.Call()
	if s.CallCount() != 3 {
		t.Error("expected call count to be 3")
	}
}

func TestStubCalled(t *testing.T) {
	s := New()
	if s.Called() {
		t.Error("expected called to be false")
	}
	s.Call()
	if !s.Called() {
		t.Error("expected called to be true")
	}
}

func TestStubNotCalled(t *testing.T) {
	s := New()
	if !s.NotCalled() {
		t.Error("expected not called to be true")
	}
	s.Call()
	if s.NotCalled() {
		t.Error("expected not called to be false")
	}
}

func TestStubCalledOnce(t *testing.T) {
	s := New()
	if s.CalledOnce() {
		t.Error("expected called once to be false when not called")
	}
	s.Call()
	if !s.CalledOnce() {
		t.Error("expected called once to be true")
	}
	s.Call()
	if s.CalledOnce() {
		t.Error("expected called once to be false when called more than once")
	}
}

func TestStubFirstCall(t *testing.T) {
	s := New()
	if s.FirstCall() != nil {
		t.Error("expected first call to be nil if never called")
	}
	s.Call(1, 2)
	s.Call(1)
	if len(s.FirstCall().Args) != 2 {
		t.Error("expected first call to be call with correct number of arguments")
	}
}

func TestStubNthCall(t *testing.T) {
	s := New()
	s.Call(1)
	if s.NthCall(2) != nil {
		t.Error("expected 2nd call to be nil if called once")
	}
	s.Call(1, 2)
	s.Call(1, 3, 3)
	if s.NthCall(3) != nil && len(s.NthCall(3).Args) != 3 {
		t.Error("expected 3rd call to be call with correct number of arguments")
	}
}

func TestStubLastCall(t *testing.T) {
	s := New()
	if s.LastCall() != nil {
		t.Error("expected last call to be nil if never called")
	}
	s.Call(1, 2)
	if len(s.LastCall().Args) != 2 {
		t.Error("expected last call to be first call if called once")
	}
	s.Call(1, 2, 3)
	if len(s.LastCall().Args) != 3 {
		t.Error("expected last call to be call with correct number of arguments")
	}
}

func TestStubCalledBefore(t *testing.T) {
	s1 := New()
	s2 := New()
	s3 := New()
	s3.Call()
	s2.Call()
	s1.Call()
	s2.Call()
	if !s1.CalledBefore(s2) {
		t.Error("expected s1 to be called before s2")
	}
	if !s2.CalledBefore(s1) {
		t.Error("expected s1 to be called before s2")
	}
	if s1.CalledBefore(s3) {
		t.Error("expected s1 to no be called before s3")
	}
}

func TestStubCalledAfter(t *testing.T) {
	s1 := New()
	s2 := New()
	s3 := New()
	s2.Call()
	s1.Call()
	s2.Call()
	s3.Call()
	if !s1.CalledAfter(s2) {
		t.Error("expected s1 to be called after s2")
	}
	if !s2.CalledAfter(s1) {
		t.Error("expected s1 to be called after s2")
	}
	if s1.CalledAfter(s3) {
		t.Error("expected s1 to no be called after s3")
	}
}

func TestStubCalledWith(t *testing.T) {
	s := New()
	s.Call(1, 2, 3)
	s.Call(3, 4)
	if !s.CalledWith(1, 2) {
		t.Error("expected stub to be called with 1, 2")
	}
	if s.CalledWith(5, 6) {
		t.Error("expected stub to not be called with 5, 6")
	}
}

func TestStubCalledWithExactly(t *testing.T) {
	s := New()
	s.Call(1, 2, 3)
	s.Call(3, 4)
	if !s.CalledWithExactly(1, 2, 3) {
		t.Error("expected stub to be called with exactly 1, 2, 3")
	}
	if s.CalledWithExactly(1, 2) {
		t.Error("expected stub to be called with exactly 1, 2")
	}
	if s.CalledWithExactly(5, 6) {
		t.Error("expected stub to not be called with 5, 6")
	}
}

func TestStubAlwaysCalledWith(t *testing.T) {
	s := New()
	s.Call(1, 2, 3)
	s.Call(1, 2)
	if !s.AlwaysCalledWith(1, 2) {
		t.Error("expected stub to always be called with 1, 2")
	}
	s.Call(3, 4)
	if s.AlwaysCalledWith(1, 2) {
		t.Error("expected stub to not always be called with 1, 2")
	}
}

func TestStubAlwaysCalledWithExactly(t *testing.T) {
	s := New()
	s.Call(1, 2)
	s.Call(1, 2)
	if !s.AlwaysCalledWithExactly(1, 2) {
		t.Error("expected stub to always be called with exactly 1, 2")
	}
	s.Call(1, 2, 3)
	if s.AlwaysCalledWithExactly(1, 2) {
		t.Error("expected stub to not always be called with exactly 1, 2")
	}
}

func TestStubNeverCalledWith(t *testing.T) {
	s := New()
	s.Call(1, 2)
	if !s.NeverCalledWith(3, 4) {
		t.Error("expected stub to never be called with 3, 4")
	}
	s.Call(3, 4, 5)
	if s.NeverCalledWith(3, 4) {
		t.Error("expected stub to not be never be called with 3, 4")
	}
}

func TestStubNeverCalledWithExactly(t *testing.T) {
	s := New()
	s.Call(1, 2)
	if !s.NeverCalledWithExactly(3, 4) {
		t.Error("expected stub to never be called with 3, 4")
	}
	s.Call(3, 4, 5)
	if !s.NeverCalledWithExactly(3, 4) {
		t.Error("expected stub to never be called with excatly 3, 4")
	}
	if s.NeverCalledWithExactly(3, 4, 5) {
		t.Error("expected stub to not be never be called with exactly 3, 4, 5")
	}
}
