package constant

import "fmt"

type ErrorCode struct {
	Code int64
	Msg  string
}

func makeErrorCode(code int64, msg string) *ErrorCode {
	return &ErrorCode{
		Code: code,
		Msg:  msg,
	}
}

func BuildErrCode(detail interface{}, errCode *ErrorCode) *ErrorCode {
	return makeErrorCode(errCode.Code, fmt.Sprintf(errCode.Msg, detail))
}

var (
	Success      = makeErrorCode(0, "成功")
	ParamErr     = makeErrorCode(1, "参数非法:%v")
	NotFoundErr  = makeErrorCode(100, "找不到需要降权的电影或者标签:%v")
	ReadRepoErr  = makeErrorCode(101, "查库错误:%v")
	WriteRepoErr = makeErrorCode(102, "写库错误:%v")
	SysErr       = makeErrorCode(999, "系统未知错误:%v")
)
