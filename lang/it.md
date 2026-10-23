# NightAuth Go SDK — Dacci una stella su GitHub ⭐

SDK ufficiale di **NightAuth** per Go e suite di client di riferimento (interfaccia grafica moderna in **Win32 GDI nativo** e CLI interattiva da terminale). Proteggi il tuo software in Go da pirateria, reverse engineering e manomissioni con gestione delle licenze cloud, blocco HWID, chat crittografata in tempo reale e download sicuro di payload in memoria in pochi minuti.

Portale Ufficiale: [https://nightauth.org](https://nightauth.org)

---

**Languages / Idiomas / Langues / Sprachen:**  
[English](../README.md) • [Español](es.md) • [Français](fr.md) • [Deutsch](de.md) • [Italiano](it.md) • [Português](pt.md) • [Русский](ru.md) • [Türkçe](tr.md) • [Polski](pl.md) • [中文](zh.md)

---

> [!NOTE]
> ### 📺 Tutorial Video e Guide di Configurazione
> I video tutorial ufficiali e le guide dettagliate passo-passo che spiegano come configurare applicazioni, piani, abbonamenti, variabili cloud, webhook e licenze nella dashboard di NightAuth **saranno pubblicati a breve** sul nostro portale ufficiale e canale YouTube.

---

## 📸 Screenshot

### Applicazione Desktop (Tema Scuro Moderno in Win32 GDI Nativo)

| Schermata Iniziale / Autenticazione | Dashboard & Gestione Licenze |
| :---: | :---: |
| ![NightAuth Desktop Login](../image/desktop.png) | ![NightAuth Desktop Dashboard](../image/desktop-dashboard.png) |

### Applicazione Console (CLI Interattiva da Terminale)

| Schermata Login da Terminale | Dashboard Interattiva da Console |
| :---: | :---: |
| ![NightAuth Console Login](../image/console.png) | ![NightAuth Console Dashboard](../image/console-dashboard.png) |

---

## 🛡️ Pratiche di Sicurezza Raccomandate

* **Offuscamento Avanzato in Go:** I binari Go compilati standard conservano nomi di simboli e tipi. Per la produzione, compila con **Garble** (`garble -literals -tiny build`) per rimuovere definizioni di tipo e cifrare le stringhe.
* **Protezioni Binarie Commerciali:** Proteggi il file compilato con soluzioni come **VMProtect** o **Themida** per difenderti da IDA Pro, Ghidra e x64dbg.
* **Integrità Memoria & Anti-Debug:** NightAuth include verifiche contro debugger (`x64dbg`, `CheatEngine`, `Wireshark`, `HTTPDebugger`) e proxy di sistema. Esegui `nightauth.CheckDebugger()` e `nightauth.CheckBlacklistedProcesses()` all'avvio (vedi [SecurityProtections](#24-securityprotections)).
* **Esecuzione Sicura in Memoria:** Durante il download di payload con `auth.DownloadFile()`, **non scrivere mai i byte sul disco fisso**. Esegui i dati direttamente in memoria RAM e azzera i buffer immediatamente dopo l'uso.

---

## 📦 Requisiti & Compilazione

Richiede **Go 1.20** o versioni successive.

L'applicazione desktop nativa Win32 utilizza le librerie di sistema Windows (`user32.dll`, `gdi32.dll`, `kernel32.dll`, `msimg32.dll`) senza dipendenze CGO esterne.

### Compila il Client Console
```bash
cd console
go build -ldflags="-s -w" -o bin/NightAuth-Console.exe .
```

### Compila il Client Desktop (Senza Finestra Console)
```bash
cd desktop
go build -ldflags="-s -w -H=windowsgui" -o bin/NightAuth-Desktop.exe .
```

---

## ⚙️ Configurazione & Inizializzazione

Configura le credenziali della tua applicazione dalla dashboard di NightAuth:

```go
package main

import (
	"fmt"
	"os"

	"nightauth"
)

// Configurazione Applicazione
const (
	AppName    = ""
	AppApiKey  = ""
	AppVersion = "1.0"
	AppToken   = ""
)

func main() {
	// Inizializza client API
	auth := nightauth.NewApi(AppName, AppApiKey, AppVersion, AppToken, "")

	// Esegui handshake crittografico (Curve25519 + Ed25519 + AES-256-GCM)
	auth.Init()

	if !auth.Response.Success {
		fmt.Printf("Inizializzazione fallita: %s\n", auth.Response.Message)
		os.Exit(1)
	}

	fmt.Println("Connesso con successo al cloud NightAuth!")
}
```

---

## 🚀 Guida all'Uso del SDK

### 1. Convalida Sessione & Verifica Blacklist

```go
if !auth.Check() {
	fmt.Printf("Sessione non valida: %s\n", auth.Response.Message)
}

hwid := nightauth.SecurityProtections{}.GetHardwareId()
if auth.CheckBlacklist(hwid) {
	fmt.Println("Accesso Negato: Macchina o IP bloccato.")
	os.Exit(1)
}
```

### 2. Accesso con Nome Utente e Password (e 2FA)

```go
if auth.Login(username, password, totpCode) {
	fmt.Printf("Bentornato, %s!\n", auth.User.Username)
} else {
	fmt.Printf("Accesso fallito: %s\n", auth.Response.Message)
}
```

### 3. Registrazione con Chiave di Licenza

```go
if auth.Register(username, password, licenseKey, email) {
	fmt.Printf("Registrazione completata: %s\n", auth.User.Username)
} else {
	fmt.Printf("Registrazione fallita: %s\n", auth.Response.Message)
}
```

### 4. Autenticazione Solo Tramite Chiave di Licenza (e 2FA)

```go
// Accesso con sola chiave di licenza
if auth.License(licenseKey) {
	fmt.Printf("Licenza valida! Accesso effettuato come: %s\n", auth.User.Username)
} else {
	fmt.Printf("Autenticazione licenza fallita: %s\n", auth.Response.Message)
}

// Accesso con chiave di licenza e codice 2FA TOTP
if auth.LicenseWithTotp(licenseKey, "123456") {
	fmt.Printf("Licenza verificata con TOTP! Accesso effettuato come: %s\n", auth.User.Username)
}
```

### 5. Aggiornamento Piano / Riscatto Nuova Chiave

```go
if auth.Upgrade(username, upgradeKey) {
	fmt.Println("Abbonamento aggiornato con successo!")
} else {
	fmt.Printf("Aggiornamento fallito: %s\n", auth.Response.Message)
}
```

### 6. Ripristino Hardware ID (Blocco HWID)

```go
if auth.ResetHwid(username, password, totpCode) {
	fmt.Println("HWID ripristinato con successo. Ora puoi accedere da questo dispositivo.")
} else {
	fmt.Printf("Ripristino HWID fallito: %s\n", auth.Response.Message)
}
```

### 7. Recupero Account / Password Dimenticata

```go
if auth.ForgotPassword(email) {
	fmt.Println("Link di recupero password inviato via e-mail.")
}

if auth.ResetPassword(resetToken, newPassword) {
	fmt.Println("Password aggiornata con successo!")
}
```

### 8. Dati del Profilo Utente (`auth.User`)

```go
user := auth.User
fmt.Printf("Utente:       %s\n", user.Username)
fmt.Printf("E-mail:       %s\n", user.Email)
fmt.Printf("HWID:         %s\n", user.Hwid)
fmt.Printf("IP Client:    %s\n", user.Ip)
fmt.Printf("Creato il:    %s\n", user.CreatedAt)
fmt.Printf("Ultimo login: %s\n", user.LastLoginAt)

for _, sub := range user.Subscriptions {
	fmt.Printf("Piano: %s | Scadenza: %s | Chiave: %s\n", sub.Subscription, sub.Expiry, sub.Key)
}
```

### 9. Variabili Cloud (Globali)

```go
motd := auth.GetGlobalVar("server_motd")
fmt.Printf("Variabile Cloud: %s\n", motd)
```

### 10. Variabili Utente

```go
auth.SetUserVar("custom_theme", "dark_neon")

val := auth.GetUserVar("custom_theme")
fmt.Printf("Variabile Utente: %s\n", val)

if auth.GetUserVariables() {
	rawVars, _ := auth.Response.Payload["variables"].([]interface{})
	for _, item := range rawVars {
		if m, ok := item.(map[string]interface{}); ok {
			fmt.Printf("Variabile %s = %s\n", m["key"], m["value"])
		}
	}
}
```

### 11. Registri di Telemetria e Audit

```go
auth.Log("Rilevato tentativo di violazione integrità.")
```

### 12. Blocco Utente e Macchina (Killswitch / Ban)

```go
auth.Ban("Violazione Anti-Cheat: Tentativo di alterazione memoria.")
```

### 13. Webhook Lato Server & Invio File

```go
payload := map[string]interface{}{
	"event": "alert_triggered",
	"level": "critical",
}
result := auth.TriggerWebhook("7kR0UedlVI", payload)
if result != nil && result.Success {
	fmt.Printf("Webhook inviato con successo (HTTP %d)\n", result.StatusCode)
}

fileBytes := []byte("dati del report")
fileResult := auth.TriggerWebhookFile("7kR0UedlVI", fileBytes, "report.txt", payload)
if fileResult != nil && fileResult.Success {
	fmt.Printf("Webhook con file inviato (HTTP %d)\n", fileResult.StatusCode)
}
```

### 14. Download Sicuro in Memoria

```go
fileBytes := auth.DownloadFile("385624")
if len(fileBytes) > 0 {
	fmt.Printf("Scaricati %d byte in memoria RAM in modo sicuro.\n", len(fileBytes))
}
```

### 15. Canali di Chat in Tempo Reale

```go
auth.SendChatMessage("Ciao a tutti da Go!", "general")

messages := auth.GetChatMessages("general")
for _, msg := range messages {
	fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.Author, msg.Message)
}

paged := auth.GetChatMessagesPage("general", 0, true)
for _, msg := range paged {
	fmt.Printf("#%s [%s] %s: %s\n", msg.Id, msg.Timestamp, msg.Author, msg.Message)
}
```

### 16. Eventi in Tempo Reale (Server-Sent Events)

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err := auth.SubscribeEvents(ctx, "", func(evt nightauth.EventFrame) {
	fmt.Printf("[EVENTO #%s] %s\n", evt.ID, evt.Data)
})
```

### 17. Autenticazione a Due Fattori (2FA)

```go
setupData := auth.Setup2FA()
if setupData != nil {
	fmt.Printf("Segreto: %s | URI: %s\n", setupData.Secret, setupData.Uri)
}

confirmData := auth.Confirm2FA("123456")
if confirmData != nil && confirmData.Success {
	fmt.Println("2FA Attivato con successo! Codici di riserva:")
	for _, code := range confirmData.BackupCodes {
		fmt.Printf(" - %s\n", code)
	}
}

if auth.Disable2FA(password, "123456") {
	fmt.Println("2FA disattivato con successo.")
}
```

### 18. Ciclo di Vita Sessione & Disconnessione

```go
auth.Heartbeat()

auth.Logout()
auth.Close()
```

### 19. Cambia Nome Utente

```go
if auth.ChangeUsername("nuovo_alias", "password_attuale") {
	fmt.Printf("Nome utente aggiornato con successo a: %s\n", auth.User.Username)
} else {
	fmt.Printf("Errore durante il cambio nome utente: %s\n", auth.Response.Message)
}
```

### 20. Statistiche dell'Applicazione e Utenti Online

```go
// Recupera le metriche aggregate dell'applicazione
stats := auth.GetStats()
if stats != nil {
	fmt.Printf("Utenti: %d totali, %d attivi, %d banditi\n", stats.TotalUsers, stats.ActiveUsers, stats.BannedUsers)
	fmt.Printf("Licenze: %d totali, %d usate, %d non usate\n", stats.TotalLicenses, stats.UsedLicenses, stats.UnusedLicenses)
	fmt.Printf("Versione server: %s\n", stats.Version)
}

// Recupera l'elenco degli utenti attualmente online
onlineUsers := auth.GetOnlineUsers()
fmt.Printf("Utenti online (%d): %v\n", len(onlineUsers), onlineUsers)
```

### 21. Verifica Hash del File

```go
// Verifica l'hash SHA-256 del binario rispetto ai valori approvati dal server
isValid := auth.VerifyFileHash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
if isValid {
	fmt.Println("Integrità del file verificata con successo.")
} else {
	fmt.Println("Controllo di integrità fallito: hash non corrispondente.")
}
```

### 22. Ticket di Accesso

```go
// Richiedi un ticket di accesso a breve termine con ambito
ticket := auth.AccessTicket("file_download", "file_abc123")
if ticket != "" {
	fmt.Printf("Ticket di accesso rilasciato: %s\n", ticket)
} else {
	fmt.Printf("Errore emissione ticket di accesso: %s\n", auth.Response.Message)
}
```

### 23. Callback Reattivi & Listener in Background

Configura handler di eventi reattivi e avvia il dispatcher in background:

```go
auth.OnSessionInvalidated = func(reason string) {
	fmt.Printf("Sessione terminata: %s\n", reason)
	os.Exit(0)
}
auth.OnAppBanned = func(reason string) {
	fmt.Printf("Applicazione bandita: %s\n", reason)
}
auth.OnAppStatusChanged = func(status string) {
	fmt.Printf("Stato app aggiornato: %s\n", status)
}
auth.OnVersionUpdated = func(version string) {
	fmt.Printf("Nuova versione disponibile: %s\n", version)
}
auth.OnChatMessageReceived = func(msg nightauth.ChatMessage, channel string) {
	fmt.Printf("[%s] %s: %s\n", channel, msg.Author, msg.Message)
}

// Avvia listener SSE in background
auth.StartRealtimeEventListener()

// Arresta listener quando non più necessario
auth.StopRealtimeEventListener()
```

### 24. `SecurityProtections`

Utilità statiche e di istanza per la protezione anti-manomissione:

```go
// Esegui controlli automatici (rilevamento proxy HTTP di sistema)
nightauth.RunSecurityChecks()

// Rileva debugger attivi (debugger Win32 locali o remoti)
nightauth.CheckDebugger()

// Rileva strumenti di analisi in blacklist (x64dbg, Cheat Engine, IDA, Wireshark, Fiddler)
// Imposta NIGHTAUTH_IGNORE_SECURITY=1 nei test locali per ignorare
nightauth.CheckBlacklistedProcesses()

// Ottieni identificatore hardware stabile del dispositivo
hwid := nightauth.GetHardwareId()
fmt.Printf("HWID dispositivo: %s\n", hwid)
```

### 25. Opzioni Standalone & Chiavi Server Personalizzate

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

### 26. Aggiornamenti Firmati Crittograficamente

Verifica e scarica release protette contro le manomissioni firmate con Ed25519:

```go
// Download verificato in unico passaggio con protezione anti-rollback
targetPath, err := nightauth.SignedUpdateInstaller.Download(
	"https://releases.example.com/manifest.json",
	"CHIAVE_PUBBLICA_PINNED_ED25519_HEX",
	"my-app",
	"updates/my-app.exe",
	5*time.Minute,
)
if err != nil {
	fmt.Printf("Aggiornamento fallito: %v\n", err)
} else {
	fmt.Printf("Aggiornamento verificato installato in: %s\n", targetPath)
}

// Oppure gestisci il flusso manualmente:
verified, err := nightauth.VerifySignedUpdate(manifestEnvelopeJSON, pinnedPublicKeyHex, "my-app", minimumSequence)
if err == nil {
	fmt.Printf("Release verificata %s v%s (seq %d)\n", verified.Application, verified.Version, verified.Sequence)
	installed, err := nightauth.DownloadSignedUpdate(verified, "updates/my-app.exe", 5*time.Minute)
}
```

## 🖥️ Esecuzione dei Client Console e Desktop

Avvia uno dei client di riferimento dalla directory `source`:

```bash
# Interactive Terminal Console Client
cd console
go run .

# Modern Pure Win32 Desktop Client (Windows)
cd desktop
go run .
```

---

## 📄 Licenza

Questo progetto è distribuito sotto licenza **Elastic License 2.0**.

* Non è consentito fornire il software a terzi come servizio gestito o in hosting.
* Non è consentito rimuovere, modificare o aggirare la funzionalità delle licenze.
* Non è consentito alterare, rimuovere o oscurare gli avvisi di copyright o licenza.

Consulta il file [LICENSE](../LICENSE) per il testo legale completo.
