package opn

import (
	"github.com/MatthiasKunnen/opn/internal/opn"
	"github.com/spf13/cobra"
)

var openUrlCmd = &cobra.Command{
	Use:   "url <URL>",
	Short: "Open the given URL",
	Long: `Looks up and presents all applications that can open this URL.
The user can then select the application to open the URL with.

Works by determining the MIME type of the URL and then finding all
applications that can open it according to the MIME Applications Associations
specification.

If --mime-type is not set (most common usage), the suggested applications
will be those that have an x-scheme-handler defined for the URL's protocol.
For http(s) URLs, the user can opt to choose to download the file to a
temporary location where the mime type will then be determined using:
1. The Content-Type header if it is set.
2. The sniffed MIME type.

Downloading is done using 'D' in the interactive prompt.

If --mime-type is set, the suggested applications will be those that support
opening that MIME type.
`,
	Example: `opn url https://example.com`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		opn.Url(args[0], opn.OpenerOpts{
			MimeOverride: mime,
			SkipCache:    skipCache,
		})
	},
}

func init() {
	openUrlCmd.SetHelpTemplate(openUrlCmd.HelpTemplate() + openHelpTemplate)
	openUrlCmd.Flags().BoolVar(
		&skipCache,
		"skip-cache",
		false,
		"Do not use the cache. Instead, all lookups are performed on the file system.",
	)
	openUrlCmd.Flags().StringVar(
		&mime,
		"mime-type",
		"",
		"Set the mime type of the resource at the URL's location and skip automatic determination.",
	)
}
