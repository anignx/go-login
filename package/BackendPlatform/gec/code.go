package gec

const (
	ErrGlobalToastCode int = 603 // 客户端全局toast错误码
)

type GlobalToastErr interface {
	error
	code() int
	message() string
}

type globalToastError struct {
	globalToastCode    int
	globalToastMessage string
}

type Global struct {
	code    int
	message string
}

var (
	ErrSuccess    = Error(int(OK), "操作成功")
	ErrBadRequest = Error(400, "无效请求")
)
