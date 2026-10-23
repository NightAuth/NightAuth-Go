# NightAuth Go SDK Console

Interactive Terminal CLI client implemented with ANSI truecolor and 256-color escape sequences.

## Features
- **Interactive Terminal UI:** ANSI-styled responsive menus and ASCII branding.
- **Full Protocol Coverage:** Authentication, license validation, cloud variables, encrypted file downloads, realtime chat, and SSE event streaming.

## Building and Running

### Run directly with Go:
```bash
go run .
```

### Compile Release Binary:
```bash
go build -ldflags="-s -w" -o bin/NightAuth-Console.exe .
```

## Configuration
Application credentials can be defined via the constants in `main.go` or passed as environment variables:
- `NIGHTAUTH_APP_NAME`
- `NIGHTAUTH_API_KEY`
- `NIGHTAUTH_APP_VERSION` (optional, defaults to `1.0`)
- `NIGHTAUTH_APP_TOKEN` (optional)
- `NIGHTAUTH_API_URL` (optional, defaults to `https://nightauth.org/api/v1`)
