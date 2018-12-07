package lb

import (
	"errors"
	"hgw/gateway/core"
)
type LoadBalance interface {
	Target() (*core.Target, error)
}
var ErrNoPointer = errors.New("no endpoints available")