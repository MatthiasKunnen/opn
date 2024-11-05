## opn query mime

Queries the applications associated with a MIME type

### Synopsis

Returns the desktop IDs of the applications associated with the given mime type.

```
opn query mime <MimeType> [flags]
```

### Examples

```
$ opn query mime application/pdf
```

### Options

```
  -h, --help   help for mime
```

### Options inherited from parent commands

```
      --format format   Sets the output format. Either json or verbose. The verbose output is not stable.
                        If the result is to be processed by a script, use the json format. (default verbose)
      --skip-cache      Do not use the cache. Instead, all lookups are performed on the file system.
```

### SEE ALSO

* [opn query](opn_query.md)	 - Query the associations and desktop IDs

