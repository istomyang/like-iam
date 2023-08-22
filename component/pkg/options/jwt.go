package options

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/spf13/pflag"
	"time"
)

type JwtOpts struct {
	// Realm name to display to the user. Required.
	Realm      string        `json:"realm" mapstructure:"realm"`
	Key        string        `json:"key" mapstructure:"key"`
	Timeout    time.Duration `json:"timeout" mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"max_refresh" mapstructure:"max_refresh"`
}

func NewJwtOpts() *JwtOpts {
	return &JwtOpts{
		Realm:      "jwt",
		Key:        "",
		Timeout:    1 * time.Hour,
		MaxRefresh: 1 * time.Hour,
	}
}

func (o *JwtOpts) Validate() []error {
	var err []error

	if !govalidator.StringLength(o.Key, "6", "32") {
		err = append(err, fmt.Errorf("--secret-key must larger than 5 and little than 33"))
	}

	return err
}

func (o *JwtOpts) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Realm, "jwt.realm", o.Realm, "Realm name to display to the user.")
	fs.StringVar(&o.Key, "jwt.key", o.Key, "Private key used to sign jwt token.")
	fs.DurationVar(&o.Timeout, "jwt.timeout", o.Timeout, "JWT token timeout.")
	fs.DurationVar(&o.MaxRefresh, "jwt.max-refresh", o.MaxRefresh, ""+
		"This field allows clients to refresh their token until MaxRefresh has passed.")
}
