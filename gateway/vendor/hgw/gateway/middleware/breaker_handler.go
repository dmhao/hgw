package middleware

import (
	"github.com/afex/hystrix-go/hystrix"
	"hgw/gateway/core"
	"net/http"
)

type BreakerMw struct {
	Baser
}

func (mw *BreakerMw) Init() func(http.Handler) http.Handler {
	handlerType := mw.GetHandlerType()
	if handlerType == DomainPathHandler {
		path := mw.GetPath()
		if path.CircuitBreakerEnabled {
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
			if path == nil {
				next.ServeHTTP(rw, r)
			} else {
				if path.CircuitBreakerForce {
					breakerResponse(rw, path)
				} else {
					if path.CircuitBreakerEnabled {
						hystrix.NewStreamHandler()
						hgwResponse := rw.(hgwResponseWriter)
						_ = hystrix.Do(path.Id, func() error {
							next.ServeHTTP(rw, r)
							select {
							case <-hgwResponse.ProxySuccessChan():
								return nil
							case err := <-hgwResponse.ProxyErrorChan():
								return err
							}
						}, func(err error) error {
							err.Error()
							breakerResponse(rw, path)
							return err
						})

					} else {
						next.ServeHTTP(rw, r)
					}
				}
			}

		})
	}
}

var defaultCircuitBreakerMsg = "CircuitBreaker"

func breakerResponse(rw http.ResponseWriter, path *core.Path) {
	breakerMsg := defaultCircuitBreakerMsg
	if path.CircuitBreakerMsg != "" {
		breakerMsg = path.CircuitBreakerMsg
	}
	rw.WriteHeader(http.StatusOK)
	_,_ = rw.Write([]byte(breakerMsg))
}



