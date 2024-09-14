package query

import (
	"github.com/MatthiasKunnen/opn/opnlib"
	"github.com/spf13/cobra"
	"log"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Queries the applications that can open a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		mime, err := opnlib.GetFileMime(filePath)
		if err != nil {
			log.Fatalf("Failed to get MIME type of file %s: %v\n", filePath, err)
		}

		queryMime(mime)
	},
}
