# Fin CLI

`fin` is the command-line interface to compile, check, and format Fin source.

## Commands
- `fin build <file.fin> [-o output.bat]`
  - Compile `.fin` to Batch. Default output is `<file>.bat` beside input.
- `fin check <file.fin>`
  - Run lexer+parser+sema+generator validation, no output file.
- `fin ast <file.fin>`
  - Print AST for debugging.
- `fin fmt [-w] <file.fin>`
  - Canonical formatting. Without `-w`, prints to stdout; with `-w`, overwrites the file.
- `fin version`
  - Print CLI version.

## Exit codes
- `0` on success
- `1` on compile/format errors
- `2` on usage errors (wrong flags/args)

## Formatting specifics
- 4-space indentation
- Blank line between top-level functions
- Range spacing: `for i in 1 .. 3`
- Stable map/list spacing: `{k: v}`, `[a, b]`

## Error reporting
Errors are printed as `error: file.fin:line:col message` with ANSI red unless `NO_COLOR` is set.

## Examples
```sh
fin build hello.fin -o hello.bat
fin check hello.fin
fin ast hello.fin
fin fmt hello.fin      # stdout
fin fmt -w hello.fin   # in-place
```
