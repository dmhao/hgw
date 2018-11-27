package lb

import (
	"hgw/gateway/def"
	"math/rand"
)

func NewRandom(s []*def.Target, seed int64) LoadBalance {
	return &random{
		s: s,
		r: rand.New(rand.NewSource(seed)),
		l: len(s),
	}
}

type random struct {
	s []*def.Target
	r *rand.Rand
	l int
}

func (r *random) Target() (*def.Target, error) {
	if len(r.s) <= 0 {
		return nil, ErrNoPointer
	}
	return r.s[r.r.Intn(r.l)], nil
}