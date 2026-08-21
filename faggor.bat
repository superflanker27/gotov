@echo off
powershell -Command "Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/superflanker27/gotov/refs/heads/main/gotovenv.ps1' -OutFile 'C:\Users\WDAGUtilityAccount\Desktop\gotov.psi'"
start powershell -NoExit -ExecutionPolicy Bypass -File "C:\Users\WDAGUtilityAccount\Desktop\gotov.ps1"
