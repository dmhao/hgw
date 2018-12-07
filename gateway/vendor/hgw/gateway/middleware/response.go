package middleware

import (
	"hgw/gateway/core"
	"net/http"
	"time"
)

type hgwResponse struct {
	rw				http.ResponseWriter
	status 			int
	size   			int
	rspBody			[]byte
	startTime		time.Time
	pErrorChan		chan error
	pSuccessChan	chan bool
	pUseTime		time.Duration
	pTarget			*core.Target
	reqIp			string
}

type hgwResponseWriter interface {
	Status() int
	Size()	int
	RspBody() []byte
	ProxyErrorChan()	chan error
	ProxySuccessChan()	chan bool
	StartTime()	time.Time
	SetProxyUseTime(time.Duration)
	ProxyUseTime() time.Duration
	SetProxyTarget(*core.Target)
	ProxyTarget() *core.Target
	SetReqIp(string)
	ReqIp() string
}

func (mw *hgwResponse) Header() http.Header {
	return mw.rw.Header()
}

func (mw *hgwResponse) Write(b []byte) (int, error) {
	size, err := mw.rw.Write(b)
	mw.size += size
	mw.rspBody = append(mw.rspBody, b...)
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

func (mw *hgwResponse) RspBody() []byte {
	return mw.rspBody
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

func (mw *hgwResponse) SetProxyTarget(t *core.Target) {
	mw.pTarget = t
}

func (mw *hgwResponse) ProxyTarget() *core.Target {
	return mw.pTarget
}

func (mw *hgwResponse) SetReqIp(ip string) {
	mw.reqIp = ip
}

func (mw *hgwResponse) ReqIp() string {
	return mw.reqIp
}