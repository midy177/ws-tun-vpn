# Define the name of the installer
OutFile "wtvc_gui_Installer.exe"

# Define the directory where the application will be installed
InstallDir $PROGRAMFILES\MyFyneApp

# Request application privileges for Windows Vista/7/8/10
RequestExecutionLevel user

# Define the default section
Section "MainSection" SEC01

  # Create the application directory
  CreateDirectory $INSTDIR

  # Copy the application executable
  SetOutPath $INSTDIR
  File "wtvc_gui.exe"
  File "wintun.dll"

  # Create a shortcut in the start menu
  CreateShortcut "$SMPROGRAMS\MyFyneApp.lnk" "$INSTDIR\wtvc_gui.exe"

  # Create a shortcut on the desktop
  CreateShortcut "$DESKTOP\MyFyneApp.lnk" "$INSTDIR\wtvc_gui.exe"

SectionEnd

# Define uninstallation script
Section "Uninstall"

  # Remove the application files
  Delete "$INSTDIR\wtvc_gui.exe"
  Delete "$INSTDIR\wintun.dll"

  # Remove shortcuts
  Delete "$SMPROGRAMS\MyFyneApp.lnk"
  Delete "$DESKTOP\MyFyneApp.lnk"

  # Remove the installation directory
  RMDir "$INSTDIR"

SectionEnd
