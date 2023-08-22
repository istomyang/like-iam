package global

import "github.com/spf13/pflag"

type Version struct {
	matchVersion bool
}

var DefaultVersion = &Version{}

// AddFlag adds:
// 1. enable match client and server version.
func (v *Version) AddFlag(fs *pflag.FlagSet) {

}

// Check determine client version and server version is matched.
func (v *Version) Check() error {
	// use http client to request and compare.
	return nil
}
