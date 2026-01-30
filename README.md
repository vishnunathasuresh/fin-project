# Fin Documentation

Welcome to the Fin docs. Fin is a Fish-inspired DSL that compiles to deterministic Windows Batch scripts.

## What’s inside
- [Language](language.md) — syntax, literals, control flow, functions.
- [CLI](cli.md) — `fin` command usage and examples.
- [Rules & Architecture](rules.md) — source code rules, layering, generator contract.

## Quick start
1. Write a `.fin` file (see [language](language.md)).
2. Run `fin build file.fin -o file.bat` to produce batch.
3. Use `fin check file.fin` to lint/validate without emitting batch.
4. Use `fin fmt file.fin` (or `fin fmt -w file.fin`) for canonical formatting.

## Support
For project policy and agent responsibilities see the repository’s `AGENTS.md` (authoritative).
