package core

type Config struct {
	Log struct {
		Level string
	}
	API struct {
		Port uint16
	}
	IP struct {
		IPV4 []string
		IPV6 []string
	}
	Service struct {
		Interval uint32
	}
	Services map[string]map[string]any
}

type ServiceConfig struct {
	ADDR4    string
	ADDR6    string
	Type     string
	Protocol string
}

type Cloudflare struct {
	ServiceConfig
	Token  string
	Zone   string
	Record string
}
