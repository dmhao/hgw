package middleware

import (
	"net"
	"net/http"
)

type BlackIpsMw struct {
	Baser
}


func (mw *BlackIpsMw) Init() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ip := realIP(r)
			hgwResponse := rw.(hgwResponseWriter)
			hgwResponse.SetReqIp(ip)

			if _,ok := mw.GetDomain().BlackIps[ip]; ok {
				rw.WriteHeader(http.StatusForbidden)
			} else {
				next.ServeHTTP(rw, r)
			}
		})
	}
}

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
)

func realIP(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

