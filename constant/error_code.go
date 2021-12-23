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
	RetSuccess      = makeErrorCode(0, "成功")
	RetParamsErr    = makeErrorCode(1, "参数非法:%v")
	RetReadRepoErr  = makeErrorCode(101, "查库错误:%v")
	RetWriteRepoErr = makeErrorCode(102, "写库错误:%v")
	RetSysErr       = makeErrorCode(999, "系统未知错误:%v")
)
