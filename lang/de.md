# NightAuth Go SDK — Gib uns einen Stern auf GitHub ⭐

Offizielles **NightAuth** Go SDK & Referenz-Client-Suite (Modernes Desktop-UI in **nativem Win32 GDI** und interaktive Terminal-CLI). Schütze deine Go-Software in wenigen Minuten vor Piraterie, Reverse Engineering und Manipulation durch Cloud-Lizenzierung, HWID-Bindung, echtzeitverschlüsselten Chat und sichere In-Memory-Nutzlasten.

Offizielles Portal: [https://nightauth.org](https://nightauth.org)

---

**Languages / Idiomas / Langues / Sprachen:**  
[English](../README.md) • [Español](es.md) • [Français](fr.md) • [Deutsch](de.md) • [Italiano](it.md) • [Português](pt.md) • [Русский](ru.md) • [Türkçe](tr.md) • [Polski](pl.md) • [中文](zh.md)

---

> [!NOTE]
> ### 📺 Video-Tutorials & System-Setup-Anleitungen
> Offizielle Video-Tutorials und Schritt-für-Schritt-Anleitungen zur Konfiguration von Anwendungen, Plänen, Abonnements, Cloud-Variablen, Webhooks und Lizenzen im NightAuth-Dashboard **werden in Kürze** auf unserem offiziellen Portal und YouTube-Kanal veröffentlicht.

---

## 📸 Screenshots

### Desktop-Anwendung (Natives Win32 GDI Modern Dark Theme)

| Start / Authentifizierungs-Bildschirm | Dashboard & Lizenzverwaltung |
| :---: | :---: |
| ![NightAuth Desktop Login](../image/desktop.png) | ![NightAuth Desktop Dashboard](../image/desktop-dashboard.png) |

### Konsolen-Anwendung (Interaktive Terminal-CLI)

| Terminal-Login-Bildschirm | Interaktives Konsolen-Dashboard |
| :---: | :---: |
| ![NightAuth Console Login](../image/console.png) | ![NightAuth Console Dashboard](../image/console-dashboard.png) |

---

## 🛡️ Empfohlene Sicherheitspraktiken

* **Erweiterte Go-Obfuskation:** Standard-Go-Binärdateien behalten Typen und Symbole. Kompiliere Releases mit **Garble** (`garble -literals -tiny build`), um Typdefinitionen zu entfernen und Strings zu verschlüsseln.
* **Kommerzielle Binary-Packer:** Schütze Binärdateien mit **VMProtect** oder **Themida**, um Analysen durch IDA Pro, Ghidra und x64dbg zu erschweren.
* **Speicherintegrität & Anti-Debugging:** NightAuth enthält Prüfungen gegen Debugger (`x64dbg`, `CheatEngine`, `Wireshark`, `HTTPDebugger`) und System-Proxys. Rufe `nightauth.CheckDebugger()` und `nightauth.CheckBlacklistedProcesses()` beim Start auf (siehe [SecurityProtections](#24-securityprotections)).
* **Sichere Ausführung im Arbeitsspeicher:** Beim Download geschützter Binärdateien mit `auth.DownloadFile()` **niemals Bytes auf die Festplatte schreiben**. Führe Payloads direkt im RAM aus und lösche Puffer sofort nach Gebrauch.

---

## 📦 Anforderungen & Kompilierung

Erfordert **Go 1.20** oder neuer.

Die native Win32-Desktop-Anwendung verwendet Windows-System-DLLs (`user32.dll`, `gdi32.dll`, `kernel32.dll`, `msimg32.dll`) ohne CGO-Abhängigkeiten.

### Konsolen-Client kompilieren
```bash
cd console
go build -ldflags="-s -w" -o bin/NightAuth-Console.exe .
```

### Desktop-Client kompilieren (Ohne Konsolenfenster)
```bash
cd desktop
go build -ldflags="-s -w -H=windowsgui" -o bin/NightAuth-Desktop.exe .
```

---

## ⚙️ Konfiguration & Initialisierung

Konfiguriere deine Anmeldedaten im NightAuth-Dashboard:

```go
package main

import (
	"fmt"
	"os"

	"nightauth"
)

// Anwendungs-Konfiguration
const (
	AppName    = ""
	AppApiKey  = ""
	AppVersion = "1.0"
	AppToken   = ""
)

func main() {
	// API-Client initialisieren
	auth := nightauth.NewApi(AppName, AppApiKey, AppVersion, AppToken, "")

	// Kryptografischen Handshake durchführen (Curve25519 + Ed25519 + AES-256-GCM)
	auth.Init()

	if !auth.Response.Success {
		fmt.Printf("Initialisierung fehlgeschlagen: %s\n", auth.Response.Message)
		os.Exit(1)
	}

	fmt.Println("Erfolgreich mit der NightAuth Cloud verbunden!")
}
```

---

## 🚀 SDK-Anleitung

### 1. Sitzungsvalidierung & Blacklist-Prüfung

```go
if !auth.Check() {
	fmt.Printf("Sitzung ungültig: %s\n", auth.Response.Message)
}

hwid := nightauth.SecurityProtections{}.GetHardwareId()
if auth.CheckBlacklist(hwid) {
	fmt.Println("Zugriff verweigert: Gerät oder IP gesperrt.")
	os.Exit(1)
}
```

### 2. Login mit Benutzername und Passwort (und 2FA)

```go
if auth.Login(username, password, totpCode) {
	fmt.Printf("Willkommen zurück, %s!\n", auth.User.Username)
} else {
	fmt.Printf("Anmeldung fehlgeschlagen: %s\n", auth.Response.Message)
}
```

### 3. Registrierung mit Lizenzschlüssel

```go
if auth.Register(username, password, licenseKey, email) {
	fmt.Printf("Erfolgreich registriert: %s\n", auth.User.Username)
} else {
	fmt.Printf("Registrierung fehlgeschlagen: %s\n", auth.Response.Message)
}
```

### 4. Authentifizierung nur mit Lizenzschlüssel (und 2FA)

```go
// Nur mit Lizenzschlüssel anmelden
if auth.License(licenseKey) {
	fmt.Printf("Lizenz gültig! Angemeldet als: %s\n", auth.User.Username)
} else {
	fmt.Printf("Lizenz-Authentifizierung fehlgeschlagen: %s\n", auth.Response.Message)
}

// Mit Lizenzschlüssel und 2FA-TOTP-Code anmelden
if auth.LicenseWithTotp(licenseKey, "123456") {
	fmt.Printf("Lizenz mit TOTP verifiziert! Angemeldet als: %s\n", auth.User.Username)
}
```

### 5. Plan-Upgrade / Lizenzschlüssel einlösen

```go
if auth.Upgrade(username, upgradeKey) {
	fmt.Println("Abonnement erfolgreich aufgewertet!")
} else {
	fmt.Printf("Upgrade fehlgeschlagen: %s\n", auth.Response.Message)
}
```

### 6. Hardware-ID zurücksetzen (HWID-Sperre)

```go
if auth.ResetHwid(username, password, totpCode) {
	fmt.Println("HWID erfolgreich zurückgesetzt. Du kannst dich jetzt auf diesem Gerät anmelden.")
} else {
	fmt.Printf("HWID-Reset fehlgeschlagen: %s\n", auth.Response.Message)
}
```

### 7. Passwort vergessen / Kontowiederherstellung

```go
if auth.ForgotPassword(email) {
	fmt.Println("Wiederherstellungslink per E-Mail versendet.")
}

if auth.ResetPassword(resetToken, newPassword) {
	fmt.Println("Passwort erfolgreich aktualisiert!")
}
```

### 8. Benutzerprofildaten (`auth.User`)

```go
user := auth.User
fmt.Printf("Benutzer:     %s\n", user.Username)
fmt.Printf("E-Mail:       %s\n", user.Email)
fmt.Printf("HWID:         %s\n", user.Hwid)
fmt.Printf("Client-IP:    %s\n", user.Ip)
fmt.Printf("Erstellt am:  %s\n", user.CreatedAt)
fmt.Printf("Letzter Login:%s\n", user.LastLoginAt)

for _, sub := range user.Subscriptions {
	fmt.Printf("Plan: %s | Ablauf: %s | Schlüssel: %s\n", sub.Subscription, sub.Expiry, sub.Key)
}
```

### 9. Cloud-Variablen (Global)

```go
motd := auth.GetGlobalVar("server_motd")
fmt.Printf("Cloud-Variable: %s\n", motd)
```

### 10. Benutzer-Variablen

```go
auth.SetUserVar("custom_theme", "dark_neon")

val := auth.GetUserVar("custom_theme")
fmt.Printf("Benutzer-Variable: %s\n", val)

if auth.GetUserVariables() {
	rawVars, _ := auth.Response.Payload["variables"].([]interface{})
	for _, item := range rawVars {
		if m, ok := item.(map[string]interface{}); ok {
			fmt.Printf("Variable %s = %s\n", m["key"], m["value"])
		}
	}
}
```

### 11. Anwendungs-Audit-Logs

```go
auth.Log("Integritätsverletzung festgestellt.")
```

### 12. Benutzer & Maschine sperren (Killswitch / Ban)

```go
auth.Ban("Anti-Cheat-Verletzung: Manipulationsversuch im Arbeitsspeicher.")
```

### 13. Server-Webhooks & Datei-Upload

```go
payload := map[string]interface{}{
	"event": "alert_triggered",
	"level": "critical",
}
result := auth.TriggerWebhook("7kR0UedlVI", payload)
if result != nil && result.Success {
	fmt.Printf("Webhook erfolgreich übermittelt (HTTP %d)\n", result.StatusCode)
}

fileBytes := []byte("Berichtsdaten")
fileResult := auth.TriggerWebhookFile("7kR0UedlVI", fileBytes, "report.txt", payload)
if fileResult != nil && fileResult.Success {
	fmt.Printf("Datei-Webhook übermittelt (HTTP %d)\n", fileResult.StatusCode)
}
```

### 14. Sicherer Download im Arbeitsspeicher

```go
fileBytes := auth.DownloadFile("385624")
if len(fileBytes) > 0 {
	fmt.Printf("%d Bytes sicher in den RAM geladen.\n", len(fileBytes))
}
```

### 15. Echtzeit-Chat-Kanäle

```go
auth.SendChatMessage("Hallo von Go!", "general")

messages := auth.GetChatMessages("general")
for _, msg := range messages {
	fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.Author, msg.Message)
}

paged := auth.GetChatMessagesPage("general", 0, true)
for _, msg := range paged {
	fmt.Printf("#%s [%s] %s: %s\n", msg.Id, msg.Timestamp, msg.Author, msg.Message)
}
```

### 16. Echtzeit-Ereignisse (Server-Sent Events)

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err := auth.SubscribeEvents(ctx, "", func(evt nightauth.EventFrame) {
	fmt.Printf("[EREIGNIS #%s] %s\n", evt.ID, evt.Data)
})
```

### 17. Zwei-Faktor-Authentifizierung (2FA)

```go
setupData := auth.Setup2FA()
if setupData != nil {
	fmt.Printf("Geheimnis: %s | URI: %s\n", setupData.Secret, setupData.Uri)
}

confirmData := auth.Confirm2FA("123456")
if confirmData != nil && confirmData.Success {
	fmt.Println("2FA erfolgreich aktiviert! Backup-Codes:")
	for _, code := range confirmData.BackupCodes {
		fmt.Printf(" - %s\n", code)
	}
}

if auth.Disable2FA(password, "123456") {
	fmt.Println("2FA erfolgreich deaktiviert.")
}
```

### 18. Sitzungslebenszyklus & Abmeldung

```go
auth.Heartbeat()

auth.Logout()
auth.Close()
```

### 19. Benutzernamen Ändern

```go
if auth.ChangeUsername("neuer_alias", "aktuelles_passwort") {
	fmt.Printf("Benutzername erfolgreich aktualisiert auf: %s\n", auth.User.Username)
} else {
	fmt.Printf("Fehler beim Ändern des Benutzernamens: %s\n", auth.Response.Message)
}
```

### 20. Anwendungsstatistiken & Online-Benutzer

```go
// Anwendungsstatistiken abrufen
stats := auth.GetStats()
if stats != nil {
	fmt.Printf("Benutzer: %d gesamt, %d aktiv, %d gesperrt\n", stats.TotalUsers, stats.ActiveUsers, stats.BannedUsers)
	fmt.Printf("Lizenzen: %d gesamt, %d genutzt, %d ungenutzt\n", stats.TotalLicenses, stats.UsedLicenses, stats.UnusedLicenses)
	fmt.Printf("Server-Version: %s\n", stats.Version)
}

// Liste der aktuell angemeldeten Benutzer abrufen
onlineUsers := auth.GetOnlineUsers()
fmt.Printf("Online-Benutzer (%d): %v\n", len(onlineUsers), onlineUsers)
```

### 21. Datei-Hash Verifizieren

```go
// SHA-256-Hash der Binärdatei gegen serverseitig registrierte Hashes validieren
isValid := auth.VerifyFileHash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
if isValid {
	fmt.Println("Dateiintegrität erfolgreich verifiziert.")
} else {
	fmt.Println("Integritätsprüfung fehlgeschlagen: Hash stimmt nicht überein.")
}
```

### 22. Zugriffstickets

```go
// Kurzlebiges Zugriffsticket mit Bereich anfordern
ticket := auth.AccessTicket("file_download", "file_abc123")
if ticket != "" {
	fmt.Printf("Zugriffsticket erteilt: %s\n", ticket)
} else {
	fmt.Printf("Fehler bei Ticketausstellung: %s\n", auth.Response.Message)
}
```

### 23. Reaktive Callbacks & Hintergrund-Listener

Konfiguriere reaktive Event-Handler und starte den Hintergrund-Dispatcher:

```go
auth.OnSessionInvalidated = func(reason string) {
	fmt.Printf("Sitzung beendet: %s\n", reason)
	os.Exit(0)
}
auth.OnAppBanned = func(reason string) {
	fmt.Printf("App gesperrt: %s\n", reason)
}
auth.OnAppStatusChanged = func(status string) {
	fmt.Printf("App-Status aktualisiert: %s\n", status)
}
auth.OnVersionUpdated = func(version string) {
	fmt.Printf("Neue Version verfügbar: %s\n", version)
}
auth.OnChatMessageReceived = func(msg nightauth.ChatMessage, channel string) {
	fmt.Printf("[%s] %s: %s\n", channel, msg.Author, msg.Message)
}

// Hintergrund-SSE-Listener starten
auth.StartRealtimeEventListener()

// Listener bei Bedarf beenden
auth.StopRealtimeEventListener()
```

### 24. `SecurityProtections`

Statische Schutzprogramme gegen Manipulation und Debugger:

```go
// Automatische Prüfungen ausführen (Erkennung von HTTP-System-Proxys)
nightauth.RunSecurityChecks()

// Aktive Debugger erkennen (lokale oder entfernte Win32-Debugger)
nightauth.CheckDebugger()

// Analysewerkzeuge auf Sperrliste erkennen (x64dbg, Cheat Engine, IDA, Wireshark, Fiddler)
// NIGHTAUTH_IGNORE_SECURITY=1 beim lokalen Testen setzen
nightauth.CheckBlacklistedProcesses()

// Stabile Hardware-ID des Geräts abrufen
hwid := nightauth.GetHardwareId()
fmt.Printf("Geräte-HWID: %s\n", hwid)
```

### 25. Erweiterte Optionen & Server-Schlüssel

```go
auth := nightauth.NewApiWithOptions(
	"my-app",
	"my-api-key",
	"1.0",
	"",
	"https://nightauth.org/api/v1",
	false,
)

auth.SetRequestTimeout(15 * time.Second)

err := auth.SetServerKeys(
	"cdb2285db481666ce40003aadb874c8751dbeb9a0672a37f740dfe881ecdc742",
	"b773fc4f6cf1a0938e78e714c3d84cb386584bea24728923008471a6cd203213",
)
```

### 26. Kryptografisch signierte Selbst-Updates

Verifiziert und lädt manipulationssichere Anwendungs-Releases herunter, die mit Ed25519 signiert sind:

```go
// Verifizierter One-Shot-Download mit Replay-Schutz
targetPath, err := nightauth.SignedUpdateInstaller.Download(
	"https://releases.example.com/manifest.json",
	"DEIN_PINNED_ED25519_SCHLUESSEL_HEX",
	"my-app",
	"updates/my-app.exe",
	5*time.Minute,
)
if err != nil {
	fmt.Printf("Update fehlgeschlagen: %v\n", err)
} else {
	fmt.Printf("Verifiziertes Update installiert unter: %s\n", targetPath)
}

// Oder den Prozess manuell steuern:
verified, err := nightauth.VerifySignedUpdate(manifestEnvelopeJSON, pinnedPublicKeyHex, "my-app", minimumSequence)
if err == nil {
	fmt.Printf("Verifiziertes Release %s v%s (Seq %d)\n", verified.Application, verified.Version, verified.Sequence)
	installed, err := nightauth.DownloadSignedUpdate(verified, "updates/my-app.exe", 5*time.Minute)
}
```

## 🖥️ Konsolen- und Desktop-Clients starten

Starte einen der Referenz-Clients aus dem `source`-Verzeichnis:

```bash
# Interactive Terminal Console Client
cd console
go run .

# Modern Pure Win32 Desktop Client (Windows)
cd desktop
go run .
```

---

## 📄 Lizenz

Dieses Projekt ist unter der **Elastic License 2.0** lizenziert.

* Du darfst die Software nicht als gehosteten oder verwalteten Dienst für Dritte bereitstellen.
* Du darfst die Lizenzschlüsselfunktionalität nicht verschieben, ändern, deaktivieren oder umgehen.
* Du darfst keine Lizenz- oder Urheberrechtshinweise ändern, entfernen oder unkenntlich machen.

Den vollständigen Text findest du in der [LICENSE](../LICENSE)-Datei.
