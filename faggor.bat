@echo off
powershell -Command "Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/main/gotovenv.ps1' -OutFile 'C:\Users\WDAGUtilityAccount\Desktop\gotov.psi'"

start powershell -NoExit -ExecutionPolicy Bypass -File "C:\Users\WDAGUtilityAccount\Desktop\gotov.psi"
