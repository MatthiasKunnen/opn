## opn

opn, a fast terminal file opener

### Synopsis

opn is a terminal program meant for opening files with the selected
associated application.

It uses xdg-mime or the file command to determine the MIME type of the
file and the Desktop Entry and MIMEApps specification to determine the
applications that can open the MIME type.

```
opn [flags]
```

### Examples

```
Open a file/URL:
$ opn resource foo.pdf

Open a file:
$ opn file /path/to/file

Open a URL:
$ opn url https://example.com

Get a list of applications that can open a file.
$ opn query file /path/to/file

Get a list of applications that can open a MIME type.
$ opn query mime text/html
```

### Options

```
  -h, --help      help for opn
      --version   Version info
```

### SEE ALSO

* [opn cache](opn_cache.md)	 - Update and view info of the cache
* [opn file](opn_file.md)	 - Open the given file
* [opn query](opn_query.md)	 - Query the associations and desktop IDs
* [opn resource](opn_resource.md)	 - Open the given resource (file or URL)
* [opn url](opn_url.md)	 - Open the given URL

