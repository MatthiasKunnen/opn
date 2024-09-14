package query

import (
	"fmt"
	"github.com/MatthiasKunnen/opn/opnlib"
	"github.com/spf13/cobra"
	"strings"
)

var mimeCmd = &cobra.Command{
	Use:   "mime",
	Short: "Queries the applications associated with a MIME type",
	Long:  `Returns the desktop IDs of the applications associated with the given mime type.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		queryMime(args[0])
	},
}

func queryMime(mimeType string) {
	index := opnlib.MustLoadIndex(skipCache, "")
	desktopIds := index.Associations[mimeType]

	if len(desktopIds) > 0 {
		println(strings.Join(desktopIds, " "))
	}

	fmt.Printf("%s: No applications associated with this MIME type.\n", mimeType)

	if mimeType = opnlib.GetBroaderMimeType(mimeType); mimeType != "" {
		fmt.Println("Trying broader mime type.")
		desktopIds := index.Associations[mimeType]

		if len(desktopIds) > 0 {
			fmt.Printf("%s: %s\n", mimeType, strings.Join(desktopIds, " "))
		} else {
			fmt.Printf("%s: No applications associated with this MIME type.\n", mimeType)
		}
	}
}
