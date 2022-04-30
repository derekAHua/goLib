package main

import (
	"fmt"
	"github.com/derekAHua/goLib/base"
)

// @Author: Derek
// @Description:
// @Date: 2022/4/30 10:27
// @Version 1.0

func main() {
	err := base.NewError(-1, "测试")
	err2 := base.NewError(-1, "测试AAA")

	fmt.Println(err.Wrap(err2))
}
