package errors

// @Author: Derek
// @Description: Err
// @Date: 2022/4/30 10:42
// @Version 1.0

type Err struct {
	ErrNo  int
	ErrMsg string
}

func (e Err) Error() string {
	return e.ErrMsg
}
