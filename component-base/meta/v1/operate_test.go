package v1

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

func TestUrlValues(t *testing.T) {
	s := &DeleteOperateMeta{
		OperateMeta: OperateMeta{ApiVersion: "v1"},
		Unscoped:    false,
	}

	out1, _ := json.Marshal(s)
	fmt.Println(string(out1))

	var v = url.Values{}
	v.Set("api-version", "v1")
	v.Set("unscoped", "false")
	p1 := v.Encode()
	//fmt.Println(p1)
	v2, _ := url.ParseQuery(p1)
	fmt.Println(v2.Encode())
}
