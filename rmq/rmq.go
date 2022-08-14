// Package rmq 提供了访问RocketMQ服务的能力
package rmq

import (
	"context"
	"errors"
	"github.com/derekAHua/goLib/zlog"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gin-gonic/gin"
)

var (
	// ErrRmqSvcConfigInvalid 服务配置无效
	ErrRmqSvcConfigInvalid = errors.New("requested rmq service is not correctly configured")
	// ErrRmqSvcNotRegistered 服务尚未被注册
	ErrRmqSvcNotRegistered = errors.New("requested rmq service is not registered")
	// ErrRmqSvcInvalidOperation 当前操作无效
	ErrRmqSvcInvalidOperation = errors.New("requested rmq service is not suitable for current operation")
)

var (
	rmqServices   = make(map[string]*client)
	rmqServicesMu sync.Mutex
)

// MessageCallback 定义业务方接收消息的回调接口
type MessageCallback func(ctx *gin.Context, msg Message) error

func (conf *ClientConfig) checkConfig() error {
	if conf.Group == "" {
		return ErrRmqSvcConfigInvalid
	}
	if conf.Topic == "" {
		return ErrRmqSvcConfigInvalid
	}
	if len(conf.NameServers) == 0 {
		return ErrRmqSvcConfigInvalid
	}
	return nil
}

func InitRmq(service string, config ClientConfig) (err error) {
	if err = config.checkConfig(); err != nil {
		return err
	}

	client := &client{
		ClientConfig: &config,
	}
	rmqServicesMu.Lock()
	defer rmqServicesMu.Unlock()

	err = client.startNamingHandler()
	if err != nil {
		return err
	}

	rmqServices[service] = client
	return nil
}

// StartProducer 启动指定已注册的RocketMQ生产服务
func StartProducer(service string) error {
	if client, ok := rmqServices[service]; ok {
		client.mu.Lock()
		defer client.mu.Unlock()
		if client.producer != nil {
			return ErrRmqSvcInvalidOperation
		}
		var err error
		var nsDomain string
		nsDomain, err = client.getNameserverDomain()
		if err != nil {
			return err
		}
		client.producer, err = newProducer(
			client.ClientConfig.Auth.AccessKey,
			client.ClientConfig.Auth.SecretKey,
			service,
			client.ClientConfig.Group,
			nsDomain,
			client.ClientConfig.Retry,
			time.Duration(client.ClientConfig.Timeout)*time.Millisecond)
		if err != nil {
			return err
		}
		return client.producer.start()
	}

	return ErrRmqSvcNotRegistered
}

// StopProducer 停止指定已注册的RocketMQ生产服务
func StopProducer(service string) error {
	if client, ok := rmqServices[service]; ok {
		client.mu.Lock()
		defer client.mu.Unlock()
		if client.producer == nil {
			return ErrRmqSvcInvalidOperation
		}
		err := client.producer.stop()
		client.producer = nil
		return err
	}
	return ErrRmqSvcNotRegistered
}

// StartConsumer 启动指定已注册的RocketMQ消费服务， 同时指定要消费的消息标签，以及消费回调
func StartConsumer(g *gin.Engine, service string, tags []string, callback MessageCallback) error {
	if _, exist := rmqServices[service]; !exist {
		return ErrRmqSvcNotRegistered
	}
	client := rmqServices[service]
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.pushConsumer != nil || callback == nil {
		return ErrRmqSvcInvalidOperation
	}
	var err error
	var nsDomain string
	nsDomain, err = client.getNameserverDomain()
	if err != nil {
		wrapLogger(zlog.ErrorLogger, nil, "invalid consumer nameServer", zap.Any("error", err))
		return err
	}

	cb := func(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, m := range msgList {
			if ctx.Err() != nil {
				wrapLogger(zlog.ErrorLogger, nil, "stop consume cause ctx cancelled", zap.Any("error", ctx.Err()))
				return consumer.SuspendCurrentQueueAMoment, ctx.Err()
			}
			if err := call(g, callback, m); err != nil {
				return consumer.SuspendCurrentQueueAMoment, nil
			}
		}
		return consumer.ConsumeSuccess, nil
	}

	client.pushConsumer, err = newPushConsumer(
		client.ClientConfig.Auth.AccessKey,
		client.ClientConfig.Auth.SecretKey,
		service,
		client.ClientConfig.Group,
		client.ClientConfig.Topic,
		client.ClientConfig.Broadcast,
		client.ClientConfig.Orderly,
		client.ClientConfig.Retry,
		tags,
		nsDomain,
		cb)
	if err != nil {
		wrapLogger(zlog.ErrorLogger, nil, "create new consumer error", zap.Any("error", err))
		return err
	}
	return client.pushConsumer.start()
}

// StopConsumer 停止指定已注册的RocketMQ消费服务
func StopConsumer(service string) error {
	if client, exist := rmqServices[service]; exist {
		client.mu.Lock()
		defer client.mu.Unlock()
		if client.pushConsumer == nil {
			return ErrRmqSvcInvalidOperation
		}
		err := client.pushConsumer.stop()
		client.pushConsumer = nil
		return err
	}
	return ErrRmqSvcNotRegistered
}

// NewMessage return a new Message.
func NewMessage(service string, content []byte) (Message, error) {
	if client, exist := rmqServices[service]; exist {
		return &messageWrapper{
			client: client,
			msg:    primitive.NewMessage(client.ClientConfig.Topic, content),
		}, nil
	}
	return nil, ErrRmqSvcNotRegistered
}

// save consumers.
var consumers []string

// Use will start the consumer of service.
func Use(g *gin.Engine, service string, tags []string, handler MessageCallback) {
	if err := StartConsumer(g, service, tags, handler); err != nil {
		panic("Start consumer  error: " + err.Error())
	}
	consumers = append(consumers, service)
}

func StopRocketMqConsume() {
	for _, svc := range consumers {
		_ = StopConsumer(svc)
	}
}
