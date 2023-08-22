package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var configFilePath string

func LoadConfig() {

}

// initConfigLoader gives app loading config from env and dir.
// rootPersistentFlags is to set config-file-path by cmdline.
func initConfigLoader(basename string, rootPersistentFlags *pflag.FlagSet) {
	cobra.OnInitialize(func() {
		// Here You know:
		// 1, assign special config file path.
		// 2, find in /HomeDir/.like-iam
		if configFilePath != "" {
			viper.SetConfigFile(configFilePath)
		} else {
			addConfigInHome()
			addConfigInProject()

			viper.SetConfigName(basename)
		}

		viper.AutomaticEnv()
		viper.SetEnvPrefix("li_")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

		if err := viper.ReadInConfig(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration file: %v\n", err)
			os.Exit(1)
		}
	})

	rootPersistentFlags.StringVarP(&configFilePath, "config", "c", "",
		"Read configuration from specified `FILE`, "+
			"support JSON, TOML, YAML, HCL, or Java properties formats.")

	if err := viper.BindPFlag("config", rootPersistentFlags.Lookup("config")); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: fail to bind pflage `config`: %v.", err)
		os.Exit(1)
	}
}

// parseConfig bind App's flags to viper and unmarshal to App's options.
func (a *App) parseConfig() error {
	if err := viper.BindPFlags(a.cmd.Flags()); err != nil {
		return err
	}
	if err := viper.Unmarshal(a.options); err != nil {
		return err
	}
	return nil
}

func addConfigInHome() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home + ".like-iam/")
}

func addConfigInProject() {
	viper.AddConfigPath("." + "config/")
}
