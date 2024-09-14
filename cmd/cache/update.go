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
		index, err := opnlib.GenerateIndex()
		if err != nil {
			log.Fatalf("Failed to generate index: %v", err)
		}

		err = index.SaveIndex("")
		if err != nil {
			log.Fatalf("Failed to save the newly generated cache: %v", err)
		}

		println("Cache successfully updated.")
	},
}
