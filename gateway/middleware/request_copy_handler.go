package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/dmhao/hgw/gateway/core"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var defaultMultipartMemory int64 = 32 << 20

type RequestCopyMw struct {
	Baser
}

func (mw *RequestCopyMw) Init() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			b := bytes.NewBuffer([]byte{})
			io.Copy(b, r.Body)
			vals, _ := url.ParseQuery(b.String())
			r.Body = ioutil.NopCloser(b)
			next.ServeHTTP(rw, r)

			reqListenMap := core.GetReqListenMap()
			hasListen := false
			if data, ok := reqListenMap[r.Host]; ok {
				for _, v := range data {
					if r.URL.Path == v.ListenPath {
						hasListen = true
					}
				}
			}
			if hasListen {
				t := time.Now()
				hgwRsp := rw.(hgwResponseWriter)
				requestCopy := new(core.RequestCopy)
				requestCopy.SerName = core.GateWaySerName()
				//时间纳秒
				requestCopy.Id = strconv.FormatInt(t.UnixNano(), 10)

				//请求时间
				requestCopy.ReqTime = t.Format("2006-01-02 15:04:05")

				//请求的路径
				requestCopy.ReqPath = r.URL.Path

				//POST 请求的数据
				requestCopy.PostForm = vals

				//请求的URI路径
				requestCopy.Get = r.RequestURI

				//请求的header
				requestCopy.ReqHeader = r.Header

				//记录请求IP
				requestCopy.ReqIp = hgwRsp.ReqIp()

				//记录放回大小
				requestCopy.RspSize = hgwRsp.Size()

				//记录返回数据头信息
				requestCopy.RspHeader = rw.Header()

				//记录返回数据  判断是否需要gzip
				needGzip := false
				if headerVal, ok := requestCopy.RspHeader["Content-Encoding"]; ok {
					for _, data := range headerVal {
						if strings.Index(data, "gzip") != -1 {
							needGzip = true
						}
					}
				}
				if needGzip {
					requestCopy.RspBody = gzipByteToString(hgwRsp.RspBody())
				} else {
					requestCopy.RspBody = string(hgwRsp.RspBody())
				}
				core.PutRequestCopy(requestCopy)
			}
		})
	}
}

func gzipByteToString(p []byte) string {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(p)
	read,_:= gzip.NewReader(buf)
	data, _:=ioutil.ReadAll(read)
	return string(data)
}