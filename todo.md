1. Figure out how xdg-mime query default gets the default for application/pdf. Does it load mimeinfo.cache? 
2. Figure how `update-desktop-database` updates `mimeinfo.cache`. Do we reimplement this functionality? Does it update any other files?
  `desktop-file-install` looks interesting as well, same util.
    - https://www.freedesktop.org/wiki/Software/desktop-file-utils/  
    - https://gitlab.freedesktop.org/xdg/desktop-file-utils
