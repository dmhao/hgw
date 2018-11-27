package middleware

import (
	"github.com/sirupsen/logrus"
	"hgw/gateway/core"
	"net/http"
	"time"
)

type LoggerMw struct {
	Baser
}

func (mw *LoggerMw) Init() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(rw, r)
			hgwResponse := rw.(hgwResponseWriter)
			target := hgwResponse.ProxyTarget()
			var pointer string
			if target != nil {
				pointer = target.Pointer
			}

			core.Proxy().WithFields(logrus.Fields{
				"req-method": r.Method,
				"req-domain": r.Host,
				"req-path": r.RequestURI,
				"proxy-target": pointer,
				"proxy-total-time": time.Now().Sub(hgwResponse.StartTime()),
				"proxy-time": hgwResponse.ProxyUseTime(),
				"rsp-status": hgwResponse.Status(),
				"rsp-size": hgwResponse.Size(),
				"handler-type": mw.GetHandlerType(),
			}).Info()
		})
	}
}


