package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var openFileCmd = &cobra.Command{
	Use:   "files",
	Short: "Open the given files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}
