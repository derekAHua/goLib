package goLib

import (
	"github.com/gin-gonic/gin"
	"testing"
)

// @Author: Derek
// @Description:
// @Date: 2022/8/14 20:24
// @Version 1.0

func TestBootstraps(t *testing.T) {
	Bootstraps(gin.New())
}
