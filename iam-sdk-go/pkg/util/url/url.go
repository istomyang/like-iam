package url

import (
	"fmt"
	"regexp"
)

func PickAddr(addr string) (string, error) {
	re := regexp.MustCompile("^https?://[A-Za-z0-9.]+:?[0-9]*")
	if re.MatchString(addr) {
		return re.FindString(addr), nil
	} else {
		return "", fmt.Errorf("your addr must satisfy regexp format: ^https?://[A-Za-z0-9.]+:?[0-9]*, but got: %s", addr)
	}
}
