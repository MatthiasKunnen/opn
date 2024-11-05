package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MatthiasKunnen/opn/pkg/opnlib"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var mimeCmd = &cobra.Command{
	Use:     "mime <MimeType>",
	Short:   "Queries the applications associated with a MIME type",
	Long:    `Returns the desktop IDs of the applications associated with the given mime type.`,
	Args:    cobra.ExactArgs(1),
	Example: `$ opn query mime application/pdf`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		queryMime(args[0])
	},
}

func queryMime(mimeType string) {
	opn := &opnlib.Opn{
		SkipCache: skipCache,
	}
	err := opn.Load()
	switch {
	case errors.Is(err, opnlib.FailedToSaveCache):
		log.Printf("%v\n", err)
	case err != nil:
		log.Fatalf("Failed to load: %v", err)
	}
	result := opn.GetDesktopIdsForBroadMime(mimeType)

	switch format {
	case outputVerbose:
		for _, item := range result {
			if len(item.DesktopIds) == 0 {
				fmt.Printf("%s: No associated applications\n", item.Mime)
			} else {
				fmt.Printf("%s: %s\n", item.Mime, strings.Join(item.DesktopIds, ", "))
			}
		}
	case outputJson:
		err := json.NewEncoder(os.Stdout).Encode(result)
		if err != nil {
			log.Fatalf("Failed to encode JSON: %v", err)
		}
	}
}
