package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component-base/base"
	"istomyang.github.com/like-iam/component-base/errors"
	"net/http"
)

// ErrorResponse defines body when send error info to client.
type ErrorResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

var _ fmt.Stringer = &ErrorResponse{}

func (er *ErrorResponse) String() string {
	var tpl = `ErrorResponse: Code(%d), Message(%s), Reference(%s);`
	return fmt.Sprintf(tpl, er.Code, er.Message, er.Reference)
}

var _ fmt.GoStringer = &ErrorResponse{}

func (er *ErrorResponse) GoString() string {
	s, _ := json.Marshal(er)
	return string(s)
}

// WriteResponse write response data or error to gin.
func WriteResponse(c *gin.Context, err error, data any) {
	if err != nil {
		coder := errors.AsCode(err)
		c.JSON(coder.HTTPCode(), ErrorResponse{
			Code:      coder.Code(),
			Message:   coder.Message(),
			Reference: coder.Reference(),
		})
		return
	}
	c.JSON(http.StatusOK, data)
}

// IsErrResponse check http.Response 's body is a ErrorResponse struct.
func IsErrResponse(body []byte, decoder base.Decoder) (res *ErrorResponse) {
	if err := decoder.Decode(body, res); err != nil {
		return nil
	}
	return
}
