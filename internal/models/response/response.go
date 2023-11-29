package response

import "strings"

type Response[T any] struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func Error[T any](msg ...string) *Response[T] {
	return &Response[T]{
		Result:  false,
		Message: strings.Join(msg, ","),
		Data:    *new(T),
	}
}

func Ok[T any](data T) *Response[T] {
	return &Response[T]{
		Result:  true,
		Message: "",
		Data:    data,
	}
}
