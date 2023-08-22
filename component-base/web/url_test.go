package web

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type OperateMeta struct {
	ApiVersion string `json:"api-version,omitempty"`
}

type DeleteOperateMeta struct {
	OperateMeta `json:",inline"`

	// Unscoped replace soft delete operation with hard delete operation.
	// Default gorm db use DeleteAt field to mark this entry need be deleted.
	// +optional
	Unscoped bool `json:"unscoped"`
}

func TestParse(t *testing.T) {
	s := &DeleteOperateMeta{
		OperateMeta: OperateMeta{ApiVersion: "v1"},
		Unscoped:    false,
	}
	var v = url.Values{}

	out1, _ := json.Marshal(s)
	str := string(out1)
	var b strings.Builder
	for _, l := range str {
		if l == '{' || l == '}' || l == '"' {
			continue
		}
		b.WriteRune(l)
	}
	sections := strings.Split(b.String(), ",")
	for _, section := range sections {
		parts := strings.SplitN(section, ":", 2)
		v.Set(parts[0], parts[1])
	}

	fmt.Println(v.Encode())
}

func TestS(t *testing.T) {
	var s = []string{"a"}
	tt := reflect.TypeOf(s)
	fmt.Println(tt.Kind())
	fmt.Println(tt.String())
}

func TestS2(t *testing.T) {
	var a any = "ssss"
	fmt.Println(a.(string))

	var _b int64 = 12323
	var b any = _b
	fmt.Println(strconv.Itoa(int(b.(int64))))
}

func TestUrlValues(t *testing.T) {
	s := &DeleteOperateMeta{
		OperateMeta: OperateMeta{ApiVersion: "v1"},
		Unscoped:    false,
	}
	uv, _ := UrlValues(s)
	fmt.Println(uv.Encode())
}

func TestUrlValuesWithTpl(t *testing.T) {
	s := &DeleteOperateMeta{
		OperateMeta: OperateMeta{ApiVersion: "v1"},
		Unscoped:    false,
	}
	_, tpl := UrlValues(s)
	uv, _ := UrlValuesWithTpl(tpl, "v2", true)
	fmt.Println(uv.Encode())
}
