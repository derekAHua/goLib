package goLib

import (
	"github.com/derekAHua/goLib/env"
	"github.com/derekAHua/goLib/middleware"
	"github.com/derekAHua/goLib/utils"
	"github.com/gin-gonic/gin"
)

func Bootstraps(router *gin.Engine) {
	utils.InitValidator()

	gin.SetMode(env.GetRunMode())
	router.Use(middleware.AccessLog())
	router.Use(gin.Recovery())
}
