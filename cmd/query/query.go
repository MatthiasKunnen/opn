package query

import (
	"github.com/spf13/cobra"
)

var skipCache bool

var QueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query the associations",
}

func init() {
	QueryCmd.AddCommand(mimeCmd)
	QueryCmd.AddCommand(fileCmd)

	QueryCmd.PersistentFlags().BoolVar(
		&skipCache,
		"skip-cache",
		false,
		"Do not use the cache. Instead, all lookups are performed on the file system.",
	)
}
