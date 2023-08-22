package options

import (
	"github.com/spf13/pflag"
	"regexp"
)

type MongoOpts struct {
	Addr           string
	DbName         string
	CollectionName string

	CAFile  string
	KeyFile string
}

func (i *MongoOpts) AddFlags(set *pflag.FlagSet) {
	//TODO implement me
	panic("implement me")
}

func (i *MongoOpts) Validate() []error {
	var errs []error
	re := regexp.MustCompile(`(.*):(.*)`)
	i.Addr = re.FindString(i.Addr)
	return errs
}

var _ ValidatableOptions = &MongoOpts{}

var _ FlagOptions = &MongoOpts{}
