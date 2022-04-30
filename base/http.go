package base

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/derekAHua/goLib/consts"
	"github.com/derekAHua/goLib/env"
	"github.com/derekAHua/goLib/utils"
	"github.com/derekAHua/goLib/zlog"
	"go.uber.org/zap"

	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

type TransportOption struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration

	CustomTransport *http.Transport
}

var globalTransport *http.Transport

// InitHttp init global transport.
func InitHttp(opts *TransportOption) {
	if opts == nil {
		globalTransport = &http.Transport{
			MaxIdleConns:        500,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     300 * time.Second,
		}
	} else if opts.CustomTransport != nil {
		globalTransport = opts.CustomTransport
	} else {
		globalTransport = &http.Transport{
			MaxIdleConns:        opts.MaxIdleConns,
			MaxIdleConnsPerHost: opts.MaxIdleConnsPerHost,
			IdleConnTimeout:     opts.IdleConnTimeout,
		}
	}
}

type HttpRequestOptions struct {
	Encode      string
	Data        map[string]string // Use when Get request
	RequestBody interface{}       // Use by Encode
	Headers     map[string]string
	Cookies     map[string]string
	BodyType    string
	RetryPolicy RetryPolicy
	// 重试间隔策略
	BackOffPolicy BackOffPolicy
}

func (o *HttpRequestOptions) GetContentType() (cType string) {
	switch o.Encode {
	case consts.EncodeJson:
		cType = "application/json"
	case consts.EncodeForm:
		fallthrough
	default:
		cType = "application/x-www-form-urlencoded"
	}
	return cType
}

func (o *HttpRequestOptions) GetRequestBody() (encodeData string, err error) {
	if o.RequestBody == nil {
		return encodeData, nil
	}

	switch o.Encode {
	case consts.EncodeJson:
		reqBody, e := json.Marshal(o.RequestBody)
		encodeData, err = string(reqBody), e
	case consts.EncodeRaw:
		ok := true
		encodeData, ok = o.RequestBody.(string)
		if !ok {
			err = errors.New("raw data need string type")
		}
	case consts.EncodeForm:
		fallthrough
	default:
		v := url.Values{}
		if data, ok := o.RequestBody.(map[string]string); ok {
			for key, value := range data {
				v.Add(key, value)
			}
		} else if data, ok := o.RequestBody.(map[string]interface{}); ok {
			for key, value := range data {
				var vStr string
				switch value.(type) {
				case string:
					vStr = value.(string)
				default:
					if tmp, err := jsoniter.Marshal(value); err != nil {
						return encodeData, err
					} else {
						vStr = string(tmp)
					}
				}
				v.Add(key, vStr)
			}
		} else {
			return encodeData, errors.New("unSupport RequestBody type")
		}
		encodeData, err = v.Encode(), nil
	}
	return encodeData, err
}

func (o *HttpRequestOptions) GetUrlData() (string, error) {
	v := url.Values{}
	for key, value := range o.Data {
		v.Add(key, value)
	}

	return v.Encode(), nil
}

func (o *HttpRequestOptions) GetRetryPolicy() RetryPolicy {
	r := defaultRetryPolicy
	if o.RetryPolicy != nil {
		r = o.RetryPolicy
	}
	return r
}

func (o *HttpRequestOptions) GetBackOffPolicy() BackOffPolicy {
	b := defaultBackOffPolicy
	if o.BackOffPolicy != nil {
		b = o.BackOffPolicy
	}

	return b
}

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

func (client *ApiClient) GetTransPort() *http.Transport {
	trans := globalTransport
	if client.Proxy != "" {
		trans.Proxy = func(_ *http.Request) (*url.URL, error) {
			return url.Parse(client.Proxy)
		}
	} else {
		trans.Proxy = nil
	}

	if client.ConnectTimeout != 0 {
		trans.DialContext = (&net.Dialer{
			Timeout: client.ConnectTimeout,
		}).DialContext
	} else {
		trans.DialContext = nil
	}

	return trans
}

