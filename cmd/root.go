package cmd

import (
	"github.com/MatthiasKunnen/opn/cmd/cache"
	"github.com/MatthiasKunnen/opn/cmd/query"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "opn",
	Short: "opn, a fast terminal file opener",
	Long:  `opn `,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	rootCmd.AddCommand(cache.CacheCmd)
	rootCmd.AddCommand(openFileCmd)
	rootCmd.AddCommand(query.QueryCmd)
}
