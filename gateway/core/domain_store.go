package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"strings"
)

const (
	domainsDataPrefix = hgwPrefix + "domain-data/"
	domainDataFormat = domainsDataPrefix + "%s/"
	domainPathsDataPrefix = hgwPrefix + "path-data/"
	domainPathsDataFormat = domainPathsDataPrefix + "%s/"

	domainPathsBakDataPrefix = hgwPrefix + "path-data-bak/"
	domainsBakDataPrefix = hgwPrefix + "domain-data-bak/"
)

//域名数据key
func domainDataK(domainId string) string {
	return fmt.Sprintf(domainDataFormat, domainId)
}

//路径key转路径备份的key
func PathKToBakPathK(dataK string) string {
	return strings.Replace(dataK, domainPathsDataPrefix, domainPathsBakDataPrefix, 1)
}

//域名key装域名备份的key
func DomainKToBakDomainK(domainK string) string {
	return strings.Replace(domainK, domainsDataPrefix, domainsBakDataPrefix, 1)
}

//域名下path数据的路径
func domainPathsDataK(domainId string) string {
	return fmt.Sprintf(domainPathsDataFormat, domainId)
}

func DomainDataByK(domainK string, getPaths bool) (*Domain, error) {
	return domainData(domainK, getPaths)
}

func DomainDataById(domainId string, getPaths bool) (*Domain, error) {
	return domainData(domainDataK(domainId), getPaths)
}

func DomainsData() ([]*Domain, error) {
	return domainsData(domainsDataPrefix)
}
//获取域名定义数据
func domainData(dataK string, getPaths bool) (*Domain, error) {
	domain := new(Domain)
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	rsp, err := cli.Get(ctx, dataK, clientv3.WithPrefix())
	cancel()
	if err != nil || rsp.Count == 0 {
		return domain, err
	}
	err = json.Unmarshal(rsp.Kvs[0].Value, domain)
	if err != nil {
		Sys().Warnf("域名数据json解析失败 key : %s val: %s err: %s", dataK, string(rsp.Kvs[0].Value), err)
	} else {
		if getPaths {
			pathsData, err := domainPathsData(domainPathsDataK(domain.Id))
			if err != nil {
				domain.Paths = pathsData
			}
		}
	}
	return domain, nil
}

//获取所有域名定义数据
func domainsData(dataK string) ([]*Domain, error) {
	var domainsData []*Domain
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	rsp, err := cli.Get(ctx, dataK, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return domainsData, err
	}
	for _, rr := range rsp.Kvs {
		domain := new(Domain)
		err := json.Unmarshal(rr.Value, domain)
		if err != nil {
			Sys().Warnf("域名数据json解析失败 key : %s val: %s err: %s", dataK, string(rr.Value), err)
		} else {
			pathsData, err := domainPathsData(domainPathsDataK(domain.Id))
			if err == nil {
				domain.Paths = append(domain.Paths, pathsData...)
			}
			domainsData = append(domainsData, domain)
		}
	}
	return domainsData, nil
}

//获取域名下所有定义路径数据
func domainPathsData(dataK string) ([]*Path, error) {
	var pathsData []*Path
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	rsp, err := cli.Get(ctx, dataK, clientv3.WithPrefix())
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
			Sys().Warnf("路径数据json解析失败 key : %s val: %s err: %s", dataK, string(rr.Value), err)
		}
	}
	return pathsData, nil
}

func DomainPathData(dataK string) (*Path, error) {
	path := new(Path)
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	rsp, err := cli.Get(ctx, dataK)
	cancel()
	if err != nil {
		return path, err
	}
	if rsp.Count > 0 {
		err = json.Unmarshal(rsp.Kvs[0].Value, path)
		if err != nil {
			Sys().Warnf("路径数据json解析失败 key : %s val: %s err: %s", dataK, string(rsp.Kvs[0].Value), err)
		}
	}
	return path, nil
}

func WatchDomains(e chan *clientv3.Event) {
	for {
		rch := cli.Watch(context.Background(), domainsDataPrefix, clientv3.WithPrefix())
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