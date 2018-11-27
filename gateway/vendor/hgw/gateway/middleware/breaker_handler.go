package middleware

import (
	"github.com/afex/hystrix-go/hystrix"
	"net/http"
)

type BreakerMw struct {
	Baser
}

func (mw *BreakerMw) Init() func(http.Handler) http.Handler {
	path := mw.GetPath()
	handlerType := mw.GetHandlerType()
	if handlerType == DomainPathHandler && path != nil {
		if path.CircuitBreakerRequest > 0 && path.CircuitBreakerPercent > 0 {
			cmdConf := hystrix.CommandConfig{
				MaxConcurrentRequests: path.CircuitBreakerRequest,
				ErrorPercentThreshold: path.CircuitBreakerPercent,
			}
			if path.CircuitBreakerTimeout > 0 {
				cmdConf.Timeout = path.CircuitBreakerTimeout
			}
			hystrix.ConfigureCommand(path.Id, cmdConf)
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			path := mw.GetPath()
			if path != nil && path.CircuitBreakerRequest > 0 && path.CircuitBreakerPercent > 0  {
				hgwResponse := rw.(hgwResponseWriter)
				hystrix.Do(path.Id, func() error {
					next.ServeHTTP(rw, r)
					select {
					case <-hgwResponse.ProxySuccessChan():
						return nil
					case err := <-hgwResponse.ProxyErrorChan():
						return err
					}
				}, func(err error) error {
					rw.WriteHeader(http.StatusOK)
					rw.Write([]byte(path.CircuitBreakerMsg))
					return err
				})

			} else {
				next.ServeHTTP(rw, r)
			}
		})
	}
}



