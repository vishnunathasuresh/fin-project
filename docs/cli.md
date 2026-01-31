# Fin CLI Reference

**Version:** 1.0.0  
**Status:** Final

The `fin` command-line interface provides tools to compile, validate, analyze, and format Fin source code.

---

## Installation

### Option 1: Installer
Download from [Releases](../../releases)

### Option 2: Build from Source
```cmd
go build -o fin.exe ./cmd/fin
```
The binary will be in the current directory.

### Option 3: Add to PATH
After installing, `fin` should be in your PATH. Verify:
```cmd
fin version
```

---

## Commands

### build
Compile a Fin script to Windows Batch.

**Syntax:**
```
fin build <file.fin> [-o output.bat]
```

**Description:**
- Lexes, parses, analyzes, and generates Batch code
- Output path defaults to `<file>.bat` in same directory as input
- Overwrites output file without warning

**Examples:**
```cmd
fin build script.fin                    # → script.bat
fin build script.fin -o out.bat         # → out.bat
fin build examples/01_variables_echo.fin # → examples/01_variables_echo.bat
```

**Exit Code:**
- `0` on success
- `1` on semantic/generation errors
- `2` on usage error (missing file, etc.)

---

### check
Validate a Fin script without generating output.

**Syntax:**
```
fin check <file.fin>
```

**Description:**
- Runs lexer, parser, semantic analysis, and generator validation
- No `.bat` file produced
- Reports all errors found

**Examples:**
```cmd
fin check script.fin
fin check examples/02_arithmetic.fin
```

**Exit Code:**
- `0` if valid
- `1` if errors found
- `2` if usage error

---

### ast
Print the abstract syntax tree (AST) of a Fin script for debugging.

**Syntax:**
```
fin ast <file.fin>
```

**Description:**
- Parses input and outputs AST in human-readable form
- Shows node types, positions, and structure
- Useful for understanding how code is parsed

**Examples:**
```cmd
fin ast script.fin
fin ast examples/03_if_else.fin
```

**Example Output:**
```
Program @1:1
  SetStmt name=x @2:1
    value: NumberLit 10 @2:7
  IfStmt @4:1
    cond: BinaryExpr op=> @4:7
      left: IdentExpr x @4:4
      right: NumberLit 5 @4:9
    then:
      EchoStmt @5:5
        value: StringLit "yes" @5:10
```

**Exit Code:**
- `0` on success
- `1` on parse errors
- `2` on usage error

---

### fmt
Format Fin code to canonical style.

**Syntax:**
```
fin fmt [-w] <file.fin>
```

**Options:**
- `-w` — Write formatted output back to file (in-place)
- Without `-w` — Print to stdout

**Description:**
- Applies canonical formatting rules
- Idempotent (formatting twice produces same result)
- Helps maintain consistent style across projects

**Examples:**
```cmd
fin fmt script.fin          # Print formatted to console
fin fmt -w script.fin       # Update file in-place
fin fmt examples/*.fin      # Format all examples (with shell glob)
```

**Formatting Rules:**
- 4-space indentation
- Blank line before each top-level function
- Space around range operator: `for i in 1 .. 5`
- Maps/lists: `{k: v}`, `[a, b]`
- Spaces around operators: `a + b`, `!flag`
- Comma-space in lists/maps: `[a, b, c]`

**Exit Code:**
- `0` on success
- `1` on parse errors
- `2` on usage error (e.g., `-w` with stdin)

---

### version
Print the Fin compiler version.

**Syntax:**
```
fin version
```

**Output:**
```
v1.0.0
```

**Exit Code:**
- `0` always

---

## Global Options

### Environment Variables

**NO_COLOR**
Disable colored output in errors:
```cmd
set NO_COLOR=1
fin check script.fin
```

---

## Error Reporting

### Format
Errors are printed to `stderr` as:
```
error: <file>:<line>:<col> <message>
```

**Example:**
```
error: script.fin:10:5 undefined variable: foo
error: script.fin:15:1 function arity mismatch: expected 2 arguments, got 1
```

### Colors
- Red for errors (unless `NO_COLOR` is set)
- Line:column info for precise location

### Multiple Errors
All errors are reported in one pass (no stopping at first error).

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Compile/semantic/generation error |
| 2 | Usage error (bad args, missing file, etc.) |

---

## Usage Examples

### Basic Workflow
```cmd
# Check syntax
fin check my_script.fin

# Compile to Batch
fin build my_script.fin

# Run generated script
my_script.bat

# View AST for debugging
fin ast my_script.fin

# Format code
fin fmt -w my_script.fin
```

### Multiple Files
```cmd
# Build all examples
for %f in (examples\*.fin) do fin build %f

# Check all files
for %f in (*.fin) do fin check %f

# Format all files in-place
for /r . %f in (*.fin) do fin fmt -w %f
```

### Integration with CI/CD
```powershell
# PowerShell: Build and check all
Get-ChildItem -Filter *.fin | ForEach-Object {
    fin check $_.FullName
    if ($LASTEXITCODE -ne 0) { exit 1 }
    fin build $_.FullName
}
```

### Debugging
```cmd
# View AST for complex script
fin ast complex_script.fin > ast.txt
type ast.txt

# Check without compiling
fin check broken_script.fin

# See formatted output before applying
fin fmt script.fin  # preview
fin fmt -w script.fin  # apply
```

---

## Troubleshooting

### "fin: command not found"
- Ensure `fin.exe` is in PATH
- Or use full path: `C:\Program Files\Fin\fin.exe build script.fin`

### Error: "undefined variable"
- Check variable is defined before use
- Fin requires explicit `set` statements

### Error: "function arity mismatch"
- Count arguments in function call
- Must match function definition exactly
- No default parameters supported

### Batch error when running generated script
- Review generated `.bat` file
- Ensure input path exists
- Check `fin ast` to verify parsing

---

## Version Information
```cmd
fin version                          # v1.0.0
```

See [language.md](language.md) for language features and [README.md](../README.md) for project overview.
