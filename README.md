# OPN - a file opener

`opn`, pronounced _open_ (creative, I know), allows you to pick the application to open a file with
in the terminal from a list of preferred applications.

![Example of opening a PDF file with opn](.github/example_open_pdf.svg)

## Installation
See [Install.md](Install.md).

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

### Which terminal to open in?
Unfortunately, there is no specification yet, though
[one is being developed](https://gitlab.freedesktop.org/terminal-wg/specifications/-/merge_requests/3),
on how to set a preferred terminal emulator across Desktop Environments or systems without one.
This requires an environment variable to make `opn` aware of your preference.

Specify the terminal to be launched using the `OPN_TERM_CMD` environment variable. E.g:
- `OPN_TERM_CMD="foot"`
- `OPN_TERM_CMD="gnome-terminal --"`

### Attaching to terminal
By default, applications are opened in a new terminal.
This behavior can be controlled using either the `OPN_TERM_TARGET` environment variable or
interactively by appending the target to the index of the application to launch.
The target is either `h`, `b`, or not set, in which case `OPN_TERM_TARGET` will be used.
`h` stands for _here_, and `b` stands for _background_.
For example, 3h will launch the application with index 3 in the current terminal.

## Usage
Open a file using `opn file /path/to/file`.

For detailed usage, see `opn --help` or view the [CLI docs](./docs/cli/opn.md).

## TODO
- Document integration with NNN
- Pacman/package manager hook to update cache
- Generate documentation from CLI help. Man page perhaps?
- Add examples on how to change preferred applications.
- Support DBUS activation for starting applications. Pending interest.
