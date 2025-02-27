# Integrations
This document describes how `opn` can be integrated with other tools.

## nnn
**nnn** is a terminal file manager. `opn` can be configured as its opener.

1. Copy [`nnn.sh`](./nnn.sh?raw=1) to your system. Suggested locations are:
   - In home, `~/.config/nnn/plugins/opn_nnn`.
   - Outside home, `/usr/lib/opn/nnn`.
1. Set the following environment variables:
   - Set the `NNN_OPENER` environment variable to the absolute path of the script you just copied.
   - Include `c` in `NNN_OPTS` or use the `-c` flag when running nnn.

   This is usually done in `.bashrc`/`.zshrc`/...
1. Done. Now, when you open a file in nnn (default Enter), `opn` will be used.
