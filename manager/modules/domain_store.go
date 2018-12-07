package modules

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

const (
	domainsDataPrefix = hgwPrefix + "domain-data/"
	domainDataFormat = domainsDataPrefix + "%s/"
	domainPathsDataFormat = hgwPrefix + "path-data/%s/"
	domainPathDataFormat = domainPathsDataFormat + "%s"

	domainsBakDataPrefix = hgwPrefix + "domain-data-bak/"
	domainBakDataFormat = domainsBakDataPrefix+"%s/"

	domainPathsBakDataPrefix = hgwPrefix + "path-data-bak/%s/"
	domainPathBakDataFormat = domainPathsBakDataPrefix+"%s"
)

//域名数据的key
func domainDataK(domainId string) string {
	return fmt.Sprintf(domainDataFormat, domainId)
}

//域名备份数据的key
func domainBakDataK(domainId string) string {
	return fmt.Sprintf(domainBakDataFormat, domainId)
}

//域名路径数据的Key
func domainPathDataK(domainId string, pathId string) string {
	return fmt.Sprintf(domainPathDataFormat, domainId, pathId)
}

//域名路径数据的key
func domainPathBakDataK(domainId string, pathId string) string {
	return fmt.Sprintf(domainPathBakDataFormat, domainId, pathId)
}

func domainPathsDataK(domainId string) string {
	return fmt.Sprintf(domainPathsDataFormat, domainId)
}

//存储域名数据
func putDomain(domainId , domainJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	_, err := cli.Put(ctx, domainDataK(domainId), domainJson)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

//删除域名数据并备份数据
func delDomain(domainId string) bool {
	dataK := domainDataK(domainId)
	dataBakK := domainBakDataK(domainId)
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	dataRsp, err  := cli.Get(ctx, dataK)
	if err != nil {
		return false
	}
	data := dataRsp.Kvs[0].Value

	txn := cli.Txn(ctx)
	lease, err := cli.Grant(ctx, bakDataTTL)
	if err != nil {
		return false
	}
	rsp, err := txn.Then(clientv3.OpDelete(dataK),
		clientv3.OpPut(dataBakK, string(data), clientv3.WithLease(lease.ID))).Commit()
	if err != nil {
		return false
	}
	return rsp.Succeeded
}

//域名数据
func domainData(domainId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, domainDataK(domainId))
	cancel()
	return rsp, err
}

//所有的域名数据
func domainsData() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, domainsDataPrefix, clientv3.WithPrefix())
	cancel()
	return rsp, err
}

//所有路径数据
func pathsData(domainId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, domainPathsDataK(domainId), clientv3.WithPrefix())
	cancel()
	return rsp, err
}

//某个路径数据
func pathData(domainId string, pathId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, domainPathDataK(domainId, pathId), clientv3.WithPrefix())
	cancel()
	return rsp, err
}

//设置路径数据
func putPath(domainId, pathId string, pathJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	_, err := cli.Put(ctx, domainPathDataK(domainId, pathId), pathJson)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

//删除域名路径配置并备份数据
func delPath(domainId string, pathId string) bool {
	dataK := domainPathDataK(domainId, pathId)
	dataBakK := domainPathBakDataK(domainId, pathId)
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	dataRsp, err  := cli.Get(ctx, dataK)
	if err != nil {
		return false
	}
	data := dataRsp.Kvs[0].Value

	txn := cli.Txn(ctx)
	lease, err := cli.Grant(ctx, bakDataTTL)
	if err != nil {
		return false
	}
	rsp, err := txn.Then(clientv3.OpDelete(dataK),
		clientv3.OpPut(dataBakK, string(data), clientv3.WithLease(lease.ID))).Commit()
	if err != nil {
		return false
	}
	return rsp.Succeeded
}
