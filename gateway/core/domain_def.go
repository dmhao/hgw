package core


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
	ReqMethod				string		`json:"req_method"`
	ReqPath					string		`json:"req_path"`
	CircuitBreakerRequest	int			`json:"circuit_breaker_request"`
	CircuitBreakerPercent	int			`json:"circuit_breaker_percent"`
	CircuitBreakerTimeout	int			`json:"circuit_breaker_timeout"`
	CircuitBreakerMsg		string		`json:"circuit_breaker_msg"`
	CircuitBreakerEnabled	bool		`json:"circuit_breaker_enabled"`
	CircuitBreakerForce		bool		`json:"circuit_breaker_force"`
	SetTime					string		`json:"set_time"`
}