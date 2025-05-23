## opn query desktop-locations

Queries the locations of a desktop ID

### Synopsis

Returns a list of all the desktop files that match a given desktop ID. The files are
returned in order from highest priority to lowest.

```
opn query desktop-locations <desktop ID> [flags]
```

### Examples

```
$ opn query desktop-locations vim.desktop
```

### Options

```
  -h, --help   help for desktop-locations
```

### Options inherited from parent commands

```
      --format format   Sets the output format. Either json or verbose. The verbose output is not stable.
                        If the result is to be processed by a script, use the json format. (default verbose)
      --skip-cache      Do not use the cache. Instead, all lookups are performed on the file system.
```

### SEE ALSO

* [opn query](opn_query.md)	 - Query the associations and desktop IDs

