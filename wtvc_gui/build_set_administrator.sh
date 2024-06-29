#/bin/sh
# windows
go build -ldflags="-H windowsgui" -o your-program.exe your-program.go
mt -manifest app.manifest -outputresource:your-program.exe;#1

# mac
cd wtvc_gui
fyne package -os darwin -icon ico/src/favicon.ico
rm *.dmg
create-dmg --volname wtvc_gui \
           --window-pos 200 120 \
           --window-size 800 400 \
           --icon-size 100 \
           --icon wtvc_gui.app 200 190 \
           --hide-extension wtvc_gui.app \
           --app-drop-link  600 185 \
           wtvc_gui_setup.dmg \
           wtvc_gui.app
