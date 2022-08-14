package rmq

import (
	"context"
	"github.com/derekAHua/goLib/zlog"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

var defaultAllMsgSelector = consumer.MessageSelector{
	Type:       consumer.TAG,
	Expression: "*",
}

type pushConsumerCallback func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)

func newPushConsumer(ak, sk string, instance, group, topic string, broadcast, orderly bool, retry int, tags []string, nsDomain string, cb pushConsumerCallback) (*rmqPushConsumer, error) {
	if broadcast {
		instance = instance + "-consumer"
	} else {
		instance = instance + "-" + strconv.Itoa(os.Getpid()) + "-consumer"
	}
	options := []consumer.Option{
		consumer.WithInstance(instance),
		consumer.WithGroupName(group),
		consumer.WithAutoCommit(true),
		consumer.WithNameServerDomain(nsDomain),
		consumer.WithConsumerOrder(orderly),
		consumer.WithConsumeMessageBatchMaxSize(1),
		consumer.WithMaxReconsumeTimes(int32(retry)),
		consumer.WithStrategy(consumer.AllocateByAveragely),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	}
	if broadcast {
		options = append(options, consumer.WithConsumerModel(consumer.BroadCasting))
	} else {
		options = append(options, consumer.WithConsumerModel(consumer.Clustering))
	}
	if ak != "" && sk != "" {
		options = append(options, consumer.WithCredentials(primitive.Credentials{
			AccessKey: ak,
			SecretKey: sk,
		}))
	}
	con, err := rocketmq.NewPushConsumer(options...)
	if err != nil {
		wrapLogger(zlog.ErrorLogger, nil, "failed to create consumer", zap.String("error", err.Error()))
		return nil, err
	}

	if len(tags) == 0 {
		err = con.Subscribe(topic, defaultAllMsgSelector, cb)
	} else if len(tags) == 1 {
		err = con.Subscribe(topic, consumer.MessageSelector{
			Type:       consumer.TAG,
			Expression: tags[0],
		}, cb)
	} else {
		err = con.Subscribe(topic, consumer.MessageSelector{
			Type:       consumer.TAG,
			Expression: strings.Join(tags, "||"),
		}, cb)
	}
	if err != nil {
		zlog.Error(nil, "failed to subscribe",
			logAgentTopic,
			zap.String("error", err.Error()))
		return nil, err
	}

	return &rmqPushConsumer{
		consumer: con,
	}, nil
}

type rmqPushConsumer struct {
	consumer rocketmq.PushConsumer
}

func (c *rmqPushConsumer) start() error {
	return c.consumer.Start()
}

func (c *rmqPushConsumer) stop() error {
	return c.consumer.Shutdown()
}
