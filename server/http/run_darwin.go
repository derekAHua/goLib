package http

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func Start(engine *gin.Engine, conf ServerConfig) error {

	appServer := endless.NewServer(conf.Address, engine)

	// 超时时间 (如果设置太小，可能导致接口响应时间超过该值，进而导致504)
	if conf.ReadTimeout > 0 {
		appServer.ReadTimeout = conf.ReadTimeout
	}

	if conf.WriteTimeout > 0 {
		appServer.WriteTimeout = conf.WriteTimeout
	}

	// 监听http端口
	if err := appServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
