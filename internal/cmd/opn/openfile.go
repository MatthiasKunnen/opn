package opn

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/MatthiasKunnen/opn/internal/util"
	"github.com/MatthiasKunnen/opn/pkg/opnlib"
	"github.com/MatthiasKunnen/xdg/desktop"
	"github.com/mattn/go-shellwords"
	"github.com/pkg/xattr"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"syscall"
)

var mime string
var skipCache bool

var appSelectRe = regexp.MustCompile(`^(\d+)(?:\.(\d+))?([ad])?$`)

type StartMode int

const (
	Unset StartMode = iota
	Attached
	Detached
)

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
		if mime == "" {
			// If not overriden by --mime-type, try to get extended file attribute
			var attrMime []byte
			if attrMime, err = xattr.Get(filePath, "user.mime"); err == nil {
				mime = string(attrMime)
			}
		}
		if mime == "" {
			mime, err = opnlib.GetFileMime(filePath)
			if err != nil {
				log.Fatalf("Failed to get MIME type of file %s: %v\n", filePath, err)
			}
		}

		type DesktopInfo struct {
			Entry    *desktop.Entry
			FilePath string
			Id       string
			Actions  []desktop.Action
		}

		desktopFiles := make([]*DesktopInfo, 0)
		desktopIdsSet := make(map[string]bool)

		for _, mimeInfo := range opn.GetDesktopIdsForBroadMime(mime) {
			for _, desktopId := range mimeInfo.DesktopIds {
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

				desktopInfo := &DesktopInfo{
					Id:       desktopId,
					FilePath: desktopFilePath,
					Entry:    entry,
					Actions:  make([]desktop.Action, 0),
				}
				desktopFiles = append(desktopFiles, desktopInfo)

				for _, action := range entry.Actions {
					if entry.Exec.CanOpenFiles() && !action.Exec.CanOpenFiles() {
						// If this subaction does not have a field code indicating file opening
						// support, but the main action does, assume this is on purpose.
						continue
					}

					desktopInfo.Actions = append(desktopInfo.Actions, action)
				}
			}
		}

		if len(desktopFiles) == 0 {
			log.Fatalf("No applications found that can open \"%s\"", mime)
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

		startModeEnv := os.Getenv("OPN_START_MODE")
		startModeTerm := Attached
		startModeGui := Detached
		startMode := Unset

		for _, tc := range strings.Split(startModeEnv, ",") {
			if tc == "" {
				continue
			}

			tcParts := strings.Split(tc, ":")
			if len(tcParts) != 2 {
				log.Fatalf(
					"Invalid value of OPN_START_MODE: '%s'. "+
						"The target conf must contain a single colon.",
					startModeEnv,
				)
			}

			switch tcParts[0] {
			case "gui":
				switch tcParts[1] {
				case "a":
					startModeGui = Attached
				case "d":
					startModeGui = Detached
				default:
					log.Fatalf(
						"Unknown start mode in OPN_START_MODE for gui: '%s'. "+
							"Either 'a' or 'd' expected",
						tcParts[1],
					)
				}
			case "term":
				switch tcParts[1] {
				case "a":
					startModeTerm = Attached
				case "d":
					startModeTerm = Detached
				default:
					log.Fatalf(
						"Unknown start mode in OPN_START_MODE for terminal: '%s'. "+
							"Either 'a' or 'd' expected",
						tcParts[1],
					)
				}
			default:
				log.Fatalf(
					"Unknown target in OPN_START_MODE: '%s'. Either 'gui' or 'term' expected",
					tcParts[0],
				)
			}
		}

		mainIndex := -1
		actionIndex := -1
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
			case text == "a":
				mainIndex = 0
				startMode = Attached
				break inputLoop
			case text == "d":
				mainIndex = 0
				startMode = Detached
				break inputLoop
			case text == "q":
				return
			case text == "?":
				var sb strings.Builder
				sb.WriteString(`Choose the application to open the file with, using the respective number.
If no number is entered, 0 is assumed.

Optionally append either a or d to control stdin/stdout behavior.
a(ttached): execute program in this terminal.
  When opening with vim, this would launch vim in the current terminal.
d(etached): launch the program detached from the terminal.
  When opening with vim, this would launch vim in a new terminal.

Current defaults:
`)
				if startModeTerm == Attached {
					sb.WriteString("Terminal: attached\n")
				} else {
					sb.WriteString("Terminal: detached\n")
				}

				if startModeGui == Attached {
					sb.WriteString("GUI: attached\n")
				} else {
					sb.WriteString("GUI: detached\n")
				}

				sb.WriteString("\nq to quit\n")
				fmt.Println(sb.String())
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
					case "":
					case "a":
						startMode = Attached
					case "d":
						startMode = Detached
					default:
						log.Fatalf("Unknown start mode: '%s'. Exepected 'a' or 'd'.", matches[3])
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

		if !execVal.CanOpenFiles() {
			// Not ideal, we don't know for sure if the program supports being launched with paths
			// in the arguments. Unfortunately, programs don't always follow the spec.
			log.Printf(
				"Warning: %s does not explicitly declare support for opening a file. "+
					"It is missing a field code in the Exec value. "+
					"The path will be added as last argument.\n", chosen.Id)
			arguments = append(arguments, filePath)
		}

		if startMode == Unset {
			if chosen.Entry.Terminal {
				startMode = startModeTerm
			} else {
				startMode = startModeGui
			}
		}

		switch startMode {
		case Attached:
			eCmd := exec.Command(arguments[0], arguments[1:]...)
			// @todo Think about using syscall.Exec as this would replace the opn process and
			//       release the resources. Gotchas are unknown.
			eCmd.Stdin = os.Stdin
			eCmd.Stdout = os.Stdout
			eCmd.Stderr = os.Stderr
			err = eCmd.Run()
			if err != nil {
				log.Fatalf("Error running command '%s': %v\n", arguments, err)
			}
		case Detached:
			startDetached(chosen.Entry.Terminal, arguments)
		default:
			log.Fatalln("Startmode not configured")
		}
	},
}

