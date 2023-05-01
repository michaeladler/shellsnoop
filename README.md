# shellsnoop

shellsnoop is a tool that spies on shells and captures user input.
Under the hood, it uses eBPF magic to do this.

It currently supports the following shells:

- zsh

Support for bash and fish is planned.

## Usage

`shellsnoop` must be run as root user. You can configure it to only spy on shells belonging to a specific uid, e.g.

```
# sudo shellsnoop -u 1000
```

From now on shellsnoop intercepts the input for all new shells.

You can use `shellsnoop-client` to ask `shellsnoop` for the last command used in a given shell, e.g.

```bash
$ shellsnoop-client 42
```

queries `shellsnoop` for the last command entered in the shell with pid 42.

## Installation

### From Source

Install `libbpf`, then run:

```bash
$ make
# make install
```

### NixOS

```nix
# flake.nix

{
  inputs.shellsnoop = {
    url = "github:michaeladler/shellsnoop";
    inputs.nixpkgs.follows = "nixpkgs";
  };

  # ...

  outputs = {nixpkgs, hyprland, ...}: {
    nixosConfigurations.HOSTNAME = nixpkgs.lib.nixosSystem {
      modules = [
        shellsnoop.nixosModules.default
      ];
    };
  };
}

# configuration.nix

{inputs, pkgs, ...}: {
  # note: this will start shellsnoop as a systemd service
  programs.shellsnoop.enable = true;
  programs.shellsnoop.uid = config.user.uid;
}
```

## Use-Cases

### tmux-resurrect

Use `shellsnoop` as a save command strategy for [tmux-resurrect](https://github.com/tmux-plugins/tmux-resurrect):

1. Create the executable file `save_command_strategies/shellsnoop.sh` (inside the `tmux-resurrect` plugin directory):

```bash
#!/usr/bin/env bash
set -eu

PID=$1
CHILDREN="/proc/$PID/task/$PID/children"
if [ -e "$CHILDREN" ]; then
    CONTENT=$(cat "$CHILDREN")
    if [ "$CONTENT" != "" ]; then
        exec shellsnoop-client "$PID"
    fi
fi
exit 0
```

2. Set `shellsnoop` as your save command strategy:
```bash
# tmux.conf

set -g @resurrect-save-command-strategy 'shellsnoop'
```
