# pj

A simple CLI tool to quickly jump into project directories and launch commands.

## Installation

```sh
curl -fsSL https://raw.githubusercontent.com/hlop3z/pj/main/install.sh | sh
```

## Installation (Source)

```bash
go build -o pj.exe .   # Windows
go build -o pj .       # Linux/macOS
```

Move the executable to a directory in your PATH, or add the build directory to your PATH.

## Configuration

Create an `apps.yaml` file in the same directory as the executable:

```yaml
project-name:
  cmd: claude.exe
  path: ~/path/to/project

another-project:
  cmd: code .
  path: ~/work/another-project
```

- `cmd`: The command to run when entering the project
- `path`: The project directory (supports `~` for home directory)

## Usage

List available projects:

```bash
pj
```

Jump to a project and run its command:

```bash
pj <project-name>
```

## Cross-Platform

Works on Windows, Linux, and macOS. Paths use Unix-style forward slashes and `~` expansion regardless of OS.
