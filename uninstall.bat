@echo off
echo.
echo  ______          _   _    _____      
echo ^|  ____^|        ^| ^| (_)  / ____^|     
echo ^| ^|__ ___  _ __ ^| ^|_ _  ^| ^|  __  ___ 
echo ^|  __/ _ \^| '_ \^| __^| ^| ^| ^| ^|_ ^|/ _ \
echo ^| ^| ^| (_) ^| ^| ^| ^| ^|_^| ^| ^| ^|__^| ^| (_) ^|
echo ^|_^|  \___/^|_^| ^|_^|\__^|_^|  \_____^|\___/
echo.
echo Fort-Go Uninstallation Script for Windows
echo.

:: Confirm uninstallation
echo Warning: This will remove Fort-Go from your system.
set /p CONFIRM="Are you sure you want to continue? (y/n): "
if /i not "%CONFIRM%" == "y" (
    echo Uninstallation cancelled.
    goto :end
)

:: Set installation directory
set INSTALL_DIR=%USERPROFILE%\.fort-go
set BIN_DIR=%USERPROFILE%\bin

:: Remove binary
if exist "%BIN_DIR%\fort.exe" (
    echo Removing Fort-Go binary...
    del /f /q "%BIN_DIR%\fort.exe" >nul 2>&1
    if exist "%BIN_DIR%\fort.exe" (
        echo Error: Failed to remove binary.
    ) else (
        echo Fort-Go binary removed successfully.
    )
) else (
    echo Fort-Go binary not found in %BIN_DIR%.
)

:: Remove from PATH
echo Updating PATH environment variable...
for /f "tokens=2*" %%a in ('reg query HKCU\Environment /v PATH 2^>nul ^| findstr PATH') do set CURRENT_PATH=%%b
if defined CURRENT_PATH (
    echo %CURRENT_PATH% | findstr /C:"%BIN_DIR%" >nul
    if %errorLevel% equ 0 (
        set NEW_PATH=%CURRENT_PATH:;%BIN_DIR%=%
        set NEW_PATH=%NEW_PATH:;%BIN_DIR%;=;%
        set NEW_PATH=%NEW_PATH:%BIN_DIR%;=%
        set NEW_PATH=%NEW_PATH:%BIN_DIR%=%
        setx PATH "%NEW_PATH%"
        echo Fort-Go removed from PATH.
    ) else (
        echo Fort-Go not found in PATH.
    )
)

:: Remove repository
if exist "%INSTALL_DIR%" (
    echo Removing Fort-Go repository...
    rmdir /s /q "%INSTALL_DIR%" >nul 2>&1
    if exist "%INSTALL_DIR%" (
        echo Error: Failed to remove repository.
    ) else (
        echo Fort-Go repository removed successfully.
    )
) else (
    echo Fort-Go repository not found.
)

echo.
echo Fort-Go has been successfully uninstalled!
echo.
echo Note: You may need to restart your command prompt for PATH changes to take effect.
goto :end

:end
echo.
echo Press any key to exit...
pause > nul 