package modules

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"strconv"
	"time"
)

type Path struct {
	Id						string		`json:"id"`
	ReqMethod				string		`json:"req_method"`
	ReqPath					string		`json:"req_path"`
	CircuitBreakerRequest	int			`json:"circuit_breaker_request"`
	CircuitBreakerPercent	int			`json:"circuit_breaker_percent"`
	CircuitBreakerTimeout	int			`json:"circuit_breaker_timeout"`
	CircuitBreakerMsg		string		`json:"circuit_breaker_msg"`
	SetTime					string		`json:"set_time"`
}

func Paths(c *gin.Context) {
	domainId := c.Param("domain_id")
	if domainId == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	rsp, err := pathsData(domainId)
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	var paths []*Path
	if rsp.Count > 0 {
		for _, kv := range rsp.Kvs {
			path := new(Path)
			err := json.Unmarshal(kv.Value, path)
			if err == nil {
				paths = append(paths, path)
			}
		}
		mContext{c}.SuccessOP(paths)
		return
	}
	mContext{c}.SuccessOP(make([]string, 0))
}

func CreatePath(c *gin.Context) {
	domainId := c.Param("domain_id")
	reqMethod := c.PostForm("req_method")
	reqPath := c.PostForm("req_path")
	cbRequest := c.PostForm("circuit_breaker_request")
	cbPercent := c.PostForm("circuit_breaker_percent")
	cbTimeout := c.PostForm("circuit_breaker_timeout")
	cbMsg := c.PostForm("circuit_breaker_msg")

	if reqMethod == "" || reqPath == "" || domainId == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	var pathId string
	pathId = c.Param("path_id")
	if pathId == "" {
		pathId = uuid.Must(uuid.NewV4()).String()
	}

	path := new(Path)
	path.Id = pathId
	path.ReqMethod = reqMethod
	path.ReqPath = reqPath
	path.CircuitBreakerRequest,_ = strconv.Atoi(cbRequest)
	path.CircuitBreakerPercent,_ = strconv.Atoi(cbPercent)
	path.CircuitBreakerTimeout,_ = strconv.Atoi(cbTimeout)
	path.CircuitBreakerMsg = cbMsg
	path.SetTime = time.Now().Format("2006/1/2 15:04:05")

	pathB, err := json.Marshal(path)
	if err != nil {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	err = putPath(domainId, path.Id, string(pathB))
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	mContext{c}.SuccessOP(path)
}

func GetPath(c *gin.Context) {
	domainId := c.Param("domain_id")
	pathId := c.Param("path_id")
	if domainId == "" || pathId == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}
	rsp, err := pathData(domainId, pathId)
	if err != nil {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	if rsp.Count > 0 {
		path := new(Path)
		err := json.Unmarshal(rsp.Kvs[0].Value, path)
		if err != nil {
			mContext{c}.ErrorOP(DataParseError)
			return
		}
		mContext{c}.SuccessOP(path)
	} else {
		mContext{c}.SuccessOP(struct {}{})
	}
}

func DelPath(c *gin.Context) {
	domainId := c.Param("domain_id")
	pathId := c.Param("path_id")
	if domainId == "" || pathId == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	deleted := delPath(domainId, pathId)
	if deleted  {
		mContext{c}.SuccessOP(nil)
		return
	}
	mContext{c}.ErrorOP(SystemError)
}