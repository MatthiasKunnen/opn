package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MatthiasKunnen/opn/opnlib"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var desktopLocationsCmd = &cobra.Command{
	Use:   "desktop-locations <desktop ID>",
	Short: "Queries the locations of a desktop ID",
	Long: `Returns a list of all the desktop files that match a given desktop ID. The files are
returned in order from highest priority to lowest.`,
	Example: `$ opn query desktop-locations vim.desktop`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		opn := &opnlib.Opn{
			SkipCache: skipCache,
		}
		err := opn.LoadAndSave()
		switch {
		case errors.Is(err, opnlib.FailedToSaveCache):
			log.Printf("%v\n", err)
		case err != nil:
			log.Fatalf("Failed to load: %v", err)
		}

		desktopId := args[0]
		result := opn.GetDesktopFileLocations(desktopId)

		if len(result) == 0 && !strings.HasSuffix(desktopId, ".desktop") {
			log.Printf(
				"No desktop file found with ID %s, but it does not end in .desktop. "+
					"Assuming it was forgotten.",
				desktopId,
			)
			result = opn.GetDesktopFileLocations(desktopId + ".desktop")
		}

		switch format {
		case outputVerbose:
			fmt.Println(strings.Join(result, "\n"))
		case outputJson:
			err := json.NewEncoder(os.Stdout).Encode(result)
			if err != nil {
				log.Fatalf("Failed to encode JSON: %v", err)
			}
		}
	},
}
