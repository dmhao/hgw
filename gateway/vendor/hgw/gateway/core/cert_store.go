package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"strings"
)

const (
	hgwCertsPath = hgwPrefix + "server-tls/"
	hgwCertFormat = hgwCertsPath + "%s"
	hgwCertBakPrefix = hgwPrefix + "server-tls-bak/"
)

func certDataPath(certId string) string {
	return fmt.Sprintf(hgwCertFormat, certId)
}

func certBakDataPath(certPath string) string {
	return strings.Replace(certPath, hgwCertsPath, hgwCertBakPrefix, 1)
}

func watchCerts(e chan *clientv3.Event) {
	for {
		rch := cli.Watch(context.Background(), hgwCertsPath, clientv3.WithPrefix())
		for rsp := range rch {
			for _, ev := range rsp.Events {
				e <- ev
			}
		}
	}
}


func certsData() ([]*Cert, error) {
	var certs []*Cert
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, hgwCertsPath, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return certs, err
	}
	for _, rr := range rsp.Kvs {
		cert := new(Cert)
		err := json.Unmarshal(rr.Value, cert)
		if err == nil {
			certs = append(certs, cert)
		} else {
			Sys().Warnf("证书数据解析失败 %s, error: %s", string(rr.Value), err)
		}
	}
	return certs, nil
}

func certData(certPath string) (*Cert, error) {
	cert := new(Cert)
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	rsp, err := cli.Get(ctx, certPath, clientv3.WithPrefix())
	cancel()
	if err != nil || rsp.Count == 0 {
		return cert, err
	}
	err = json.Unmarshal(rsp.Kvs[0].Value, cert)
	if err == nil {
		return cert, nil
	}
	Sys().Warnf("store [%s] json parse error, value : %q err: %s", certPath, rsp.Kvs[0].Value, err)
	return nil, err
}
