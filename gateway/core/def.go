package core

import "net/http"

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
	Paths				[]*Path				`json:"paths"`
}

type Target struct {
	Pointer			string		`json:"pointer"`
	Weight			int8		`json:"weight"`
	CurrentWeight	int8		`json:"current_weight"`
}

type Path struct {
	Id						string		`json:"id"`
	DomainId				string		`json:"domain_id"`
	ReqMethod				string		`json:"req_method"`
	ReqPath					string		`json:"req_path"`
	SearchPath				string		`json:"search_path"`
	ReplacePath				string		`json:"replace_path"`
	CircuitBreakerRequest	int			`json:"circuit_breaker_request"`
	CircuitBreakerPercent	int			`json:"circuit_breaker_percent"`
	CircuitBreakerTimeout	int			`json:"circuit_breaker_timeout"`
	CircuitBreakerMsg		string		`json:"circuit_breaker_msg"`
	CircuitBreakerEnabled	bool		`json:"circuit_breaker_enabled"`
	CircuitBreakerForce		bool		`json:"circuit_breaker_force"`
	PrivateProxyEnabled		bool		`json:"private_proxy_enabled"`
	LbType					string		`json:"lb_type"`
	Targets					[]*Target	`json:"targets"`
	SetTime					string		`json:"set_time"`
}

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