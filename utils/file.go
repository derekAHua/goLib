package utils

import "os"

// @Author: Derek
// @Description:
// @Date: 2022/8/14 09:05
// @Version 1.0

func WriteFile(filename string) (f *os.File, err error) {
	if CheckFileIsExist(filename) {
		return os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666)
	}

	return os.Create(filename)
}

// CheckFileIsExist 检查文件是否存在
func CheckFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
