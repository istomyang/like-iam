package options

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"istomyang.github.com/like-iam/component-base/helper"
	"istomyang.github.com/like-iam/component/pkg/middleware"
	"time"
)

type ServerOpts struct {
	Mode    string `json:"mode"        mapstructure:"mode"`
	Healthz bool   `json:"healthz"     mapstructure:"healthz"`

	// PresetMiddlewares save preset middlewares in middlewares.PresetMiddlewares.
	PresetMiddlewares []string `json:"preset-middlewares" mapstructure:"preset-middlewares"`

	// ShutdownTime gives times to program do some clean work after close signal reached.
	ShutdownTime time.Duration `json:"shutdown-time" mapstructure:"shutdown-time"`
}

func NewServerOpts() *ServerOpts {
	var ms []string
	for s, _ := range middleware.PresetMiddlewares {
		ms = append(ms, s)
	}
	return &ServerOpts{
		Mode:              gin.DebugMode,
		Healthz:           true,
		PresetMiddlewares: ms,
		ShutdownTime:      time.Second * 10,
	}
}

func (o *ServerOpts) Validate() []error {
	var err []error

	if !helper.InStrings(o.Mode, []string{gin.DebugMode, gin.ReleaseMode, gin.TestMode}) {
		err = append(err, fmt.Errorf("http server option mode wrong: %s", o.Mode))
	}

	if o.PresetMiddlewares != nil {
		var presets []string
		for k, _ := range middleware.PresetMiddlewares {
			presets = append(presets, k)
		}
		for _, name := range o.PresetMiddlewares {
			if !helper.InStrings(name, presets) {
				err = append(err, fmt.Errorf("must use defined preset middleware, got: %s", name))
			}
		}
	}

	return err
}

func (o *ServerOpts) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Mode, "server.mode", o.Mode, ""+
		"Start the server in a specified server mode. Supported server mode: debug, test, release.")

	fs.BoolVar(&o.Healthz, "server.healthz", o.Healthz, ""+
		"Add self readiness check and install /healthz router.")

	fs.StringSliceVar(&o.PresetMiddlewares, "server.preset-middlewares", o.PresetMiddlewares, ""+
		"List of allowed preset-middlewares for server, comma separated. If this list is empty default middlewares will be used.")

	fs.DurationVar(&o.ShutdownTime, "server.shutdown-time", o.ShutdownTime, "Gives times to program do some "+
		"clean work after close signal reached.")
}
