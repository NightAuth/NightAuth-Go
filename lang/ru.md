# NightAuth Go SDK — Поставьте нам звезду на GitHub ⭐

Официальный **NightAuth** Go SDK и набор эталонных клиентов (современный графический интерфейс на **нативном Win32 GDI** и интерактивный консольный CLI). Защитите свое программное обеспечение на Go от пиратства, реверс-инжиниринга и взлома с помощью облачного лицензирования, привязки к HWID, зашифрованного чата в реальном времени и безопасной загрузки модулей в память за считанные минуты.

Официальный портал: [https://nightauth.org](https://nightauth.org)

---

**Languages / Idiomas / Langues / Sprachen:**  
[English](../README.md) • [Español](es.md) • [Français](fr.md) • [Deutsch](de.md) • [Italiano](it.md) • [Português](pt.md) • [Русский](ru.md) • [Türkçe](tr.md) • [Polski](pl.md) • [中文](zh.md)

---

> [!NOTE]
> ### 📺 Видеоуроки и руководства по настройке системы
> Официальные видеоуроки и пошаговые руководства по настройке приложений, тарифов, подписок, облачных переменных, вебхуков и лицензий в панели управления NightAuth **скоро будут опубликованы** на нашем официальном портале и YouTube-канале.

---

## 📸 Скриншоты

### Десктопное приложение (Нативная темная тема Win32 GDI)

| Экран входа / Аутентификация | Панель управления & Подписки |
| :---: | :---: |
| ![NightAuth Desktop Login](../image/desktop.png) | ![NightAuth Desktop Dashboard](../image/desktop-dashboard.png) |

### Консольное приложение (Интерактивный терминальный CLI)

| Терминальный экран входа | Интерактивная панель консоли |
| :---: | :---: |
| ![NightAuth Console Login](../image/console.png) | ![NightAuth Console Dashboard](../image/console-dashboard.png) |

---

## 🛡️ Рекомендуемые практики безопасности

* **Продвинутая обфускация Go:** Стандартные бинарные файлы Go сохраняют имена символов и данные рефлексии. Для релизов компилируйте с помощью **Garble** (`garble -literals -tiny build`), чтобы удалить определения типов и зашифровать строки.
* **Коммерческие протекторы:** Дополнительно защитите исполняемый файл с помощью **VMProtect** или **Themida** от анализа в IDA Pro, Ghidra и x64dbg.
* **Целостность памяти и антиотладка:** NightAuth включает встроенные проверки отладчиков (`x64dbg`, `CheatEngine`, `Wireshark`, `HTTPDebugger`) и системных прокси. Вызывайте `nightauth.CheckDebugger()` и `nightauth.CheckBlacklistedProcesses()` при запуске (см. [SecurityProtections](#24-securityprotections)).
* **Безопасное выполнение в оперативной памяти:** При загрузке защищенных модулей с помощью `auth.DownloadFile()` **никогда не сохраняйте байты на жесткий диск**. Загружайте и выполняйте код напрямую в RAM и очищайте буферы сразу после использования.

---

## 📦 Требования и сборка

Требуется **Go 1.20** или выше.

Нативное приложение Win32 использует стандартные системные библиотеки Windows (`user32.dll`, `gdi32.dll`, `kernel32.dll`, `msimg32.dll`) без сторонних зависимостей CGO.

### Сборка консольного клиента
```bash
cd console
go build -ldflags="-s -w" -o bin/NightAuth-Console.exe .
```

### Сборка десктопного клиента (Без окна консоли)
```bash
cd desktop
go build -ldflags="-s -w -H=windowsgui" -o bin/NightAuth-Desktop.exe .
```

---

## ⚙️ Конфигурация и инициализация

Настройте учетные данные приложения в панели управления NightAuth:

```go
package main

import (
	"fmt"
	"os"

	"nightauth"
)

// Конфигурация приложения
const (
	AppName    = ""
	AppApiKey  = ""
	AppVersion = "1.0"
	AppToken   = ""
)

func main() {
	// Инициализация API клиента
	auth := nightauth.NewApi(AppName, AppApiKey, AppVersion, AppToken, "")

	// Выполнение криптографического рукопожатия (Curve25519 + Ed25519 + AES-256-GCM)
	auth.Init()

	if !auth.Response.Success {
		fmt.Printf("Ошибка инициализации: %s\n", auth.Response.Message)
		os.Exit(1)
	}

	fmt.Println("Успешное подключение к облаку NightAuth!")
}
```

---

## 🚀 Руководство по использованию SDK

### 1. Проверка валидности сессии и черного списка

```go
if !auth.Check() {
	fmt.Printf("Сессия недействительна: %s\n", auth.Response.Message)
}

hwid := nightauth.SecurityProtections{}.GetHardwareId()
if auth.CheckBlacklist(hwid) {
	fmt.Println("Доступ запрещен: компьютер или IP заблокированы.")
	os.Exit(1)
}
```

### 2. Вход по логину и паролю (и 2FA)

```go
if auth.Login(username, password, totpCode) {
	fmt.Printf("С возвращением, %s!\n", auth.User.Username)
} else {
	fmt.Printf("Ошибка входа: %s\n", auth.Response.Message)
}
```

### 3. Регистрация с лицензионным ключом

```go
if auth.Register(username, password, licenseKey, email) {
	fmt.Printf("Регистрация успешна: %s\n", auth.User.Username)
} else {
	fmt.Printf("Ошибка регистрации: %s\n", auth.Response.Message)
}
```

### 4. Вход только по лицензионному ключу (и 2FA)

```go
// Вход только по лицензионному ключу
if auth.License(licenseKey) {
	fmt.Printf("Лицензия действительна! Вход выполнен как: %s\n", auth.User.Username)
} else {
	fmt.Printf("Ошибка аутентификации по лицензии: %s\n", auth.Response.Message)
}

// Вход по лицензионному ключу с кодом 2FA TOTP
if auth.LicenseWithTotp(licenseKey, "123456") {
	fmt.Printf("Лицензия подтверждена через TOTP! Вход выполнен как: %s\n", auth.User.Username)
}
```

### 5. Обновление тарифа / Активация ключа

```go
if auth.Upgrade(username, upgradeKey) {
	fmt.Println("Подписка успешно обновлена!")
} else {
	fmt.Printf("Ошибка обновления: %s\n", auth.Response.Message)
}
```

### 6. Сброс аппаратного идентификатора (HWID Lock)

```go
if auth.ResetHwid(username, password, totpCode) {
	fmt.Println("HWID сброшен. Теперь вы можете войти на этом устройстве.")
} else {
	fmt.Printf("Ошибка сброса HWID: %s\n", auth.Response.Message)
}
```

### 7. Восстановление аккаунта / Сброс пароля

```go
if auth.ForgotPassword(email) {
	fmt.Println("Ссылка для сброса пароля отправлена на почту.")
}

if auth.ResetPassword(resetToken, newPassword) {
	fmt.Println("Пароль успешно обновлен!")
}
```

### 8. Данные профиля пользователя (`auth.User`)

```go
user := auth.User
fmt.Printf("Пользователь: %s\n", user.Username)
fmt.Printf("Email:        %s\n", user.Email)
fmt.Printf("HWID:         %s\n", user.Hwid)
fmt.Printf("IP-адрес:     %s\n", user.Ip)
fmt.Printf("Регистрация:  %s\n", user.CreatedAt)
fmt.Printf("Вход:         %s\n", user.LastLoginAt)

for _, sub := range user.Subscriptions {
	fmt.Printf("План: %s | Истекает: %s | Ключ: %s\n", sub.Subscription, sub.Expiry, sub.Key)
}
```

### 9. Облачные переменные (Глобальные)

```go
motd := auth.GetGlobalVar("server_motd")
fmt.Printf("Облачная переменная: %s\n", motd)
```

### 10. Переменные пользователя

```go
auth.SetUserVar("custom_theme", "dark_neon")

val := auth.GetUserVar("custom_theme")
fmt.Printf("Переменная пользователя: %s\n", val)

if auth.GetUserVariables() {
	rawVars, _ := auth.Response.Payload["variables"].([]interface{})
	for _, item := range rawVars {
		if m, ok := item.(map[string]interface{}); ok {
			fmt.Printf("Переменная %s = %s\n", m["key"], m["value"])
		}
	}
}
```

### 11. Журнал аудита приложения

```go
auth.Log("Обнаружена подозрительная активность.")
```

### 12. Блокировка пользователя и ПК (Killswitch / Ban)

```go
auth.Ban("Нарушение безопасности античита: попытка модификации памяти.")
```

### 13. Серверные вебхуки и загрузка файлов

```go
payload := map[string]interface{}{
	"event": "alert_triggered",
	"level": "critical",
}
result := auth.TriggerWebhook("7kR0UedlVI", payload)
if result != nil && result.Success {
	fmt.Printf("Вебхук доставлен (HTTP %d)\n", result.StatusCode)
}

fileBytes := []byte("данные отчета")
fileResult := auth.TriggerWebhookFile("7kR0UedlVI", fileBytes, "report.txt", payload)
if fileResult != nil && fileResult.Success {
	fmt.Printf("Вебхук с файлом доставлен (HTTP %d)\n", fileResult.StatusCode)
}
```

### 14. Безопасная загрузка в оперативную память

```go
fileBytes := auth.DownloadFile("385624")
if len(fileBytes) > 0 {
	fmt.Printf("Загружено %d байт напрямую в RAM.\n", len(fileBytes))
}
```

### 15. Чаты в реальном времени

```go
auth.SendChatMessage("Привет сообществу из Go!", "general")

messages := auth.GetChatMessages("general")
for _, msg := range messages {
	fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.Author, msg.Message)
}

paged := auth.GetChatMessagesPage("general", 0, true)
for _, msg := range paged {
	fmt.Printf("#%s [%s] %s: %s\n", msg.Id, msg.Timestamp, msg.Author, msg.Message)
}
```

### 16. События в реальном времени (Server-Sent Events)

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err := auth.SubscribeEvents(ctx, "", func(evt nightauth.EventFrame) {
	fmt.Printf("[СОБЫТИЕ #%s] %s\n", evt.ID, evt.Data)
})
```

### 17. Двухфакторная аутентификация (2FA)

```go
setupData := auth.Setup2FA()
if setupData != nil {
	fmt.Printf("Секрет: %s | URI: %s\n", setupData.Secret, setupData.Uri)
}

confirmData := auth.Confirm2FA("123456")
if confirmData != nil && confirmData.Success {
	fmt.Println("2FA успешно активирована! Резервные коды:")
	for _, code := range confirmData.BackupCodes {
		fmt.Printf(" - %s\n", code)
	}
}

if auth.Disable2FA(password, "123456") {
	fmt.Println("2FA успешно отключена.")
}
```

### 18. Жизненный цикл сессии и выход

```go
auth.Heartbeat()

auth.Logout()
auth.Close()
```

### 19. Смена Имени Пользователя

```go
if auth.ChangeUsername("novyi_alias", "tekushchii_parol") {
	fmt.Printf("Имя пользователя успешно обновлено на: %s\n", auth.User.Username)
} else {
	fmt.Printf("Ошибка при смене имени пользователя: %s\n", auth.Response.Message)
}
```

### 20. Статистика Приложения и Пользователи Онлайн

```go
// Получение статистики приложения
stats := auth.GetStats()
if stats != nil {
	fmt.Printf("Пользователи: %d всего, %d активных, %d заблокированных\n", stats.TotalUsers, stats.ActiveUsers, stats.BannedUsers)
	fmt.Printf("Лицензии: %d всего, %d использовано, %d свободно\n", stats.TotalLicenses, stats.UsedLicenses, stats.UnusedLicenses)
	fmt.Printf("Версия сервера: %s\n", stats.Version)
}

// Получение списка пользователей онлайн
onlineUsers := auth.GetOnlineUsers()
fmt.Printf("Пользователи онлайн (%d): %v\n", len(onlineUsers), onlineUsers)
```

### 21. Проверка Хеша Файла

```go
// Проверка SHA-256 хеша бинарного файла против разрешенных сервером
isValid := auth.VerifyFileHash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
if isValid {
	fmt.Println("Целостность файла успешно подтверждена.")
} else {
	fmt.Println("Сбой проверки целостности: несовпадение хеша.")
}
```

### 22. Билеты Доступа

```go
// Запрос временного билета доступа с областью действия
ticket := auth.AccessTicket("file_download", "file_abc123")
if ticket != "" {
	fmt.Printf("Билет доступа получен: %s\n", ticket)
} else {
	fmt.Printf("Ошибка выдачи билета доступа: %s\n", auth.Response.Message)
}
```

### 23. Реактивные обработчики и фоновый слушатель

Настройте реактивные обработчики событий и запустите фоновый диспетчер:

```go
auth.OnSessionInvalidated = func(reason string) {
	fmt.Printf("Сессия завершена: %s\n", reason)
	os.Exit(0)
}
auth.OnAppBanned = func(reason string) {
	fmt.Printf("Приложение заблокировано: %s\n", reason)
}
auth.OnAppStatusChanged = func(status string) {
	fmt.Printf("Статус приложения обновлен: %s\n", status)
}
auth.OnVersionUpdated = func(version string) {
	fmt.Printf("Доступна новая версия: %s\n", version)
}
auth.OnChatMessageReceived = func(msg nightauth.ChatMessage, channel string) {
	fmt.Printf("[%s] %s: %s\n", channel, msg.Author, msg.Message)
}

// Запуск фонового слушателя SSE событий
auth.StartRealtimeEventListener()

// Остановка слушателя, когда он больше не требуется
auth.StopRealtimeEventListener()
```

### 24. `SecurityProtections`

Статические и объектные утилиты защиты от реверс-инжиниринга:

```go
// Автоматические проверки (обнаружение системного HTTP-прокси)
nightauth.RunSecurityChecks()

// Обнаружение активных отладчиков (локальных и удаленных Win32)
nightauth.CheckDebugger()

// Обнаружение заблокированных инструментов анализа (x64dbg, Cheat Engine, IDA, Wireshark, Fiddler)
// Установите NIGHTAUTH_IGNORE_SECURITY=1 при локальном тестировании
nightauth.CheckBlacklistedProcesses()

// Получение стабильного аппаратного идентификатора (HWID)
hwid := nightauth.GetHardwareId()
fmt.Printf("HWID устройства: %s\n", hwid)
```

### 25. Дополнительные опции и ключи сервера

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

### 26. Криптографически подписанные обновления

Проверяйте и скачивайте защищенные релизы приложения, подписанные Ed25519:

```go
// Прямая проверенная загрузка с защитой от отката версий
targetPath, err := nightauth.SignedUpdateInstaller.Download(
	"https://releases.example.com/manifest.json",
	"ВАШ_PINNED_ED25519_ПУБЛИЧНЫЙ_КЛЮЧ_HEX",
	"my-app",
	"updates/my-app.exe",
	5*time.Minute,
)
if err != nil {
	fmt.Printf("Ошибка обновления: %v\n", err)
} else {
	fmt.Printf("Проверенное обновление установлено в: %s\n", targetPath)
}

// Или выполните пайплайн вручную:
verified, err := nightauth.VerifySignedUpdate(manifestEnvelopeJSON, pinnedPublicKeyHex, "my-app", minimumSequence)
if err == nil {
	fmt.Printf("Проверенный релиз %s v%s (секв %d)\n", verified.Application, verified.Version, verified.Sequence)
	installed, err := nightauth.DownloadSignedUpdate(verified, "updates/my-app.exe", 5*time.Minute)
}
```

## 🖥️ Запуск консольного и десктопного клиентов

Запустите любой из эталонных клиентов из каталога `source`:

```bash
# Interactive Terminal Console Client
cd console
go run .

# Modern Pure Win32 Desktop Client (Windows)
cd desktop
go run .
```

---

## 📄 Лицензия

Этот проект лицензирован в соответствии с **Elastic License 2.0**.

* Вы не имеете права предоставлять программное обеспечение третьим лицам в качестве размещенного или управляемого сервиса.
* Вы не имеете права изменять, удалять или обходить функциональность лицензионных ключей.
* Вы не имеете права изменять, удалять или скрывать уведомления об авторских правах или лицензиях.

Полный текст лицензии доступен в файле [LICENSE](../LICENSE).
