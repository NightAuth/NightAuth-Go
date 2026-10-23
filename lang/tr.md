# NightAuth Go SDK — GitHub'da bize yıldız verin ⭐

Resmi **NightAuth** Go SDK ve referans istemci paketi (Native Win32 GDI tabanlı modern masaüstü arayüzü ve interaktif terminal CLI). Go yazılımlarınızı tersine mühendisliğe, korsana ve kurcalamaya karşı bulut lisanslama, HWID kilitleme, uçtan uca şifreli sohbet ve bellekte güvenli çalıştırma ile dakikalar içinde koruyun.

Resmi Portal: [https://nightauth.org](https://nightauth.org)

---

**Languages / Idiomas / Langues / Sprachen:**  
[English](../README.md) • [Español](es.md) • [Français](fr.md) • [Deutsch](de.md) • [Italiano](it.md) • [Português](pt.md) • [Русский](ru.md) • [Türkçe](tr.md) • [Polski](pl.md) • [中文](zh.md)

---

> [!NOTE]
> ### 📺 Video Eğitimler ve Kurulum Kılavuzları
> NightAuth kontrol panelinde uygulamaların, planların, aboneliklerin, bulut değişkenlerinin, webhook'ların ve lisansların nasıl yapılandırılacağını açıklayan resmi video eğitimler ve adım adım kurulum kılavuzları **yakında yayınlanacaktır** resmi portalımızda ve YouTube kanalımızda.

---

## 📸 Ekran Görüntüleri

### Masaüstü Uygulaması (Yerel Win32 GDI Modern Koyu Tema)

| Başlangıç / Kimlik Doğrulama Ekranı | Kontrol Paneli ve Abonelik Yönetimi |
| :---: | :---: |
| ![NightAuth Masaüstü Giriş](../image/desktop.png) | ![NightAuth Masaüstü Kontrol Paneli](../image/desktop-dashboard.png) |

### Konsol Uygulaması (İnteraktif Terminal CLI)

| Terminal Giriş Ekranı | İnteraktif Konsol Paneli |
| :---: | :---: |
| ![NightAuth Konsol Giriş](../image/console.png) | ![NightAuth Konsol Kontrol Paneli](../image/console-dashboard.png) |

---

## 🛡️ Önerilen Güvenlik Uygulamaları

* **Gelişmiş Go Karartma (Obfuscation):** Standart Go ikili dosyaları çalışma zamanı tip bilgilerini ve sembol isimlerini saklar. Canlı sürümler için **Garble** (`garble -literals -tiny build`) kullanarak tipleri kaldırın, denetim akışını karartın ve dize değişmezlerini şifreleyin.
* **Ticari İkili Paketleme:** Derlenmiş ikili dosyayı IDA Pro, Ghidra ve x64dbg'ye karşı **VMProtect** veya **Themida** gibi ticari koruyucularla ek olarak koruyun.
* **Bellek Bütünlüğü ve Anti-Debug:** NightAuth, hata ayıklayıcıları (`x64dbg`, `CheatEngine`, `Wireshark`, `HTTPDebugger`) ve sistem proxy'lerini tespit eden yerleşik güvenlik korumaları içerir. Başlangıçta `nightauth.CheckDebugger()` ve `nightauth.CheckBlacklistedProcesses()` fonksiyonlarını çağırarak maksimum koruma sağlayın (bkz. [SecurityProtections](#24-securityprotections)).
* **Bellek İçi Güvenli Yük Çalıştırma:** Korunan ikili modülleri veya kabuk kodlarını `auth.DownloadFile()` ile alırken **asla baytları fiziksel diske kaydetmeyin**. Yükleri doğrudan sanal belleğe yükleyip çalıştırın ve işlem tamamlandıktan hemen sonra arabellekleri sıfırlayın.

---

## 📦 Gereksinimler ve Derleme

**Go 1.20** veya üzeri gereklidir.

Yerel Win32 masaüstü uygulaması CGO bağımlılıkları olmadan saf Windows sistem DLL'lerini (`user32.dll`, `gdi32.dll`, `kernel32.dll`, `msimg32.dll`) kullanır.

### Konsol İstemcisini Derleme
```bash
cd console
go build -ldflags="-s -w" -o bin/NightAuth-Console.exe .
```

### Masaüstü İstemcisini Derleme (Pencereli / Konsol Penceresiz)
```bash
cd desktop
go build -ldflags="-s -w -H=windowsgui" -o bin/NightAuth-Desktop.exe .
```

---

## ⚙️ Yapılandırma ve Başlatma

NightAuth Kontrol Panelindeki **Manage Apps** bölümünden uygulama kimlik bilgilerinizi yapılandırın:

```go
package main

import (
	"fmt"
	"os"

	"nightauth"
)

// Uygulama Yapılandırması (NightAuth Kontrol Paneli)
const (
	AppName    = ""
	AppApiKey  = ""
	AppVersion = "1.0"
	AppToken   = ""
)

func main() {
	// API istemcisini başlat
	auth := nightauth.NewApi(AppName, AppApiKey, AppVersion, AppToken, "")

	// Kriptografik el sıkışmayı gerçekleştir (Curve25519 + Ed25519 + AES-256-GCM)
	auth.Init()

	if !auth.Response.Success {
		fmt.Printf("Başlatma başarısız: %s\n", auth.Response.Message)
		os.Exit(1)
	}

	fmt.Println("NightAuth bulutuna başarıyla bağlanıldı!")
}
```

---

## 🚀 SDK Kullanım Kılavuzu

### 1. Oturum Doğrulama ve Kara Liste Durumu Kontrolü

```go
// Sunucu ile aktif oturum durumunu doğrula
if !auth.Check() {
	fmt.Printf("Geçersiz oturum: %s\n", auth.Response.Message)
}

// İstemci HWID veya IP adresinin kara listede olup olmadığını denetle
hwid := nightauth.SecurityProtections{}.GetHardwareId()
if auth.CheckBlacklist(hwid) {
	fmt.Println("Erişim Reddedildi: Cihaz veya IP kara listeye alınmış.")
	os.Exit(1)
}
```

### 2. Kullanıcı Adı ve Şifre ile Giriş (ve 2FA)

```go
if auth.Login(username, password, totpCode) {
	fmt.Printf("Tekrar hoş geldiniz, %s!\n", auth.User.Username)
} else {
	fmt.Printf("Giriş başarısız: %s\n", auth.Response.Message)
}
```

### 3. Lisans Anahtarı ile Kayıt Olma

```go
if auth.Register(username, password, licenseKey, email) {
	fmt.Printf("Kayıt başarılı: %s\n", auth.User.Username)
} else {
	fmt.Printf("Kayıt başarısız: %s\n", auth.Response.Message)
}
```

### 4. Yalnızca Lisans Anahtarı ile Kimlik Doğrulama (ve 2FA)

```go
// Yalnızca lisans anahtarı ile giriş
if auth.License(licenseKey) {
	fmt.Printf("Lisans geçerli! Giriş yapıldı: %s\n", auth.User.Username)
} else {
	fmt.Printf("Lisans kimlik doğrulaması başarısız: %s\n", auth.Response.Message)
}

// Lisans anahtarı ve 2FA TOTP kodu ile giriş
if auth.LicenseWithTotp(licenseKey, "123456") {
	fmt.Printf("Lisans TOTP ile doğrulandı! Giriş yapıldı: %s\n", auth.User.Username)
}
```

### 5. Hesap Yükseltme / Lisans Anahtarı Tanımlama

```go
if auth.Upgrade(username, upgradeKey) {
	fmt.Println("Abonelik başarıyla yükseltildi!")
} else {
	fmt.Printf("Yükseltme başarısız: %s\n", auth.Response.Message)
}
```

### 6. Donanım Kimliği Sıfırlama (HWID Kilidi)

```go
if auth.ResetHwid(username, password, totpCode) {
	fmt.Println("HWID başarıyla sıfırlandı. Artık bu cihazda giriş yapabilirsiniz.")
} else {
	fmt.Printf("HWID sıfırlama başarısız: %s\n", auth.Response.Message)
}
```

### 7. Şifremi Unuttum / Hesap Kurtarma

```go
// Adım 1: Şifre sıfırlama e-postası iste
if auth.ForgotPassword(email) {
	fmt.Println("Şifre sıfırlama bağlantısı e-postanıza gönderildi.")
}

// Adım 2: Belirteç ile yeni şifreyi onayla
if auth.ResetPassword(resetToken, newPassword) {
	fmt.Println("Şifre başarıyla güncellendi!")
}
```

### 8. Kullanıcı Profil Verileri (`auth.User`)

```go
user := auth.User
fmt.Printf("Kullanıcı Adı: %s\n", user.Username)
fmt.Printf("E-posta:       %s\n", user.Email)
fmt.Printf("Donanım ID:    %s\n", user.Hwid)
fmt.Printf("IP Adresi:     %s\n", user.Ip)
fmt.Printf("Oluşturulma:   %s\n", user.CreatedAt)
fmt.Printf("Son Giriş:     %s\n", user.LastLoginAt)

for _, sub := range user.Subscriptions {
	fmt.Printf("Plan: %s | Bitiş: %s | Anahtar: %s\n", sub.Subscription, sub.Expiry, sub.Key)
}
```

### 9. Bulut Değişkenleri (Genel)

```go
motd := auth.GetGlobalVar("server_motd")
fmt.Printf("Bulut Değişkeni (Global): %s\n", motd)
```

### 10. Kullanıcı Değişkenleri

```go
// Özel kullanıcı değişkeni ata
auth.SetUserVar("custom_theme", "dark_neon")

// Kayıtlı değişkeni oku
val := auth.GetUserVar("custom_theme")
fmt.Printf("Kullanıcı Değişkeni: %s\n", val)

// Tüm kullanıcı değişkenlerini listele
if auth.GetUserVariables() {
	rawVars, _ := auth.Response.Payload["variables"].([]interface{})
	for _, item := range rawVars {
		if m, ok := item.(map[string]interface{}); ok {
			fmt.Printf("Değişken %s = %s\n", m["key"], m["value"])
		}
	}
}
```

### 11. Uygulama Denetim Günlükleri

```go
auth.Log("Şüpheli bütünlük ihlali yakalandı.")
```

### 12. Kullanıcıyı Yasaklama ve Cihazı Kara Listeye Alma

```go
auth.Ban("Anti-Cheat bütünlük hatası: Bellek değiştirme girişimi tespit edildi.")
```

### 13. Sunucu Taraflı Webhook'lar ve Dosya Yükleme

```go
// Standart JSON webhook
payload := map[string]interface{}{
	"event": "alert_triggered",
	"level": "critical",
}
result := auth.TriggerWebhook("7kR0UedlVI", payload)
if result != nil && result.Success {
	fmt.Printf("Webhook başarıyla iletildi (HTTP %d)\n", result.StatusCode)
}

// Dosya eki ile çok parçalı (multipart) webhook
fileBytes := []byte("örnek rapor verisi")
fileResult := auth.TriggerWebhookFile("7kR0UedlVI", fileBytes, "report.txt", payload)
if fileResult != nil && fileResult.Success {
	fmt.Printf("Dosya webhook'u başarıyla iletildi (HTTP %d)\n", fileResult.StatusCode)
}
```

### 14. Bellek İçi Güvenli Dosya İndirme

```go
fileBytes := auth.DownloadFile("385624")
if len(fileBytes) > 0 {
	fmt.Printf("%d bayt güvenli şekilde RAM belleğe indirildi.\n", len(fileBytes))
	// Fiziksel diske yazmadan doğrudan bellek üzerinden çalıştırın veya işleyin
}
```

### 15. Gerçek Zamanlı Sohbet Kanalları

```go
// Kanala mesaj gönder
auth.SendChatMessage("Go SDK üzerinden topluluğa selamlar!", "general")

// Kanaldaki son mesajları çek
messages := auth.GetChatMessages("general")
for _, msg := range messages {
	fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.Author, msg.Message)
}

// Belirli bir kimlikten önceki sayfalanmış mesajları çek
paged := auth.GetChatMessagesPage("general", 0, true)
for _, msg := range paged {
	fmt.Printf("#%s [%s] %s: %s\n", msg.Id, msg.Timestamp, msg.Author, msg.Message)
}
```

### 16. Gerçek Zamanlı Olay Akışı (Server-Sent Events)

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err := auth.SubscribeEvents(ctx, "", func(evt nightauth.EventFrame) {
	fmt.Printf("[OLAY #%s] %s\n", evt.ID, evt.Data)
})
if err != nil {
	fmt.Printf("Olay akışı sonlandı: %v\n", err)
}
```

### 17. İki Faktörlü Kimlik Doğrulama (2FA)

```go
// TOTP 2FA kurulumu başlat
setupData := auth.Setup2FA()
if setupData != nil {
	fmt.Printf("Gizli Anahtar: %s | URI: %s\n", setupData.Secret, setupData.Uri)
}

// Kimlik doğrulama kodu ile 2FA kurulumunu onayla
confirmData := auth.Confirm2FA("123456")
if confirmData != nil && confirmData.Success {
	fmt.Println("2FA başarıyla etkinleştirildi! Yedek Kodlar:")
	for _, code := range confirmData.BackupCodes {
		fmt.Printf(" - %s\n", code)
	}
}

// 2FA devre dışı bırak
if auth.Disable2FA(password, "123456") {
	fmt.Println("2FA başarıyla devre dışı bırakıldı.")
}
```

### 18. Oturum Yaşam Döngüsü ve Çıkış

```go
// Oturumu aktif tutmak için periyodik sinyal (heartbeat) gönder
auth.Heartbeat()

// Bulutta oturumu sonlandır ve olay dinleyicisini kapat
auth.Logout()
auth.Close()
```

### 19. Kullanıcı Adını Değiştirme

```go
if auth.ChangeUsername("yeni_takma_ad", "gecerli_sifre") {
	fmt.Printf("Kullanıcı adı başarıyla güncellendi: %s\n", auth.User.Username)
} else {
	fmt.Printf("Kullanıcı adı değiştirilemedi: %s\n", auth.Response.Message)
}
```

### 20. Uygulama İstatistikleri ve Çevrimiçi Kullanıcılar

```go
// Uygulama istatistiklerini getir
stats := auth.GetStats()
if stats != nil {
	fmt.Printf("Kullanıcılar: %d toplam, %d aktif, %d yasaklı\n", stats.TotalUsers, stats.ActiveUsers, stats.BannedUsers)
	fmt.Printf("Lisanslar: %d toplam, %d kullanılan, %d boşta\n", stats.TotalLicenses, stats.UsedLicenses, stats.UnusedLicenses)
	fmt.Printf("Sunucu Sürümü: %s\n", stats.Version)
}

// Çevrimiçi kullanıcıların listesini getir
onlineUsers := auth.GetOnlineUsers()
fmt.Printf("Çevrimiçi kullanıcılar (%d): %v\n", len(onlineUsers), onlineUsers)
```

### 21. Dosya Hash Doğrulama

```go
// Sunucuda kayıtlı SHA-256 hash değerine karşı binary doğrulaması
isValid := auth.VerifyFileHash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
if isValid {
	fmt.Println("Dosya bütünlüğü başarıyla doğrulandı.")
} else {
	fmt.Println("Bütünlük doğrulaması başarısız: hash uyuşmazlığı.")
}
```

### 22. Erişim Biletleri

```go
// Kapsamlı kısa ömürlü erişim bileti talep et
ticket := auth.AccessTicket("file_download", "file_abc123")
if ticket != "" {
	fmt.Printf("Erişim bileti verildi: %s\n", ticket)
} else {
	fmt.Printf("Erişim bileti verilemedi: %s\n", auth.Response.Message)
}
```

### 23. Gerçek Zamanlı Olay Geri Çağrımları ve Arka Plan Dinleyicisi

Reaktif olay işleyicilerini yapılandırın ve arka plan dağıtıcısını başlatın:

```go
auth.OnSessionInvalidated = func(reason string) {
	fmt.Printf("Oturum sonlandırıldı: %s\n", reason)
	os.Exit(0)
}
auth.OnAppBanned = func(reason string) {
	fmt.Printf("Uygulama yasaklandı: %s\n", reason)
}
auth.OnAppStatusChanged = func(status string) {
	fmt.Printf("Uygulama durumu güncellendi: %s\n", status)
}
auth.OnVersionUpdated = func(version string) {
	fmt.Printf("Yeni sürüm mevcut: %s\n", version)
}
auth.OnChatMessageReceived = func(msg nightauth.ChatMessage, channel string) {
	fmt.Printf("[%s] %s: %s\n", channel, msg.Author, msg.Message)
}

// Arka plan SSE olay dinleyicisini başlatın
auth.StartRealtimeEventListener()

// Artık gerekmediğinde dinleyiciyi durdurun
auth.StopRealtimeEventListener()
```

### 24. `SecurityProtections`

Statik ve örnek kurcalama önleme güvenlik yardımcıları:

```go
// Otomatik denetimleri çalıştırın (sistem HTTP proxy algılama)
nightauth.RunSecurityChecks()

// Etkin hata ayıklayıcıları algılayın (yerel veya uzak Win32 hata ayıklayıcıları)
nightauth.CheckDebugger()

// Kara listedeki analiz araçlarını algılayın (x64dbg, Cheat Engine, IDA, Wireshark, Fiddler)
// Yerel test sırasında atlamak için NIGHTAUTH_IGNORE_SECURITY=1 ayarlayın
nightauth.CheckBlacklistedProcesses()

// Sabit donanım parmak izi tanımlayıcısını alın
hwid := nightauth.GetHardwareId()
fmt.Printf("Cihaz HWID: %s\n", hwid)
```

### 25. Bağımsız Seçenekler ve Özel Sunucu Anahtarları

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

### 26. İmzalı Otomatik Güncellemeler

Ed25519 ile imzalanmış kurcalamaya karşı korumalı uygulama sürümlerini doğrulayın ve indirin:

```go
// Yeniden oynatma geri alma korumasıyla tek seferlik doğrulanmış güncelleme indirmesi
targetPath, err := nightauth.SignedUpdateInstaller.Download(
	"https://releases.example.com/manifest.json",
	"ÖZEL_PINNED_ED25519_GENEL_ANAHTAR_HEX",
	"my-app",
	"updates/my-app.exe",
	5*time.Minute,
)
if err != nil {
	fmt.Printf("Güncelleme başarısız: %v\n", err)
} else {
	fmt.Printf("Doğrulanmış güncelleme yüklendi: %s\n", targetPath)
}

// Veya ardışık düzeni manuel olarak yürütün:
verified, err := nightauth.VerifySignedUpdate(manifestEnvelopeJSON, pinnedPublicKeyHex, "my-app", minimumSequence)
if err == nil {
	fmt.Printf("Doğrulanan sürüm %s v%s (sıra %d)\n", verified.Application, verified.Version, verified.Sequence)
	installed, err := nightauth.DownloadSignedUpdate(verified, "updates/my-app.exe", 5*time.Minute)
}
```

## 🖥️ Konsol ve Masaüstü İstemcilerini Çalıştırma

`source` dizininden referans istemcilerden birini başlatın:

```bash
# Interactive Terminal Console Client
cd console
go run .

# Modern Pure Win32 Desktop Client (Windows)
cd desktop
go run .
```

---

## 📄 Lisans

Bu proje **Elastic License 2.0** kapsamında lisanslanmıştır.

* Yazılımı üçüncü taraflara barındırılan veya yönetilen bir hizmet olarak sağlayamazsınız.
* Lisans anahtarı işlevselliğini taşıyamaz, değiştiremez, devre dışı bırakamaz veya atlatamazsınız.
* Lisanslama veya telif hakkı bildirimlerini değiştiremez, kaldıramaz veya gizleyemezsiniz.

Tam yasal metin için [LICENSE](../LICENSE) dosyasına bakın.
