package base

import (
	"google.golang.org/grpc"
)

type GrpcClient struct {
	Service string `yaml:"service"`
	Domain  string `yaml:"domain"`

	//Timeout        time.Duration `yaml:"timeout"`
	//ConnectTimeout time.Duration `yaml:"connectTimeout"`
	//Retry          int           `yaml:"retry"`
	//HttpStat       bool          `yaml:"httpStat"`
	//Host           string        `yaml:"host"`
	//Proxy          string        `yaml:"proxy"`

	*grpc.ClientConn
}
