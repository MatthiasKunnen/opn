package cmd

import (
	"github.com/spf13/cobra"
)

var openFileCmd = &cobra.Command{
	Use:   "files",
	Short: "Open the given files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// @todo, how do we handle multiple files that could potentially be of a different mime
		//   type.
	},
}
