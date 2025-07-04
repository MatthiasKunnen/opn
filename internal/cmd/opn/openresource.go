package opn

import (
	"github.com/MatthiasKunnen/opn/internal/opn"
	"github.com/spf13/cobra"
)

var openResourceCmd = &cobra.Command{
	Use:     "resource <File or URL>",
	Aliases: []string{"file-or-url", "r"},
	Short:   "Open the given resource (file or URL)",
	Long: `Looks up and presents all applications that can open this URL/file.
The user can then select the application to open the URL/file with.

For details, see:
- For files: opn file --help.
- For URLs: opn url --help.
`,
	Example: `With file:
$ opn resource foo.pdf

With URL:
$ opn resource https://example.com`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		opn.FileOrUrl(args[0], opn.OpenerOpts{
			MimeOverride: mime,
			SkipCache:    skipCache,
		})
	},
}

func init() {
	openResourceCmd.SetHelpTemplate(openUrlCmd.HelpTemplate() + openHelpTemplate)
	openResourceCmd.Flags().BoolVar(
		&skipCache,
		"skip-cache",
		false,
		"Do not use the cache. Instead, all lookups are performed on the file system.",
	)
	openResourceCmd.Flags().StringVar(
		&mime,
		"mime-type",
		"",
		"Set the mime type of the file/resource at the URL's location and skip automatic determination.",
	)
}