func (client *ApiClient) makeRequest(ctx *gin.Context, method, url string, data io.Reader, opts HttpRequestOptions) (*http.Request, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	if opts.Headers != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}

	if client.Host != "" {
		req.Host = client.Host
	} else if h := req.Header.Get("host"); h != "" {
		req.Host = h
	}

	for k, v := range opts.Cookies {
		req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	if client.BasicAuth.Username != "" {
		req.SetBasicAuth(client.BasicAuth.Username, client.BasicAuth.Password)
	}

	cType := opts.BodyType
	if cType == "" {
		cType = opts.GetContentType()
	}
	req.Header.Set("Content-Type", cType)

	req.Header.Set(consts.HttpHeaderService, env.GetAppName())
	req.Header.Set(consts.TraceHeaderKey, zlog.GetRequestId(ctx))
	req.Header.Set(consts.LogIdHeaderKey, zlog.GetLogId(ctx))

	return req, nil
}

func (client *ApiClient) HttpGet(ctx *gin.Context, path string, opts HttpRequestOptions) (*ApiResult, error) {
	urlData, err := opts.GetUrlData()
	if err != nil {
		zlog.ErrorLogger(ctx, zlog.LogNameServer, fmt.Sprintf("Http client make GetUrlData error: %v", err), zlog.WithTopicField(zlog.LogNameModule))
		return nil, err
	}

	var requestUrl string
	if urlData == "" {
		requestUrl = fmt.Sprintf("%s%s", client.Domain, path)
	} else {
		requestUrl = fmt.Sprintf("%s%s?%s", client.Domain, path, urlData)
	}

	req, err := client.makeRequest(ctx, "GET", requestUrl, nil, opts)
	if err != nil {
		zlog.ErrorLogger(ctx, zlog.LogNameServer, fmt.Sprintf("http client makeRequest error: %v", err), zlog.WithTopicField(zlog.LogNameModule))
		return nil, err
	}

	t := client.beforeHttpStat(ctx, req)
	body, fields, err := client.httpDo(ctx, req, &opts)
	client.afterHttpStat(ctx, req.URL.Scheme, t)

	zlog.InfoLogger(ctx, zlog.LogNameServer, "http get request",
		zlog.WithTopicField(zlog.LogNameModule),
		zap.String("url", requestUrl),
		zap.Int("responseCode", body.HttpCode),
		zap.String("responseBody", string(body.Response)),
	)

	msg := "http request success"
	if err != nil {
		msg = err.Error()
	}

	zlog.InfoLogger(ctx, zlog.LogNameModule, msg, fields...)

	return &body, err
}

func (client *ApiClient) HttpPost(ctx *gin.Context, path string, opts HttpRequestOptions) (*ApiResult, error) {
	// http request
	urlData, err := opts.GetRequestBody()
	if err != nil {
		zlog.WarnLogger(ctx, zlog.LogNameServer, fmt.Sprintf("http client make data error: %v", err), zlog.WithTopicField(zlog.LogNameModule))
		return nil, err
	}

	u := fmt.Sprintf("%s%s", client.Domain, path)

	req, err := client.makeRequest(ctx, "POST", u, strings.NewReader(urlData), opts)
	if err != nil {
		zlog.WarnLogger(ctx, zlog.LogNameServer, fmt.Sprintf("http client makeRequest error: %v", err), zlog.WithTopicField(zlog.LogNameModule))
		return nil, err
	}

	t := client.beforeHttpStat(ctx, req)
	body, fields, err := client.httpDo(ctx, req, &opts)
	client.afterHttpStat(ctx, req.URL.Scheme, t)

	zlog.InfoLogger(ctx, zlog.LogNameServer, "http post request",
		zlog.WithTopicField(zlog.LogNameModule),
		zap.String("url", u),
		zap.String("params", urlData),
		zap.Int("responseCode", body.HttpCode),
		zap.String("responseBody", string(body.Response)),
	)

	msg := "http request success"
	if err != nil {
		msg = err.Error()
	}

	zlog.InfoLogger(ctx, zlog.LogNameModule, msg, fields...)

	return &body, err
}

type ApiResult struct {
	HttpCode int
	Response []byte
	Ctx      *gin.Context
}

