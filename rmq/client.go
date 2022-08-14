package rmq

import (
	"github.com/derekAHua/goLib/zlog"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/apache/rocketmq-client-go/v2/rlog"
)

// auth 提供链接到Broker所需要的验证信息（按需配置）
type auth struct {
	AccessKey string `json:"ak,omitempty" yaml:"ak,omitempty"`
	SecretKey string `json:"sk,omitempty" yaml:"sk,omitempty"`
}

// ClientConfig 包含链接到RocketMQ服务所需要的各配置项
type ClientConfig struct {
	// 集群名字
	Service string `json:"-" yaml:"service"`
	// 提供名字服务器的地址列表，例如: [ "127.0.0.1:9876" ]
	NameServers []string `json:"nameservers" yaml:"nameservers"`
	// 生产/消费者组名称，各业务线间需要保持唯一
	Group string `json:"group" yaml:"group"`
	// 要消费/订阅的主题
	Topic string `json:"topic" yaml:"topic"`
	// 如果配置了ACL，需提供验证信息
	Auth auth `json:"auth" yaml:"auth"`
	// 是否是广播消费模式
	Broadcast bool `json:"broadcast" yaml:"broadcast"`
	// 是否是顺序消费模式
	Orderly bool `json:"orderly" yaml:"orderly"`
	// 生产失败时的重试次数
	Retry int `json:"retry" yaml:"retry"`
	// 生产超时时间
	Timeout int `json:"timeout" yaml:"timeout"`
}

// Client 为客户端主体结构
type client struct {
	*ClientConfig
	mu sync.RWMutex

	producer       *rmqProducer
	pushConsumer   *rmqPushConsumer
	namingListener net.Listener
}

func (c *client) startNamingHandler() error {
	var err error
	c.namingListener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		wrapLogger(zlog.ErrorLogger, nil, "failed to create naming listener")
		return err
	}

	go func() {
		err = http.Serve(c.namingListener, c.createNamingHandler())
		wrapLogger(zlog.ErrorLogger, nil, "naming handler stopped", zap.String("error", err.Error()))
	}()

	return nil
}

func (c *client) getNameserverDomain() (string, error) {
	if c.namingListener != nil {
		return "https://" + c.namingListener.Addr().String(), nil
	}
	return "", ErrRmqSvcInvalidOperation
}

func (c *client) createNamingHandler() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		hostList := c.getHostListByDns()

		if hostList == "" {
			// no ns available
			resp.WriteHeader(http.StatusNotFound)
			return
		}

		_, err := resp.Write([]byte(hostList))
		if err != nil {
			wrapLogger(zlog.ErrorLogger, nil, "write response failed", zap.Any("error", err))
		}
		return
	}
}

func (c *client) getHostListByDns() (hostList string) {
	wrapLogger(zlog.DebugLogger, nil, "try serve through static config, nameServer",
		zap.Strings("ns", c.ClientConfig.NameServers),
	)

	var firstItem = true
	for _, ns := range c.ClientConfig.NameServers {
		var parts = strings.Split(ns, ":")
		if len(parts) != 2 {
			wrapLogger(zlog.WarnLogger, nil, "invalid nameserver config", zap.String("ns", ns))
			continue
		}

		var host = parts[0]
		var port = parts[1]
		// have to resolve the domain name to ips
		addressList, err := net.LookupHost(host)
		if err != nil {
			wrapLogger(zlog.WarnLogger, nil, "failed to lookup nameserver", zap.String("host", host))
			continue
		}

		for _, addr := range addressList {
			if !firstItem {
				hostList += ";"
			}

			hostList += addr + ":" + port
			firstItem = false
		}
	}

	return hostList
}

// DelayLevel 定义消息延迟发送的级别
type DelayLevel int

const (
	Second = DelayLevel(iota + 1)
	Seconds5
	Seconds10
	Seconds30
	Minute1
	Minutes2
	Minutes3
	Minutes4
	Minutes5
	Minutes6
	Minutes7
	Minutes8
	Minutes9
	Minutes10
	Minutes20
	Minutes30
	Hour1
	Hours2
)

func init() {
	rlog.SetLogger(&rmqLogger{})
}
