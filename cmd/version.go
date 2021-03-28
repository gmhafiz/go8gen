package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of go8",
	Long:  `All software has versions. This is go8's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("go8 golang API scaffolder v0.9.0")
	},
}
