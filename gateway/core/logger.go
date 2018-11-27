package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var sysLog *logrus.Logger
var proxyLog *logrus.Logger


func init() {
	sysLog = logrus.New()
	proxyLog = logrus.New()
	proxyLog.Formatter = &logrus.JSONFormatter{}
}

func setNull(logger *logrus.Logger) {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	logger.Out = src
}

func Sys() *logrus.Logger {
	return sysLog
}

func Proxy() *logrus.Logger {
	return proxyLog
}