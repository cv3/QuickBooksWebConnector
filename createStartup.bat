REM Batch file to put a shortcut of the qbwcServer.exe into the startup folder

set SCRIPT="%TEMP%\%RANDOM%-%RANDOM%-%RANDOM%-%RANDOM%.vbs"

echo Set oWS = WScript.CreateObject("WScript.Shell") >> %SCRIPT%
echo sLinkFile = "%appdata%\Microsoft\Windows\Start Menu\Programs\Startup\qbwcServer.lnk" >> %SCRIPT%
echo Set oLink = oWS.CreateShortcut(sLinkFile) >> %SCRIPT%
echo oLink.TargetPath = "%~dp0\qbwcServer.exe" >> %SCRIPT%
echo oLink.WorkingDirectory = "%~dp0" >> %SCRIPT%
echo oLink.WindowStyle = 7 >> %SCRIPT%  && 7=Minimized 0=Maximized  4=Normal 
echo oLink.Description = "QBWC Server" >> %SCRIPT%
echo oLink.Save >> %SCRIPT%

cscript /nologo %SCRIPT%
del %SCRIPT%