package modules

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

const (
	domainDataPrefix = hgwPrefix + "domain-data/"
	domainDataFormat = domainDataPrefix + "%s/"
	domainPathsDataFormat = hgwPrefix + "path-data/%s/"
	domainPathDataFormat = domainPathsDataFormat + "%s"

	domainBakDataPrefix = hgwPrefix + "domain-data-bak/"
	domainBakDataFormat = domainBakDataPrefix+"%s/"
)

//域名数据的路径
func domainDataPath(domainId string) string {
	return fmt.Sprintf(domainDataFormat, domainId)
}

//域名备份数据的路径
func domainBakDataPath(domainId string) string {
	return fmt.Sprintf(domainBakDataFormat, domainId)
}

//存储域名数据
func putDomain(domainId , domainJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	_, err := cli.Put(ctx, domainDataPath(domainId), domainJson)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

//删除域名并在备份目录备份数据
func delDomain(domainId string) bool {
	dataPath := domainDataPath(domainId)
	dataBakPath := domainBakDataPath(domainId)
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

func domainData(domainId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, domainDataPath(domainId))
	cancel()
	return rsp, err
}

//所有的域名数据
func domainsData() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, domainDataPrefix, clientv3.WithPrefix())
	cancel()
	return rsp, err
}

//域名路径数据的路径
func domainPathDataPath(domainId string, pathId string) string {
	return fmt.Sprintf(domainPathDataFormat, domainId, pathId)
}

func domainPathsDataPath(domainId string) string {
	return fmt.Sprintf(domainPathsDataFormat, domainId)
}

func pathsData(domainId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, domainPathsDataPath(domainId), clientv3.WithPrefix())
	cancel()
	return rsp, err
}

func pathData(domainId string, pathId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, domainPathDataPath(domainId, pathId), clientv3.WithPrefix())
	cancel()
	return rsp, err
}

func putPath(domainId, pathId string, pathJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	_, err := cli.Put(ctx, domainPathDataPath(domainId, pathId), pathJson)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

//删除域名并在备份目录备份数据
func delPath(domainId string, pathId string) bool {
	dataPath := domainPathDataPath(domainId, pathId)
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)

	rsp, err := cli.Delete(ctx, dataPath)
	cancel()
	if err != nil {
		return false
	}
	if rsp.Deleted == 0 {
		return false
	}
	return true
}
