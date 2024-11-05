## opn query

Query the associations and desktop IDs

### Examples

```
Get all desktop IDs associated with a file:
$ opn query file <file>

Get all desktop IDs associated with a MIME type:
$ opn query mime <MIME type>

Get the locations of the .desktop files for a given desktop ID:
$ opn query desktop-locations <desktop ID>

```

### Options

```
      --format format   Sets the output format. Either json or verbose. The verbose output is not stable.
                        If the result is to be processed by a script, use the json format. (default verbose)
  -h, --help            help for query
      --skip-cache      Do not use the cache. Instead, all lookups are performed on the file system.
```

### SEE ALSO

* [opn](opn.md)	 - opn, a fast terminal file opener
* [opn query desktop-locations](opn_query_desktop-locations.md)	 - Queries the locations of a desktop ID
* [opn query file](opn_query_file.md)	 - Queries the applications that can open a file
* [opn query mime](opn_query_mime.md)	 - Queries the applications associated with a MIME type

