package middleware

import (
	"github.com/didip/tollbooth"
	"net/http"
)

type RateLimiterMw struct {
	Baser
}

func (mw *RateLimiterMw) Init() func(http.Handler) http.Handler {
	if mw.GetLmt() == nil {
		if mw.GetDomain().RateLimiterNum > 0 {
			mw.SetLmt(tollbooth.NewLimiter(mw.GetDomain().RateLimiterNum, nil))
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if mw.GetLmt() != nil {
				httpError := tollbooth.LimitByRequest(mw.GetLmt(), rw, r)
				if httpError != nil {
					if mw.GetDomain().RateLimiterMsg != "" {
						rw.WriteHeader(http.StatusOK)
						rw.Write([]byte(mw.GetDomain().RateLimiterMsg))
					} else {
						rw.Header().Add("Content-Type", mw.GetLmt().GetMessageContentType())
						rw.WriteHeader(httpError.StatusCode)
						rw.Write([]byte(httpError.Message))
					}
				} else {
					next.ServeHTTP(rw, r)
				}
			} else {
				next.ServeHTTP(rw, r)
			}
		})
	}
}

