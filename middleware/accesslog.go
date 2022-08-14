package middleware

import (
	"bytes"
	"fmt"
	"github.com/derekAHua/goLib/utils"
	"github.com/derekAHua/goLib/zlog"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	printRequestLen  = 10240
	printResponseLen = 10240
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	s = strings.Replace(s, "\n", "", -1)
	if w.body != nil {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

// AccessLog print access log.
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		var requestBody []byte
		if c.Request.Body != nil {
			var err error
			requestBody, err = c.GetRawData()
			if err != nil {
				zlog.WarnF(c, "Get http request body error: %s", err.Error())
			}
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}

		c.Writer = blw

		c.Set(zlog.ContextKeyUri, path)
		_ = zlog.GetLogId(c)
		_ = zlog.GetRequestId(c)
		// 处理请求
		c.Next()

		response := ""
		if blw.body != nil {
			if len(blw.body.String()) <= printResponseLen {
				response = blw.body.String()
			} else {
				response = blw.body.String()[:printResponseLen]
			}
		}

		bodyStr := string(requestBody)

		if c.Request.URL.RawQuery != "" {
			bodyStr += "&" + c.Request.URL.RawQuery
		}

		if len(bodyStr) > printRequestLen {
			bodyStr = bodyStr[:printRequestLen]
		}

		// 结束时间
		end := time.Now()

		// 固定notice
		commonFields := []zlog.Field{
			zap.String("os", getReqValueByKey(c, "os")),
			zap.String("userId", getReqValueByKey(c, "userId")),
			zap.String("uri", c.Request.RequestURI),
			zap.String("host", c.Request.Host),
			zap.String("method", c.Request.Method),
			zap.String("httpProto", c.Request.Proto),
			zap.String("handle", c.HandlerName()),
			zap.String("userAgent", c.Request.UserAgent()),
			zap.String("refer", c.Request.Referer()),
			zap.String("clientIp", utils.GetClientIp(c)),
			zap.String("cookie", getCookie(c)),
			zap.String("requestStartTime", utils.GetFormatRequestTime(start)),
			zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
			zap.Float64("cost", utils.GetRequestCost(start, end)),
			zap.String("requestParam", bodyStr),
			zap.Int("responseStatus", c.Writer.Status()),
			zap.String("response", response),
		}

		zlog.InfoLogger(c, zlog.LogNameAccess, "access", commonFields...)
	}
}

func getReqValueByKey(ctx *gin.Context, k string) string {
	if vs, exist := ctx.Request.Form[k]; exist && len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func getCookie(ctx *gin.Context) string {
	cStr := ""
	for _, c := range ctx.Request.Cookies() {
		cStr += fmt.Sprintf("%s=%s&", c.Name, c.Value)
	}
	return strings.TrimRight(cStr, "&")
}
