package log

import (
	"fmt"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

const (
	EncodingJson    = "json"
	EncodingConsole = "console"
)

const (
	flagName              = "log.name"
	flagLevel             = "log.level"
	flagDevelopment       = "log.development"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagEncoding          = "log.encoding"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagEnableColor       = "log.enable-color"
)

// Options provides config to create a logger.
// You can add others props as you need.
type Options struct {
	Name              string   `json:"name"              mapstructure:"name"`
	Level             string   `json:"level"              mapstructure:"level"`
	Development       bool     `json:"development"        mapstructure:"development"`
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	Encoding          string   `json:"encoding"           mapstructure:"encoding"`
	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	EnableColor       bool     `json:"enable-color"       mapstructure:"enable-color"`
}

// NewOptions use preset flag to build Options.
func NewOptions(basename string, development *bool) *Options {
	if development == nil {
		development = new(bool)
		env := os.Getenv("APP_ENV")
		if env == "development" {
			*development = true
		} else if env == "production" {
			*development = false
		} else {
			*development = true
		}
	}

	if *development {
		return NewDevelopmentOptions(basename)
	} else {
		return NewProductionOptions(basename)
	}
}

// NewDevelopmentOptions is a recommended config in development environment.
func NewDevelopmentOptions(basename string) *Options {
	return &Options{
		Name:              basename,
		Level:             DebugLevel.string(),
		Development:       true,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		Encoding:          EncodingConsole,
		DisableCaller:     false,
		DisableStacktrace: false,
		EnableColor:       true,
	}
}

// NewProductionOptions is a recommended config in development environment.
func NewProductionOptions(basename string) *Options {
	return &Options{
		Name:              basename,
		Level:             InfoLevel.string(),
		Development:       false,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		Encoding:          EncodingJson,
		DisableCaller:     false,
		DisableStacktrace: false,
		EnableColor:       true,
	}
}

func (o *Options) AddFlags(set *pflag.FlagSet) {
	set.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
	set.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	set.BoolVar(&o.Development, flagDevelopment, o.Development,
		"Development puts the logger in development mode, which changes"+
			"the behavior of DPanicLevel and takes stacktraces more liberally.")
	set.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	set.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
	set.StringVar(&o.Encoding, flagEncoding, o.Encoding, "Log output `FORMAT`, support plain or json format, "+
		"value you can pass `console` and `json`.")
	set.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	set.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	set.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
}

func (o *Options) Validate() []error {
	var errors []error
	// level string
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(o.Level)); err != nil {
		errors = append(errors, err)
	}

	// encoding
	lower := strings.ToLower(o.Encoding)
	if lower != EncodingConsole && lower != EncodingJson {
		errors = append(errors, fmt.Errorf("option 'Encoding' mistake, got %s", o.Encoding))
	}

	return errors
}

// Build uses Options to create Logger.
func (o *Options) Build() (*zap.Logger, error) {
	options := o
	if errs := options.Validate(); len(errs) > 0 {
		for _, err := range errs {
			panic(fmt.Sprintf("error: build options get wrong string, please check: %v", err))
		}
	}

	// refer to zap.NewProductionConfig
	var sampling *zap.SamplingConfig = nil
	if !options.Development {
		sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
			Hook:       nil,
		}
	}

	// lv, _ := zap.ParseAtomicLevel(options.Level)
	// Use level limitation in this package layer.
	lv := zap.NewAtomicLevelAt(zapcore.DebugLevel)

	config := zap.Config{
		Level:             lv,
		Development:       options.Development,
		DisableCaller:     options.DisableCaller,
		DisableStacktrace: options.DisableStacktrace,
		Sampling:          sampling,
		Encoding:          options.Encoding,
		EncoderConfig:     buildEncoderConfig(options),
		OutputPaths:       options.OutputPaths,
		ErrorOutputPaths:  options.ErrorOutputPaths,
	}

	build, err := config.Build()
	if err != nil {
		return nil, err
	}

	build.Named(options.Name)
	build.WithOptions(zap.WithCaller(true))

	zap.RedirectStdLog(build)
	zap.ReplaceGlobals(build)

	return build, nil
}

func buildEncoderConfig(options *Options) (encoderConfig zapcore.EncoderConfig) {

	// You can refer to zap src.
	if options.Development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	// You can use others like `time.RFC3339`
	const timeFormat = "2022-10-01 18:00:00.000"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(timeFormat)

	encoderConfig.EncodeName = zapcore.FullNameEncoder // default
	encoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	if options.Encoding == EncodingConsole && options.EnableColor {
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}

	return
}
