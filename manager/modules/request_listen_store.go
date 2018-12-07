package modules

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"net/http"
)

type RequestListen struct {
	DomainUrl		string		`json:"domain_url"`
	ListenPath		string		`json:"listen_path"`
}


type RequestCopy struct {
	SerName		string			`json:"ser_name"`
	Id			string			`json:"id"`
	ReqTime		string			`json:"req_time"`
	ReqIp		string			`json:"req_ip"`
	ReqPath		string			`json:"req_path"`
	PostForm	interface{}		`json:"post_form"`
	Get			string			`json:"get"`
	ReqHeader	interface{}		`json:"req_header"`
	RspSize		int				`json:"rsp_size"`
	RspHeader	http.Header		`json:"rsp_header"`
	RspBody		string			`json:"rsp_body"`
}


const (
	requestsListenDataPrefix = hgwPrefix + "requests-listen/"
	requestListenDataFormat = requestsListenDataPrefix + "%s"

	requestsCopyDataPrefix = hgwPrefix + "requests-copy/"
)

var requestListenTTL int64 = 60
func requestListenDataK(listenId string) string {
	return fmt.Sprintf(requestListenDataFormat, listenId)
}
func putRequestListen (listenId string, copyJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	lease, err := cli.Grant(ctx, requestListenTTL)
	if err != nil {
		return err
	}
	_, err = cli.Put(ctx, requestListenDataK(listenId), copyJson, clientv3.WithLease(lease.ID))
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func requestsCopy() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	op := clientv3.WithLastKey()
	op = append(op, clientv3.WithLimit(500))
	rsp, err := cli.Get(ctx, requestsCopyDataPrefix, op...)
	cancel()
	return rsp, err
}
