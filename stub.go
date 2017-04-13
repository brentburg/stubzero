package stubzero

import (
	"container/list"
)

type Stub struct {
	calls         *list.List
	returns       *list.List
	defaultReturn []interface{}
}

func New() *Stub {
	s := &Stub{
		calls:   list.New(),
		returns: list.New(),
	}
	s.Reset()
	return s
}

func (s *Stub) Reset() {
	s.calls.Init()
	s.returns.Init()
	s.defaultReturn = make([]interface{}, 0)
}

func (s *Stub) Call(args ...interface{}) []interface{} {
	s.calls.PushBack(newCall(args))
	if s.returns.Len() > 0 {
		return s.returns.Remove(s.returns.Front()).([]interface{})
	}
	return s.defaultReturn
}

func (s *Stub) ReturnsOnce(vals ...interface{}) {
	s.returns.PushBack(vals)
}

func (s *Stub) Returns(vals ...interface{}) {
	s.defaultReturn = vals
}

func (s *Stub) CallCount() int {
	return s.calls.Len()
}

func (s *Stub) Called() bool {
	return s.CallCount() > 0
}

func (s *Stub) NotCalled() bool {
	return s.CallCount() == 0
}

func (s *Stub) CalledOnce() bool {
	return s.CallCount() == 1
}

func (s *Stub) FirstCall() *Call {
	return s.NthCall(1)
}

func (s *Stub) NthCall(n int) *Call {
	if s.CallCount() < n {
		return nil
	}
	e := s.calls.Front()
	for i := 1; i < n; n++ {
		e = e.Next()
	}
	return e.Value.(*Call)
}

func (s *Stub) LastCall() *Call {
	if !s.Called() {
		return nil
	}
	return s.calls.Back().Value.(*Call)
}

func (s *Stub) CalledBefore(t *Stub) bool {
	if s.NotCalled() || t.NotCalled() {
		return false
	}
	return s.FirstCall().CalledBefore(t.LastCall())
}

func (s *Stub) CalledAfter(t *Stub) bool {
	if s.NotCalled() || t.NotCalled() {
		return false
	}
	return s.LastCall().CalledAfter(t.FirstCall())
}

func (s *Stub) CalledWith(args ...interface{}) bool {
	for e := s.calls.Front(); e != nil; e = e.Next() {
		if e.Value.(*Call).CalledWith(args...) {
			return true
		}
	}
	return false
}

func (s *Stub) CalledWithExactly(args ...interface{}) bool {
	for e := s.calls.Front(); e != nil; e = e.Next() {
		if e.Value.(*Call).CalledWithExactly(args...) {
			return true
		}
	}
	return false
}

func (s *Stub) AlwaysCalledWith(args ...interface{}) bool {
	for e := s.calls.Front(); e != nil; e = e.Next() {
		if !e.Value.(*Call).CalledWith(args...) {
			return false
		}
	}
	return true
}

func (s *Stub) AlwaysCalledWithExactly(args ...interface{}) bool {
	for e := s.calls.Front(); e != nil; e = e.Next() {
		if !e.Value.(*Call).CalledWithExactly(args...) {
			return false
		}
	}
	return true
}

func (s *Stub) NeverCalledWith(args ...interface{}) bool {
	return !s.CalledWith(args...)
}

func (s *Stub) NeverCalledWithExactly(args ...interface{}) bool {
	return !s.CalledWithExactly(args...)
}
