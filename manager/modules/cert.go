package modules

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"time"
)

type Cert struct {
	Id				string			`json:"id"`
	SerName			string			`json:"ser_name"`
	CertBlock		string			`json:"cert_block"`
	CertKeyBlock	string			`json:"cert_key_block"`
	SetTime			string			`json:"set_time"`
}



func Certs(c *gin.Context) {
	rsp, err := certsData()
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	var certs []*Cert
	if rsp.Count > 0 {
		for _, kv := range rsp.Kvs {
			cert := new(Cert)
			err := json.Unmarshal(kv.Value, cert)
			if err == nil {
				certs = append(certs, cert)
			}
		}
		mContext{c}.SuccessOP(certs)
		return
	}
	mContext{c}.SuccessOP(make([]string, 0))
}

func PutCert(c *gin.Context) {
	certBlock := c.PostForm("cert_block")
	certKeyBlock := c.PostForm("cert_key_block")
	serName := c.PostForm("ser_name")

	if serName == "" || certKeyBlock == "" || certBlock == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	_, err := tls.X509KeyPair([]byte(certBlock), []byte(certKeyBlock))
	if err != nil {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	//有接收到certId 就是修改操作， 否则就是新增
	var certId string
	certId = c.Param("cert_id")
	if certId == "" {
		certId = uuid.Must(uuid.NewV4()).String()
	}
	cert := new(Cert)
	cert.Id = certId
	cert.SerName = serName
	cert.CertBlock = certBlock
	cert.CertKeyBlock = certKeyBlock
	cert.SetTime = time.Now().Format("2006/1/2 15:04:05")

	certB, err := json.Marshal(cert)
	if err != nil {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	err = putCert(cert.Id, string(certB))
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	mContext{c}.SuccessOP(cert)
}

func DelCert(c *gin.Context) {
	certId := c.Param("cert_id")
	if certId == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	deleted := delCert(certId)
	if deleted  {
		mContext{c}.SuccessOP(nil)
		return
	}
	mContext{c}.ErrorOP(SystemError)
}