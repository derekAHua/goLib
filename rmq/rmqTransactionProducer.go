package rmq

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// @Author: Derek
// @Description:
// @Date: 2022/5/1 00:08
// @Version 1.0

type (
	RocketMqTransactionProducer interface {
		rocketmq.TransactionProducer
	}
)

func NewRocketMqTransactionProducer(listener primitive.TransactionListener, opts ...producer.Option) (RocketMqTransactionProducer, error) {
	// TODO Implement RocketMqTransactionProducer.
	return producer.NewTransactionProducer(listener, opts...)
}
