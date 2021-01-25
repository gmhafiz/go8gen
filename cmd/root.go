package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// https://gist.github.com/ik5/d8ecde700972d4378d87#file-colors-go
const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

var rootCmd = &cobra.Command{
	Use:   "go8",
	Short: "go8 is a golang API scaffolder",
	Long: `Quickly scaffold a new Go API site - https://github.com/gmhafiz/go8`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}