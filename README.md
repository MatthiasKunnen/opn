# OPN - a file opener

`opn`, pronounced _open_ (creative, I know), allows you to pick the application to open a file with
in the terminal from a list of preferred applications.


## How it works
To determine the list of programs that are suggested for a file, the 
[MIME Applications Specification](https://specifications.freedesktop.org/mime-apps-spec/1.0.1/)
is used.

A simplified explanation:
- `.desktop` files, located in `$XDG_DATA_HOME/applications` and `$XDG_DATA_DIRS/applications`,
  specify the name of applications, their path, and optionally some MIME types they can open.
  These files are typically included when installing applications.  
  The name of these files equals the _desktop ID_.
  See the [Desktop Entry spec](https://specifications.freedesktop.org/desktop-entry-spec/1.5/index.html#).
- `mimeapps.list` files, located in one of the [specified directories](https://specifications.freedesktop.org/mime-apps-spec/1.0.1/file.html),
  associates a MIME type with desktop IDs.
- `opn` reads these desktop and `mimeapps.list` files to determine the MIME-application association.
- When `opn file path/to/file` is executed, `xdg-mime` or `file` is used to determine the MIME type
  of the file. Then, the MIME is looked up in the MIME-application index to find all associated
  applications.

## Requirements
- `xdg-mime` (preferred, better accuracy, `xdg-utils` package) or `file`

## Terminal applications
Some applications require a terminal to launch in. Examples of this include; vim, nano, nnn, and opn
itself. These applications have `Terminal=true` in their respective desktop files.
Unfortunately, there is no specification yet on how to set a preferred terminal emulator across
Desktop Environments or systems without one so we use an environment variable to make `opn` aware of
your preference.

Specify the terminal to be launched using the `OPN_TERMINAL_COMMAND` environment variable. E.g:
- `OPN_TERMINAL_COMMAND="foot"`
- `OPN_TERMINAL_COMMAND="gnome-terminal --"`

## Usage
For detailed usage, see `opn --help` or view [cli_docs/opn.md](./cli_docs/opn.md).

### Open file
`opn file /path/to/file`

## TODO
- Document integration with NNN
- Publish to AUR
- Installation instructions
- Pacman/package manager hook to update cache
- Update cache if older than x days?
- Generate documentation from CLI help. Man page perhaps?
- Add examples on how to change preferred applications.
- Add compact format to `query` commands for use in scripts. Pending interest.
- Localization, currently, all output is in English. Pending interest.
- Support DBUS activation for starting applications. Pending interest.
