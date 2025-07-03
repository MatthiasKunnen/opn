package opn

import (
	"github.com/MatthiasKunnen/opn/internal/opn"
	"github.com/spf13/cobra"
)

var openResourceCmd = &cobra.Command{
	Use:     "resource <URL>",
	Aliases: []string{"file-or-url", "r"},
	Short:   "Open the given resource (file or URL)",
	Long: `Looks up and presents all applications that can open this URL/file.
The user can then select the application to open the URL/file with.

Works by first determining the MIME type of the URL and then finding all
applications that can open it according to the MIME Applications Associations
specification.

The MIME type is determined in this order:
1. The value specified using the --mime-type option.
2. The value reported by the relevant utility, xdg-query or file.`,
	Example: `opn file foo.pdf`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		opn.File(args[0], opn.OpenerOpts{
			MimeOverride: mime,
			SkipCache:    skipCache,
		})
	},
}

func init() {
	openResourceCmd.SetHelpTemplate(openUrlCmd.HelpTemplate() + fileHelpTemplate)
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
