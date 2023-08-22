package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// formatInfo is a model for withCode.Format which transfer by WithCode and others external error object.
// It describes the known data in output string.
type formatInfo struct {
	// caller is a stack caller with the latest file path and line-number.
	// i represent depth of codes.
	caller func(i int) string

	// error always comes from error.Error.
	error string

	// message has the same value with formatInfo.error generally, if came across WithCode, use Coder.Message.
	message string
}

func buildFormatInfo(err error) *formatInfo {
	var ret *formatInfo
	switch err := err.(type) {
	case *withCode:
		coder, ok := codes[err.code]
		if !ok {
			coder = unknownCode
		}

		extMsg := coder.Message()
		if extMsg == "" {
			extMsg = err.Error()
		}

		ret = &formatInfo{
			error:   err.Error(),
			caller:  callerStringFactory(s2f(err.stack)),
			message: messageString(err.code, extMsg),
		}
	case *withMessage:
		ret = &formatInfo{
			error:   err.Error(),
			caller:  nil,
			message: messageString(unknownCode.Code(), err.msg),
		}
	case *withStack:
		ret = &formatInfo{
			error:   err.Error(),
			caller:  callerStringFactory(s2f(err.stack)),
			message: messageString(unknownCode.Code(), err.Error()),
		}
	case *fundamental:
		ret = &formatInfo{
			error:   err.Error(),
			caller:  callerStringFactory(s2f(err.stack)),
			message: messageString(unknownCode.Code(), err.msg),
		}
	default:
		ret = &formatInfo{
			error:   err.Error(),
			caller:  nil,
			message: messageString(unknownCode.Code(), err.Error()),
		}
	}
	return ret
}

func s2f(s *stack) Frame {
	return Frame((*s)[0])
}

// callerStringFactory provides string for formatInfo.caller.
// E.g. #0 [/home/yangyang/workspace/golang/src/github.com/istomyang/like-iam/main.go:12 (main.main)]
func callerStringFactory(f Frame) func(int) string {
	return func(i int) string {
		return fmt.Sprintf("#%d [%s:%d (%s)]", i, f.file(), f.line(), f.name())
	}
}

// messageString provides string for formatInfo.message.
// E.g. (#100102) Internal Server Error
func messageString(code int, message string) string {
	return fmt.Sprintf("(#%d) %s", code, message)
}

