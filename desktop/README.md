# NightAuth Go SDK Desktop Client

Native Windows 64-bit Desktop GUI client implemented in pure Go using Windows system APIs (`user32.dll`, `gdi32.dll`, `kernel32.dll`, `msimg32.dll`) without external CGO dependencies.

## Features
- **Native Pure Win32 GDI Rendering:** Custom dark-themed card UI, anti-aliased typographic hierarchy, and custom vector icons.
- **Authentication Flows:** Login with TOTP 2FA, License Key Only activation, User Registration with License, HWID Reset, Account Upgrades, and Password Recovery.
- **Console-less Build:** Run as a pure GUI application without launching an auxiliary cmd console window.

## Building and Running

### Run directly with Go:
```bash
go run .
```

### Compile windowed executable:
```bash
go build -ldflags="-s -w -H=windowsgui" -o bin/NightAuth-Desktop.exe .
```

## Configuration
Application credentials can be defined via the constants in `main.go` or passed as environment variables:
- `NIGHTAUTH_APP_NAME`
- `NIGHTAUTH_API_KEY`
- `NIGHTAUTH_APP_VERSION` (optional, defaults to `1.0`)
- `NIGHTAUTH_APP_TOKEN` (optional)
- `NIGHTAUTH_API_URL` (optional, defaults to `https://nightauth.org/api/v1`)
