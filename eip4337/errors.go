package eip4337

type RPCError struct {
	code    int
	message string
	data    interface{}
}

const (
	ErrorCodeInvalidFields                  = -32602
	ErrorCodeSimulateValidation             = -32500
	ErrorCodeSimulatePaymasterValidation    = -32501
	ErrorCodeOpcodeValidation               = -32502
	ErrorCodeExpiresShortly                 = -32503
	ErrorCodeReputation                     = -32504
	ErrorCodeInsufficientStake              = -32505
	ErrorCodeUnsupportedSignatureAggregator = -32506
	ErrorCodeInvalidSignature               = -32507

	ErrorCodeUserOperationReverted = -32521
	ErrorCodeUnknown               = -39999 // FIXME: find a better number
)

func (err *RPCError) ErrorCode() int {
	return err.code
}

func (err *RPCError) Error() string {
	return err.message
}

func (err *RPCError) ErrorData() interface{} {
	return err.data
}

func NewRPCError(code int, message string, data interface{}) *RPCError {
	return &RPCError{
		code:    code,
		message: message,
		data:    data,
	}
}
