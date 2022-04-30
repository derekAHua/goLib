package env

import (
	"github.com/derekAHua/goLib/consts"
	"github.com/derekAHua/goLib/utils"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"sync"
)

var (
	LocalIP string
	runMode string

	runEnv int

	rootPath string

	appName     string
	onceAppName sync.Once
)

func SetAppName(name string) {
	onceAppName.Do(func() {
		appName = name
	})
}

func GetAppName() string {
	return appName
}

func init() {
	LocalIP = utils.GetLocalIp()

	gin.SetMode(gin.DebugMode) // 运行环境

	initDBSecret()
}

func SetRootPath(r string) {
	rootPath = r
}

func GetRootPath() string {
	if rootPath != "" {
		return rootPath
	} else {
		return consts.DefaultRootPath
	}
}

// GetConfDirPath 返回配置文件目录绝对地址
func GetConfDirPath() string {
	return filepath.Join(GetRootPath(), "conf")
}

// GetLogDirPath 返回log目录的绝对地址
func GetLogDirPath() string {
	return filepath.Join(GetRootPath(), "log")
}

func GetRunEnv() int {
	return runEnv
}

func GetLocalIP() string {
	return LocalIP
}

func GetRunMode() string {
	return runMode
}
