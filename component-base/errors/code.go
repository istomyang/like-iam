package errors

import (
	"fmt"
	"net/http"
	"sync"
)

// Coder save a error detail information.
type Coder interface {
	// Code return project-related error code, like 1011001
	Code() int

	// HTTPCode associates this error with http status code.
	HTTPCode() int

	// 	Message return a user-friendly text about this error.
	Message() string

	// Reference return a document-link to guide.
	Reference() string
}

type defaultCoder struct {
	code      int
	httpCode  int
	message   string
	reference string
}

// Code implements Coder's Coder.Code.
func (c defaultCoder) Code() int {
	return c.code
}

// HTTPCode implements Coder's Coder.HTTPCode.
func (c defaultCoder) HTTPCode() int {
	return c.httpCode
}

// Message implements Coder's Coder.Message.
func (c defaultCoder) Message() string {
	return c.message
}

// Reference implements Coder's Coder.Reference.
func (c defaultCoder) Reference() string {
	return c.reference
}

var mu sync.Mutex

// codes saves all register Coder in memory.
var codes map[int]Coder

var unknownCode Coder = &defaultCoder{
	code:      0,
	httpCode:  http.StatusInternalServerError,
	message:   "An internal server error occurs.",
	reference: "",
}

// Register registers a coder with override strategy.
func Register(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved by this package's `unknownCode`.")
	}

	mu.Lock()
	defer mu.Unlock()

	codes[coder.Code()] = coder
}

// MustRegister registers a coder with override-panic strategy.
func MustRegister(coder Coder) {
	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %v already exists.", coder.Code()))
	}
	Register(coder)
}

// AsCode parse err created by WithCode into Coder registered by Register.
func AsCode(err error) Coder {
	if err == nil {
		return nil
	}

	if c, ok := err.(*withCode); ok {
		if coder, ok := codes[c.code]; ok {
			return coder
		}
	}
	return unknownCode
}

// IsCode checks whether err is associated with coder deeply.
func IsCode(err error, coder Coder) bool {
	if err == nil {
		return false
	}
	if c, ok := err.(*withCode); ok {

		if c.code != coder.Code() {
			IsCode(c.Unwrap(), coder)
		} else {
			return true
		}
	}
	return false
}

// Registered check whether special code has registered.
func Registered(code int) bool {
	_, has := codes[code]
	return has
}

func init() {
	codes = make(map[int]Coder)
	codes[unknownCode.Code()] = unknownCode
}
