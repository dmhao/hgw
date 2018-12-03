package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"strings"
)

const (
	domainDataPrefix = hgwPrefix + "domain-data/"
	domainDataFormat = domainDataPrefix + "%s/"
	domainPathsDataPrefix = hgwPrefix + "path-data/"
	domainPathsDataFormat = domainPathsDataPrefix + "%s/"

	domainBakDataPrefix = hgwPrefix + "domain-data-bak/"
)

// 组装域名数据路径
func domainDataPath(domainId string) string {
	return fmt.Sprintf(domainDataFormat, domainId)
}

// 通过替换得到备份域名数据的路径
func DomainToBakDomainPath(domainPath string) string {
	return strings.Replace(domainPath, domainDataPrefix, domainBakDataPrefix, 1)
}

//域名下path数据的路径
func domainPathsDataPath(domainId string) string {
	return fmt.Sprintf(domainPathsDataFormat, domainId)
}

//获取域名定义数据
func domainData(dataPath string, getPaths bool) (*Domain, error) {
	domain := new(Domain)
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	rsp, err := cli.Get(ctx, dataPath, clientv3.WithPrefix())
	cancel()
	if err != nil || rsp.Count == 0 {
		return domain, err
	}
	err = json.Unmarshal(rsp.Kvs[0].Value, domain)
	if err == nil {
		if getPaths {
			pathsData, err := domainPathsData(domainPathsDataPath(domain.Id))
			if err != nil {
				domain.Paths = pathsData
			}
		}
	} else {
		Sys().Warnf("store [%s] json parse error, value : %q err: %s", dataPath, rsp.Kvs[0].Value, err)
	}
	return domain, nil
}

//获取所有域名定义数据
func domainsData(dataPath string) ([]*Domain, error) {
	var domainsData []*Domain
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	rsp, err := cli.Get(ctx, dataPath, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return domainsData, err
	}
	for _, rr := range rsp.Kvs {
		domain := new(Domain)
		err := json.Unmarshal(rr.Value, domain)
		if err == nil {
			pathsData, err := domainPathsData(domainPathsDataPath(domain.Id))
			if err == nil {
				domain.Paths = append(domain.Paths, pathsData...)
			}
			domainsData = append(domainsData, domain)
		} else {
			Sys().Warnf("store [%s] json parse error, value : %q err: %s", dataPath, rr.Value, err)
		}
	}
	return domainsData, nil
}

//获取域名下所有定义路径数据
func domainPathsData(dataPath string) ([]*Path, error) {
	var pathsData []*Path
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	rsp, err := cli.Get(ctx, dataPath, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return pathsData, err
	}
	for _, rr := range rsp.Kvs {
		path := new(Path)
		err := json.Unmarshal(rr.Value, path)
		if err == nil {
			pathsData = append(pathsData, path)
		} else {
			Sys().Warnf("store [%s] json parse error, value : %q err: %s", dataPath, rr.Value, err)
		}
	}
	return pathsData, nil
}

func DomainDataByPath(domainPath string, getPaths bool) (*Domain, error) {
	return domainData(domainPath, getPaths)
}

func DomainDataById(domainId string, getPaths bool) (*Domain, error) {
	return domainData(domainDataPath(domainId), getPaths)
}

func DomainsData() ([]*Domain, error) {
	return domainsData(domainDataPrefix)
}


func WatchDomains(e chan *clientv3.Event) {
	for {
		rch := cli.Watch(context.Background(), domainDataPrefix, clientv3.WithPrefix())
		for rsp := range rch {
			for _, ev := range rsp.Events {
				e <- ev
			}
		}
	}
}

func WatchPaths(e chan *clientv3.Event) {
	for {
		rch := cli.Watch(context.Background(), domainPathsDataPrefix, clientv3.WithPrefix())
		for rsp := range rch {
			for _, ev := range rsp.Events {
				e <- ev
			}
		}
	}
}