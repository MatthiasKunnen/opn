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
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"syscall"
)

var appSelectRe = regexp.MustCompile(`^(\d+)(?:\.(\d+))?([ad])?$`)

type StartMode int

const (
	Unset StartMode = iota
	Attached
	Detached
)

type valueType int

const (
	valueTypeFile valueType = iota
	valueTypeUrl
	valueTypeUnknown
)

type desktopInfo struct {
	Entry    *desktop.Entry
	FilePath string
	Id       string
	Actions  []desktop.Action
}

type OpenerOpts struct {
	MimeOverride string
	SkipCache    bool

	fileOrUrl string
	valueType valueType
}

func Url(url string, opts OpenerOpts) {
	opts.fileOrUrl = url
	opts.valueType = valueTypeUrl
	openFileOrUrl(opts)
}

func File(filePath string, opts OpenerOpts) {
	opts.fileOrUrl = filePath
	opts.valueType = valueTypeFile
	openFileOrUrl(opts)
}

func FileOrUrl(fileOrUrl string, opts OpenerOpts) {
	opts.fileOrUrl = fileOrUrl
	opts.valueType = valueTypeUnknown
	openFileOrUrl(opts)
}

func openFileOrUrl(opts OpenerOpts) {
	newOpener(opts).run()
}

type opener struct {
	mimeOverride          string
	localFile             string
	localFileMime         string
	localFileIsDownloaded bool
	url                   string
	urlIsDownloadable     bool
	urlScheme             string
	opn                   *opnlib.Opn
}

func newOpener(opts OpenerOpts) *opener {
	opn := &opnlib.Opn{
		SkipCache: opts.SkipCache,
	}
	err := opn.LoadAndSave()
	switch {
	case errors.Is(err, opnlib.FailedToSaveCache):
		log.Printf("%v\n", err)
	case err != nil:
		log.Fatalf("Error loading: %v", err)
	}

	o := &opener{
		mimeOverride: opts.MimeOverride,
		opn:          opn,
	}

	switch opts.valueType {
	case valueTypeFile:
		o.localFile = opts.fileOrUrl
	case valueTypeUrl:
		parsedUrl, err := url.Parse(opts.fileOrUrl)
		if err != nil {
			log.Fatalf("Error parsing URL: %v", err)
		}

		if !parsedUrl.IsAbs() {
			log.Fatalf("URL must be absolute: %v", opts.fileOrUrl)
		}
		o.url = opts.fileOrUrl
		o.urlScheme = parsedUrl.Scheme
		o.urlIsDownloadable = isDownloadSupported(parsedUrl.Scheme)
	case valueTypeUnknown:
		parsedUrl, err := url.Parse(opts.fileOrUrl)
		if err != nil || !parsedUrl.IsAbs() {
			o.localFile = opts.fileOrUrl
			break
		}

		o.url = opts.fileOrUrl
		o.urlScheme = parsedUrl.Scheme
		o.urlIsDownloadable = isDownloadSupported(parsedUrl.Scheme)
	}

	return o
}

func (o *opener) run() {
	if o.localFile == "" && (o.url == "" || o.urlScheme == "") {
		panic("Either localFile or both o.url && o.urlScheme must be set")
	}

	o.updateLocalFileMime()
	desktopFiles := o.mustGetOptions()

	var err error
	mainIndex := -1
	actionIndex := -1
	startMode, startModeGui, startModeTerm := getDefaultStartMode()
	scanner := bufio.NewScanner(os.Stdin)
inputLoop:
	for {
		printOptions(desktopFiles)
		fmt.Printf(
			"Open %s with (?=help)[0]: ",
			o.getPrintHint(),
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
		case text == "D":
			if o.url == "" {
				fmt.Println("D(ownload) is only supported for URL inputs")
				break
			}

			if o.localFileIsDownloaded {
				fmt.Println("File is already downloaded")
				break
			}

			if !o.urlIsDownloadable {
				fmt.Println("Download is not supported for this protocol/scheme")
				break
			}

			o.download()
			desktopFiles = o.mustGetOptions()
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
			return o.getExecArg(true)
		},
		GetFiles: func() []string {
			return []string{o.getExecArg(true)}
		},
		GetName: func() string {
			return chosen.Entry.Name.Default
		},
		GetUrl: func() string {
			return o.getExecArg(false)
		},
		GetUrls: func() []string {
			return []string{o.getExecArg(false)}
		},
	})

	if !execVal.CanOpenFiles() {
		// Not ideal, we don't know for sure if the program supports being launched with paths
		// in the arguments. Unfortunately, programs don't always follow the spec.
		log.Printf(
			"Warning: %s does not explicitly declare support for opening a file. "+
				"It is missing a field code in the Exec value. "+
				"The path will be added as last argument.\n", chosen.Id)
		arguments = append(arguments, o.getExecArg(false))
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

}

