package middleware

import (
	"hgw/gateway/def"
	"net/http"
	"time"
)

type hgwResponse struct {
	rw				http.ResponseWriter
	status 			int
	size   			int
	startTime		time.Time
	pErrorChan		chan error
	pSuccessChan	chan bool
	pUseTime		time.Duration
	pTarget			*def.Target
}

type hgwResponseWriter interface {
	Status() int
	Size()	int
	ProxyErrorChan()	chan error
	ProxySuccessChan()	chan bool
	StartTime()	time.Time
	SetProxyUseTime(time.Duration)
	ProxyUseTime() time.Duration
	SetProxyTarget(*def.Target)
	ProxyTarget() *def.Target
}

func (mw *hgwResponse) Header() http.Header {
	return mw.rw.Header()
}

func (mw *hgwResponse) Write(b []byte) (int, error) {
	size, err := mw.rw.Write(b)
	mw.size += size
	return size, err
}

func (mw *hgwResponse) WriteHeader(s int) {
	mw.rw.WriteHeader(s)
	mw.status = s
}

func (mw *hgwResponse) Status() int {
	return mw.status
}

func (mw *hgwResponse) Size() int {
	return mw.size
}

func (mw *hgwResponse) ProxyErrorChan() chan error {
	return mw.pErrorChan
}

func (mw *hgwResponse) ProxySuccessChan() chan bool {
	return mw.pSuccessChan
}

func (mw *hgwResponse) StartTime() time.Time {
	return mw.startTime
}

func (mw *hgwResponse) SetProxyUseTime(pUse time.Duration) {
	mw.pUseTime = pUse
}

func (mw *hgwResponse) ProxyUseTime() time.Duration {
	return mw.pUseTime
}

func (mw *hgwResponse) SetProxyTarget(t *def.Target) {
	mw.pTarget = t
}

func (mw *hgwResponse) ProxyTarget() *def.Target {
	return mw.pTarget
}