package consts

import "github.com/derekAHua/goLib/base"

// @Author: Derek
// @Description: Error Code.
// @Date: 2022/4/30 11:45
// @Version 1.0

var (
	ParamUnValid = base.NewError(4000, "参数错误！")
)

// JWT Error. [1000-1100)
var (
	TokenExpired     = base.NewError(1000, "Token is expired.")
	TokenNotValidYet = base.NewError(1001, "Token not active yet.")
	TokenMalformed   = base.NewError(1002, "That's not even a token.")
	TokenInvalid     = base.NewError(1003, "Couldn't handle this token:")
)
