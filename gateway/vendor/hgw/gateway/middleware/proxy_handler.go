package middleware

import (
	"hgw/gateway/lb"
	"net/http"
	"net/url"
	"strings"
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
		lbType := p.GetDomain().LbType
		switch strings.ToUpper(lbType) {
		case LbRandom:
			p.SetLb(lb.NewRandom(p.GetDomain().Targets, time.Now().UnixNano()))
		case LbRoundRobin:
			p.SetLb(lb.NewRoundRobin(p.GetDomain().Targets))
		case LbWeightRoundRobin:
			p.SetLb(lb.NewWeightRoundRobin(p.GetDomain().Targets))
		default:
			p.SetLb(lb.NewRoundRobin(p.GetDomain().Targets))
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
	proxy := NewSingleHostReverseProxy(remote)
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

