```
 ███████╗██╗███╗   ██╗
 ██╔════╝██║████╗  ██║
 █████╗  ██║██╔██╗ ██║
 ██╔══╝  ██║██║╚██╗██║
 ██║     ██║██║ ╚████║
 ╚═╝     ╚═╝╚═╝  ╚═══╝
```

# Fin — Windows Shell DSL Compiler

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![Windows](https://img.shields.io/badge/platform-Windows-0078D4.svg)](https://www.microsoft.com/windows)
[![Version](https://img.shields.io/badge/version-v1.0.0-brightgreen.svg)](releases)

A production-grade, **Fish-inspired DSL** that compiles deterministically into **Windows Batch** scripts. Write readable shell code, generate correct batch automation—zero runtime dependencies.

---

## Features

- **Fish-inspired syntax** — Variables, functions, conditionals, loops with intuitive semantics
- **Deterministic compilation** — Same source always produces identical batch output
- **No dependencies** — Generated `.bat` files run on any Windows system without additional tools
- **Safe by design** — Semantic analysis prevents common shell errors
- **Type-aware** — Handles lists, maps, and string interpolation elegantly
- **Recursive functions** — Full support for function recursion with local scoping
- **Proper delayed expansion** — Correct batch variable expansion in all contexts

---

## Quick Start

### Install

Download the latest installer from [Releases](../../releases):
```
.\scripts\Fin-v1.0.0-Setup.exe
```

Or build from source (see [Building](#-building-from-source)).

### CLI Usage

```bash
# Compile a Fin script
fin build script.fin                 # → script.bat
fin build script.fin -o output.bat   # Custom output

# Validate without compiling
fin check script.fin

# View AST
fin ast script.fin

# Format code
fin fmt script.fin                   # Print formatted
fin fmt -w script.fin                # Write formatted

# Version
fin version
```

---

## Language Features

### Variables
```fin
set name "Alice"
set count 42
set pi 3.14159
echo "Hello, $name! Count: $count"
```

### Arithmetic
```fin
set x 10
set y 5
set result $x + $y        # 15
set squared $x ** 2       # 100
set quotient $x / $y      # 2
```

### Control Flow
```fin
# If/Else
if $x > 5
    echo "x is large"
else
    echo "x is small"
end

# For Loop
for i in 1..5
    echo "Count: $i"
end

# While Loop
set n 10
while $n > 0
    echo "Countdown: $n"
    set n $n - 1
end
```

### Functions
```fin
fn greet name times
    echo "Hello, $name!"
end

greet "World" 1
```

### Lists & Maps
```fin
set colors [red, green, blue]
echo "First: $colors[0]"
echo "Length: $colors_len"

set person {name: "Bob", age: 30, city: "NYC"}
echo "Name: $person.name, Age: $person.age"
```

### String Interpolation
```fin
set user "admin"
set port 8080
echo "Server: $user@localhost:$port"
```

---

## Building from Source

### Prerequisites
- **Go** 1.21+ ([download](https://golang.org/dl))
- **NSIS** (for installer) ([download](https://nsis.sourceforge.io))
- **Windows SDK** (optional, for code signing)

### Build

Using **PowerShell**:
```powershell
.\build.ps1              # Build fin.exe and installer
.\build.ps1 -Sign        # Build and sign installer
.\build.ps1 -Help        # Show options
```

Using **Batch**:
```cmd
build.bat
```

Or manually:
```bash
go build -o fin.exe ./cmd/fin
cd scripts
makensis fin_installer.nsi
```

Output:
- `fin.exe` — Compiler binary
- `scripts/Fin-v1.0.0-Setup.exe` — Installer

---

## Documentation

| File | Purpose |
|------|---------|
| [language.md](docs/language.md) | Complete language specification |
| [cli.md](docs/cli.md) | CLI usage and commands |
| [tests.md](docs/tests.md) | Testing guide |
| [AGENTS.md](AGENTS.md) | Compiler architecture & design rules |

---

## Examples

Complete examples in [examples/](examples/) directory:

| File | Feature |
|------|---------|
| `01_variables_echo.fin` | Variables and echo |
| `02_arithmetic.fin` | Arithmetic operations |
| `03_if_else.fin` | Conditionals |
| `04_for_loop.fin` | For loops |
| `05_while_loop.fin` | While loops |
| `06_functions.fin` | Function definitions |
| `07_lists.fin` | List operations |
| `08_maps.fin` | Map (object) operations |
| `09_nested_control.fin` | Nested loops/conditionals |
| `10_comparisons.fin` | All comparison operators |

Try them:
```bash
fin build examples/01_variables_echo.fin
01_variables_echo.bat
```

---

## Development

### Running Tests
```bash
go test ./...                    # All tests
go test ./internal/generator/... # Generator tests only
```

### Project Structure
```
fin/
 ├─ cmd/fin/               # CLI entry point
 ├─ internal/
 │   ├─ token/             # Token definitions
 │   ├─ lexer/             # Scanner
 │   ├─ parser/            # Parser (recursive descent)
 │   ├─ ast/               # Abstract syntax tree
 │   ├─ sema/              # Semantic analysis
 │   ├─ generator/         # Batch code generation
 │   └─ diagnostics/       # Error reporting
 ├─ examples/              # Example .fin files
 ├─ tests/                 # Integration tests
 ├─ scripts/               # Build scripts & installer
 └─ docs/                  # Documentation
```

### Compiler Pipeline
```
Source Code (.fin)
    ↓
Lexer (Scanner)
    ↓
Parser (AST)
    ↓
Semantic Analysis
    ↓
Code Generator
    ↓
Batch Script (.bat)
```

---

## Supported Operators

| Category | Operators |
|----------|-----------|
| Arithmetic | `+` `-` `*` `/` `%` `**` |
| Comparison | `<` `<=` `>` `>=` `==` `!=` |
| Logical | `&&` `\|\|` `!` |
| Range | `..` |

---

## Design Principles

From [AGENTS.md](AGENTS.md):

1. **Determinism** — Same input always produces identical output
2. **Explicitness** — No magic; all behavior is predictable
3. **Backward compatibility** — Breaking changes require version bump
4. **Tooling-friendly** — Clean errors, AST dumping, formatting

---

## License

MIT License — See [LICENSE](LICENSE) for details.

---

## Author

**bfrovrflw**

---

## Version History

### v1.0.0 (2026-02-01)
- Core language features (variables, functions, control flow)
- List and map support
- Recursive functions
- String interpolation with property/index access
- All comparison operators in if/while
- Proper batch delayed expansion handling
- Semantic analysis and error reporting
- CLI tools (build, check, ast, fmt)
- NSIS installer with PATH updates
- Complete test suite (100+ tests)

---

## Contributing

Contributions welcome! Please:
1. Follow [AGENTS.md](AGENTS.md) architecture rules
2. Add tests for new features
3. Update documentation
4. Run `go test ./...` before submitting

---

## Known Limitations

- **No return values** — Functions don't return values (planned for v1.1)
- **No imports** — All code in single file (planned for v1.1)
- **No closures** — Functions don't capture outer scope
- **Windows only** — Generates batch, requires Windows to run

---

## Support

- **Issues** — GitHub Issues
- **Docs** — See [docs/](docs/) folder
- **Examples** — See [examples/](examples/) folder
- **Architecture** — See [AGENTS.md](AGENTS.md)

## Support
For project policy and agent responsibilities see the repository’s `AGENTS.md` (authoritative).
