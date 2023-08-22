package client

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"testing"
)

func TestCmd(t *testing.T) {
	fmt.Println(cases.Title(language.Und).String(runtime.GOARCH))
}

func TestURL(t *testing.T) {
	var l = "https://localhost:443/a/b/c/d/你好?k1=v1&k2=1,2,3,4,5,6&h=你好#hash你好"
	u, _ := url.Parse(l)
	s, _ := json.Marshal(u)
	fmt.Println(u.String())
	fmt.Println(string(s))
	fmt.Println(u.EscapedFragment())
	fmt.Println(u.EscapedPath())
}

func TestTLS(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://www.baidu.com", nil)
	res, _ := http.DefaultClient.Do(req)
	fmt.Println(res.Status)
}

func TestRegExp(t *testing.T) {
	u := regexp.MustCompile("^https?://[A-Za-z0-9.]+:?[0-9]*/?")
	fmt.Println(u.MatchString("http://192.168.2.1:3000/"))
	fmt.Println(u.MatchString("http://192.168.2.1/"))
	fmt.Println(u.MatchString("https://localhost.com:3000"))
	fmt.Println(u.MatchString("https://localhost.com:3000/a/b/c"))

	u2 := regexp.MustCompile("^https?://[A-Za-z0-9.]+:?[0-9]*")
	fmt.Println(u2.FindString("https://localhost.com:3000/a/b/c"))
}

type Test1 struct {
	a string
}

func TestCopy(t *testing.T) {
	var p = &Test1{}
	var c = *p
	c.a = "123"
	var p2 = &c
	fmt.Println(p.a == p2.a)

	var p3 = &(*(&Test1{}))
	p3.a = "456"
	fmt.Println(p.a == p3.a)
}
