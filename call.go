package stubzero

import (
	"time"

	"github.com/brentburg/stubzero/match"
)

type Call struct {
	Args []interface{}
	Time time.Time
}

func newCall(args ...interface{}) *Call {
	return &Call{
		Args: args,
		Time: time.Now(),
	}
}

func (c *Call) CalledWith(args ...interface{}) bool {
	for i, arg := range args {
		if !match.Match(arg, c.Args[i]) {
			return false
		}
	}
	return true
}

func (c *Call) CalledWithExactly(args ...interface{}) bool {
	if len(args) != len(c.Args) {
		return false
	}
	return c.CalledWith(args...)
}

func (c *Call) CalledBefore(d *Call) bool {
	return c.Time.Before(d.Time)
}

func (c *Call) CalledAfter(d *Call) bool {
	return c.Time.After(d.Time)
}
