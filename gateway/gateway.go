package main

import (
	"crypto/tls"
	"gopkg.in/alecthomas/kingpin.v2"
	"hgw/gateway/core"
	"hgw/gateway/router"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type HgwGw struct {
}

func (h *HgwGw) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	hs := router.GetHs()
	hs.ServeHTTP(rw, r)
}

var (
	version = "0.1"
/*	serName = kingpin.Flag("ser-name", "SerName: gateway listen addr").Default("gateway-1").String()
	addr = kingpin.Flag("addr", "Addr: gateway listen addr").Default(":80").String()
	tlsAddr = kingpin.Flag("tls-addr", "Tls-Addr: gateway tls listen addr").Default(":443").String()
	etcd = kingpin.Flag("etcd", "Addr: etcd server addr").Default("127.0.0.1:2379").String()
	username = kingpin.Flag("u", "Username: etcd username").Default("").String()
	password = kingpin.Flag("p", "Password: etcd password").Default("").String()*/

	serName = kingpin.Flag("ser-name", "SerName: gateway listen addr").Required().String()
	addr = kingpin.Flag("addr", "Addr: gateway listen addr").Required().String()
	tlsAddr = kingpin.Flag("tls-addr", "Tls-Addr: gateway tls listen addr").Default("").String()
	etcd = kingpin.Flag("etcd", "Addr: etcd server addr").Required().String()
	username = kingpin.Flag("u", "Username: etcd username").Default("").String()
	password = kingpin.Flag("p", "Password: etcd password").Default("").String()
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version)
	kingpin.Parse()
	err := core.ConnectStore([]string{*etcd}, *username, *password)
	if err != nil {
		core.Sys().Fatalln("etcd 连接失败", err)
	}

	//初始化路由
	err = router.InitHs()
	if err != nil {
		core.Sys().Fatalln("etcd 读取路由失败", err)
	}

	//启动配置变动监听
	router.WatchDomainChange()


	if *addr != "" {
		core.Sys().Infof("GateWay启动 服务地址：%s etcd服务地址：%s", *addr, *etcd)
		go func() {
			err = http.ListenAndServe(*addr, &HgwGw{})
			if err != nil {
				core.Sys().Fatalln("端口80启动监听失败", err)
			}
		}()
	}

	if *tlsAddr != "" {
		go runTls(*tlsAddr, &HgwGw{})
	}

	go core.RecordMetrics(*serName)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(twg *sync.WaitGroup) {
		sig := make(chan os.Signal, 2)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		core.Sys().Infoln("服务关闭信号")
		unRegisterGateWay()
		twg.Done()
	}(wg)

	registerGateWay()
	wg.Wait()
	core.Sys().Infoln("服务已停止")
}


func runTls(addr string, handler http.Handler) {
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		core.Sys().Fatalln("tls ",addr,"监听失败", err)
	}
	core.InitCerts()
	go core.CertsChangeListen()

	//设置获取证书的方法
	tlsListener := tls.NewListener(ls, &tls.Config{GetCertificate: core.GetCert})
	//设置处理的handler方法
	ser := http.Server{Handler: handler}
	core.Sys().Infof("GateWay启动 服务地址：%s", addr)
	err = ser.Serve(tlsListener)
	if err != nil {
		core.Sys().Fatalln("tls ",addr,"关闭", err)
	}
}

func registerGateWay() {
}

func unRegisterGateWay() {
}