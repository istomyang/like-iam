package v1

import (
	"encoding/json"
)

// Extend provides a field type for database storing extended data.
// Using DB Hook feature to exchange string for Extend and vice versa.
type Extend map[string]interface{}

// String transfers Extend to String, and store in persistent media.
func (e Extend) String() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// Load init empty Extend with string coming from persistent media.
func (e Extend) Load(extentShadow string) error {
	var extent Extend
	if err := json.Unmarshal([]byte(extentShadow), &extent); err != nil {
		return err
	}
	for s, i := range extent {
		e[s] = i
	}
	return nil
}
