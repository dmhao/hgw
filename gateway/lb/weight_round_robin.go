package lb

import (
	"hgw/gateway/def"
	"sync"
)



func NewWeightRoundRobin(s []*def.Target) LoadBalance {
	for _, target := range s {
		if target.Weight == 0 {
			target.Weight = 1
		}
	}
	ss := &safeTargets{
		mtx: &sync.Mutex{},
		targets: s,
	}
	return &weightRoundRobin{
		s: ss,
		c: 0,
		l: len(s),
	}
}

type safeTargets struct {
	mtx			*sync.Mutex
	targets		[]*def.Target
}

type weightRoundRobin struct {
	s *safeTargets
	c uint64
	l int
}

func (wrr *weightRoundRobin) getNextPointerIndex() int {
	index := -1
	var total int8 = 0
	wrr.s.mtx.Lock()
	defer wrr.s.mtx.Unlock()

	for i := 0; i < wrr.l; i++ {
		wrr.s.targets[i].CurrentWeight += wrr.s.targets[i].Weight

		total += wrr.s.targets[i].Weight
		if index == -1 || wrr.s.targets[index].CurrentWeight < wrr.s.targets[i].CurrentWeight {
			index = i
		}
	}
	wrr.s.targets[index].CurrentWeight -= total
	return index
}

func (wrr *weightRoundRobin) Target() (*def.Target, error) {
	index := wrr.getNextPointerIndex()
	return wrr.s.targets[index], nil
}