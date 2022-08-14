package middleware

import (
	"github.com/derekAHua/goLib/function"
	"github.com/gin-gonic/gin"
)

func Recover(ctx *gin.Context) {
	defer function.CatchPanic(ctx)
	ctx.Next()
}
