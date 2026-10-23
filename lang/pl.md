# NightAuth Go SDK — Daj nam gwiazdkę na GitHubie ⭐

Oficjalny zestaw **NightAuth** Go SDK i aplikacje referencyjne (nowoczesny interfejs graficzny w **natywnym Win32 GDI** oraz interaktywny interfejs wiersza poleceń CLI). Chroń swoje oprogramowanie Go przed piractwem, inżynierią wsteczną i manipulacjami dzięki licencjonowaniu w chmurze, blokadzie HWID, szyfrowanemu czatowi w czasie rzeczywistym i bezpiecznemu przesyłaniu danych w pamięci RAM w kilka minut.

Oficjalny Portal: [https://nightauth.org](https://nightauth.org)

---

**Languages / Idiomas / Langues / Sprachen:**  
[English](../README.md) • [Español](es.md) • [Français](fr.md) • [Deutsch](de.md) • [Italiano](it.md) • [Português](pt.md) • [Русский](ru.md) • [Türkçe](tr.md) • [Polski](pl.md) • [中文](zh.md)

---

> [!NOTE]
> ### 📺 Samouczki wideo i przewodniki konfiguracji
> Oficjalne poradniki wideo i szczegółowe przewodniki krok po kroku wyjaśniające, jak konfigurować aplikacje, plany, subskrypcje, zmienne w chmurze, webhooki i licencje w panelu NightAuth **zostaną wkrótce opublikowane** na naszym oficjalnym portalu i kanale YouTube.

---

## 📸 Zrzuty Ekranu

### Aplikacja Desktopowa (Natywny Ciemny Motyw Win32 GDI)

| Ekran Początkowy / Uwierzytelnianie | Panel Główny & Zarządzanie Licencjami |
| :---: | :---: |
| ![NightAuth Desktop Login](../image/desktop.png) | ![NightAuth Desktop Dashboard](../image/desktop-dashboard.png) |

### Aplikacja Konsolowa (Interaktywne Terminalowe CLI)

| Ekran Logowania w Terminalu | Interaktywny Pulpit Konsoli |
| :---: | :---: |
| ![NightAuth Console Login](../image/console.png) | ![NightAuth Console Dashboard](../image/console-dashboard.png) |

---

## 🛡️ Zalecane Praktyki Bezpieczeństwa

* **Zaawansowana Obfuskacja Go:** Standardowe pliki binarne Go zachowują nazwy symboli i typy refleksji. W wersjach produkcyjnych kompiluj za pomocą **Garble** (`garble -literals -tiny build`), aby usunąć definicje typów i zaszyfrować literały tekstowe.
* **Komercyjne Pakery:** Zabezpiecz plik binarny za pomocą **VMProtect** lub **Themida**, aby chronić kod przed deasemblacją w IDA Pro, Ghidra i x64dbg.
* **Integralność Pamięci & Anti-Debug:** NightAuth zawiera wbudowane zabezpieczenia wykrywające debuggery (`x64dbg`, `CheatEngine`, `Wireshark`, `HTTPDebugger`) i systemowe proxy. Wywołaj `nightauth.CheckDebugger()` i `nightauth.CheckBlacklistedProcesses()` przy starcie programu (patrz [SecurityProtections](#24-securityprotections)).
* **Bezpieczne Wykonywanie w Pamięci RAM:** Pobierając zabezpieczone pliki binarne za pomocą `auth.DownloadFile()`, **nigdy nie zapisuj bajtów na dysku twardym**. Ładuj dane bezpośrednio do pamięci RAM i zeruj bufory natychmiast po użyciu.

---

## 📦 Wymagania i Kompilacja

Wymaga **Go 1.20** lub nowszego.

Natywna aplikacja desktopowa Win32 korzysta z bibliotek systemowych Windows (`user32.dll`, `gdi32.dll`, `kernel32.dll`, `msimg32.dll`) bez zależności CGO.

### Kompilacja Klienta Konsolowego
```bash
cd console
go build -ldflags="-s -w" -o bin/NightAuth-Console.exe .
```

### Kompilacja Klienta Desktopowego (Bez Okna Konsoli)
```bash
cd desktop
go build -ldflags="-s -w -H=windowsgui" -o bin/NightAuth-Desktop.exe .
```

---

## ⚙️ Konfiguracja i Inicjalizacja

Skonfiguruj dane uwierzytelniające w panelu NightAuth:

```go
package main

import (
	"fmt"
	"os"

	"nightauth"
)

// Konfiguracja Aplikacji
const (
	AppName    = ""
	AppApiKey  = ""
	AppVersion = "1.0"
	AppToken   = ""
)

func main() {
	// Inicjalizacja klienta API
	auth := nightauth.NewApi(AppName, AppApiKey, AppVersion, AppToken, "")

	// Wykonanie kryptograficznego uzgadniania (Curve25519 + Ed25519 + AES-256-GCM)
	auth.Init()

	if !auth.Response.Success {
		fmt.Printf("Błąd inicjalizacji: %s\n", auth.Response.Message)
		os.Exit(1)
	}

	fmt.Println("Pomyślnie połączono z chmurą NightAuth!")
}
```

---

## 🚀 Przewodnik po SDK

### 1. Walidacja Sesji & Sprawdzenie Czarnej Listy

```go
if !auth.Check() {
	fmt.Printf("Nieprawidłowa sesja: %s\n", auth.Response.Message)
}

hwid := nightauth.SecurityProtections{}.GetHardwareId()
if auth.CheckBlacklist(hwid) {
	fmt.Println("Odmowa dostępu: Urządzenie lub IP zablokowane.")
	os.Exit(1)
}
```

### 2. Logowanie Loginem i Hasłem (oraz 2FA)

```go
if auth.Login(username, password, totpCode) {
	fmt.Printf("Witaj ponownie, %s!\n", auth.User.Username)
} else {
	fmt.Printf("Błąd logowania: %s\n", auth.Response.Message)
}
```

### 3. Rejestracja z Kluczem Licencyjnym

```go
if auth.Register(username, password, licenseKey, email) {
	fmt.Printf("Rejestracja pomyślna: %s\n", auth.User.Username)
} else {
	fmt.Printf("Błąd rejestracji: %s\n", auth.Response.Message)
}
```

### 4. Uwierzytelnianie Samym Kluczem Licencji (i 2FA)

```go
// Logowanie tylko za pomocą klucza licencji
if auth.License(licenseKey) {
	fmt.Printf("Licencja poprawna! Zalogowano jako: %s\n", auth.User.Username)
} else {
	fmt.Printf("Błąd uwierzytelniania licencji: %s\n", auth.Response.Message)
}

// Logowanie za pomocą klucza licencji i kodu 2FA TOTP
if auth.LicenseWithTotp(licenseKey, "123456") {
	fmt.Printf("Licencja zweryfikowana z TOTP! Zalogowano jako: %s\n", auth.User.Username)
}
```

### 5. Uaktualnienie Planu / Nowy Klucz

```go
if auth.Upgrade(username, upgradeKey) {
	fmt.Println("Pomyślnie uaktualniono subskrypcję!")
} else {
	fmt.Printf("Błąd uaktualnienia: %s\n", auth.Response.Message)
}
```

### 6. Reset Identyfikatora Sprzętu (Blokada HWID)

```go
if auth.ResetHwid(username, password, totpCode) {
	fmt.Println("HWID zresetowany pomyślnie. Możesz zalogować się na tym urządzeniu.")
} else {
	fmt.Printf("Błąd resetowania HWID: %s\n", auth.Response.Message)
}
```

### 7. Odzyskiwanie Konta / Zapomniane Hasło

```go
if auth.ForgotPassword(email) {
	fmt.Println("Link do resetowania hasła wysłany na e-mail.")
}

if auth.ResetPassword(resetToken, newPassword) {
	fmt.Println("Hasło zaktualizowane pomyślnie!")
}
```

### 8. Dane Profilu Użytkownika (`auth.User`)

```go
user := auth.User
fmt.Printf("Użytkownik:  %s\n", user.Username)
fmt.Printf("E-mail:      %s\n", user.Email)
fmt.Printf("HWID:        %s\n", user.Hwid)
fmt.Printf("IP Klienta:  %s\n", user.Ip)
fmt.Printf("Utworzono:   %s\n", user.CreatedAt)
fmt.Printf("Ostatni log: %s\n", user.LastLoginAt)

for _, sub := range user.Subscriptions {
	fmt.Printf("Plan: %s | Wygasa: %s | Klucz: %s\n", sub.Subscription, sub.Expiry, sub.Key)
}
```

### 9. Zmienne Chmurowe (Globalne)

```go
motd := auth.GetGlobalVar("server_motd")
fmt.Printf("Zmienna w chmurze: %s\n", motd)
```

### 10. Zmienne Użytkownika

```go
auth.SetUserVar("custom_theme", "dark_neon")

val := auth.GetUserVar("custom_theme")
fmt.Printf("Zmienna użytkownika: %s\n", val)

if auth.GetUserVariables() {
	rawVars, _ := auth.Response.Payload["variables"].([]interface{})
	for _, item := range rawVars {
		if m, ok := item.(map[string]interface{}); ok {
			fmt.Printf("Zmienna %s = %s\n", m["key"], m["value"])
		}
	}
}
```

### 11. Dzienniki Zdarzeń i Telemetria

```go
auth.Log("Wykryto próbę naruszenia integralności.")
```

### 12. Blokada Użytkownika i Maszyny (Killswitch / Ban)

```go
auth.Ban("Naruszenie Anti-Cheat: Próba modyfikacji pamięci.")
```

### 13. Webhooki Serwera & Przesyłanie Plików

```go
payload := map[string]interface{}{
	"event": "alert_triggered",
	"level": "critical",
}
result := auth.TriggerWebhook("7kR0UedlVI", payload)
if result != nil && result.Success {
	fmt.Printf("Webhook dostarczony (HTTP %d)\n", result.StatusCode)
}

fileBytes := []byte("dane raportu")
fileResult := auth.TriggerWebhookFile("7kR0UedlVI", fileBytes, "report.txt", payload)
if fileResult != nil && fileResult.Success {
	fmt.Printf("Webhook z plikiem dostarczony (HTTP %d)\n", fileResult.StatusCode)
}
```

### 14. Bezpieczne Pobieranie do Pamięci RAM

```go
fileBytes := auth.DownloadFile("385624")
if len(fileBytes) > 0 {
	fmt.Printf("Pobrano %d bajtów bezpośrednio do pamięci RAM.\n", len(fileBytes))
}
```

### 15. Kanały Czatu w Czasie Rzeczywistym

```go
auth.SendChatMessage("Cześć wszystkim z Go!", "general")

messages := auth.GetChatMessages("general")
for _, msg := range messages {
	fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.Author, msg.Message)
}

paged := auth.GetChatMessagesPage("general", 0, true)
for _, msg := range paged {
	fmt.Printf("#%s [%s] %s: %s\n", msg.Id, msg.Timestamp, msg.Author, msg.Message)
}
```

### 16. Strumień Zdarzeń w Czasie Rzeczywistym (Server-Sent Events)

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err := auth.SubscribeEvents(ctx, "", func(evt nightauth.EventFrame) {
	fmt.Printf("[ZDARZENIE #%s] %s\n", evt.ID, evt.Data)
})
```

### 17. Uwierzytelnianie Dwuetapowe (2FA)

```go
setupData := auth.Setup2FA()
if setupData != nil {
	fmt.Printf("Sekret: %s | URI: %s\n", setupData.Secret, setupData.Uri)
}

confirmData := auth.Confirm2FA("123456")
if confirmData != nil && confirmData.Success {
	fmt.Println("2FA Włączone pomyślnie! Kody zapasowe:")
	for _, code := range confirmData.BackupCodes {
		fmt.Printf(" - %s\n", code)
	}
}

if auth.Disable2FA(password, "123456") {
	fmt.Println("2FA wyłączone pomyślnie.")
}
```

### 18. Cykl Życia Sesji & Wylogowanie

```go
auth.Heartbeat()

auth.Logout()
auth.Close()
```

### 19. Zmiana Nazwy Użytkownika

```go
if auth.ChangeUsername("nowy_alias", "aktualne_haslo") {
	fmt.Printf("Nazwa użytkownika pomyślnie zmieniona na: %s\n", auth.User.Username)
} else {
	fmt.Printf("Błąd zmiany nazwy użytkownika: %s\n", auth.Response.Message)
}
```

### 20. Statystyki Aplikacji i Użytkownicy Online

```go
// Pobierz statystyki aplikacji
stats := auth.GetStats()
if stats != nil {
	fmt.Printf("Użytkownicy: %d łącznie, %d aktywnych, %d zablokowanych\n", stats.TotalUsers, stats.ActiveUsers, stats.BannedUsers)
	fmt.Printf("Licencje: %d łącznie, %d użytych, %d nieużytych\n", stats.TotalLicenses, stats.UsedLicenses, stats.UnusedLicenses)
	fmt.Printf("Wersja serwera: %s\n", stats.Version)
}

// Pobierz listę użytkowników aktualnie online
onlineUsers := auth.GetOnlineUsers()
fmt.Printf("Użytkownicy online (%d): %v\n", len(onlineUsers), onlineUsers)
```

### 21. Weryfikacja Hasha Pliku

```go
// Weryfikacja hasha SHA-256 pliku binarnego z zatwierdzonymi hashami serwera
isValid := auth.VerifyFileHash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
if isValid {
	fmt.Println("Integralność pliku zweryfikowana pomyślnie.")
} else {
	fmt.Println("Błąd weryfikacji integralności: niezgodny hash pliku.")
}
```

### 22. Bilety Dostępu

```go
// Pobierz krótkotrwały bilet dostępu z zakresem
ticket := auth.AccessTicket("file_download", "file_abc123")
if ticket != "" {
	fmt.Printf("Bilet dostępu przyznany: %s\n", ticket)
} else {
	fmt.Printf("Błąd wydawania biletu dostępu: %s\n", auth.Response.Message)
}
```

### 23. Reaktywne Zdarzenia & Odbiornik w Tle

Skonfiguruj reaktywne procedury obsługi zdarzeń i uruchom dyspozytora w tle:

```go
auth.OnSessionInvalidated = func(reason string) {
	fmt.Printf("Sesja zakończona: %s\n", reason)
	os.Exit(0)
}
auth.OnAppBanned = func(reason string) {
	fmt.Printf("Aplikacja zablokowana: %s\n", reason)
}
auth.OnAppStatusChanged = func(status string) {
	fmt.Printf("Status aplikacji zaktualizowany: %s\n", status)
}
auth.OnVersionUpdated = func(version string) {
	fmt.Printf("Dostępna nowa wersja: %s\n", version)
}
auth.OnChatMessageReceived = func(msg nightauth.ChatMessage, channel string) {
	fmt.Printf("[%s] %s: %s\n", channel, msg.Author, msg.Message)
}

// Rozpocznij nasłuchiwanie zdarzeń SSE w tle
auth.StartRealtimeEventListener()

// Zatrzymaj nasłuchiwanie, gdy nie jest potrzebne
auth.StopRealtimeEventListener()
```

### 24. `SecurityProtections`

Statyczne i instancyjne mechanizmy ochrony przed manipulacją:

```go
// Uruchom automatyczne kontrole (wykrywanie proxy HTTP w systemie)
nightauth.RunSecurityChecks()

// Wykryj aktywne debuggery (lokalne i zdalne Win32)
nightauth.CheckDebugger()

// Wykryj zablokowane narzędzia analityczne (x64dbg, Cheat Engine, IDA, Wireshark, Fiddler)
// Ustaw NIGHTAUTH_IGNORE_SECURITY=1 podczas testów lokalnych
nightauth.CheckBlacklistedProcesses()

// Pobierz unikalny identyfikator sprzętowy (HWID)
hwid := nightauth.GetHardwareId()
fmt.Printf("HWID urządzenia: %s\n", hwid)
```

### 25. Opcje Zaawansowane & Własne Klucze Serwera

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

### 26. Kryptograficznie Podpisane Aktualizacje

Weryfikuj i pobieraj bezpieczne wydania aplikacji podpisane kluczem Ed25519:

```go
// Bezpośrednie, bezpieczne pobieranie aktualizacji z ochroną przed rollbackiem
targetPath, err := nightauth.SignedUpdateInstaller.Download(
	"https://releases.example.com/manifest.json",
	"TWOJ_KLUCZ_PUBLICZNY_ED25519_HEX",
	"my-app",
	"updates/my-app.exe",
	5*time.Minute,
)
if err != nil {
	fmt.Printf("Błąd aktualizacji: %v\n", err)
} else {
	fmt.Printf("Zweryfikowana aktualizacja zainstalowana w: %s\n", targetPath)
}

// Lub przeprowadź proces ręcznie:
verified, err := nightauth.VerifySignedUpdate(manifestEnvelopeJSON, pinnedPublicKeyHex, "my-app", minimumSequence)
if err == nil {
	fmt.Printf("Zweryfikowane wydanie %s v%s (sekwencja %d)\n", verified.Application, verified.Version, verified.Sequence)
	installed, err := nightauth.DownloadSignedUpdate(verified, "updates/my-app.exe", 5*time.Minute)
}
```

## 🖥️ Uruchamianie Klientów Konsoli i Desktopu

Uruchom dowolnego klienta referencyjnego z katalogu `source`:

```bash
# Interactive Terminal Console Client
cd console
go run .

# Modern Pure Win32 Desktop Client (Windows)
cd desktop
go run .
```

---

## 📄 Licencja

Projekt jest objęty licencją **Elastic License 2.0**.

* Nie możesz udostępniać oprogramowania stronom trzecim jako usługi hostowanej lub zarządzanej.
* Nie możesz modyfikować, usuwać ani obchodzić funkcji klucza licencyjnego.
* Nie możesz modyfikować, usuwać ani ukrywać powiadomień o prawach autorskich lub licencji.

Pełny tekst prawny znajduje się w pliku [LICENSE](../LICENSE).
