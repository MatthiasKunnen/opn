## opn file

Open the given file

### Synopsis

Looks up and presents all applications that can open this file.
The user can then select the application to open the file with.

Works by first obtaining the MIME type of the file and then finding all
applications that can open it according to the MIME Applications Associations
specification.

```
opn file <filename> [flags]
```

### Examples

```
opn file foo.pdf
```

### Options

```
  -h, --help         help for file
      --skip-cache   Do not use the cache. Instead, all lookups are performed on the file system.
```

### Attaching to terminal

Applications that need a terminal can be launched in the current terminal or be opened in a new
terminal. By default, applications are opened in a new terminal. This behavior can be controlled
using the OPN_TERM_TARGET environment variable or, interactively, by appending the target to the
index of the application to launch. The target is either 'h', 'b', or not set, in which case
'OPN_TERM_TARGET' will be used. 'h' stands for _here_, and 'b' stands for _background_.
For example, 3h will launch the application with index 3 in the current terminal.

### Environment

#### OPN_TERM_CMD
The command to use when starting an application that has Terminal=true.
The arguments will be appended to this command.
E.g. "foot", "gnome-terminal --".

#### OPN_TERM_TARGET
The default target to open terminal applications in:
`b`, background (default), a new terminal will be spawned based on `OPN_TERM_CMD`.
`h`, here, the application will be opened in the current terminal.
The target can still be overwritten by appending the target to the application's index.

### SEE ALSO

* [opn](opn.md)	 - opn, a fast terminal file opener

