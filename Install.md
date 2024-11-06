# Installation

## Arch Linux
`opn` is available on the AUR: <https://aur.archlinux.org/packages/opn>.

## Manual install
1. `go install github.com/MatthiasKunnen/opn@latest`
1. Add shell completions (optional)  
   You may need to change these destinations based on the distro you are using and the permissions 
   you have.
    ```shell
    opn completion bash > "/usr/share/bash-completion/completions/opn"
    opn completion fish > "/usr/share/fish/vendor_completions.d/opn.fish"
    opn completion zsh > "/usr/share/zsh/site-functions/_opn"
    ```
