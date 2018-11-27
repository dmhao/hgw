package middleware

import (
	"github.com/justinas/alice"
	"hgw/gateway/core"
	"hgw/gateway/def"
	"net/http"
)


type HgwMiddleWare interface {
	Init() func(http.Handler) http.Handler
}

func mwList(chain *[]alice.Constructor, hd func(http.Handler) http.Handler) bool {
	*chain = append(*chain, hd)
	return true
}

func CreateMwChain(domain *def.Domain) http.Handler {
	baseMw := &Base{Domain: domain, HandlerType: DomainHandler}
	baseMw.SetMt(core.NewDomainMetrics(domain))
	return createMwChain(baseMw)
}

func CreatePathMwChain(domain *def.Domain, path *def.Path) http.Handler {
	baseMw := &Base{Domain: domain, Path: path, HandlerType: DomainPathHandler}
	baseMw.SetMt(core.NewDomainPathMetrics(domain, path))
	return createMwChain(baseMw)
}

func createMwChain(base *Base) http.Handler {
	var chainArray []alice.Constructor
	mwList(&chainArray, base.Init())
	mwList(&chainArray, (&RecoverMw{base}).Init())
	mwList(&chainArray, (&BlackIpsMw{base}).Init())
	mwList(&chainArray, (&MetricsMw{base}).Init())
	mwList(&chainArray, (&LoggerMw{base}).Init())
	mwList(&chainArray, (&RateLimiterMw{base}).Init())
	mwList(&chainArray, (&BreakerMw{base}).Init())
	gw := &GateWay{base}
	gw.Init()
	hl := alice.New(chainArray...).Then(gw)
	return hl
}

