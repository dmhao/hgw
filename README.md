# hgw
【http-reverse_proxy, http-gateway，hot-reload】

hgw是由gateway网关服务、manager控制服务构成的一套轻量级网关系统。目前支持http协议的服务控制。

hgw通过【ETCD】存储服务数据，状态监控。
1. 反向代理
2. 黑名单
3. 访问日志
4. 限流
5. 熔断
6. 负载均衡
7. 访问统计

## 安装
#### 1. 获取hgw代码
```
go get github.com/dmhao/hgw
```

gateway文件夹是网关服务的核心代码

manager文件夹是控制服务的核心代码


#### 2. 编译gateway
```
go build gateway.go
```

#### 3. 编译manager
```
go build manager.go
```

#### 4. 运行gateway
```
[root@localhost gateway]# ./gateway -h
usage: gateway --ser-name=SER-NAME --addr=ADDR --etcd=ETCD [<flags>]

Flags:
  -h, --help               Show context-sensitive help (also try --help-long and --help-man).
      --ser-name=SER-NAME  SerName: gateway listen addr
      --addr=ADDR          Addr: gateway listen addr
      --tls-addr=""        Tls-Addr: gateway tls listen addr
      --etcd=ETCD          Addr: etcd server addr
      --u=""               Username: etcd username
      --p=""               Password: etcd password
      --version            Show application version.
```

##### 参数
--ser-name [gateway服务的识别名称]

--addr[gateway http服务的监听地址]

--tls-addr[gateway https服务的监听地址]

--etcd[连接etcd服务的地址]

--u[连接etcd服务的账户]

--p[连接etcd服务的密码]

```
./gateway --ser-name=gateway-1 --addr=0.0.0.0:80 --etcd=127.0.0.1:2379
```

#### 5. 运行manager
```
[root@localhost manager]# ./manager -h
usage: manager --addr=ADDR --etcd=ETCD [<flags>]

Flags:
  -h, --help       Show context-sensitive help (also try --help-long and --help-man).
      --addr=ADDR  gateway listen addr
      --etcd=ETCD  etcd server addr
      --u=""       Username: etcd username
      --p=""       Password: etcd password
```
##### 参数
--addr[manager 服务的监听地址]

--etcd[连接etcd服务的地址]

--u[连接etcd服务的账户]

--p[连接etcd服务的密码]

```
./manager --addr=0.0.0.0:8080 --etcd=127.0.0.1:2379
```

#### 6. 使用
访问 manager监听的服务地址+/admin/ 管理gateway服务。

<font color="red">**注**</font>： 第一次访问未初始化后台账户请先访问 /admin/init.html 初始化。