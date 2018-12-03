package modules

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/url"
	"strconv"
	"time"
)

var LbMap = map[string] bool {"roundRobin": true, "random": true}

type Domain struct {
	Id					string				`json:"id"`
	DomainName			string				`json:"domain_name"`
	DomainUrl			string				`json:"domain_url"`
	LbType				string				`json:"lb_type"`
	Targets				[]*Target			`json:"targets"`
	BlackIps			map[string]bool 	`json:"black_ips"`
	RateLimiterNum		float64				`json:"rate_limiter_num"`
	RateLimiterMsg		string				`json:"rate_limiter_msg"`
	RateLimiterEnabled	bool				`json:"rate_limiter_enabled"`
	SetTime				string				`json:"set_time"`
}

type Target struct {
	Pointer			string		`json:"pointer"`
	Weight			int8		`json:"weight"`
	CurrentWeight	int8		`json:"current_weight"`
}

func Domains(c *gin.Context) {
	rsp, err := domainsData()
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	var domains []*Domain
	if rsp.Count > 0 {
		for _, kv := range rsp.Kvs {
			domain := new(Domain)
			err := json.Unmarshal(kv.Value, domain)
			if err == nil {
				domains = append(domains, domain)
			}
		}
		mContext{c}.SuccessOP(domains)
		return
	}
	mContext{c}.SuccessOP(make([]string, 0))
}

func GetDomain(c *gin.Context) {
	domainId := c.Param("domain_id")
	if domainId == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}
	rsp, err := domainData(domainId)
	if err != nil {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	if rsp.Count > 0 {
		domain := new(Domain)
		err := json.Unmarshal(rsp.Kvs[0].Value, domain)
		if err != nil {
			mContext{c}.ErrorOP(DataParseError)
			return
		}
		mContext{c}.SuccessOP(domain)
	} else {
		mContext{c}.SuccessOP(struct {}{})
	}
}

func DelDomain(c *gin.Context) {
	domainId := c.Param("domain_id")
	if domainId == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	deleted := delDomain(domainId)
	if deleted  {
		mContext{c}.SuccessOP(nil)
		return
	}
	mContext{c}.ErrorOP(SystemError)
}

func PutDomain(c *gin.Context) {
	domainUrl := c.PostForm("domain_url")
	domainName := c.PostForm("domain_name")
	lbType := c.PostForm("lb_type")
	proxyTargets := c.PostForm("proxy_targets")
	blackIpsJson := c.PostForm("black_ips")
	rateLimiterNum := c.PostForm("rate_limiter_num")
	rateLimiterMsg := c.PostForm("rate_limiter_msg")
	rateLimiterEnabled := c.PostForm("rate_limiter_enabled")

	if domainUrl == "" || domainName == "" || lbType == "" || proxyTargets == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	urlParse, err := url.ParseRequestURI(domainUrl)
	if err != nil {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	//检查负载均衡模式
	if _, ok := LbMap[lbType]; !ok {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	//检查代理目标数据
	var targets []*Target
	err = json.Unmarshal([]byte(proxyTargets), &targets)
	if err != nil {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	//黑名单解析
	blackIps := make(map[string]bool)
	if blackIpsJson != "" {
		err := json.Unmarshal([]byte(blackIpsJson), &blackIps)
		if err != nil {
			mContext{c}.ErrorOP(DataParseError)
			return
		}
	}


	//有接收到domainId 就是修改操作， 否则就是新增
	var domainId string
	domainId = c.Param("domain_id")
	if domainId == "" {
		domainId = uuid.Must(uuid.NewV4()).String()
	}
	domain := new(Domain)
	domain.Id = domainId
	domain.DomainName = domainName
	domain.DomainUrl = urlParse.Host
	domain.LbType = lbType
	domain.Targets = targets
	domain.BlackIps = blackIps
	domain.RateLimiterNum,_ = strconv.ParseFloat(rateLimiterNum, 10)
	domain.RateLimiterMsg = rateLimiterMsg
	domain.RateLimiterEnabled,_ = strconv.ParseBool(rateLimiterEnabled)
	domain.SetTime = time.Now().Format("2006/1/2 15:04:05")

	domainB, err := json.Marshal(domain)
	if err != nil {
		mContext{c}.ErrorOP(DataParseError)
		return
	}

	err = putDomain(domain.Id, string(domainB))
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	mContext{c}.SuccessOP(domain)
}