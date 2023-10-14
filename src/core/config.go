package core

type Config struct {
	Log struct {
		Level string
	}
	Store struct {
		Limit uint16
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

type Alibaba struct {
	ServiceConfig
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Domain    string
}

type Tencent struct {
	ServiceConfig
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Domain    string
}