func (o *opener) updateLocalFileMime() {
	if o.localFile == "" {
		o.localFileMime = ""
		return
	}

	if o.localFileMime == "" && !o.localFileIsDownloaded {
		// If not overriden by --mime-type, try to get extended file attribute
		if attrMime, err := xattr.Get(o.localFile, "user.mime"); err == nil {
			o.localFileMime = string(attrMime)
		}
	}

	if o.localFileMime == "" {
		mime, err := opnlib.GetFileMime(o.localFile)
		if err != nil {
			log.Fatalf("Failed to get MIME type of file %s: %v\n", o.localFile, err)
		}
		o.localFileMime = mime
	}
}

func (o *opener) mustGetOptions() []*desktopInfo {
	desktopFiles := make([]*desktopInfo, 0)
	desktopIdsSet := make(map[string]bool)

	var mimes []string
	if o.mimeOverride != "" {
		mimes = append(mimes, o.mimeOverride)
	} else if o.localFileMime != "" {
		mimes = append(mimes, o.localFileMime)
	}

	if o.urlScheme != "" && !o.localFileIsDownloaded {
		mimes = append(mimes, "x-scheme-handler/"+o.urlScheme)
	}

	var desktopIdsToSuggest []opnlib.MimeDesktopIds
	for _, mime := range mimes {
		desktopIdsToSuggest = append(desktopIdsToSuggest, o.opn.GetDesktopIdsForBroadMime(mime)...)
	}

	for _, mimeInfo := range desktopIdsToSuggest {
		for _, desktopId := range mimeInfo.DesktopIds {
			if desktopIdsSet[desktopId] {
				continue
			}

			desktopIdsSet[desktopId] = true
			var desktopParseError error
			var entry *desktop.Entry
			var desktopFilePath string
			for _, desktopFilePath = range o.opn.GetDesktopFileLocations(desktopId) {
				entry, desktopParseError = desktop.ParseFile(desktopFilePath)
				if desktopParseError == nil {
					break
				}

				log.Printf("Error parsing desktop file %s: %v\n", desktopFilePath, desktopParseError)
			}

			if desktopParseError != nil || entry == nil {
				continue
			}

			if o.localFile == "" && !entry.Exec.CanOpenUrls() && !o.urlIsDownloadable {
				// If opening a URL that is not downloadable, and the desktop entry cannot open
				// URLs, exclude it.
				// If the user downloads the file, it will become localFile and this exclusion will
				// be skipped.
				continue
			}

			desktopInfo := &desktopInfo{
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

				if o.localFile == "" && !action.Exec.CanOpenUrls() && !o.urlIsDownloadable {
					// If opening a URL that is not downloadable, and the desktop entry cannot open
					// URLs, exclude it.
					// If the user downloads the file, it will become localFile and this exclusion will
					// be skipped.
					continue
				}

				desktopInfo.Actions = append(desktopInfo.Actions, action)
			}
		}
	}

	if len(desktopFiles) == 0 {
		log.Fatalf("No applications found that can open %v", mimes)
	}

	return desktopFiles
}

func (o *opener) getExecArg(mustBeLocal bool) string {
	if o.localFile != "" {
		return o.localFile
	}

	if !mustBeLocal && o.url != "" {
		return o.url
	}

	o.download()

	return o.localFile
}

func isDownloadSupported(urlScheme string) bool {
	switch urlScheme {
	case "http", "https":
		return true
	default:
		return false
	}
}

func (o *opener) download() {
	if o.url == "" {
		log.Fatalln("Could not download, URL is not set.")
	}

	if !o.urlIsDownloadable {
		log.Fatalf("Downloading is not supported for the scheme %s\n", o.urlScheme)
	}

	log.Println("Downloading...")

	temp, err := os.CreateTemp("", "opn_download")
	if err != nil {
		log.Fatalf("Error creating temporary file: %v\n", err)
	}
	defer temp.Close()

	resp, err := http.Get(o.url)
	if err != nil {
		log.Fatalf("Error downloading %s: %v\n", o.url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		log.Fatalf("Bad status when dowloading: %s", resp.Status)
	}

	_, err = io.Copy(temp, resp.Body)
	if err != nil {
		log.Fatalf("Error writing downloaded content to file: %v\n", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		o.localFileMime = contentType
	}

	o.localFile = temp.Name()
	o.localFileIsDownloaded = true
	o.updateLocalFileMime()
}

func (o *opener) getPrintHint() string {
	if o.localFile != "" {
		return filepath.Base(o.localFile)
	}

	return o.url
}

func getDefaultStartMode() (StartMode, StartMode, StartMode) {
	startModeEnv := os.Getenv("OPN_START_MODE")
	startMode := Unset
	startModeGui := Detached
	startModeTerm := Attached

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

	return startMode, startModeGui, startModeTerm
}

func printOptions(desktopFiles []*desktopInfo) {
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
