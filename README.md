# NightAuth Go SDK — Star us on GitHub ⭐

Official **NightAuth** Go SDK & Reference Client Suite (Modern Desktop UI in **Native Win32 GDI** and Interactive Terminal CLI). Protect your Go software against piracy, reverse engineering, and tampering with cloud licensing, HWID locking, realtime encrypted chat, and secure in-memory payload streaming in minutes.

Official Portal: [https://nightauth.org](https://nightauth.org)  
Documentation: [https://nightauth.org/docs](https://nightauth.org/docs)  
Technical Architecture Benchmark: [NightAuth vs KeyAuth vs Auth.gg](https://nightauth.org/blog/nightauth-vs-keyauth-authgg-comparison)

---

**Languages / Idiomas / Langues / Sprachen:**  
[English](README.md) • [Español](lang/es.md) • [Français](lang/fr.md) • [Deutsch](lang/de.md) • [Italiano](lang/it.md) • [Português](lang/pt.md) • [Русский](lang/ru.md) • [Türkçe](lang/tr.md) • [Polski](lang/pl.md) • [中文](lang/zh.md)

---

> [!NOTE]
> ### 📺 Video Tutorials & System Setup Guides
> Official video tutorials and detailed step-by-step guides explaining how to configure applications, plans, subscriptions, cloud variables, webhooks, and licenses in the NightAuth Dashboard **will be published soon** on our official portal and YouTube channel.

---

## 📸 Screenshots

### Desktop Application (Native Win32 GDI Modern Dark Theme)

| Start / Authentication Screen | Dashboard & Subscription Management |
| :---: | :---: |
| ![NightAuth Desktop Login](image/desktop.png) | ![NightAuth Desktop Dashboard](image/desktop-dashboard.png) |

### Console Application (Interactive Terminal CLI)

| Terminal Login Screen | Interactive Console Dashboard |
| :---: | :---: |
| ![NightAuth Console Login](image/console.png) | ![NightAuth Console Dashboard](image/console-dashboard.png) |

---

## 🛡️ Recommended Security Practices

* **Advanced Go Obfuscation:** Standard Go binaries retain runtime type reflections and symbol names. For production releases, build using **Garble** (`garble -literals -tiny build`) to strip type definitions, obfuscate control flow, and encrypt string literals.
* **Commercial Binary Packing:** Further protect the compiled binary with commercial protectors such as **VMProtect** or **Themida** to safeguard against IDA Pro, Ghidra, and x64dbg.
* **Memory Integrity & Anti-Debugging:** NightAuth includes built-in security protections detecting debuggers (`x64dbg`, `CheatEngine`, `Wireshark`, `HTTPDebugger`) and system proxies. Invoke `nightauth.CheckDebugger()` and `nightauth.CheckBlacklistedProcesses()` at startup for maximum protection (see [SecurityProtections](#24-securityprotections)).
* **In-Memory Secure Payload Execution:** When retrieving protected binary modules or shellcodes with `auth.DownloadFile()`, **never write bytes to physical storage or disk**. Execute payloads directly in memory (e.g. via virtual memory allocation and execution) and zero out buffers immediately after execution.

---

## 📦 Requirements & Building

Requires **Go 1.20** or later.

Native Win32 desktop application uses pure Windows system DLLs (`user32.dll`, `gdi32.dll`, `kernel32.dll`, `msimg32.dll`) without CGO dependencies.

### Build Console Client
```bash
cd console
go build -ldflags="-s -w" -o bin/NightAuth-Console.exe .
```

### Build Desktop Client (Windowed / No Console Window)
```bash
cd desktop
go build -ldflags="-s -w -H=windowsgui" -o bin/NightAuth-Desktop.exe .
```

---

## ⚙️ Configuration & Initialization

Configure your application credentials from the **Manage Apps** section of the NightAuth Dashboard:

```go
package main

import (
	"fmt"
	"os"

	"nightauth"
)

// Application Configuration (from NightAuth Dashboard)
const (
	AppName    = ""
	AppApiKey  = ""
	AppVersion = "1.0"
	AppToken   = ""
)

func main() {
	// Initialize API client
	auth := nightauth.NewApi(AppName, AppApiKey, AppVersion, AppToken, "")

	// Perform cryptographic handshake (Curve25519 + Ed25519 + AES-256-GCM)
	auth.Init()

	if !auth.Response.Success {
		fmt.Printf("Initialization failed: %s\n", auth.Response.Message)
		os.Exit(1)
	}

	fmt.Println("Connected to NightAuth cloud successfully!")
}
```

---

## 🚀 SDK Usage Guide

### 1. Check Session Validation & Blacklist Status

```go
// Validate active session status with the server
if !auth.Check() {
	fmt.Printf("Session invalid: %s\n", auth.Response.Message)
}

// Check if client HWID or IP address is blacklisted
hwid := nightauth.SecurityProtections{}.GetHardwareId()
if auth.CheckBlacklist(hwid) {
	fmt.Println("Access Denied: Machine or IP is blacklisted.")
	os.Exit(1)
}
```

### 2. Login with Username & Password (and 2FA)

```go
if auth.Login(username, password, totpCode) {
	fmt.Printf("Welcome back, %s!\n", auth.User.Username)
} else {
	fmt.Printf("Login failed: %s\n", auth.Response.Message)
}
```

### 3. Register with License Key

```go
if auth.Register(username, password, licenseKey, email) {
	fmt.Printf("Registration successful: %s\n", auth.User.Username)
} else {
	fmt.Printf("Registration failed: %s\n", auth.Response.Message)
}
```

### 4. Authenticate with License Key Only (and 2FA)

```go
// Login with license key only
if auth.License(licenseKey) {
	fmt.Printf("License valid! Logged in as: %s\n", auth.User.Username)
} else {
	fmt.Printf("License authentication failed: %s\n", auth.Response.Message)
}

// Login with license key and 2FA TOTP code
if auth.LicenseWithTotp(licenseKey, "123456") {
	fmt.Printf("License verified with TOTP! Logged in as: %s\n", auth.User.Username)
}
```

### 5. Upgrade Account / Claim License Key

```go
if auth.Upgrade(username, upgradeKey) {
	fmt.Println("Subscription upgraded successfully!")
} else {
	fmt.Printf("Upgrade failed: %s\n", auth.Response.Message)
}
```

### 6. Reset Hardware ID (HWID Lock)

```go
if auth.ResetHwid(username, password, totpCode) {
	fmt.Println("HWID reset successfully. You can now login on this device.")
} else {
	fmt.Printf("HWID reset failed: %s\n", auth.Response.Message)
}
```

### 7. Forgot Password / Account Recovery

```go
// Step 1: Request password reset email
if auth.ForgotPassword(email) {
	fmt.Println("Password reset link sent to your email.")
}

// Step 2: Confirm password reset with token
if auth.ResetPassword(resetToken, newPassword) {
	fmt.Println("Password updated successfully!")
}
```

### 8. User Profile Data (`auth.User`)

```go
user := auth.User
fmt.Printf("Username:    %s\n", user.Username)
fmt.Printf("Email:       %s\n", user.Email)
fmt.Printf("Hardware ID: %s\n", user.Hwid)
fmt.Printf("IP Address:  %s\n", user.Ip)
fmt.Printf("Created At:  %s\n", user.CreatedAt)
fmt.Printf("Last Login:  %s\n", user.LastLoginAt)

for _, sub := range user.Subscriptions {
	fmt.Printf("Plan: %s | Expiry: %s | Key: %s\n", sub.Subscription, sub.Expiry, sub.Key)
}
```

### 9. Cloud Variables (Global)

```go
motd := auth.GetGlobalVar("server_motd")
fmt.Printf("Cloud Variable (Global): %s\n", motd)
```

### 10. User Variables

```go
// Set a custom user variable
auth.SetUserVar("custom_theme", "dark_neon")

// Retrieve the stored variable
val := auth.GetUserVar("custom_theme")
fmt.Printf("User Custom Variable: %s\n", val)

// List all user variables
if auth.GetUserVariables() {
	rawVars, _ := auth.Response.Payload["variables"].([]interface{})
	for _, item := range rawVars {
		if m, ok := item.(map[string]interface{}); ok {
			fmt.Printf("Variable %s = %s\n", m["key"], m["value"])
		}
	}
}
```

### 11. Application Audit Logs

```go
auth.Log("Suspicious integrity violation intercepted.")
```

### 12. Ban User & Blacklist Machine

```go
auth.Ban("Anti-Cheat integrity failure: Memory modification attempt.")
```

### 13. Server-Sided Webhooks & File Upload

```go
// Standard JSON webhook
payload := map[string]interface{}{
	"event": "alert_triggered",
	"level": "critical",
}
result := auth.TriggerWebhook("7kR0UedlVI", payload)
if result != nil && result.Success {
	fmt.Printf("Webhook delivered successfully (HTTP %d)\n", result.StatusCode)
}

// Multipart webhook with file attachment
fileBytes := []byte("example payload data")
fileResult := auth.TriggerWebhookFile("7kR0UedlVI", fileBytes, "report.txt", payload)
if fileResult != nil && fileResult.Success {
	fmt.Printf("File webhook delivered (HTTP %d)\n", fileResult.StatusCode)
}
```

### 14. Secure In-Memory File Download

```go
fileBytes := auth.DownloadFile("385624")
if len(fileBytes) > 0 {
	fmt.Printf("Downloaded %d bytes securely into RAM.\n", len(fileBytes))
	// Execute or process directly from memory without writing to physical disk
}
```

### 15. Realtime Chat Channels

```go
// Send message to channel
auth.SendChatMessage("Hello community from Go!", "general")

// Fetch latest messages from channel
messages := auth.GetChatMessages("general")
for _, msg := range messages {
	fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.Author, msg.Message)
}

// Fetch paginated messages before a specific ID
paged := auth.GetChatMessagesPage("general", 0, true)
for _, msg := range paged {
	fmt.Printf("#%s [%s] %s: %s\n", msg.Id, msg.Timestamp, msg.Author, msg.Message)
}
```

### 16. Realtime Event Streaming (Server-Sent Events)

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err := auth.SubscribeEvents(ctx, "", func(evt nightauth.EventFrame) {
	fmt.Printf("[EVENT #%s] %s\n", evt.ID, evt.Data)
})
if err != nil {
	fmt.Printf("Event stream ended: %v\n", err)
}
```

### 17. Two-Factor Authentication (2FA)

```go
// Setup TOTP 2FA
setupData := auth.Setup2FA()
if setupData != nil {
	fmt.Printf("Secret: %s | URI: %s\n", setupData.Secret, setupData.Uri)
}

// Confirm 2FA setup with authenticator code
confirmData := auth.Confirm2FA("123456")
if confirmData != nil && confirmData.Success {
	fmt.Println("2FA Enabled successfully! Backup Codes:")
	for _, code := range confirmData.BackupCodes {
		fmt.Printf(" - %s\n", code)
	}
}

// Disable 2FA
if auth.Disable2FA(password, "123456") {
	fmt.Println("2FA disabled successfully.")
}
```

### 18. Session Lifecycle & Logout

```go
// Send periodic heartbeat to keep session active
auth.Heartbeat()

// Revoke session on cloud and terminate listener
auth.Logout()
auth.Close()
```

### 19. Change Account Username

```go
if auth.ChangeUsername("new_alias", "current_password") {
	fmt.Printf("Username successfully updated to: %s\n", auth.User.Username)
} else {
	fmt.Printf("Failed to change username: %s\n", auth.Response.Message)
}
```

### 20. Application Telemetry & Online Users

```go
// Fetch application statistics
stats := auth.GetStats()
if stats != nil {
	fmt.Printf("Users: %d total, %d active, %d banned\n", stats.TotalUsers, stats.ActiveUsers, stats.BannedUsers)
	fmt.Printf("Licenses: %d total, %d used, %d unused\n", stats.TotalLicenses, stats.UsedLicenses, stats.UnusedLicenses)
	fmt.Printf("Server Version: %s\n", stats.Version)
}

// Fetch list of currently online usernames
onlineUsers := auth.GetOnlineUsers()
fmt.Printf("Online users (%d): %v\n", len(onlineUsers), onlineUsers)
```

### 21. Verify File Hash

```go
// Verify SHA-256 binary hash against server-approved release hashes
isValid := auth.VerifyFileHash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
if isValid {
	fmt.Println("File integrity verified successfully.")
} else {
	fmt.Println("Integrity check failed: file hash mismatch or unauthorized build.")
}
```

### 22. Access Tickets

```go
// Request an authenticated scoped access ticket (e.g. for secure asset streaming or download)
ticket := auth.AccessTicket("file_download", "file_abc123")
if ticket != "" {
	fmt.Printf("Access ticket granted: %s\n", ticket)
} else {
	fmt.Printf("Failed to issue access ticket: %s\n", auth.Response.Message)
}
```

### 23. Realtime Event Callbacks & Background Listener

Configure reactive event handlers and start the background event dispatcher:

```go
auth.OnSessionInvalidated = func(reason string) {
	fmt.Printf("Session terminated: %s\n", reason)
	os.Exit(0)
}
auth.OnAppBanned = func(reason string) {
	fmt.Printf("Application banned: %s\n", reason)
}
auth.OnAppStatusChanged = func(status string) {
	fmt.Printf("App status updated: %s\n", status)
}
auth.OnVersionUpdated = func(version string) {
	fmt.Printf("New version available: %s\n", version)
}
auth.OnChatMessageReceived = func(msg nightauth.ChatMessage, channel string) {
	fmt.Printf("[%s] %s: %s\n", channel, msg.Author, msg.Message)
}

// Start background SSE event listener goroutine
auth.StartRealtimeEventListener()

// Stop listener when no longer needed
auth.StopRealtimeEventListener()
```

### 24. `SecurityProtections`

Static and instance anti-tampering security helpers:

```go
// Run automatic checks (system HTTP proxy detection)
nightauth.RunSecurityChecks()

// Detect active debuggers (local or remote Win32 debuggers)
nightauth.CheckDebugger()

// Detect blacklisted analysis tools (x64dbg, Cheat Engine, IDA, Wireshark, Fiddler)
// Set NIGHTAUTH_IGNORE_SECURITY=1 during local testing to bypass
nightauth.CheckBlacklistedProcesses()

// Retrieve stable hardware fingerprint identifier
hwid := nightauth.GetHardwareId()
fmt.Printf("Device HWID: %s\n", hwid)
```

### 25. Standalone Options & Custom Server Keys

```go
// Initialize with custom API endpoint and optional insecure loopback for testing
auth := nightauth.NewApiWithOptions(
	"my-app",
	"my-api-key",
	"1.0",
	"",
	"https://nightauth.org/api/v1",
	false,
)

// Configure custom request timeout
auth.SetRequestTimeout(15 * time.Second)

// Configure pinned Ed25519 signing key and X25519 sealed box key (64 hex characters)
err := auth.SetServerKeys(
	"cdb2285db481666ce40003aadb874c8751dbeb9a0672a37f740dfe881ecdc742",
	"b773fc4f6cf1a0938e78e714c3d84cb386584bea24728923008471a6cd203213",
)
```

### 26. Signed Self-Updates

Verify and download tamper-evident application releases signed with Ed25519:

```go
// One-shot verified update download with replay rollback protection
targetPath, err := nightauth.SignedUpdateInstaller.Download(
	"https://releases.example.com/manifest.json",
	"YOUR_PINNED_ED25519_PUBLIC_KEY_HEX",
	"my-app",
	"updates/my-app.exe",
	5*time.Minute,
)
if err != nil {
	fmt.Printf("Update failed: %v\n", err)
} else {
	fmt.Printf("Verified update installed at: %s\n", targetPath)
}

// Or drive the pipeline manually:
verified, err := nightauth.VerifySignedUpdate(manifestEnvelopeJSON, pinnedPublicKeyHex, "my-app", minimumSequence)
if err == nil {
	fmt.Printf("Verified release %s v%s (seq %d)\n", verified.Application, verified.Version, verified.Sequence)
	installed, err := nightauth.DownloadSignedUpdate(verified, "updates/my-app.exe", 5*time.Minute)
}
```

## 🖥️ Run the Console and Desktop Clients

Launch either reference client from the `source` directory:

```bash
# Interactive Terminal Console Client
cd console
go run .

# Modern Pure Win32 Desktop Client (Windows)
cd desktop
go run .
```

---

---

## 🆚 Why NightAuth? (KeyAuth & Auth.gg Comparison)

| Feature | NightAuth | KeyAuth | Auth.gg |
| :--- | :---: | :---: | :---: |
| **Cryptography** | **End-to-End Curve25519 (X25519) + ChaCha20-Poly1305** | Static AES-CBC (Vulnerable to MITM & Proxy Interception) | Legacy HTTPS / Plain JSON |
| **Hardware ID (HWID)** | **Multi-Component Cryptographic Hash (Anti-Spoofing)** | Raw WMI Strings (Trivially Spoofed) | Simple Registry / MAC |
| **Edge Network Latency** | **< 35ms Global Anycast Edge Network** | Single-Origin VPS (Frequent Network Outages) | Unmaintained / High Latency |
| **Reseller System & Bots** | **Built-in Multi-Tier Credit Balances + Discord & Telegram Bots** | Manual API / Third-Party Scripts | Incomplete / Deprecated |
| **In-Memory Streaming** | **Zero-Disk Decryption with Direct Memory Execution & Buffer Wiping** | Disk Writes Required | Not Supported |
| **Internationalization** | **10 Languages Supported Out of the Box across Documentation & SDKs** | English Only | English Only |

Read the full architectural benchmark: [NightAuth vs KeyAuth vs Auth.gg](https://nightauth.org/blog/nightauth-vs-keyauth-authgg-comparison)

---

## 📄 License

This project is licensed under the **Elastic License 2.0**.

* You may not provide the software to third parties as a hosted or managed service.
* You may not move, change, disable, or circumvent the license key functionality.
* You may not alter, remove, or obscure any licensing or copyright notices.

See the [LICENSE](LICENSE) file for the full legal text.