func init() {
	openFileCmd.SetHelpTemplate(openFileCmd.HelpTemplate() + `
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
`)
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

func startDetached(isTerminal bool, arguments []string) {
	if isTerminal {
		terminalEnvVars := []string{"OPN_TERM_CMD", "TERMINAL_COMMAND"}
		var terminalArgs []string

		for _, envVar := range terminalEnvVars {
			envVal := os.Getenv(envVar)
			if envVal == "" {
				continue
			}

			parsedArgs, err := shellwords.Parse(envVal)
			if err != nil {
				log.Fatalf("Failed to parse %s=%s: %v", envVar, envVal, err)
			}
			terminalArgs = parsedArgs
			break
		}

		if terminalArgs == nil {
			log.Fatalf(
				"Program needs to be opened in a new terminal but none of these environment variables are set: %s. See opn file --help. \n",
				strings.Join(terminalEnvVars, ", "),
			)
		}

		if util.ParentIsShell() {
			arguments = append(terminalArgs, arguments...)
		} else {
			// If the parent process is not a shell, assume it is a terminal that will close
			// immediately after opn exits. This risks taking out the newly launched detached
			// program.
			// To prevent this, we need to make sure the program is launched before exiting.
			startDetachedWithStartSignaling(terminalArgs, arguments)
			return
		}
	}

	eCmd := exec.Command(arguments[0], arguments[1:]...)
	eCmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // Start new session
	}

	err := eCmd.Start()
	if err != nil {
		log.Fatalf("Error starting command '%s': %v\n", arguments, err)
	}

	err = eCmd.Process.Release()
	if err != nil {
		log.Fatalf("Failed to release process: %v\n", err)
	}
}

// startDetachedWithStartSignaling will start a terminal program in a new terminal.
// See the description in openWithSignalCmd.
func startDetachedWithStartSignaling(terminalArgs []string, launchArgs []string) {
	fifoPath, err := createFifo()
	if err != nil {
		log.Fatalf("Error creating FIFO: %v\n", err)
	}

	selfExe, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to determine opn's path: %v\n", err)
	}

	terminalProgram := terminalArgs[0]
	args := make([]string, 0, len(terminalArgs)+len(launchArgs)+2)
	args = append(args, terminalArgs[1:]...)
	args = append(args, selfExe, "openwithsig", fifoPath)
	args = append(args, launchArgs...)

	eCmd := exec.Command(terminalProgram, args...)
	eCmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // Start new session
	}

	err = eCmd.Start()
	if err != nil {
		log.Fatalf("Error starting command '%s': %v\n", launchArgs, err)
	}

	fifo, err := os.Open(fifoPath)
	if err != nil {
		log.Printf("failed to open fifo file: %v\n", err)
	}
	defer fifo.Close()

	buf := make([]byte, 1)
	if _, err := fifo.Read(buf); err != nil {
		log.Printf("Failed to receive start signal: %v\n", err)
	}

	err = eCmd.Process.Release()
	if err != nil {
		log.Fatalf("Failed to release process: %v\n", err)
	}
}

func createFifo() (string, error) {
	var fifoPath string
existsLoop:
	for i := 0; i < 5; i++ {
		filename := util.RandString(10)
		testPath := filepath.Join(os.TempDir(), filename)
		err := syscall.Mkfifo(testPath, 0600)
		switch {
		case err == nil:
			fifoPath = testPath
			break existsLoop
		case errors.Is(err, os.ErrExist):
			continue
		default:
			return "", fmt.Errorf("Error creating fifo: %w\n", err)
		}
	}

	if fifoPath == "" {
		return "", fmt.Errorf("failed to create fifo, files already exist")
	}

	return fifoPath, nil
}
