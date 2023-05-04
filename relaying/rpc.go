package relaying

import "fmt"

const (
	ErrorCodeInternal       = -99999
	ErrorCodeSystemBusy     = -50000
	ErrorCodeInvalidRequest = -40000
	ErrorCodeInvalidRelayTx = -40001
)

type RpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (err *RpcError) Error() string {
	return err.Message
}

func (err *RpcError) ErrorCode() int {
	return err.Code
}

func (err *RpcError) ErrorData() interface{} {
	return err.Data
}

func NewRpcError(code int, data interface{}, msg string) *RpcError {
	return &RpcError{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func NewRpcErrorf(code int, data interface{}, msgFmt string, msgArgs ...interface{}) *RpcError {
	return NewRpcError(code, data, fmt.Sprintf(msgFmt, msgArgs...))
}
