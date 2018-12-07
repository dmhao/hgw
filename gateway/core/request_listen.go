package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"sync"
	"time"
)

const (
	requestsListenDataPrefix = hgwPrefix + "requests-listen/"
	requestsCopyDataPrefix = hgwPrefix + "requests-copy/"
	requestsCopyDataFormat = requestsCopyDataPrefix + "%s"
)

const (
	copyChLen = 1000
	copyDataTTL = 120
)

var listenCopyCh chan *RequestCopy
var reqListenMap map[string][]*RequestListen
var reqListenLk *sync.RWMutex

func init() {
	listenCopyCh = make(chan *RequestCopy, copyChLen)
	reqListenMap = make(map[string][]*RequestListen)
	reqListenLk = new(sync.RWMutex)
}

func RequestListenRun() {
	//启动请求拷贝数据存储队列程序
	go requestCopyChRun()
	//加载请求监听数据
	reloadRequestsListen()
	//监听请求拷贝设置
	watchRequestListen()
}

func GetReqListenMap() map[string][]*RequestListen {
	return reqListenMap
}

func watchRequestListen() {
	for {
		rch := cli.Watch(context.Background(), requestsListenDataPrefix, clientv3.WithPrefix())
		for rsp := range rch {
			for _, ev := range rsp.Events {
				if ev.Type == clientv3.EventTypePut {
					fmt.Println(123)
					reqListen := new(RequestListen)
					err := json.Unmarshal(ev.Kv.Value, reqListen)
					if err != nil || reqListen == nil {
						continue
					}
					reqListenLk.Lock()
					if data, ok := reqListenMap[reqListen.DomainUrl]; ok {
						data = append(data, reqListen)
						reqListenMap[reqListen.DomainUrl] = data
					} else {
						var data []*RequestListen
						data = append(data, reqListen)
						reqListenMap[reqListen.DomainUrl] = data
					}
					reqListenLk.Unlock()
				} else if ev.Type == clientv3.EventTypeDelete {
					reloadRequestsListen()
				}
			}
		}
	}
}


func reloadRequestsListen() error {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, requestsListenDataPrefix, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return err
	}

	tmpMap := make(map[string][]*RequestListen)
	for _, val := range rsp.Kvs {
		reqListen := new(RequestListen)
		err := json.Unmarshal(val.Value, reqListen)
		if err != nil || reqListen == nil {
			continue
		}
		if data, ok := tmpMap[reqListen.DomainUrl]; ok {
			data = append(data, reqListen)
			tmpMap[reqListen.DomainUrl] = data
		} else {
			var data []*RequestListen
			data = append(data, reqListen)
			tmpMap[reqListen.DomainUrl] = data
		}
	}
	reqListenLk.Lock()
	reqListenMap = tmpMap
	reqListenLk.Unlock()
	return nil
}
func requestCopyDataK(listenId string) string {
	return fmt.Sprintf(requestsCopyDataFormat, listenId)
}

func putRequestCopy(copyId, copyJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	lease, err := cli.Grant(ctx, copyDataTTL)
	if err != nil {
		return err
	}
	_, err = cli.Put(ctx, requestCopyDataK(copyId), copyJson, clientv3.WithLease(lease.ID))
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func requestCopyChRun() {
	for {
		select {
			case data := <-listenCopyCh:
				dataB, err := json.Marshal(data)
				if err == nil {
					putRequestCopy(data.Id, string(dataB))
				}
		}
	}
}

func PutRequestCopy(requestCopy *RequestCopy) bool {
	t := time.NewTimer(time.Second)
	select {
	case listenCopyCh <- requestCopy:
		t.Stop()
		return true
	case <-t.C:
		return false
	}
}