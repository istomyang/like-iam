package helper

// InStrings check one string in strings.
func InStrings(s string, ss []string) bool {
	for _, s2 := range ss {
		if s == s2 {
			return true
		}
	}
	return false
}
