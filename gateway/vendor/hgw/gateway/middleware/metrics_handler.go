package middleware

import (
	"net/http"
	"time"
)

type MetricsMw struct {
	Baser
}

func (mw *MetricsMw) Init() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			n := time.Now()
			next.ServeHTTP(rw, r)
			mt := mw.GetMt()
			if mt != nil {
				if mt.Histograms != nil {
					mt.Histograms.Update(int64(time.Since(n)))
				}
			}
		})
	}
}