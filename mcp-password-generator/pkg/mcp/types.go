package mcp

import (
	"encoding/json"
)

const (
	JSONRPCVersion = "2.0"
)

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
	ID      interface{}     `json:"id"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

type ErrorCode int

const (
	ParseError     ErrorCode = -32700
	InvalidRequest ErrorCode = -32600
	MethodNotFound ErrorCode = -32601
	InvalidParams  ErrorCode = -32602
	InternalError  ErrorCode = -32603
)

const (
	ToolNotFound      ErrorCode = -32000
	ToolExecutionFail ErrorCode = -32001
	ConsentRequired   ErrorCode = -32002
)

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
	Annotations interface{} `json:"annotations,omitempty"`
}

type ToolsListParams struct{}

type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

type ToolsCallParams struct {
	Name    string          `json:"name"`
	Input   json.RawMessage `json:"input"`
	Consent *bool           `json:"consent,omitempty"`
}

type TextContent struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Title string `json:"title,omitempty"`
}

type ToolResult struct {
	Content []interface{} `json:"content"`
}

func NewResponse(id interface{}, result interface{}, err error) (*Response, error) {
	var rawResult json.RawMessage
	if result != nil {
		var marshalErr error
		rawResult, marshalErr = json.Marshal(result)
		if marshalErr != nil {
			return nil, marshalErr
		}
	}

	var jsonRPCError *Error
	if err != nil {
		jsonRPCError = &Error{
			Code:    int(InternalError),
			Message: err.Error(),
		}
	}

	return &Response{
		JSONRPC: JSONRPCVersion,
		Result:  rawResult,
		Error:   jsonRPCError,
		ID:      id,
	}, nil
}

func NewError(code ErrorCode, message string, data interface{}) *Error {
	return &Error{
		Code:    int(code),
		Message: message,
		Data:    data,
	}
}
