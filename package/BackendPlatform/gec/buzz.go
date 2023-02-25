package gec

import (
	"encoding/json"
	"reflect"

	"github.com/gin-gonic/gin"
)

var (
	emptySuccessRsp = Success(nil)
	jsRawType       = reflect.TypeOf(json.RawMessage{})
)

type BuzzResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Debug   string      `json:"debug,omitempty"`
	Data    interface{} `json:"data"`
}

// err 类型是 ecode.Codes 时：
// 		Response.Code = ecode.Codes.Code(),
// 		Response.Message = ecode.Codes.Message(),
// err 类型是 GlobalToastError 时：
// 		Response.Code = 603,
// 		Response.Message = GlobalToastError.Message(),
// err == nil 时：
// 		Response.Code = 0
// 其他情况：
//		Response.Code = 500
// ⚠️ data 允许的类型为：struct, map, json.RawMessage
func NewGmResponse(data interface{}, err error, msgs ...string) *BuzzResponse {
	var (
		m     BuzzResponse
		valid bool
		v     reflect.Value
	)

	if data == nil {
		goto DONE
	}

	v = reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			goto DONE
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct, reflect.Map:
		valid = true
	default:
		switch v.Type() {
		case jsRawType:
			valid = true
		}
	}

DONE:
	if !valid {
		data = json.RawMessage("{}")
	}
	m.SetData(data).WithError(err, msgs...)
	return &m
}

func Success(data interface{}) *BuzzResponse {
	return NewGmResponse(data, ErrSuccess)
}

func EmptyBody(err error) *BuzzResponse {
	return NewGmResponse(nil, err)
}

func EmptySuccess() *BuzzResponse {
	return emptySuccessRsp
}

func (r *BuzzResponse) SetData(data interface{}) *BuzzResponse {
	r.Data = data
	return r
}

func (r *BuzzResponse) WithError(err error, msgs ...string) *BuzzResponse {
	// 检查Error类型
	// GlobalToastErr
	e, ok := err.(GlobalToastErr)
	if ok {
		r.Code = e.code()
		r.Message = e.message()
		return r
	}
	// ecode.Codes
	code := Cause(err)
	r.Code = code.Code()
	var msg string
	if len(msgs) > 0 {
		msg = msgs[0]
	} else {
		msg = code.Message()
	}
	r.Message = msg
	return r
}

func (r *BuzzResponse) WithDebug(debug string) *BuzzResponse {
	r.Debug = debug
	return r
}

func (r *BuzzResponse) JSON(c *gin.Context) {
	c.JSON(200, r)
}

// err 类型是 ecode.Codes 时：
// 		Response.Code = ecode.Codes.Code(),
// 		Response.Message = ecode.Codes.Message(),
// err == nil 时：
// 		Response.Code = 0
// 其他情况：
//		Response.Code = 500
// ⚠️ data 允许的类型为：struct, map, json.RawMessage
func JSON(c *gin.Context, data interface{}, err error, msgs ...string) {
	NewGmResponse(data, err, msgs...).JSON(c)
}
