package lb

import (
	"github.com/dmhao/hgw/gateway/core"
	"sync/atomic"
)

func NewRoundRobin(s []*core.Target) LoadBalance {
	return &roundRobin{
		s: s,
		c: new(uint64),
		l: len(s),
	}
}

type roundRobin struct {
	s []*core.Target
	c *uint64
	l int
}


func (rr *roundRobin) Target() (*core.Target, error) {
	if rr.l <= 0 {
		return nil, ErrNoPointer
	}
	old := atomic.AddUint64(rr.c, 1) - 1
	idx := old % uint64(len(rr.s))
	return rr.s[idx], nil
}
