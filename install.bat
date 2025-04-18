@echo off
echo.
echo  ______          _   _    _____      
echo ^|  ____^|        ^| ^| (_)  / ____^|     
echo ^| ^|__ ___  _ __ ^| ^|_ _  ^| ^|  __  ___ 
echo ^|  __/ _ \^| '_ \^| __^| ^| ^| ^| ^|_ ^|/ _ \
echo ^| ^| ^| (_) ^| ^| ^| ^| ^|_^| ^| ^| ^|__^| ^| (_) ^|
echo ^|_^|  \___/^|_^| ^|_^|\__^|_^|  \_____^|\___/
echo.
echo Fort-Go Installation Script for Windows
echo.

:: Check for admin privileges
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo Warning: Not running with administrator privileges.
    echo Some operations may fail, consider running this script as administrator.
    echo.
)

:: Check for Go installation
where go >nul 2>&1
if %errorLevel% neq 0 (
    echo Error: Go is not installed or not in your PATH.
    echo Please download and install Go from https://golang.org/dl/
    echo After installation, run this script again.
    goto :error
) else (
    for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
    echo ✅ Go is installed: %GO_VERSION%
)

:: Check for Git installation
where git >nul 2>&1
if %errorLevel% neq 0 (
    echo Error: Git is not installed or not in your PATH.
    echo Please download and install Git from https://git-scm.com/download/win
    echo After installation, run this script again.
    goto :error
) else (
    for /f "tokens=3" %%i in ('git --version') do set GIT_VERSION=%%i
    echo ✅ Git is installed: %GIT_VERSION%
)

:: Set installation directory
set INSTALL_DIR=%USERPROFILE%\.fort-go

:: Clone or update repository
echo.
echo Setting up Fort-Go repository...
if exist "%INSTALL_DIR%" (
    echo Fort-Go repository already exists. Updating...
    cd /d "%INSTALL_DIR%"
    git pull
) else (
    echo Cloning Fort-Go repository...
    git clone https://github.com/princebabou/fort-go.git "%INSTALL_DIR%"
    if %errorLevel% neq 0 (
        echo Error: Failed to clone repository
        goto :error
    )
    cd /d "%INSTALL_DIR%"
)

:: Build the application
echo.
echo Building Fort-Go...
go build -o fort.exe .\cmd\fort
if %errorLevel% neq 0 (
    echo Error: Build failed
    goto :error
)

:: Add to PATH
echo.
echo Installing Fort-Go...

:: Create bin directory if it doesn't exist
set BIN_DIR=%USERPROFILE%\bin
if not exist "%BIN_DIR%" (
    mkdir "%BIN_DIR%"
)

:: Copy the binary
copy /Y "%INSTALL_DIR%\fort.exe" "%BIN_DIR%\fort.exe" >nul
if %errorLevel% neq 0 (
    echo Error: Failed to copy Fort-Go binary
    goto :error
)

:: Add to user PATH if not already there
echo Updating PATH environment variable...
for /f "tokens=2*" %%a in ('reg query HKCU\Environment /v PATH 2^>nul ^| findstr PATH') do set CURRENT_PATH=%%b
if not defined CURRENT_PATH (
    setx PATH "%BIN_DIR%"
) else (
    echo %CURRENT_PATH% | findstr /C:"%BIN_DIR%" >nul
    if %errorLevel% neq 0 (
        setx PATH "%CURRENT_PATH%;%BIN_DIR%"
    ) else (
        echo %BIN_DIR% is already in your PATH
    )
)

echo.
echo Fort-Go has been successfully installed!
echo.
echo You will need to restart your command prompt to use the 'fort' command.
echo.
echo Example commands:
echo   fort scan -t example.com              # Perform a full scan
echo   fort exploit -t example.com           # Attempt safe exploitation
echo   fort report -i results.json -f pdf    # Generate a PDF report
echo.
echo For more information, run: fort --help
goto :end

:error
echo.
echo Installation failed. Please resolve the errors and try again.
exit /b 1

:end
echo.
echo Press any key to exit...
pause > nul 