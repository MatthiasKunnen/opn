#!/usr/bin/env bash

# Use the CLI file opener program "opn" which allows selecting the application to open the
# file with based on the XDG Desktop and MimeApps Specification.
# See https://github.com/MatthiasKunnen/opn
#
# 1. Copy this file to your system. Suggested locations are:
#    - In home, `~/.config/nnn/plugins/opn_nnn`.
#    - Outside home, `/usr/lib/opn/nnn`.
# 2. Set the following environment variables:
#    - Set the `NNN_OPENER` environment variable to the absolute path of the script you just copied.
#    - Include `c` in `NNN_OPTS` or use the `-c` flag when running nnn. This disables `-e`.
#
#    This is usually done in `.bashrc`/`.zshrc`/...
# 3. Done. Now, when you open a file (default Enter), `opn` will be used.
#
# Required configuration of opn:
#   1. Install opn, see https://github.com/MatthiasKunnen/opn/blob/master/Install.md
#   2. Set OPN_TERM_CMD to configure your preferred terminal.
#      E.g. "foot" or "gnome-terminal --".
#      See https://github.com/MatthiasKunnen/opn?tab=readme-ov-file#terminal-applications

set -uo pipefail

FPATH="$1"

clear -x

opn file "${FPATH}"
status=$?

if [ $status -eq 0 ]; then
	exit 0
fi

# Show any error message to the user
read -rsn1 -p"Press any key to return to nnn";echo
