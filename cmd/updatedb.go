package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var updateDbCmd = &cobra.Command{
	Use:   "updatedb",
	Short: "Updates the MIME associations",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}
