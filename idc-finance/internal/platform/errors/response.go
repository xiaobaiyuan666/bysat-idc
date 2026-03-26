package errors

type APIResponse[T any] struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"requestId,omitempty"`
	Data      T      `json:"data"`
}

func Ok[T any](data T, requestID string) APIResponse[T] {
	return APIResponse[T]{
		Code:      "OK",
		Message:   "success",
		RequestID: requestID,
		Data:      data,
	}
}

func BadRequest(message, requestID string) APIResponse[any] {
	return APIResponse[any]{
		Code:      "BAD_REQUEST",
		Message:   message,
		RequestID: requestID,
		Data:      nil,
	}
}
