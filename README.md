# SkillStack CLI

The official command-line tool for building, testing, and submitting SkillStack
Guided Projects from your local workstation.

The CLI helps you:

- authenticate with a SkillStack account;
- initialize Guided Project starter files;
- run local verification tests; and
- submit verified checkpoints and receive progress credit.

## Table of contents

- [Installation](#installation)
  - [1. Install Go](#1-install-go)
  - [2. Install the SkillStack CLI](#2-install-the-skillstack-cli)
  - [3. Verify the installation](#3-verify-the-installation)
- [Getting started](#getting-started)
  - [Log in](#log-in)
  - [Initialize a project](#initialize-a-project)
  - [Run and submit checkpoints](#run-and-submit-checkpoints)
- [Command reference](#command-reference)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
- [Developing locally](#developing-locally)

## Installation

### 1. Install Go

The SkillStack CLI requires **Go 1.22 or newer**.

For Linux, WSL, and macOS, the [Webi Go installer](https://webinstall.dev/golang/)
is a convenient option:

```sh
curl -sS https://webi.sh/golang | sh
```

You can also follow the [official Go installation instructions](https://go.dev/doc/install).
After installing Go, open a new terminal and verify it:

```sh
go version
```

### 2. Install the SkillStack CLI

Install the latest CLI with `go install`:

```sh
go install github.com/bagelbrantley/skillstack@latest
```

Go installs the executable in `GOBIN` (usually `~/go/bin`). Make sure that
directory is on your `PATH`.

### 3. Verify the installation

```sh
skillstack --help
skillstack version
```

If your shell cannot find `skillstack`, see
[Troubleshooting](#command-not-found).

## Getting started

### Log in

Start the login flow:

```sh
skillstack login
```

The CLI displays a URL where you can request a single-use login code. Paste
that code into the terminal when prompted.

For scripts and CI environments, pass the code explicitly and request JSON
output:

```sh
skillstack login --code YOUR_LOGIN_CODE --json
```

Credentials are stored locally at
`~/.config/skillstack/credentials.json` with owner-only permissions.

### Initialize a project

Create the starter workspace for the default Guided Project:

```sh
skillstack init
cd engine-kv
```

You can provide a project ID explicitly. `clone` and `fetch` are aliases for
`init`:

```sh
skillstack init project-kv-store
skillstack clone project-kv-store
```

### Run and submit checkpoints

Run local verification tests without submitting:

```sh
skillstack run
```

Pass a lesson ID when needed:

```sh
skillstack run proj-lesson-1-2
```

After the checks pass, submit the checkpoint:

```sh
skillstack submit
```

Or run the checks and submit in one step:

```sh
skillstack run --submit
# -s is the short form
skillstack run -s
```

## Command reference

| Command | Description |
| --- | --- |
| `skillstack login` | Authenticate with SkillStack |
| `skillstack logout` | Remove locally stored credentials |
| `skillstack status` | Check authentication and session status |
| `skillstack list` | List available Guided Projects |
| `skillstack show <project>` | Show project or lesson details |
| `skillstack init [project_id]` | Initialize a local project workspace |
| `skillstack run [lesson_id]` | Run local verification tests |
| `skillstack run -s [lesson_id]` | Run tests and submit on success |
| `skillstack submit [lesson_id]` | Run tests and submit a checkpoint |
| `skillstack check` | Check Go, Git, and environment readiness |
| `skillstack version` | Display version and build information |

Every command supports `--help`:

```sh
skillstack run --help
skillstack submit --help
```

## Configuration

### Server URL

The CLI uses `http://localhost:3000` by default. Override it for an individual
command with the global `--server` (`-S`) flag:

```sh
skillstack --server https://your-skillstack-server.example login
```

You can also set the `SKILLSTACK_SERVER` environment variable:

```sh
export SKILLSTACK_SERVER=https://your-skillstack-server.example
```

### JSON output

Use `--json` for machine-readable output. The flag is available globally:

```sh
skillstack --json status
skillstack --json run
```

### Configuration files

Authentication credentials are saved to:

```text
~/.config/skillstack/credentials.json
```

The CLI also reads optional Viper settings from
`~/.config/skillstack/.skillstack.yaml`. Use `--config PATH` to select a
different configuration file.

## Troubleshooting

### Command not found

Find Go's binary directory:

```sh
go env GOBIN
go env GOPATH
```

If `GOBIN` is empty, the default executable directory is usually
`$GOPATH/bin`. Add it to your shell's `PATH`, then start a new shell:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Check the local environment

Run the built-in diagnostics command:

```sh
skillstack check
```

### Reset authentication

Log out and authenticate again:

```sh
skillstack logout
skillstack login
```

## Developing locally

Clone the repository and build from its root:

```sh
git clone https://github.com/bagelbrantley/skillstack.git
cd skillstack
make build
./bin/skillstack version
```
