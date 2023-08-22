package ctl

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"istomyang.github.com/like-iam/iam/internal/ctl/global"
	"istomyang.github.com/like-iam/iam/internal/ctl/util"
	"istomyang.github.com/like-iam/log"
	"os"
	"strings"
)

type IO struct {
	Out   io.Reader
	In    io.Writer
	Error io.Writer
}

// Run gets new ctl command and start to run and should be called in main function.
func Run(cmd *cobra.Command) error {
	return cmd.Execute()
}

// NewStd creates a ctl with standard out, in and error.
func NewStd() *cobra.Command {
	return New(IO{
		Out:   os.Stdout,
		In:    os.Stdin,
		Error: os.Stderr,
	})
}

// New creates a ctl command with specific in, out and error streams.
func New(io IO) *cobra.Command {
	var appCmd = &cobra.Command{
		Use:   "iamctl",
		Short: "iamctl controls the iam platform",
		Long: util.NewNormalize(`
			iamctl controls the iam platform, is the client side tool for iam platform.
			Find more information at:
				https://github.com/marmotedu/iam/blob/master/docs/guide/en-US/cmd/iamctl/iamctl.md`).Heredoc().String(),
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
		// children command can run.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return global.DefaultProfile.Run()
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return global.DefaultProfile.Close()
		},
	}

	appCmd.SetGlobalNormalizationFunc(normalizeDash2Underline())
	setUsageFunc(appCmd)
	setHelpFunc(appCmd)

	var fs = appCmd.PersistentFlags()
	fs.SetNormalizeFunc(normalizeDash2Underline())

	defaultFactory.AddFlagTo(fs)

	global.DefaultProfile.AddFlag(fs)
	global.DefaultVersion.AddFlag(fs)

	configViper(fs)

	initCommands(fs)

	return appCmd
}

func normalizeDash2Underline() func(f *pflag.FlagSet, name string) pflag.NormalizedName {
	return func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		if strings.Contains(name, "_") {
			strings.ReplaceAll(name, "_", "-")
			log.Warnf("config field should use dash other than underline, got: %s", name)
		}
		return pflag.NormalizedName(name)
	}
}

// configViper config load config from file or other places using viper.
func configViper(fs *pflag.FlagSet) {
	var cfgFile = viper.GetString("config")
	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)

			// Search config in home directory with name ".cobra" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigType("yaml")
			viper.SetConfigName(".cobra")
		}

		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	})

	_ = viper.BindPFlags(fs)
}

func initCommands(fs *pflag.FlagSet) {

}

func setHelpFunc(command *cobra.Command) {
	//TODO implement me
}

func setUsageFunc(command *cobra.Command) {
	//TODO implement me
}
