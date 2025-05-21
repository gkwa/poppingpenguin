package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of poppingpenguin`,
	Run: func(cmd *cobra.Command, args []string) {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			fmt.Println("poppingpenguin version: development")
			return
		}
		fmt.Println("poppingpenguin version:", info.Main.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
