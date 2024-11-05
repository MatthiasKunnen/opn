package query

import (
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
	"log"
)

type outputMode enumflag.Flag

const (
	outputVerbose = iota
	outputJson
)

var skipCache bool
var format outputMode

var outputModeMap = map[outputMode][]string{
	outputJson:    {"json"},
	outputVerbose: {"verbose"},
}

var QueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query the associations and desktop IDs",
	Example: `Get all desktop IDs associated with a file:
$ opn query file <file>

Get all desktop IDs associated with a MIME type:
$ opn query mime <MIME type>

Get the locations of the .desktop files for a given desktop ID:
$ opn query desktop-locations <desktop ID>
`,
}

func init() {
	QueryCmd.AddCommand(mimeCmd)
	QueryCmd.AddCommand(fileCmd)
	QueryCmd.AddCommand(desktopLocationsCmd)

	QueryCmd.PersistentFlags().BoolVar(
		&skipCache,
		"skip-cache",
		false,
		"Do not use the cache. Instead, all lookups are performed on the file system.",
	)

	formatFlag := enumflag.New(&format, "format", outputModeMap, enumflag.EnumCaseInsensitive)
	const formatFlagName = "format"
	QueryCmd.PersistentFlags().Var(
		formatFlag,
		formatFlagName,
		`Sets the output format. Either json or verbose. The verbose output is not stable.
If the result is to be processed by a script, use the json format.`)

	err := formatFlag.RegisterCompletion(QueryCmd, formatFlagName, enumflag.Help[outputMode]{
		outputJson:    "Stdout is JSON.",
		outputVerbose: "Stdout is verbose text and not stable.",
	})
	if err != nil {
		log.Printf("failed to register shell completion of query --format flag: %v\n", err)
	}

	// @todo Add compact option with separator. 0 => nul delimited, all the rest, simple
	//   replacement.
}
