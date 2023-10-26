package modelUtils

type JsonResponseStruct struct {
	ResponseCode int
	Detail       JsonResponseStructDetail
}

type JsonResponseStructDetail struct {
	Timestamp string      `json:"timestamp"`
	Success   bool        `json:"success"`
	ErrorCode string      `json:"error_code"`
	Data      interface{} `json:"data"`
	Error     interface{} `json:"error"`
}

type ApiResponseStruct struct {
	ResponseCode int
	ErrorCode    string
	Data         interface{}
	Error        interface{}
}

type ApiResponseConfigStruct struct {
	ErrorCode string
	Error     interface{}
}
