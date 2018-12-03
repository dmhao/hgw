package modules

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

const (
	hgwCertsPrefix = hgwPrefix + "server-tls/"
	hgwCertFormat = hgwCertsPrefix + "%s"
	hgwCertBakPrefix = hgwPrefix + "server-tls-bak/"
	hgwCertBakFormat = hgwCertBakPrefix + "%s"
)

func certDataPath(certId string) string {
	return fmt.Sprintf(hgwCertFormat, certId)
}

func certBakDataPath(certId string) string {
	return fmt.Sprintf(hgwCertBakFormat, certId)
}


func certsData() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, hgwCertsPrefix, clientv3.WithPrefix())
	cancel()
	return rsp, err
}


func putCert(certId , certJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	_, err := cli.Put(ctx, certDataPath(certId), certJson)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func delCert(certId string) bool {
	dataPath := certDataPath(certId)
	dataBakPath := certBakDataPath(certId)
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	dataRsp, err  := cli.Get(ctx, dataPath)
	if err != nil {
		return false
	}
	data := dataRsp.Kvs[0].Value

	txn := cli.Txn(ctx)
	lease, err := cli.Grant(ctx, bakDataTTL)
	if err != nil {
		return false
	}
	rsp, err := txn.Then(clientv3.OpDelete(dataPath),
		clientv3.OpPut(dataBakPath, string(data), clientv3.WithLease(lease.ID))).Commit()
	if err != nil {
		return false
	}
	return rsp.Succeeded
}