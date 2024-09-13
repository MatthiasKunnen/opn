package query

import (
	"errors"
	"fmt"
	"github.com/MatthiasKunnen/opn/opnlib"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var mimeCmd = &cobra.Command{
	Use:   "mime",
	Short: "Queries the applications associated with a MIME type",
	Long:  `Returns the desktop IDs of the applications associated with the given mime type.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mimeType := args[0]
		db, err := opnlib.LoadDb("")
		switch {
		case errors.Is(err, opnlib.ErrDbNotFound):
			log.Printf("%s", err)
			log.Printf("Generate the database using opn gendb")
		}

		if err != nil {
			log.Printf("Error loading database: %v\n", err)
			os.Exit(1)
		}

		desktopIds := db.Associations[mimeType]

		if len(desktopIds) == 0 {
			fmt.Println("No applications associated with this MIME type")
		} else {
			println(strings.Join(desktopIds, " "))
		}
	},
}
