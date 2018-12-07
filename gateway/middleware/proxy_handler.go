package middleware

import (
	"github.com/dmhao/hgw/gateway/lb"
	"net/http"
	"net/url"
	"time"
)

type GateWay struct {
	Baser
}
const (
	LbRandom = "RANDOM"
	LbRoundRobin = "ROUNDROBIN"
	LbWeightRoundRobin = "WEIGHTROUNDROBIN"
)

func (p *GateWay) Init() {
	if p.GetLb() == nil {
		path := p.GetPath()
		lbType := p.GetDomain().LbType
		targets := p.GetDomain().Targets

		if p.GetHandlerType() == DomainPathHandler && path != nil {
			if path.PrivateProxyEnabled {
				lbType = path.LbType
				targets = path.Targets
			}
		}

		switch lbType {
		case LbRandom:
			p.SetLb(lb.NewRandom(targets, time.Now().UnixNano()))
		case LbRoundRobin:
			p.SetLb(lb.NewRoundRobin(targets))
		case LbWeightRoundRobin:
			p.SetLb(lb.NewWeightRoundRobin(targets))
		default:
			p.SetLb(lb.NewRoundRobin(targets))
		}
	}
}

func(p *GateWay) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	hgwResponse := rw.(hgwResponseWriter)
	target, err := p.GetLb().Target()
	hgwResponse.SetProxyTarget(target)

	if err != nil {
		hgwResponse.ProxyErrorChan() <- err
	}
	remote, err := url.Parse(target.Pointer)
	if err != nil {
		hgwResponse.ProxyErrorChan() <- err
	}
	proxy := NewSingleHostReverseProxy(remote, p.GetPath())
	proxy.ErrorHandler = func(rw http.ResponseWriter, request *http.Request, e error) {
		hgwResponse.ProxyErrorChan() <- e
	}
	proxy.SuccessHandler = func(rw http.ResponseWriter, request *http.Request) {
		hgwResponse.ProxySuccessChan() <- true
	}

	startTime := time.Now()
	proxy.ServeHTTP(rw, r)
	hgwResponse.SetProxyUseTime(time.Now().Sub(startTime))
}

