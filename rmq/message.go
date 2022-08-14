package rmq

import (
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/derekAHua/goLib/zlog"
	"go.uber.org/zap"
)

// Message 消息提供的接口定义
type (
	Message interface {
		// WithTag 设置消息的标签Tag
		WithTag(string) Message
		// WithShard 设置消息的分片键
		WithShard(string) Message
		// WithDelay 设置消息的延迟等级
		WithDelay(DelayLevel) Message
		// Send 发送消息
		Send() (msgID string, err error)
		// GetContent 获取消息体内容
		GetContent() []byte
		// GetTag 获取消息标签
		GetTag() string
		// GetShard 获取消息分片键
		GetShard() string
		// GetID 获取消息ID
		GetID() string
	}

	messageWrapper struct {
		msg      *primitive.Message
		client   *client
		offsetID string
		msgID    string
	}
)

func (m *messageWrapper) WithTag(tag string) Message {
	m.msg = m.msg.WithTag(tag)
	return m
}

func (m *messageWrapper) WithShard(shard string) Message {
	m.msg = m.msg.WithShardingKey(shard)
	return m
}

func (m *messageWrapper) WithDelay(lvl DelayLevel) Message {
	m.msg = m.msg.WithDelayTimeLevel(int(lvl))
	return m
}

func (m *messageWrapper) Send() (msgID string, err error) {
	if m.client == nil {
		wrapLogger(zlog.ErrorLogger, nil, "client is not specified")
		return "", ErrRmqSvcInvalidOperation
	}
	m.client.mu.Lock()
	prod := m.client.producer
	m.client.mu.Unlock()
	if prod == nil {
		wrapLogger(zlog.ErrorLogger, nil, "producer not started")
		return "", ErrRmqSvcInvalidOperation
	}

	queue, id, offset, err := m.client.producer.SendMessage(m.msg)
	if err != nil {
		wrapLogger(zlog.ErrorLogger, nil, "failed to send message",
			zap.String("error", err.Error()),
			zap.String("message", m.msg.String()),
		)
		return "", err
	}

	fields := []zlog.Field{
		zap.String("message", m.msg.String()),
		zap.String("queue", queue),
		zap.String("msgId", id),
		zap.String("offsetId", offset),
	}
	wrapLogger(zlog.InfoLogger, nil, "rmq sent message", fields...)

	return offset, nil
}

func (m *messageWrapper) GetContent() []byte {
	return m.msg.Body
}

func (m *messageWrapper) GetTag() string {
	return m.msg.GetTags()
}

func (m *messageWrapper) GetShard() string {
	return m.msg.GetShardingKey()
}

func (m *messageWrapper) GetID() string {
	return m.offsetID
}

type MessageBatch []Message

func (batch MessageBatch) Send() (msgID string, err error) {
	var msgList = make([]*primitive.Message, 0)
	for _, m := range batch {
		msgList = append(msgList, m.(*messageWrapper).msg)
	}
	if len(msgList) < 1 {
		return "", ErrRmqSvcInvalidOperation
	}

	queue, id, offset, err := batch[0].(*messageWrapper).client.producer.SendMessage(msgList...)

	if err != nil {
		wrapLogger(zlog.ErrorLogger, nil, "failed to send message batch",
			zap.String("error", err.Error()),
		)
		return "", err
	}

	fields := []zlog.Field{
		zap.String("queue", queue),
		zap.String("msgId", id),
		zap.String("offsetId", offset),
	}
	wrapLogger(zlog.InfoLogger, nil, "sent message batch", fields...)

	return offset, nil
}
