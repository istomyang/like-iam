package apiserver

import (
	"context"
	"istomyang.github.com/like-iam/component/pkg/app"
	"istomyang.github.com/like-iam/log"
)

const description = `The IAM API server validates and configures data
for the api objects which include users, policies, secrets, and
others. The API Server services REST operations to do the api objects management.

Find more iam-apiserver information at:
    https://github.com/istomyang/iam/blob/master/docs/guide/en-US/cmd/iam-apiserver.md`

// NewApp provides an entry to run, with basename depends on binary name you choose.
// Then, you should call app.App's Run().
func NewApp(basename string) *app.App {
	// Create a new options with default value, and then app.New will
	// put config and flags value to options, and then do validate.
	options := NewOptions(basename)
	newApp := app.New("IAM ApiServer",
		basename,
		Run(options),
		app.WithBrief("IAM ApiServer is a authn app."),
		app.WithOptions(options),
		app.WithDescription(description))
	return newApp
}

func Run(options *Options) func() error {
	return func() error {
		// In this context, options is initialized by cobra flags and viper.
		log.Init(context.Background(), options.Log)
		defer log.Sync()

		if err := newApiServer(options).run(); err != nil {
			return err
		}

		return nil
	}
}
