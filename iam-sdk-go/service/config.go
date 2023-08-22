package service

import (
	"gopkg.in/yaml.v3"
	"istomyang.github.com/like-iam/iam-sdk-go/service/iam"
	"os"
)

type Config struct {
	IAM *iam.Config `yaml:"iam,omitempty"`
}

func (c *Config) Correct() []error {
	var errs []error

	return errs
}

// LoadFromFile load config from file which ext name is .yaml
// filePath can use relative and absolute path.
func LoadFromFile(filePath string) (*Config, []error) {
	var errs []error

	bs, err := os.ReadFile(filePath)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	var config Config
	err = yaml.Unmarshal(bs, &config)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	errs = config.Correct()
	if errs != nil {
		return nil, errs
	}

	return &config, nil
}
