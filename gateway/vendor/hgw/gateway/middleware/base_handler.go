package middleware

import (
	"github.com/didip/tollbooth/limiter"
	"hgw/gateway/core"
	"hgw/gateway/lb"
	"net/http"
	"time"
)

const (
	DomainHandler	= 1
	DomainPathHandler = 2
)

type Base struct {
	HandlerType		int8
	Domain			*core.Domain
	Path			*core.Path
	lb				lb.LoadBalance
	lmt				*limiter.Limiter
	mt				*core.Metrics
}

type Baser interface {
	GetHandlerType() int8
	GetDomain() *core.Domain
	GetPath()	*core.Path
	GetLb() 	lb.LoadBalance
	GetLmt()	*limiter.Limiter
	GetMt()		*core.Metrics
	SetHandlerType(int8)
	SetDomain(*core.Domain)
	SetPath(*core.Path)
	SetLb(lb.LoadBalance)
	SetLmt(*limiter.Limiter)
	SetMt(*core.Metrics)
}

func (mw *Base) GetHandlerType() int8 {
	return mw.HandlerType
}

func(mw *Base) GetDomain() *core.Domain {
	return mw.Domain
}

func (mw *Base) GetPath() *core.Path {
	return mw.Path
}

func (mw *Base) GetLb() lb.LoadBalance {
	return mw.lb
}

func (mw *Base) GetLmt() *limiter.Limiter {
	return mw.lmt
}

func (mw *Base) GetMt() *core.Metrics {
	return mw.mt
}

func (mw *Base) SetHandlerType(i int8) {
	mw.HandlerType = i
}

func (mw *Base) SetDomain(d *core.Domain) {
	mw.Domain = d
}

func (mw *Base) SetPath(p *core.Path) {
	mw.Path = p
}

func (mw *Base) SetLb(l lb.LoadBalance) {
	mw.lb = l
}

func (mw *Base) SetLmt(lmt *limiter.Limiter) {
	mw.lmt = lmt
}

func (mw *Base) SetMt(mt *core.Metrics) {
	mw.mt = mt
}

func (mw *Base) Init() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			hgwRw := &hgwResponse{}
			hgwRw.rw = rw
			hgwRw.startTime = time.Now()
			hgwRw.pErrorChan = make(chan error, 1)
			hgwRw.pSuccessChan = make(chan bool, 1)
			next.ServeHTTP(hgwRw, r)
		})
	}
}