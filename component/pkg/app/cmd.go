package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func (a *App) buildCommand(run func(cmd *cobra.Command, args []string)) {
	appCmd := &cobra.Command{
		Use:   a.basename,
		Short: a.brief,
		Long:  a.description,
		Run:   run,
	}

	a.cmd = appCmd
}

func (a *App) ExecuteCommand() {
	if err := a.cmd.Execute(); err != nil {
		// TODO: output colorful text
		// fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		fmt.Printf("%v: %v\n", "Error", err)
		os.Exit(1)
	}
}
