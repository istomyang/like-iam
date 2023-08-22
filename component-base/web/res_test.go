package web

import (
	"fmt"
	"testing"
)

func TestErrorResponse_String(t *testing.T) {
	var er = ErrorResponse{
		Code:      1010100,
		Message:   "This is a message",
		Reference: "You can refer: https://iam.istomyang.github.com/reference.md",
	}

	fmt.Println(er)
	fmt.Println(er.String())
}
