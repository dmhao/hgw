package router

import (
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi"
	"hgw/gateway/core"
	"hgw/gateway/def"
	"hgw/gateway/middleware"
	"net/http"
	"strings"
)

type HostSwitch map[string]*chi.Mux

func (hs HostSwitch) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(rw, r)
	} else {
		http.Error(rw, "Forbidden", 403)
	}
}


var hsMux HostSwitch
//创建路由
func init() {
	hsMux = make(HostSwitch)
}

//获取当前路由
func GetHs() HostSwitch {
	return hsMux
}

//初始化路由
func InitHs() error {
	core.Sys().Infof("读取数据初始化配置路由")
	domainsData, err := core.DomainsData()
	if err != nil {
		return err
	}
	reloadHs(domainsData)
	return nil
}

//重新加载路由数据
func reloadHs(domainsData []*def.Domain) {
	newHsMux := make(HostSwitch)
	for _,domain := range domainsData {
		if _,ok := newHsMux[domain.DomainUrl]; !ok {
			rt := chi.NewMux()
			newHsMux[domain.DomainUrl] = rt
		}
		mux := newHsMux[domain.DomainUrl]
		mux.Handle("/*", middleware.CreateMwChain(domain))
		if len(domain.Paths) > 0 {
			for _,path := range domain.Paths {
				mux.Handle(path.ReqPath, middleware.CreatePathMwChain(domain, path))
				core.Sys().Infof("【域名%s】【路径%s】-配置完毕", domain.DomainUrl, path.ReqPath)
			}
		}
		core.Sys().Infof("【域名%s】-配置完毕", domain.DomainUrl)
	}
	hsMux = newHsMux
}

func initRouter(domainId string) error {
	gin.Recovery()
	domainData, err := core.DomainDataById(domainId, true)
	if err != nil {
		return err
	}
	if domainData != nil {
		reloadRouter(domainData)
	}
	return nil
}

func reloadRouter(domain *def.Domain) {
	mux := chi.NewMux()
	mux.Handle("/*", middleware.CreateMwChain(domain))
	if len(domain.Paths) > 0 {
		for _,path := range domain.Paths {
			mux.Handle(path.ReqPath, middleware.CreatePathMwChain(domain, path))
		}
	}
	hs := GetHs()
	hs[domain.DomainUrl] = mux
}

func WatchDomainChange() {
	core.Sys().Infof("启动域名路径变动监听")
	go pathChange()
	core.Sys().Infof("启动域名变动监听")
	go domainChange()
}

//路径变动更新路由
func pathChange() {
	ech := make(chan *clientv3.Event, 100)
	go core.WatchPaths(ech)
	for{
		select {
		case ev := <-ech:
			var domainId string
			keySplit := strings.Split(string(ev.Kv.Key), "/")
			domainId = keySplit[3]

			if ev.Type == clientv3.EventTypeDelete{
				core.Sys().Infof("【域名%s】下路径删除，重新加载域名路由")
				initRouter(domainId)
				continue
			} else if ev.Type == clientv3.EventTypePut {
				domain,err  := core.DomainDataById(domainId, false)
				if err != nil {
					continue
				}
				hsMux := GetHs()
				if _,ok := hsMux[domain.DomainUrl]; !ok {
					tmpMux := chi.NewMux()
					hsMux[domain.DomainUrl] = tmpMux
				}
				if mux, ok := hsMux[domain.DomainUrl]; ok {
					pathDef := new(def.Path)
					err := json.Unmarshal(ev.Kv.Value, pathDef)
					core.Sys().Warnf("【域名%s】路径更新，json解析失败 %q", domainId, ev.Kv.Value)
					if err != nil {
						continue
					}
					mux.Handle(pathDef.ReqPath, middleware.CreatePathMwChain(domain, pathDef))
				}
			}
		}
	}
}

//域名变动更新路由
func domainChange() {
	ech := make(chan *clientv3.Event, 100)
	go core.WatchDomains(ech)
	for{
		select {
		case ev := <-ech:
			if ev.Type == clientv3.EventTypeDelete{
				domainBakPath := core.DomainToBakDomainPath(string(ev.Kv.Key))
				domain,err  := core.DomainDataByPath(domainBakPath, false)
				if err != nil {
					core.Sys().Warnf("【域名%s】删除-获取备份数据失败")
					continue
				}
				hsMux := GetHs()
				delete(hsMux, domain.DomainUrl)
				core.Sys().Infof("【域名%s】删除", domain.DomainUrl)
			} else if ev.Type == clientv3.EventTypePut {
				var domainId string
				keySplit := strings.Split(string(ev.Kv.Key), "/")
				domainId = keySplit[3]
				domain := new(def.Domain)
				err := json.Unmarshal(ev.Kv.Value, domain)
				if err != nil {
					core.Sys().Warnf("【域名%s】更新，json解析失败 %q", domainId, ev.Kv.Value)
					continue
				}
				hsMux := GetHs()
				if _,ok := hsMux[domain.DomainUrl]; !ok {
					tmpMux := chi.NewMux()
					hsMux[domain.DomainUrl] = tmpMux
				}
				if mux, ok := hsMux[domain.DomainUrl]; ok {
					mux.Handle("/*", middleware.CreateMwChain(domain))
					core.Sys().Infof("【域名%s】更新", domain.DomainUrl)
				}
			}
		}
	}
}