func (client *ApiClient) httpDo(ctx *gin.Context, req *http.Request, opts *HttpRequestOptions) (res ApiResult, field []zlog.Field, err error) {
	start := time.Now()
	fields := []zlog.Field{
		zlog.WithTopicField(zlog.LogNameModule),
		zap.String("prot", "http"),
		zap.String("service", client.Service),
		zap.String("method", req.Method),
		zap.String("domain", client.Domain),
		zap.String("requestUri", req.URL.Path),
		zap.String("proxy", client.Proxy),
		zap.Duration("timeout", client.Timeout),
		zap.String("requestStartTime", utils.GetFormatRequestTime(start)),
	}

	client.clientInit.Do(func() {
		if client.HTTPClient == nil {
			timeout := 3 * time.Second
			if client.Timeout > 0 {
				timeout = client.Timeout
			}

			trans := client.GetTransPort()
			client.HTTPClient = &http.Client{
				Timeout:   timeout,
				Transport: trans,
			}
		}
	})

	var (
		resp         *http.Response
		dataBuffer   *bytes.Reader
		maxAttempts  int
		attemptCount int
		doErr        error
		shouldRetry  bool
	)

	attemptCount, maxAttempts = 0, client.Retry

	retryPolicy := opts.GetRetryPolicy()
	backOffPolicy := opts.GetBackOffPolicy()

	for {
		if req.GetBody != nil {
			bodyReadCloser, _ := req.GetBody()
			req.Body = bodyReadCloser
		} else if req.Body != nil {
			if dataBuffer == nil {
				data, err := ioutil.ReadAll(req.Body)
				_ = req.Body.Close()
				if err != nil {
					return res, fields, err
				}
				dataBuffer = bytes.NewReader(data)
				req.ContentLength = int64(dataBuffer.Len())
				req.Body = ioutil.NopCloser(dataBuffer)
			}
			_, _ = dataBuffer.Seek(0, io.SeekStart)
		}

		attemptCount++
		resp, doErr = client.HTTPClient.Do(req)
		if doErr != nil {
			f := []zlog.Field{
				zlog.WithTopicField(zlog.LogNameModule),
				zap.String("prot", "http"),
				zap.String("service", client.Service),
				zap.String("requestUri", req.URL.Path),
				zap.Duration("timeout", client.Timeout),
				zap.Int("attemptCount", attemptCount),
			}
			zlog.WarnLogger(ctx, zlog.LogNameModule, doErr.Error(), f...)
		}

		shouldRetry = retryPolicy(resp, doErr)
		if !shouldRetry {
			break
		}

		if attemptCount > maxAttempts {
			break
		}

		if doErr == nil {
			drainAndCloseBody(resp, 16384)
		}
		wait := backOffPolicy(attemptCount)
		select {
		case <-req.Context().Done():
			return res, fields, req.Context().Err()
		case <-time.After(wait):
		}
	}

	if resp != nil {
		res.HttpCode = resp.StatusCode
		res.Response, err = ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
	}

	err = doErr
	if err == nil && shouldRetry {
		err = fmt.Errorf("hit retry policy")
	}

	end := time.Now()
	if err != nil {
		err = fmt.Errorf("giving up after %d attempt(s): %w", attemptCount, err)
	}

	fields = append(fields,
		zap.String("retry", fmt.Sprintf("%d/%d", attemptCount-1, client.Retry)),
		zap.Int("httpCode", res.HttpCode),
		zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
		zap.Float64("cost", utils.GetRequestCost(start, end)),
		zap.Int("requestCode", client.requestCode(resp, err)),
	)

	return res, fields, err
}

func (client *ApiClient) requestCode(resp *http.Response, err error) int {
	if err != nil || resp == nil || resp.StatusCode >= 400 || resp.StatusCode == 0 {
		return -1
	}
	return 0
}

type timeTrace struct {
	dnsStartTime,
	dnsDoneTime,
	connectDoneTime,
	gotConnTime,
	gotFirstRespTime,
	tlsHandshakeStartTime,
	tlsHandshakeDoneTime,
	finishTime time.Time
}

