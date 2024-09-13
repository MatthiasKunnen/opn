package query

import (
	"fmt"
	"github.com/spf13/cobra"
)

var QueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query the associations",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}

func init() {
	QueryCmd.AddCommand(mimeCmd)
	QueryCmd.AddCommand(fileCmd)
}
