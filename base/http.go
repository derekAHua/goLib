package base

import (
	"net/http"
	"sync"
	"time"
)

type ApiClient struct {
	Service        string        `yaml:"service"`
	AppKey         string        `yaml:"appKey"`
	AppSecret      string        `yaml:"appSecret"`
	Domain         string        `yaml:"domain"`
	Timeout        time.Duration `yaml:"timeout"`
	ConnectTimeout time.Duration `yaml:"connectTimeout"`
	Retry          int           `yaml:"retry"`
	HttpStat       bool          `yaml:"httpStat"`
	Host           string        `yaml:"host"`
	Proxy          string        `yaml:"proxy"`
	BasicAuth      struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}

	HTTPClient *http.Client
	clientInit sync.Once
}
