# OPN - a file opener

`opn` allows you to pick the program to open a file with in the terminal.

The following approaches are used to determine the list of programs that are suggested for a file:
- A `mimemap.conf` generated from the desktop files on your system. See [desktop files](#desktop-files)
- Your own custom `.mimemap.conf` files. These config files map MIME types to a program. See [mimemap](#mimemap)
- `mimeapps.list`


It uses `xdg-mime` from the xdg-utils package to determine the MIME type of the file and maps
it to programs that can open it.

It uses the following ways to determine which programs are suggested:
- The desktop files on your system. See [desktop files](#desktop-files).
- The `mimemap.conf` file. See [desktop files](#desktop-files).

## Desktop files
Desktop files specify how to launch a program and the mime types it can open.
`opn` uses the desktop files to generate a map between

## Usage

### Open file
`opn /path/to/file`

### Reload desktop entries
Searches for, and reads, .desktop files

1. Read custom mapping and add to list.
2. Read `mimeapps.list` and add to list.  
   See https://specifications.freedesktop.org/mime-apps-spec/latest/file.html
3. Read desktop files and add to list.  
   See https://specifications.freedesktop.org/menu-spec/latest/paths.html
4. Store in `mimemap.json`

`opn update-desktop`

### List programs
`opn list /path/to/file`
