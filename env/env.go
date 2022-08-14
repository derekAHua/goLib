package env

import (
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

func init() {
	LocalIP = utils.GetLocalIp()

	gin.SetMode(gin.DebugMode) // 运行环境
}

func GetAppName() string {
	return appName
}

func SetAppName(name string) {
	onceAppName.Do(func() {
		appName = name
	})
}

func GetRootPath() string {
	if rootPath != "" {
		return rootPath
	} else {
		return DefaultRootPath
	}
}

func SetRootPath(r string) {
	rootPath = r
}

// GetConfDirPath 返回配置文件目录绝对地址
func GetConfDirPath() string {
	return filepath.Join(GetRootPath(), ConfDir)
}

// GetLogDirPath 返回log目录的绝对地址
func GetLogDirPath() string {
	return filepath.Join(GetRootPath(), LogDir)
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
