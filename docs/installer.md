# Windows Installer (NSIS)

We provide an NSIS script to package `fin.exe` for Windows.

## Files
- `scripts/fin_installer.nsi` — NSIS script

## Prerequisites
- Go build of the CLI: `go build -o fin.exe ./cmd/fin`
- NSIS installed (provides `makensis`)

## Build the installer
```
# from repo root
makensis scripts/fin_installer.nsi
```
This produces `Fin-v1.0.0-Setup.exe` in the current directory.

## Install behavior
- Installs to `C:\Program Files\Fin` (64-bit Program Files)
- Adds the install directory to PATH for current user and machine (append)
- Creates a Start Menu shortcut `Fin.lnk`
- Uninstall removes binaries and shortcut; PATH cleanup is **not** performed automatically to avoid clobbering user edits (remove manually if desired).

## Using the installer
Run the generated `Fin-v1.0.0-Setup.exe` and follow prompts. After install, open a new shell and run `fin version` to verify.
