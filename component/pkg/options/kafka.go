package options

import (
	"github.com/spf13/pflag"
)

type KafkaOpts struct {
	Addrs         string         `json:"addrs,omitempty" mapstructure:"addrs"`
	Topic         string         `json:"topic,omitempty" mapstructure:"topic"`
	MessageExtend map[string]any `json:"message-extend" mapstructure:"message-extend"`
}

func (i *KafkaOpts) AddFlags(set *pflag.FlagSet) {
}

func (i *KafkaOpts) Validate() []error {
	return nil
}

var _ ValidatableOptions = &KafkaOpts{}

var _ FlagOptions = &KafkaOpts{}
