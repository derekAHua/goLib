package utils

import (
	"os"
	"path"
	"runtime"
)

// @Author: Derek
// @Description:
// @Date: 2022/8/14 11:36
// @Version 1.0

func MakeDirIfNo(logDir string) (err error) {
	if _, err = os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0777)
		if err != nil {
			return
		}
	}

	return
}

// GetSourcePath returns the directory containing the source code that is calling this function.
func GetSourcePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
