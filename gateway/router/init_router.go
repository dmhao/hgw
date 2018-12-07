package router

import (
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/go-chi/chi"
	"hgw/gateway/core"
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
//创建路由选择器
func init() {
	hsMux = make(HostSwitch)
}

//获取当前路由选择器
func GetHs() HostSwitch {
	return hsMux
}

//初始化路由选择器
func InitHs() error {
	core.Sys().Infof("初始化配置路由")
	domainsData, err := core.DomainsData()
	if err != nil {
		return err
	}
	loadHs(domainsData)
	return nil
}

//加载所有路由数据
func loadHs(domainsData []*core.Domain) {
	newHsMux := make(HostSwitch)
	for _,domain := range domainsData {
		if _,ok := hsMux[domain.DomainUrl]; !ok {
			tmpMux := chi.NewMux()
			newHsMux[domain.DomainUrl] = tmpMux
		}
		mux := newHsMux[domain.DomainUrl]
		loadDomainHandler(mux, domain)
		core.Sys().Infof("【域名%s】-配置完毕", domain.DomainUrl)
	}
	hsMux = newHsMux
}

func loadDomainHandler(mux *chi.Mux,domain *core.Domain) {
	mux.Handle("/*", middleware.CreateMwChain(domain))
	if len(domain.Paths) > 0 {
		for _,path := range domain.Paths {
			mapMethodHandler(mux, path.ReqMethod, path.ReqPath, middleware.CreatePathMwChain(domain, path))
		}
	}
}

//重新加载域名数据的所有handler
func reloadDomainHandler(domain *core.Domain) {
	mux := chi.NewMux()
	loadDomainHandler(mux, domain)
	hs := GetHs()
	hs[domain.DomainUrl] = mux
}

func checkMuxExists(url string) {
	hsMux := GetHs()
	if _,ok := hsMux[url]; !ok {
		tmpMux := chi.NewMux()
		hsMux[url] = tmpMux
	}
}

func WatchDomainChange() {
	core.Sys().Infof("启动域名路径变动监听")
	go pathChangeListen()
	core.Sys().Infof("启动域名变动监听")
	go domainChangeListen()
}

//路径变动更新路由
func pathChangeListen() {
	ech := make(chan *clientv3.Event, 100)
	go core.WatchPaths(ech)
	for{
		select {
		case ev := <-ech:
			if ev.Type == clientv3.EventTypeDelete{
				dataPath := core.PathKToBakPathK(string(ev.Kv.Key))
				path, err := core.DomainPathData(dataPath)
				if err != nil || path == nil {
					core.Sys().Warnf("【域名-路径%s】备份数据获取失败", string(ev.Kv.Key))
					continue
				}

				domain, err := core.DomainDataById(path.DomainId, true)
				if err != nil {
					continue
				}
				if domain != nil {
					reloadDomainHandler(domain)
				}
				core.Sys().Infof("【域名%s】下路径删除，重新加载域名路由", domain.DomainUrl)
			} else if ev.Type == clientv3.EventTypePut {
				path := new(core.Path)
				err := json.Unmarshal(ev.Kv.Value, path)
				if err != nil {
					keySplit := strings.Split(string(ev.Kv.Key), "/")
					domainId := keySplit[3]
					core.Sys().Warnf("【域名%s】路径更新，json解析失败 %q", domainId, ev.Kv.Value)
					continue
				}
				domain,err  := core.DomainDataById(path.DomainId, false)
				if err != nil {
					continue
				}
				hsMux := GetHs()
				checkMuxExists(domain.DomainUrl)
				if mux, ok := hsMux[domain.DomainUrl]; ok {
					mapMethodHandler(mux, path.ReqMethod, path.ReqPath, middleware.CreatePathMwChain(domain, path))
				}
			}
		}
	}
}

//域名变动更新路由
func domainChangeListen() {
	ech := make(chan *clientv3.Event, 100)
	go core.WatchDomains(ech)
	for{
		select {
		case ev := <-ech:
			if ev.Type == clientv3.EventTypeDelete{
				domainBakPath := core.DomainKToBakDomainK(string(ev.Kv.Key))
				domain,err  := core.DomainDataByK(domainBakPath, false)
				if err != nil {
					core.Sys().Warnf("【域名路径%s】删除-获取备份数据失败", string(ev.Kv.Key))
					continue
				}
				delDomainClean(domain)
			} else if ev.Type == clientv3.EventTypePut {
				domain := new(core.Domain)
				err := json.Unmarshal(ev.Kv.Value, domain)
				if err != nil {
					keySplit := strings.Split(string(ev.Kv.Key), "/")
					domainId := keySplit[3]
					core.Sys().Warnf("【域名%s】更新，json解析失败 %s", domainId, string(ev.Kv.Value))
					continue
				}
				hsMux := GetHs()
				checkMuxExists(domain.DomainUrl)
				if mux, ok := hsMux[domain.DomainUrl]; ok {
					mux.Handle("/*", middleware.CreateMwChain(domain))
					core.Sys().Infof("【域名%s】更新", domain.DomainUrl)
				}
			}
		}
	}
}

func delDomainClean(domain *core.Domain) {
	hsMux := GetHs()
	delete(hsMux, domain.DomainUrl)
	core.DelDomainMetrics(domain)
	core.Sys().Infof("【域名%s】删除", domain.DomainUrl)
}

func mapMethodHandler(mux *chi.Mux, method string, pattern string, handler http.Handler) {
	met := strings.ToUpper(method)
	switch met {
	case "ALL":
		mux.Handle(pattern, handler)
	case "GET":
		mux.Method(method, pattern, handler)
	case "POST":
		mux.Method(method, pattern, handler)
	case "PUT":
		mux.Method(method, pattern, handler)
	case "PATCH":
		mux.Method(method, pattern, handler)
	case "DELETE":
		mux.Method(method, pattern, handler)
	case "OPTIONS":
		mux.Method(method, pattern, handler)
	case "HEAD":
		mux.Method(method, pattern, handler)
	default:
		mux.Handle(pattern, handler)
	}
}