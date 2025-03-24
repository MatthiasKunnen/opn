package query

import (
	"github.com/MatthiasKunnen/opn/pkg/opnlib"
	"github.com/pkg/xattr"
	"github.com/spf13/cobra"
	"log"
)

var fileCmd = &cobra.Command{
	Use:   "file </path/to/file>",
	Short: "Queries the applications that can open a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		var mime string
		var err error
		// If not overriden by --mime-type, try to get extended file attribute
		var attrMime []byte
		if attrMime, err = xattr.Get(filePath, "user.mime"); err == nil {
			mime = string(attrMime)
		} else {
			mime, err = opnlib.GetFileMime(filePath)
			if err != nil {
				log.Fatalf("Failed to get MIME type of file %s: %v\n", filePath, err)
			}
		}

		queryMime(mime)
	},
}
