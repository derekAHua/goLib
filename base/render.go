package base

import (
	"encoding/json"
	"fmt"
	"github.com/derekAHua/goLib/errors"
	"github.com/derekAHua/goLib/zlog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type DefaultRender struct {
	ErrNo  int         `json:"errNo"`
	ErrMsg string      `json:"errMsg"`
	Data   interface{} `json:"data"`
}

func RenderJson(ctx *gin.Context, code int, msg string, data interface{}) {
	renderJson := DefaultRender{code, msg, data}
	ctx.JSON(http.StatusOK, renderJson)
	return
}

func RenderJsonSuc(ctx *gin.Context, data interface{}) {
	RenderJson(ctx, 0, "success", data)
	return
}

func RenderJsonFail(ctx *gin.Context, err error) {
	errNo := -1
	if v, ok := err.(errors.Err); ok {
		errNo = v.ErrNo
	}
	RenderJson(ctx, errNo, err.Error(), nil)
	StackLogger(ctx, err)
	return
}

func RenderJsonAbort(ctx *gin.Context, err errors.Err) {
	ctx.AbortWithStatusJSON(http.StatusOK, DefaultRender{ErrNo: err.ErrNo, ErrMsg: err.ErrMsg})
	return
}

// StackLogger 打印错误栈
func StackLogger(ctx *gin.Context, err error) {
	if !strings.Contains(fmt.Sprintf("%+v", err), "\n") {
		return
	}

	var info []byte
	if ctx != nil {
		info, _ = json.Marshal(map[string]interface{}{"time": time.Now().Format("2006-01-02 15:04:05"), "level": "error", "module": "errorStack", "requestId": zlog.GetRequestId(ctx)})
	} else {
		info, _ = json.Marshal(map[string]interface{}{"time": time.Now().Format("2006-01-02 15:04:05"), "level": "error", "module": "errorStack"})
	}

	fmt.Printf("%s\n-------------------stack-start-------------------\n%+v\n-------------------stack-end-------------------\n", string(info), err)
}
