package lb

import (
	"errors"
	"hgw/gateway/def"
)
type LoadBalance interface {
	Target() (*def.Target, error)
}
var ErrNoPointer = errors.New("no endpoints available")