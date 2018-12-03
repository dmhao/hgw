package core

import (
	"crypto/tls"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
	"strings"
)

type Cert struct {
	Id				string			`json:"id"`
	SerName			string			`json:"ser_name"`
	CertBlock		string			`json:"cert_block"`
	CertKeyBlock	string			`json:"cert_key_block"`
}

var certMap map[string] *tls.Certificate

func init() {
	certMap = make(map[string] *tls.Certificate)
}

func InitCerts() {
	certs, err := certsData()
	if err == nil {
		for _, cert := range certs {
			certificate, err := tls.X509KeyPair([]byte(cert.CertBlock), []byte(cert.CertKeyBlock))
			if err != nil {
				Sys().Warnf("证书生成失败 %s", string(cert.SerName))
			}
			certMap[cert.SerName] = &certificate
		}
		Sys().Warnf("所有域名证书设置完成")
	}
}

func GetCert(info *tls.ClientHelloInfo) (certificate *tls.Certificate, e error) {
	if certMap == nil {
		return nil, errors.New("tls: no certificates configured")
	}
	name := strings.ToLower(info.ServerName)
	if cert, ok := certMap[name]; ok {
		return cert, nil
	}

	labels := strings.Split(name, ".")
	for i := range labels {
		labels[i] = "*"
		candidate := strings.Join(labels, ".")
		if cert, ok := certMap[candidate]; ok {
			return cert, nil
		}
	}
	return nil, errors.New("tls: no certificates configured")
}


func CertsChangeListen() {
	ech := make(chan *clientv3.Event, 100)
	go watchCerts(ech)
	for{
		select {
		case ev := <-ech:
			if ev.Type == clientv3.EventTypePut{
				cert := new(Cert)
				err := json.Unmarshal(ev.Kv.Value, cert)
				if err != nil {
					Sys().Warnf("证书数据解析失败 %s", string(ev.Kv.Value))
					continue
				}
				certificate, err := tls.X509KeyPair([]byte(cert.CertBlock), []byte(cert.CertKeyBlock))
				if err != nil {
					Sys().Warnf("证书生成失败 %s", string(ev.Kv.Value))
				}
				certMap[cert.SerName] = &certificate
				Sys().Infof("域名%s证书更新完成", cert.SerName)
			} else if ev.Type == clientv3.EventTypeDelete {
				certBak, err := certData(certBakDataPath(string(ev.Kv.Key)))
				if err != nil {
					Sys().Warnf("【域名证书路径%s】删除-获取备份数据失败", string(ev.Kv.Key))
					continue
				}
				delete(certMap, certBak.SerName)
			}
		}
	}

}