package options

import "github.com/spf13/pflag"

// ValidatableOptions abstracts options can be validated after set.
type ValidatableOptions interface {
	Validate() []error
}

// FlagOptions abstracts options can be used with pflag.FlagSet.
type FlagOptions interface {
	// AddFlags adds this options to special flag set.
	AddFlags(*pflag.FlagSet)
}
