package lb

import (
	"errors"
	"github.com/dmhao/hgw/gateway/core"
)
type LoadBalance interface {
	Target() (*core.Target, error)
}
var ErrNoPointer = errors.New("no endpoints available")