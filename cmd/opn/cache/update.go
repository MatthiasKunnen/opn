package cache

import (
	"github.com/MatthiasKunnen/opn/opnlib"
	"github.com/spf13/cobra"
	"log"
)

var updateCacheCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the index that is used to look up MIME/application association",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		opn := &opnlib.Opn{
			SkipCache: true,
		}
		err := opn.Load()
		if err != nil {
			log.Fatalf("Failed generate index: %v", err)
		}
		err = opn.SaveIndex()
		if err != nil {
			log.Fatalf("Failed to save index: %v", err)
		}

		println("Cache successfully updated.")
	},
}
