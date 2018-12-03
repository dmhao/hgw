package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/alecthomas/kingpin.v2"
	"hgw/manager/modules"
)

var (
/*	addr = kingpin.Flag("addr", "Addr: gateway listen addr").Default("127.0.0.1:8080").String()
	etcd = kingpin.Flag("etcd", "Addr: etcd server addr").Default("127.0.0.1:2379").String()
	username = kingpin.Flag("u", "Username: etcd username").Default("").String()
	password = kingpin.Flag("p", "Password: etcd password").Default("").String()*/
	addr = kingpin.Flag("addr", "gateway listen addr").Required().String()
	etcd = kingpin.Flag("etcd", "etcd server addr").Required().String()
	username = kingpin.Flag("u", "Username: etcd username").Default("").String()
	password = kingpin.Flag("p", "Password: etcd password").Default("").String()
)
func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	err := modules.ConnectStore([]string{*etcd}, *username, *password)
	if err != nil {
		panic(err)
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Static("/admin/", "admin")
	r.Static("/static/", "admin/static")

	v1 := r.Group("/v1", modules.AuthHandler)

	v1.POST("/domains/", modules.PutDomain)
	v1.GET("/domains/", modules.Domains)
	v1.POST("/domains/:domain_id", modules.PutDomain)
	v1.GET("/domains/:domain_id", modules.GetDomain)
	v1.DELETE("/domains/:domain_id", modules.DelDomain)

	v1.POST("/domains/:domain_id/paths/", modules.PutPath)
	v1.POST("/domains/:domain_id/paths/:path_id", modules.PutPath)
	v1.GET("/domains/:domain_id/paths/", modules.Paths)
	v1.GET("/domains/:domain_id/paths/:path_id", modules.GetPath)
	v1.DELETE("/domains/:domain_id/paths/:path_id", modules.DelPath)


	v1.POST("/certs/", modules.PutCert)
	v1.GET("/certs/", modules.Certs)
	v1.POST("/certs/:cert_id", modules.PutCert)
	v1.DELETE("/certs/:cert_id", modules.DelCert)

	v1.GET("/gateways/", modules.Gateways)
	v1.GET("/gateways/:server_name", modules.Gateway)


	r.GET("/index", modules.Index, modules.AuthHandler)
	r.POST("/init", modules.AuthInit)
	r.POST("/login", modules.Login)
	r.GET("/logout", modules.Logout)

	fmt.Printf("GW manager启动 服务地址：%s etcd服务地址：%s", *addr , *etcd)
	err = r.Run(*addr)
	if err != nil {
		fmt.Errorf("启动失败，端口监听 %v", err)
	}
}