<div align="center">
<img src="assets/beautree-banner.svg" alt="beautree" width="600" /><br /><br />

![Version](https://img.shields.io/github/v/release/patppuccin/beautree?style=flat-square&color=e5c07b&labelColor=1e1e2e)
![License](https://img.shields.io/github/license/patppuccin/beautree?style=flat-square&color=61afef&labelColor=1e1e2e)
![Build](https://img.shields.io/github/actions/workflow/status/patppuccin/beautree/release.yml?style=flat-square&color=98c379&labelColor=1e1e2e)
![Go version](https://img.shields.io/github/go-mod/go-version/patppuccin/beautree?style=flat-square&color=56b6c2&labelColor=1e1e2e)

<img src="assets/beautree-demo.gif" alt="beautree demo" width="900" />
</div>

## About

`beautree` is a directory tree viewer for people who find the classic `tree` command lacking. It renders clean Unicode connectors, respects your `.gitignore`, handles symlinks honestly, and gets out of the way.

No subcommands. No magic. A single static binary you can alias to `tree` and forget about.

**What it does:**

- Unicode box-drawing connectors (`╭─ ├─ ╰─`) with automatic ASCII fallback on dumb terminals and pipes
- Reads nested `.gitignore` files at every level, including negation patterns and dir-only rules
- Detects working and broken symlinks and shows their targets inline
- Identifies named pipes and sockets
- Emits JSON output for piping into other tools
- Supports a three-layer config system: CLI flags, env vars, and `config.toml`

**What it doesn't do:**

- No icons, no colors that need a theme, no opinions about your terminal
- No global ignore files, no hidden magic
- No subcommands

## Installation

**Linux / macOS**

```sh
curl -fsSL https://raw.githubusercontent.com/patppuccin/beautree/main/scripts/install.sh | sh
```

**Windows (PowerShell)**

```powershell
irm https://raw.githubusercontent.com/patppuccin/beautree/main/scripts/install.ps1 | iex
```

**Manual:** grab a binary from [releases](https://github.com/patppuccin/beautree/releases) and drop it on your `PATH`.

**Alias to `tree`**

```sh
# bash / zsh / fish
alias tree=beautree

# powershell
Set-Alias tree beautree

# nushell
alias tree = beautree
```

## Usage

```
beautree [path] [flags]
```

`path` defaults to the current directory.

| Flag            | Short | Description                                | Default   |
| --------------- | ----- | ------------------------------------------ | --------- |
| `--depth N`     | `-L`  | Max depth to recurse, 0 = unlimited        | `0`       |
| `--all`         | `-a`  | Include hidden files and dirs              | `false`   |
| `--dirs-only`   | `-d`  | Show directories only                      | `false`   |
| `--size`        |       | Show human-readable size beside each entry | `false`   |
| `--ignore GLOB` | `-I`  | Exclude entries matching GLOB, repeatable  |           |
| `--no-ignore`   |       | Disable `.gitignore` parsing               | `false`   |
| `--format`      | `-f`  | Output format: `unicode` `ascii` `json`    | `unicode` |
| `--no-color`    |       | Disable color output                       | `false`   |
| `--no-summary`  |       | Disable summary footer                     | `false`   |
| `--output FILE` | `-o`  | Write output to file                       |           |
| `--version`     | `-v`  | Print version                              |           |
| `--help`        | `-h`  | Show help                                  |           |

```sh
# two levels deep
beautree -L 2

# dirs only, hidden included
beautree -da

# exclude patterns
beautree -I node_modules -I dist -I "*.log"

# show sizes
beautree --size

# pipe to jq
beautree --format json | jq .

# write to file
beautree -o tree.txt
```

## Config

Persistent preferences live at `~/.config/beautree/config.toml`. The `XDG_CONFIG_HOME` environment variable is respected if set.

```toml
# ~/.config/beautree/config.toml

depth      = 4
all        = false
no_summary = false
format     = "unicode"
ignore     = ["node_modules", "dist", "*.log"]
```

Environment variables are also supported, prefixed with `BEAUTREE_`:

```sh
BEAUTREE_DEPTH=3
BEAUTREE_FORMAT=ascii
BEAUTREE_NO_SUMMARY=true
```

**Precedence order**: `CLI flags > env vars > config file`

## Contributing

Pull requests are not being accepted at this time, but issues are always welcome. Bug reports, feature requests, and general feedback can be submitted [here](https://github.com/patppuccin/beautree/issues).

## License

The project is released under the [Apache-2.0](LICENSE) license.
