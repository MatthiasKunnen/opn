package opn

import (
	"github.com/MatthiasKunnen/opn/internal/opn"
	"github.com/spf13/cobra"
)

var mime string
var skipCache bool

var openFileCmd = &cobra.Command{
	Use:   "file <filename>",
	Short: "Open the given file",
	Long: `Looks up and presents all applications that can open this file.
The user can then select the application to open the file with.

Works by first determining the MIME type of the file and then finding all
applications that can open it according to the MIME Applications Associations
specification.

The MIME type is determined in this order:
1. The value specified using the --mime-type option.
2. The value of the extended file attribute user.mime, if it exists.
3. The value reported by the relevant utility, xdg-query or file.`,
	Example: `opn file foo.pdf`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		opn.File(args[0], opn.Opts{
			MimeOverride: mime,
			SkipCache:    skipCache,
		})
	},
}

const fileHelpTemplate = `
ATTACHING TO TERMINAL:
  Applications that need a terminal can be launched in the current terminal or be opened in a new
  terminal. By default, GUI applications are started detached from the terminal and terminal
  applications are opened in the current terminal. This behavior can be controlled interactively or
  using an environment variable.
  Interactively, when choosing the application, optionally append the start mode to the index:
    a attached, the application will be opened in the current terminal.
    d detached. GUI application will be detached, terminal applications will be opened in
      a new terminal based on 'OPN_TERM_CMD'.
  For example, 3h will launch the application with index 3 in the current terminal.
  If no start mode is specified, 'OPN_START_MODE' is used to determine the default.

ENVIRONMENT:
  OPN_START_MODE
    Configures where to open applications.
    Examples:
      OPN_START_MODE="gui:d,term:a", the default, GUI applications are detached and terminal
        applications will be opened in the current terminal.
      OPN_START_MODE="gui:d,term:d", always detach.
    The start mode can be overwritten by appending it to the application's index.
  OPN_TERM_CMD
    The command to use when starting an application that has Terminal=true.
    The arguments will be appended to this command.
    E.g. "foot", "gnome-terminal --".
  TERMINAL_COMMAND
    Lower priority alias for OPN_TERM_CMD.
`

func init() {
	openFileCmd.SetHelpTemplate(openFileCmd.HelpTemplate() + fileHelpTemplate)
	openFileCmd.Flags().BoolVar(
		&skipCache,
		"skip-cache",
		false,
		"Do not use the cache. Instead, all lookups are performed on the file system.",
	)
	openFileCmd.Flags().StringVar(
		&mime,
		"mime-type",
		"",
		"Set the mime type of the file and skip automatic determination.",
	)
}
