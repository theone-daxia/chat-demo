package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/theone-daxia/chat-demo/config"
	"github.com/theone-daxia/chat-demo/serializer"
)

func ErrorResponse(err error) serializer.Response {
	if es, ok := err.(validator.ValidationErrors); ok {
		for _, e := range es {
			field := config.T(fmt.Sprintf("Field.%s", e.Field()))
			tag := config.T(fmt.Sprintf("Tag.Valid.%s", e.Tag()))
			return serializer.Response{
				Status: 400,
				Msg:    fmt.Sprintf("%s%s", field, tag),
				Error:  fmt.Sprint(err),
			}
		}
	}

	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.Response{
			Status: 400,
			Msg:    "JSON类型不匹配",
			Error:  fmt.Sprint(err),
		}
	}

	return serializer.Response{
		Status: 400,
		Msg:    "参数错误",
		Error:  fmt.Sprint(err),
	}
}
