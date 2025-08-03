# AGENT.md

## Purpose
This is the repository for `wsj`, a scripting tool for Worldographer.
The syntax for this DSL (domain specific language) is based on a subset of Javascript and is specific to Worldographer's data structures.

## Project Structure
    project-root/
    ├── deploy/              # Deployment scripts and configs
    │   ├── deploy.sh
    │   └── ansible/         # If we use Ansible
    │       └── playbook.yml
    │
    ├── dist/                # Build artifacts, one directory per deploy target
    │   ├── linux/           # Linux (production target)
    │   └── local/           # Local development
    │
    ├── docs/                # Documentation (Markdown, Diagrams, etc.)
    │
    ├── testdata/            # Data for testing the application
    │
    ├── tmp/                 # Temporary files; used by Amp agent when experimenting
    │
    ├── tools/               # Development scripts and tools
    │   └── ... (dev tools, bash scripts, etc.)
    │
    ├── ast/                 # AST code
    ├── parser/              # Parser code, grammar uses Pigeon
    ├── vm/                  # VM code, executes the AST directly
    │
    ├── .gitattributes
    ├── .gitignore
    ├── go.mod
    ├── go.work              # Development may use local repositories
    ├── main.go              # Application code
    ├── LICENSE
    ├── README.md
    └── ... (CI/CD configs, etc.)

## Commands
* Generate parser: `pigeon -o parser/grammar.go parser/grammar.peg`
* WSJ runner: `go build -o dist/local/wsj`
* Version: `dist/local/wsj --version`

## Code Style
- Version information for the project is in `main.go`, in the `version` variable
- Standard Go formatting using `gofmt`
- Imports organized by stdlib first, then external packages
- Error handling: return errors to caller, log.Fatal only in main
- Function comments use Go standard format `// FunctionName does X`
- Variable naming follows camelCase
- File structure follows standard Go package conventions

## Packages
* github.com/maloquacious/wxx: read and write Worldographer files.
* github.com/maloquacious/hexg: working with coordinates in the files.
