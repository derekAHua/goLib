package utils

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Int64sContain(a []int64, x int64) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}

func GetFunctionName(i interface{}, seps ...rune) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}

// RandNum 获取随机数
//  不传参：0-100
//  传1个参数：0-指定参数
//  传2个参数：第1个参数-第2个参数
func RandNum(num ...int) int {
	var start, end int
	if len(num) == 0 {
		start = 0
		end = 100
	} else if len(num) == 1 {
		start = 0
		end = num[0]
	} else {
		start = num[0]
		end = num[1]
	}

	rRandNumUtils := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rRandNumUtils.Intn(end-start+1) + start
}

func GetHandler(ctx *gin.Context) (handler string) {
	if ctx != nil {
		handler = ctx.HandlerName()
	}
	return handler
}

func JoinArgs(showByte int, args ...interface{}) string {
	cnt := len(args)
	f := "%v"
	for cnt > 1 {
		f += " %v"
	}

	argVal := fmt.Sprintf(f, args...)

	l := len(argVal)
	if l > showByte {
		l = showByte
		argVal = argVal[:l] + " ..."
	}
	return argVal
}
