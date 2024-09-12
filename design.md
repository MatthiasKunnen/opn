## After changing .desktop file or installing a package
1. Run [`update-desktop-database`](https://www.freedesktop.org/wiki/Software/desktop-file-utils/)  
   Usually executed by the package manager on package installation.
   This updates the `mimeinfo.cache` file containing a lookup of MIME=>desktop IDS.
   This file does not store the preference, this is handled by `mimeapps.list`.
2. Run `opn update-desktop-database`. This makes opn update `db.json`, which stores a map with
   key=MIME type, value=list of desktop IDs in order of most preferred to least.
