package authzserver

import (
	"context"
	"istomyang.github.com/like-iam/component/pkg/app"
	"istomyang.github.com/like-iam/log"
)

const description = `Authorization server to run ladon policies which can protecting your resources.
It is written inspired by AWS IAM policies.

Find more ladon information at:
    https://github.com/ory/ladon`

// NewApp provides an entry to run, with basename depends on binary name you choose.
// Then, you should call app.App's Run().
func NewApp(basename string) *app.App {
	// Create a new options with default value, and then app.New will
	// put config and flags value to options, and then do validate.
	options := NewOptions(basename)
	newApp := app.New("IAM Authorization Server",
		basename,
		Run(options),
		app.WithBrief("IAM Authorization Server is a authz app."),
		app.WithOptions(options),
		app.WithDescription(description))
	return newApp
}

func Run(options *Options) func() error {
	return func() error {
		// In this context, options is initialized by cobra flags and viper.
		log.Init(context.Background(), options.Log)
		defer log.Sync()

		if err := newAuthzServer(options).run(); err != nil {
			return err
		}

		return nil
	}
}
