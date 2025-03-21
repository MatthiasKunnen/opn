package opn

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"strings"
)

var openWithSignalCmd = &cobra.Command{
	Use:    "openwithsig fifo_path cmd_args...",
	Short:  "Executed the given command and send a signal when done",
	Hidden: true,
	Long: `This internal command only exists because of timing issues that can cause programs
to never show when opn closes too quickly after starting the program.
This only applies to detached terminal programs launched from opn running from one-off
terminals.
Check to see if your terminal is affected by running the following command:
  $TERM_CMD setsid gnome-terminal vim /path/to/file
Replace $TERM_CMD with your terminal's specific command, e.g. foot, gnome-terminal --, ...

The reproduction for opn is as follows:
1. Start a terminal and pass opn file /path/to/file,
   e.g. gnome-terminal -- opn file /path/to/file
2. Choose a terminal program (e.g. vim) to open the file with and open it detached.
   This should spawn another terminal with vim.
   However, opn now exits followed by terminal one, which might kill the newly spawned
   terminal before it has a chance to establish itself.

To combat this, and to avoid having to add a dirty, arbitrary, delay, when this scenario is
detected, the following is done:
1. Instead of launching e.g. gnome-terminal -- vim /path/to/file, opn starts
   gnome-terminal -- opn openwithsig fifo_path vim /path/to/file
2. Now there are two opn instances running.
3. The first one will read from the FIFO and block
4. The second opn will start the given command in attached mode.
5. After the command is started, opn two will write to the FIFO, signaling that the process has
   started and opn one will exit after cleaning up.
   This should make sure that the launched program had time to establish itself.

The only downside is that as long as the launched terminal program stays running, opn two will keep
running and using a certain amount of memory.
`,
	Example: `opn openwithsig fifo_path gnome-terminal -- vim /path/to/file`,
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fifoPath := args[0]
		fifo, err := os.OpenFile(fifoPath, os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("Failed to open fifo at %s: %v\n", fifoPath, err)
		}
		defer os.Remove(fifoPath)
		defer fifo.Close()

		arguments := args[1:]
		eCmd := exec.Command(arguments[0], arguments[1:]...)
		eCmd.Stdin = os.Stdin
		eCmd.Stdout = os.Stdout
		eCmd.Stderr = os.Stderr
		err = eCmd.Start()
		if err != nil {
			log.Fatalf(
				"Error running command '%s': %v\n",
				strings.Join(arguments, " "),
				err,
			)
		}

		_, err = fifo.Write([]byte("1")) // Signal that the program has started
		if err != nil {
			log.Printf("Error writing to FIFO: %v\n", err)
		}

		err = eCmd.Wait()
		if err != nil {
			log.Fatalf("Error waiting for command to finish: %v\n", err)
		}
	},
}
