## opn resource

Open the given resource (file or URL)

### Synopsis

Looks up and presents all applications that can open this URL/file.
The user can then select the application to open the URL/file with.

For details, see:
- For files: `opn file --help`.
- For URLs: `opn url --help`.


```
opn resource <File or URL> [flags]
```

### Examples

```
With file:
$ opn resource foo.pdf

With URL:
$ opn resource https://example.com
```

### Options

```
  -h, --help               help for resource
      --mime-type string   Set the mime type of the file/resource at the URL's location and skip automatic determination.
      --skip-cache         Do not use the cache. Instead, all lookups are performed on the file system.
```


### Attaching to terminal

Applications that need a terminal can be launched in the current terminal or be opened in a new
terminal. By default, GUI applications are started detached from the terminal and terminal
applications are opened in the current terminal.
This behavior can be controlled interactively or using an environment variable.
Interactively, when choosing the application, optionally append the start mode to the index:

- `a`, attached, the application will be opened in the current terminal.
- `d`, detached. GUI application will be detached, terminal applications will be opened in a new
  terminal based on [`OPN_TERM_CMD`](#opn_term_cmd).

For example, 3h will launch the application with index 3 in the current terminal.
If no start mode is specified, [`OPN_START_MODE`](#opn_start_mode) is used to determine the
default.

### Environment

#### OPN_START_MODE
Configures where to open applications.

```shell
# The default, GUI applications are detached and terminal applications will be opened in the
# current terminal.
OPN_START_MODE="gui:d,term:a"

# Open both GUI and terminal applications are detached from the terminal.
OPN_START_MODE="gui:d,term:d"
```

The start mode can be overwritten by appending it to the application's index.

#### OPN_TERM_CMD
The command to use when starting an application that has Terminal=true.
The arguments will be appended to this command.
E.g. `foot`, `gnome-terminal --`.

#### TERMINAL_COMMAND
Lower priority alias for [OPN_TERM_CMD](#opn_term_cmd).


### SEE ALSO

* [opn](opn.md)	 - opn, a fast terminal file opener

