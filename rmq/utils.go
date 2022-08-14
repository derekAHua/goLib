package rmq

import (
	"encoding/binary"
	"encoding/json"
	"github.com/derekAHua/goLib/env"
	"github.com/derekAHua/goLib/zlog"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logSDKTopic   = zlog.WithTopicField("rocketmq-sdk")
	logAgentTopic = zlog.WithTopicField("rmq-agent")
)

var _node *snowflake.Node

func init() {
	// first try to get env config
	id := getIDBasedOnEnviron()
	if id == 0 {
		// then try to get id based on ip
		id = binary.BigEndian.Uint16(net.ParseIP(env.LocalIP).To4()[2:])
	}

	snowflake.NodeBits = 16 // to hold the last half of the IP, eg. 172.29.240.120 -> 240.120 -> 61560
	snowflake.StepBits = 6  // 6bit for 64 IDs per millisecond, yielding 64000 qps, should be enough for nmq proxy

	var err error
	_node, err = snowflake.NewNode(int64(id))
	if err != nil {
		wrapLogger(zlog.ErrorLogger, nil, "failed to set worker id")
	}
}

func getIDBasedOnEnviron() uint16 {
	if os.Getenv("SNOWFLAKE_ID") != "" {
		wid, err := strconv.ParseUint(os.Getenv("SNOWFLAKE_ID"), 10, 16)
		if err != nil {
			return 0
		}

		return uint16(wid)
	}

	return 0
}

func generateSnowflake() int64 {
	return _node.Generate().Int64()
}

func wrapLogger(actualLogger func(*gin.Context, string, string, ...zapcore.Field), ctx *gin.Context, msg string, fields ...zapcore.Field) {
	logFields := []zapcore.Field{
		zlog.WithTopicField(zlog.LogNameRMQ),
	}

	if ctx != nil {
		logFields = append(logFields,
			zap.String("logId", zlog.GetLogId(ctx)),
			zap.String("requestId", zlog.GetRequestId(ctx)),
			zap.String("module", env.GetAppName()),
			zap.String("localIp", env.LocalIP),
		)
	}
	logFields = append(logFields, fields...)
	actualLogger(nil, zlog.LogNameRMQ, msg, logFields...)
}

func call(g *gin.Engine, fn MessageCallback, m *primitive.MessageExt) (err error) {
	ctx := &gin.Context{}
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]

			info, _ := json.Marshal(map[string]interface{}{
				"time":   time.Now().Format("2006-01-02 15:04:05"),
				"level":  "error",
				"module": "stack",
			})
			log.Printf("%s\n-------------------stack-start-------------------\n%v\n%s\n-------------------stack-end-------------------\n", string(info), r, buf)
		}
		ctx.Done()
	}()

	err = fn(ctx, &messageWrapper{
		msg:      &m.Message,
		offsetID: m.OffsetMsgId,
		msgID:    m.MsgId,
	})
	if err != nil {
		wrapLogger(zlog.ErrorLogger, ctx, "failed to consume message: "+err.Error())
	}

	var fields []zlog.Field
	fields = append(fields,
		zap.String("topic", m.Topic),
		zap.ByteString("body", m.Body),
		zap.String("msgID", m.MsgId),
		zap.String("offsetMsgID", m.OffsetMsgId),
		zap.String("consumeDelay", strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond)-m.BornTimestamp, 10)),
	)

	wrapLogger(zlog.InfoLogger, ctx, "rmq-access", fields...)

	return err
}
