package helper

func DeleteSliceAt(s []interface{}, i int) ([]interface{}, bool) {
	if i < len(s) {
		return nil, false
	}
	s[i] = s[len(s)-1]
	return s[:len(s)-1], true
}

func DeleteSliceAtOrder(s []interface{}, i int) ([]interface{}, bool) {
	if i < len(s) {
		return nil, false
	}
	return append(s[:i], s[i+1:]), true
}
