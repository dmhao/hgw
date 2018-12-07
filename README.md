# hgw

hgw 是一套支持http/https协议的网关系统，由gateway服务、manager服务构成。

### 开发初衷
产品开发过程中版本多次迭代，缺乏服务管理，通过网关系统的负载均衡转发Api请求
可以清晰了解到各个模块的请求情况，对接口细节进行细致管理，提高服务的可用性。
有了网关还可以增加一些参数绑定、路径替换、权限验证、copy请求数据等功能，方便我们开发调试。

### 功能特性
+ 反向代理 
+ 黑名单
+ 访问日志
+ 限流
+ 熔断 【错误熔断、超时熔断、强制熔断】
+ 负载均衡 【轮询、随机、权重】
+ 请求路径过滤、路径级私有负载均衡
+ 访问统计 【访问量、访问时长】
+ 支持https, 动态设置https证书 【https证书热加载，动态更新】
+ 访问拷贝 【复制请求数据、返回数据】
#### 测试地址
网关服务：https://test.articlechain.cn/
控制服务：http://test.articlechain.cn:8080/admin/  账号：admin 密码：admin

<font color=red>尽量新增数据，请不要删除测试数据</font>

#### 域名配置
![hgw](https://github.com/dmhao/hgw/blob/master/img/hgw.gif)

#### 路径配置 请求拦截
![hgw-path](https://github.com/dmhao/hgw/blob/master/img/hgw-path.gif)

#### 强制熔断
![hgw-breaker](https://github.com/dmhao/hgw/blob/master/img/hgw-breaker.gif)

#### 访问统计
![domain-metrics](https://github.com/dmhao/hgw/blob/master/img/metrics.png)

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
--ser-name 【gateway服务的识别名称】

--addr 【gateway http服务的监听地址】

--tls-addr 【gateway https服务的监听地址】

--etcd 【连接etcd服务的地址】

--u 【连接etcd服务的账户】

--p 【连接etcd服务的密码】

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
--addr 【manager 服务的监听地址】

--etcd 【连接etcd服务的地址】

--u 【连接etcd服务的账户】

--p 【连接etcd服务的密码】

```
./manager --addr=0.0.0.0:8080 --etcd=127.0.0.1:2379
```

#### 6. 使用
访问 manager监听的服务地址+/admin/ 管理gateway服务。

<font color="red">**注**</font>： 自己搭建服务，第一次访问请先访问 /admin/init.html 初始化管理账号密码。


## 感谢
[x-admin](http://x.xuebingsi.com/) 后台管理html模板

[gin](https://github.com/gin-gonic/gin) Golang Web框架

[go-chi](https://github.com/go-chi/chi) Golang Router

[jwt-go](https://github.com/dgrijalva/jwt-go) Golang Jwt

[logrus](https://github.com/sirupsen/logrus) Golang Log

[hystrix-go](https://github.com/afex/hystrix-go) Golang CircuitBreaker