func (client *ApiClient) beforeHttpStat(ctx *gin.Context, req *http.Request) *timeTrace {
	if client.HttpStat == false {
		return nil
	}

	var t = &timeTrace{}
	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { t.dnsStartTime = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { t.dnsDoneTime = time.Now() },
		ConnectStart: func(_, _ string) {
			if t.dnsDoneTime.IsZero() {
				t.dnsDoneTime = time.Now()
			}
		},
		ConnectDone: func(net, addr string, err error) {
			t.connectDoneTime = time.Now()
		},
		GotConn:              func(_ httptrace.GotConnInfo) { t.gotConnTime = time.Now() },
		GotFirstResponseByte: func() { t.gotFirstRespTime = time.Now() },
		TLSHandshakeStart:    func() { t.tlsHandshakeStartTime = time.Now() },
		TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) { t.tlsHandshakeDoneTime = time.Now() },
	}
	*req = *req.WithContext(httptrace.WithClientTrace(context.Background(), trace))
	return t
}

func (client *ApiClient) afterHttpStat(ctx *gin.Context, scheme string, t *timeTrace) {
	if client.HttpStat == false {
		return
	}

	t.finishTime = time.Now() // after read body

	if t.dnsStartTime.IsZero() {
		t.dnsStartTime = t.dnsDoneTime
	}

	cost := func(d time.Duration) float64 {
		if d < 0 {
			return -1
		}
		return float64(d.Nanoseconds()/1e4) / 100.0
	}

	switch scheme {
	case "https":
		f := []zlog.Field{
			zlog.WithTopicField(zlog.LogNameModule),
			zap.Float64("dnsLookupCost", cost(t.dnsDoneTime.Sub(t.dnsStartTime))),                       // dns lookup
			zap.Float64("tcpConnectCost", cost(t.connectDoneTime.Sub(t.dnsDoneTime))),                   // tcp connection
			zap.Float64("tlsHandshakeCost", cost(t.tlsHandshakeStartTime.Sub(t.tlsHandshakeStartTime))), // tls handshake
			zap.Float64("serverProcessCost", cost(t.gotFirstRespTime.Sub(t.gotConnTime))),               // server processing
			zap.Float64("contentTransferCost", cost(t.finishTime.Sub(t.gotFirstRespTime))),              // content transfer
			zap.Float64("totalCost", cost(t.finishTime.Sub(t.dnsStartTime))),                            // total cost
		}
		zlog.InfoLogger(ctx, zlog.LogNameModule, "time trace", f...)
	case "http":
		f := []zlog.Field{
			zlog.WithTopicField(zlog.LogNameModule),
			zap.Float64("dnsLookupCost", cost(t.dnsDoneTime.Sub(t.dnsStartTime))),          // dns lookup
			zap.Float64("tcpConnectCost", cost(t.gotConnTime.Sub(t.dnsDoneTime))),          // tcp connection
			zap.Float64("serverProcessCost", cost(t.gotFirstRespTime.Sub(t.gotConnTime))),  // server processing
			zap.Float64("contentTransferCost", cost(t.finishTime.Sub(t.gotFirstRespTime))), // content transfer
			zap.Float64("totalCost", cost(t.finishTime.Sub(t.dnsStartTime))),               // total cost
		}
		zlog.InfoLogger(ctx, zlog.LogNameModule, "time trace", f...)
	}
}

func drainAndCloseBody(resp *http.Response, maxBytes int64) {
	if resp != nil {
		_, _ = io.CopyN(ioutil.Discard, resp.Body, maxBytes)
		_ = resp.Body.Close()
	}
}

// RetryPolicy retry policy.
type RetryPolicy func(resp *http.Response, err error) bool

var defaultRetryPolicy RetryPolicy = func(resp *http.Response, err error) bool {
	return err != nil || resp == nil || resp.StatusCode >= 500 || resp.StatusCode == 0
}

// BackOffPolicy set policy of retry wait time.
type BackOffPolicy func(attemptCount int) time.Duration

var defaultBackOffPolicy = func(attemptNum int) time.Duration { // retry immediately
	return 0
}
