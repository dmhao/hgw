package modules

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	DataParseError = -4001
	DataCannotDeleteError = -3001
)

const (
	SystemErrorNotInit = -2000
	SystemErrorNotLogin = -2001
	LoginParamsError = -2002
)

const (
	SystemError = -1001
)

type SuccessOutPut struct {
	Status		int				`json:"status"`
	Data 		interface{}		`json:"data"`
}

type ErrorOutPut struct {
	Status		int				`json:"status"`
	ErrorCode	int				`json:"error_code"`
}

type mContext struct {
	*gin.Context
}

func (c mContext) SuccessOP(data interface{}) {
	c.JSON(http.StatusOK, SuccessOutPut{Status: 1, Data: data})
	c.Next()
}

func (c mContext) ErrorOP(errorCode int) {
	c.AbortWithStatusJSON(http.StatusOK, ErrorOutPut{Status: 0, ErrorCode: errorCode})
}