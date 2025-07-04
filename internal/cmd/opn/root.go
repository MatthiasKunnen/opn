package opn

import (
	"fmt"
	"github.com/MatthiasKunnen/opn/internal/cmd/opn/cache"
	"github.com/MatthiasKunnen/opn/internal/cmd/opn/query"
	"github.com/spf13/cobra"
	"log"
)

var versionRequested = false

var rootCmd = &cobra.Command{
	Use:   "opn",
	Short: "opn, a fast terminal file opener",
	Long: `opn is a terminal program meant for opening files with the selected
associated application.

It uses xdg-mime or the file command to determine the MIME type of the
file and the Desktop Entry and MIMEApps specification to determine the
applications that can open the MIME type.`,
	Example: `Open a file/URL:
$ opn resource foo.pdf

Open a file:
$ opn file /path/to/file

Open a URL:
$ opn url https://example.com

Get a list of applications that can open a file.
$ opn query file /path/to/file

Get a list of applications that can open a MIME type.
$ opn query mime text/html`,
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		if versionRequested {
			fmt.Println("opn version 0.4.0")
			return
		}

		err := cmd.Help()
		if err != nil {
			log.Fatalf("Error printing help information: %v\n", err)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func GetCommand() *cobra.Command {
	return rootCmd
}

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	rootCmd.AddCommand(cache.CacheCmd)
	rootCmd.AddCommand(openFileCmd)
	rootCmd.AddCommand(openResourceCmd)
	rootCmd.AddCommand(openUrlCmd)
	rootCmd.AddCommand(openWithSignalCmd)
	rootCmd.AddCommand(query.QueryCmd)
	rootCmd.Flags().BoolVar(&versionRequested, "version", false, "Version info")
}
