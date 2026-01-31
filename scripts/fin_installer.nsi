; Fin NSIS Installer
; Requires NSIS (makensis). Builds an installer that places fin.exe into Program Files and updates PATH.

!include "LogicLib.nsh"

!define APP_NAME "Fin"
!define APP_EXE "..\\fin.exe"
!define APP_VERSION "v1.0.0"
!define INSTALL_DIR "$PROGRAMFILES64\${APP_NAME}"

OutFile "Fin-${APP_VERSION}-Setup.exe"
InstallDir "${INSTALL_DIR}"
RequestExecutionLevel admin
ShowInstDetails show

Page directory
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

Section "Install"
  SetOutPath "$INSTDIR"
  File "${APP_EXE}"

  ; Add to PATH for current user and machine (append if not already present)
  ReadEnvStr $0 "PATH"
  ${If} $0 == ""
    StrCpy $0 "$INSTDIR"
  ${Else}
    StrCpy $0 "$0;$INSTDIR"
  ${EndIf}
  WriteRegExpandStr HKLM "SYSTEM\CurrentControlSet\Control\Session Manager\Environment" "Path" $0
  WriteRegExpandStr HKCU "Environment" "Path" $0

  ; Uninstaller
  WriteUninstaller "$INSTDIR\Uninstall.exe"

  ; Start menu shortcut
  CreateShortcut "$SMPROGRAMS\${APP_NAME}.lnk" "$INSTDIR\${APP_EXE}"
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\${APP_EXE}"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"
  Delete "$SMPROGRAMS\${APP_NAME}.lnk"

  ; PATH cleanup is not performed to avoid clobbering user edits. Remove manually if needed.
SectionEnd
