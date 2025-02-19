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

### SEE ALSO

* [opn](opn.md)	 - opn, a fast terminal file opener

