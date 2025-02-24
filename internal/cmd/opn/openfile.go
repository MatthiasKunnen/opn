package opn

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/MatthiasKunnen/opn/pkg/opnlib"
	"github.com/MatthiasKunnen/xdg/desktop"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"slices"
	"strconv"

	"github.com/mattn/go-shellwords"
)

var skipCache bool

var appSelectRe = regexp.MustCompile(`^(\d+)(?:\.(\d+))?([bh])?$`)

var openFileCmd = &cobra.Command{
	Use:   "file <filename>",
	Short: "Open the given file",
	Long: `Looks up and presents all applications that can open this file.
The user can then select the application to open the file with.

Works by first obtaining the MIME type of the file and then finding all
applications that can open it according to the MIME Applications Associations
specification.`,
	Example: `opn file foo.pdf`,
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
			log.Fatalf("Error loading: %v", err)
		}

		filePath := args[0]
		mime, err := opnlib.GetFileMime(filePath)
		if err != nil {
			log.Fatalf("Failed to get MIME type of file %s: %v\n", filePath, err)
		}

		type DesktopInfo struct {
			Entry    *desktop.Entry
			FilePath string
			Id       string
			Actions  []desktop.Action
		}

		desktopFiles := make([]*DesktopInfo, 0)
		desktopIdsSet := make(map[string]bool)

		for mime != "" {
			for _, desktopId := range opn.GetDesktopIdsForMime(mime) {
				if desktopIdsSet[desktopId] {
					continue
				}

				desktopIdsSet[desktopId] = true
				var desktopParseError error
				var entry *desktop.Entry
				var desktopFilePath string
				for _, desktopFilePath = range opn.GetDesktopFileLocations(desktopId) {
					entry, desktopParseError = desktop.ParseFile(desktopFilePath)
					if err == nil {
						break
					}

					log.Printf("Error parsing desktop file %s: %v\n", desktopFilePath, err)
				}

				if desktopParseError != nil {
					continue
				}

				if !entry.Exec.CanOpenFiles() {
					continue
				}

				desktopInfo := &DesktopInfo{
					Id:       desktopId,
					FilePath: desktopFilePath,
					Entry:    entry,
					Actions:  make([]desktop.Action, 0),
				}
				desktopFiles = append(desktopFiles, desktopInfo)

				for _, action := range entry.Actions {
					if !action.Exec.CanOpenFiles() {
						continue
					}

					desktopInfo.Actions = append(desktopInfo.Actions, action)
				}
			}

			mime = opnlib.GetBroaderMimeType(mime)
		}

		for index, desktopFile := range slices.Backward(desktopFiles) {
			fmt.Printf("%d) %s\n", index, desktopFile.Entry.Name.Default)

			for actionIndex, action := range desktopFile.Actions {
				fmt.Printf(
					"  %d.%d) %s\n",
					index,
					actionIndex+1,
					action.Name.Default,
				)
			}
		}

		mainIndex := -1
		actionIndex := -1
		openHere := false
		scanner := bufio.NewScanner(os.Stdin)
	inputLoop:
		for {
			fmt.Printf(
				"Open %s with (?=help)[0]: ",
				path.Base(filePath),
			)
			scanner.Scan()
			text := scanner.Text()

			switch {
			case text == "":
				mainIndex = 0
				break inputLoop
			case text == "h":
				mainIndex = 0
				openHere = true
				break inputLoop
			case text == "b":
				mainIndex = 0
				openHere = false
				break inputLoop
			case text == "?":
				fmt.Println(`
Choose the application to open the file with using the respective number.
Optionally append either h or b to control stdin/stdout behavior.
h(ere): execute program in this terminal. E.g. when opening with vim, this would launch vim in the
current terminal.
b(ackground) (default): launch the program in the background. When opening with vim, this would
launch vim in a new terminal.

If no number is entered, 0 is assumed.`)
			case appSelectRe.MatchString(text):
				matches := appSelectRe.FindStringSubmatch(text)
				mainIndex, err = strconv.Atoi(matches[1])
				maxMainIndex := len(desktopFiles) - 1
				actionIndex = -1
				if err != nil {
					log.Printf("Error converting %s to int: %v.\n", text, err)
					continue
				} else if mainIndex < 0 {
					log.Printf("Number cannot be less than 0, got %d.\n", mainIndex)
					continue
				} else if mainIndex > maxMainIndex {
					log.Printf(
						"Number cannot be greater than %d, got %d.\n",
						maxMainIndex,
						mainIndex,
					)
					continue
				}

				if len(matches) > 1 && matches[2] != "" {
					newActionIndex, err := strconv.Atoi(matches[2])
					maxSubIndex := len(desktopFiles[mainIndex].Actions)

					if err != nil {
						log.Printf("Error converting %s to int: %v.\n", text, err)
						continue
					} else if newActionIndex < 1 {
						log.Printf("Sub index cannot be less than 1, got %d.\n", newActionIndex)
						continue
					} else if newActionIndex > maxSubIndex {
						log.Printf(
							"Sub index cannot be greater than %d, got %d.\n",
							maxSubIndex,
							newActionIndex,
						)
						continue
					} else {
						actionIndex = newActionIndex - 1
					}
				}

				if len(matches) > 2 {
					switch matches[3] {
					case "h":
						openHere = true
					case "b":
						openHere = false
					}
				}

				break inputLoop
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("error reading standard input: %v\n", err)
		}

		chosen := desktopFiles[mainIndex]
		var execVal desktop.ExecValue
		if actionIndex > -1 {
			execVal = chosen.Actions[actionIndex].Exec
		} else {
			execVal = chosen.Entry.Exec
		}

		arguments := execVal.ToArguments(desktop.FieldCodeProvider{
			GetDesktopFileLocation: func() string {
				return chosen.FilePath
			},
			GetFile: func() string {
				return filePath
			},
			GetFiles: func() []string {
				return []string{filePath}
			},
			GetName: func() string {
				return chosen.Entry.Name.Default
			},
			GetUrl: func() string {
				return filePath
			},
			GetUrls: func() []string {
				return []string{filePath}
			},
		})

		if chosen.Entry.Terminal && !openHere {
			terminalCommand := os.Getenv("OPN_TERMINAL_COMMAND")
			if terminalCommand == "" {
				log.Fatal("Program needs to be opened in a terminal but OPN_TERMINAL_COMMAND" +
					" is not set. See opn file --help.")
			}

			terminalArgs, err := shellwords.Parse(terminalCommand)
			if err != nil {
				log.Fatalf("Failed to parse OPN_TERMINAL_COMMAND=%s: %v", terminalCommand, err)
			}

			arguments = append(terminalArgs, arguments...)
		}

		eCmd := exec.Command(arguments[0], arguments[1:]...)
		if openHere {
			// @todo Think about using syscall.Exec as this would replace the opn process and
			//       release the resources. Gotchas are unknown.
			eCmd.Stdin = os.Stdin
			eCmd.Stdout = os.Stdout
			eCmd.Stderr = os.Stderr
			err = eCmd.Run()
			if err != nil {
				log.Fatalf("Error running command '%s': %v\n", arguments, err)
			}
		} else {
			err = eCmd.Start()
			if err != nil {
				log.Fatalf("Error starting command '%s': %v\n", arguments, err)
			}
		}
	},
}

func init() {
	openFileCmd.SetHelpTemplate(openFileCmd.HelpTemplate() + `
ATTACHING TO TERMINAL:
  Applications that need a terminal can be launched in the current terminal or be opened in a new
  terminal. By default, applications are opened in a new terminal. This behavior can be changed by
  appending 'h' to the index of the application to launch. For example, 3h will launch the
  application with index 3 in the same terminal. As a mnemonic, the h stands for here.

ENVIRONMENT:
  OPN_TERMINAL_COMMAND
    The command to use when starting an application that has Terminal=true.
    The arguments will be appended to this command.
    E.g. "foot", "gnome-terminal --".
`)
	openFileCmd.Flags().BoolVar(
		&skipCache,
		"skip-cache",
		false,
		"Do not use the cache. Instead, all lookups are performed on the file system.",
	)
}
