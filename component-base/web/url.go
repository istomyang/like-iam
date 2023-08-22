package web

import (
	"fmt"
	"istomyang.github.com/like-iam/component-base/json"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// UrlValues parses a struct having json tag into url.Values.
//
// The second return vale is a tpl you can use in UrlValuesWithTpl to reduce cost of CPU in some condition.
func UrlValues(s any) (*url.Values, any) {
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
	var tpl = make([]string, len(sections))
	for i, section := range sections {
		parts := strings.SplitN(section, ":", 2)
		v.Set(parts[0], parts[1])
		tpl[i] = parts[0]
	}
	return &v, tpl
}

// UrlValuesWithTpl fill a url.Values tpl with values.
//
// You must ensure values be in order.
func UrlValuesWithTpl(tpl any, values ...any) (*url.Values, error) {
	t := reflect.TypeOf(tpl)
	if t.Kind() != reflect.Slice {
		return nil, fmt.Errorf("tpl is not a slice type, got: %s", t.Kind())
	}
	var uv = url.Values{}
	for i, k := range tpl.([]string) {
		value := values[i]
		kind := reflect.TypeOf(value).Kind()
		switch kind {
		case reflect.String:
			uv.Set(k, value.(string))
		case reflect.Int:
			uv.Set(k, strconv.Itoa(value.(int)))
		case reflect.Int64:
			uv.Set(k, strconv.FormatInt(value.(int64), 10))
		case reflect.Bool:
			if v := value.(bool); v {
				uv.Set(k, "true")
			} else {
				uv.Set(k, "false")
			}
		default:
			return nil, fmt.Errorf("this kind of value can't be implement, got: %s", kind)
		}
	}
	return &uv, nil
}