// Format implements fmt.Formatter. https://golang.org/pkg/fmt/#hdr-Printing
//
// Verbs:
//
//	%s  - Returns the user-safe error string mapped to the error code or
//	  ┊   the error message if none is specified.
//	%v      Alias for %s
//
// Flags:
//
//	#      JSON formatted output, useful for logging
//	-      Output caller details, useful for troubleshooting
//	+      Output full error stack details, useful for debugging
//
// Examples:
//
//	%s:    error for internal read B
//	%v:    error for internal read B
//	%-v:   error for internal read B - #0 [/home/yangyang/workspace/golang/src/github.com/istomyang/like-iam/main.go:12 (main.main)] (#100102) Internal Server Error
//	%+v:   error for internal read B - #0 [/home/yangyang/workspace/golang/src/github.com/istomyang/like-iam/main.go:12 (main.main)] (#100102) Internal Server Error; error for internal read A - #1 [/home/yangyang/workspace/golang/src/github.com/istomyang/like-iam/main.go:35 (main.newErrorB)] (#100104) Validation failed
//	%#v:   [{"error":"error for internal read B"}]
//	%#-v:  [{"caller":"#0 /home/yangyang/workspace/golang/src/github.com/istomyang/like-iam/main.go:12 (main.main)","error":"error for internal read B","message":"(#100102) Internal Server Error"}]
//	%#+v:  [{"caller":"#0 /home/yangyang/workspace/golang/src/github.com/istomyang/like-iam/main.go:12 (main.main)","error":"error for internal read B","message":"(#100102) Internal Server Error"},{"caller":"#1 /home/yangyang/workspace/golang/src/github.com/istomyang/like-iam/main.go:35 (main.newErrorB)","error":"error for internal read A","message":"(#100104) Validation failed"}]
func (c *withCode) Format(s fmt.State, verb rune) {
	var flagJson = s.Flag('#')
	var flagDetail = s.Flag('-')
	var flagStack = s.Flag('+')

	switch verb {
	case 'v':
		// %v
		if !flagDetail && !flagStack && !flagJson {
			infos := []*formatInfo{buildFormatInfo(c)}
			_, _ = fmt.Fprint(s, outputTextString(infos, true))
			return
		}

		//// %-v
		//if flagDetail && !flagStack && !flagJson {
		//	infos := []*formatInfo{buildFormatInfo(c)}
		//	_, _ = fmt.Fprint(s, outputTextString(infos, false))
		//	return
		//}
		//
		//// %+v
		//if !flagDetail && flagStack && !flagJson {
		//	errs := list(c)
		//	infos := make([]*formatInfo, len(errs))
		//	for _, err := range errs {
		//		infos = append(infos, buildFormatInfo(err))
		//	}
		//	_, _ = fmt.Fprint(s, outputTextString(infos, false))
		//	return
		//}
		//
		//// %#v
		//if !flagDetail && !flagStack && flagJson {
		//	infos := []*formatInfo{buildFormatInfo(c)}
		//	_, _ = fmt.Fprint(s, outputJsonString(infos, true))
		//	return
		//}
		//
		//// %#-v
		//if flagDetail && !flagStack && flagJson {
		//	infos := []*formatInfo{buildFormatInfo(c)}
		//	_, _ = fmt.Fprint(s, outputJsonString(infos, false))
		//	return
		//}
		//
		//// %#+v
		//if !flagDetail && flagStack && flagJson {
		//	errs := list(c)
		//	infos := make([]*formatInfo, len(errs))
		//	for _, err := range errs {
		//		infos = append(infos, buildFormatInfo(err))
		//	}
		//	_, _ = fmt.Fprint(s, outputJsonString(infos, false))
		//	return
		//}

		var infos []*formatInfo

		if flagDetail {
			infos = []*formatInfo{buildFormatInfo(c)}
		}

		if flagStack {
			errs := list(c)
			infos = make([]*formatInfo, len(errs))
			for _, err := range errs {
				infos = append(infos, buildFormatInfo(err))
			}
		}

		var brief bool
		if !flagDetail && !flagStack {
			brief = true
		}

		var o string
		if flagJson {
			o = outputJsonString(infos, brief)
		} else {
			o = outputTextString(infos, brief)
		}
		_, _ = fmt.Fprint(s, o)

	default:
		//	%s
		errs := []*formatInfo{buildFormatInfo(c)}
		_, _ = fmt.Fprint(s, outputTextString(errs, true))
	}
}

// collectFormatInfos return formatInfo collections according to flag.
func collectFormatInfos(e error, stackFlag bool) []*formatInfo {
	var es []error
	if stackFlag {
		es = []error{e}
	} else {
		es = list(e)
	}
	var ret = make([]*formatInfo, len(es))
	for _, err := range es {
		ret = append(ret, buildFormatInfo(err))
	}
	return ret
}

// outputTextString output text format.
func outputTextString(infos []*formatInfo, brief bool) string {
	var b strings.Builder
	const sep = ';'

	for i, info := range infos {
		if brief {
			b.WriteString(info.error)
		} else {
			b.WriteString(info.error)
			b.WriteString(" - ")
			b.WriteString(info.caller(i))
			b.WriteRune(' ')
			b.WriteString(info.message)
			if i != len(infos)-1 {
				b.WriteRune(sep)
			}
		}
	}

	return b.String()
}

// outputJsonString outputs two kinds of format.
// if brief is false, includes caller, error and message.
// if brief is true, only includes error.
func outputJsonString(infos []*formatInfo, brief bool) string {
	var jsonMap = make([]map[string]string, len(infos))

	for i, info := range infos {
		var m map[string]string
		if brief {
			m = map[string]string{
				"error": info.error,
			}
		} else {
			m = map[string]string{
				"caller":  info.caller(i),
				"error":   info.error,
				"message": info.message,
			}
		}

		jsonMap = append(jsonMap, m)
	}

	b, _ := json.Marshal(jsonMap)
	return string(b)
}

// list collects all codes with Unwrap.
func list(e error) []error {
	var ret []error
	if e != nil {
		ret = append(ret, e)
		if ue, ok := e.(interface {
			Unwrap() error
		}); ok {
			ret = append(ret, list(ue.Unwrap())...)
		}
	}
	return ret
}
