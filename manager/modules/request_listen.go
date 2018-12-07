package modules

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func AddRequestListen(c *gin.Context) {
	domainId := c.Param("domain_id")
	listenPath := c.PostForm("listen_path")
	if domainId == "" || listenPath == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	domainRsp, err := domainData(domainId)
	if err != nil || domainRsp.Count == 0 {
		mContext{c}.ErrorOP(DataParseError)
		return
	}
	domain := new(Domain)
	err = json.Unmarshal(domainRsp.Kvs[0].Value, domain)
	if err != nil || domain == nil {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	reqListen := new(RequestListen)
	reqListen.DomainUrl = domain.DomainUrl
	reqListen.ListenPath = listenPath
	reqCopyB, _ := json.Marshal(reqListen)
	listenId := uuid.Must(uuid.NewV4()).String()
	err = putRequestListen(listenId, string(reqCopyB))
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	mContext{c}.SuccessOP(make([]string, 0))
}

func RequestsCopy(c *gin.Context) {
	var data []*RequestCopy
	rsp, err := requestsCopy()
	if err != nil || rsp.Count == 0 {
		mContext{c}.SuccessOP(make([]string, 0))
		return
	}

	for _,kv := range rsp.Kvs {
		reqCopy := new(RequestCopy)
		err := json.Unmarshal(kv.Value, reqCopy)
		if err != nil {
			continue
		}
		data = append(data, reqCopy)
	}
	mContext{c}.SuccessOP(data)
}
