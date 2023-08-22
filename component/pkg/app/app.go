package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type Option func(app *App)

// App is a main struct to define cmdline application.
// It's recommended use app.NewApp() to create your application.
// You can add properties you need by referring to cobra.Command.
// https://github.com/spf13/cobra/blob/main/user_guide.md.
type App struct {
	basename    string
	name        string
	brief       string
	description string

	options FlagOptions

	noConfig bool // load options from config file.

	// TODO:
	quiet     bool // remove log
	noVersion bool // show version flag.

	run func() error // run after App has initialized.

	cmd *cobra.Command
}

// WithBrief gives App a brief introduction.
// E.g. "Hugo is a very fast static site generator"
func WithBrief(brief string) Option {
	return func(a *App) {
		a.brief = brief
	}
}

// WithDescription gives App a long introduction text.
// E.g. `A Fast and Flexible Static Site Generator built with \n
// love by spf13 and friends in Go. \n
// Complete documentation is available at https://gohugo.io/documentation/`
func WithDescription(brief string) Option {
	return func(a *App) {
		a.brief = brief
	}
}

func WithOptions(options FlagOptions) Option {
	return func(a *App) {
		a.options = options
	}
}

// WithNoConfig forbids App using config file.
func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

// WithQuiet reduces App's log.
func WithQuiet() Option {
	return func(a *App) {
		a.quiet = true
	}
}

// New is recommended to create cmdline application.
// basename is App's binary name, required, E.g. "hugo".
// name is App's title name, required, E.g. "Tom's Hugo"
func New(basename string, name string, run func() error, options ...Option) *App {
	app := &App{
		basename: basename,
		name:     name,
		run:      run,
	}

	for _, option := range options {
		option(app)
	}

	app.buildCommand(func(cmd *cobra.Command, args []string) {
		if !app.noConfig {
			if err := app.parseConfig(); err != nil {
				_ = fmt.Errorf("app parse config error: %s", err.Error())
				os.Exit(1)
			}
		}

		// In this Step, options has initialized by flags or config.

		if app.options != nil {
			cobra.CheckErr(app.options.Validate())
		}

		if err := app.run(); err != nil {
			_ = fmt.Errorf("app run got error: %s", err.Error())
			os.Exit(1)
		}
	})

	app.addFlags()

	if !app.noConfig {
		initConfigLoader(app.basename, app.cmd.PersistentFlags())
	}

	return app
}

func (a *App) Run() {
	a.ExecuteCommand()
}
