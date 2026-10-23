# NightAuth Go SDK — Deixe uma estrela no GitHub ⭐

SDK oficial de **NightAuth** para Go e suíte de clientes de referência (interface gráfica moderna em **Win32 GDI nativo** e CLI interativa para terminal). Proteja seu software em Go contra pirataria, engenharia reversa e adulterações com gerenciamento de licenças na nuvem, bloqueio por HWID, chat criptografado em tempo real e entrega segura de módulos em memória em poucos minutos.

Portal Oficial: [https://nightauth.org](https://nightauth.org)

---

**Languages / Idiomas / Langues / Sprachen:**  
[English](../README.md) • [Español](es.md) • [Français](fr.md) • [Deutsch](de.md) • [Italiano](it.md) • [Português](pt.md) • [Русский](ru.md) • [Türkçe](tr.md) • [Polski](pl.md) • [中文](zh.md)

---

> [!NOTE]
> ### 📺 Tutoriais em Vídeo e Guias de Configuração
> Vídeos tutoriais oficiais e guías passo a passo detalhando como configurar aplicativos, planos, assinaturas, variáveis na nuvem, webhooks e licenças no painel da NightAuth **serão publicados em breve** em nosso portal e canal oficial.

---

## 📸 Capturas de Tela

### Aplicativo Desktop (Tema Escuro Moderno Nativo Win32 GDI)

| Tela Inicial / Autenticação | Painel de Controle & Gestão de Licenças |
| :---: | :---: |
| ![NightAuth Desktop Login](../image/desktop.png) | ![NightAuth Desktop Dashboard](../image/desktop-dashboard.png) |

### Aplicativo de Console (CLI Interativa para Terminal)

| Tela de Login no Terminal | Painel Interativo de Console |
| :---: | :---: |
| ![NightAuth Console Login](../image/console.png) | ![NightAuth Console Dashboard](../image/console-dashboard.png) |

---

## 🛡️ Práticas Recomendadas de Segurança

* **Ofuscação Avançada em Go:** Binários Go padrão preservam nomes de símbolos e tipos de reflexão. Para versões de produção, compile usando **Garble** (`garble -literals -tiny build`) para remover tipos e criptografar strings.
* **Empacotamento Comercial:** Proteja o executável com ferramentas como **VMProtect** ou **Themida** para salvaguardar contra IDA Pro, Ghidra e x64dbg.
* **Integridade de Memória & Anti-Depuração:** NightAuth inclui proteções nativas contra depuradores (`x64dbg`, `CheatEngine`, `Wireshark`, `HTTPDebugger`) e proxies de sistema. Invoque `nightauth.CheckDebugger()` e `nightauth.CheckBlacklistedProcesses()` na inicialização (veja [SecurityProtections](#24-securityprotections)).
* **Execução Segura em Memória:** Ao baixar módulos protegidos com `auth.DownloadFile()`, **nunca grave os bytes no disco físico**. Execute-os diretamente na memória RAM e limpe os buffers imediatamente após o uso.

---

## 📦 Requisitos & Compilação

Requer **Go 1.20** ou superior.

O aplicativo de desktop nativo Win32 utiliza DLLs padrão do Windows (`user32.dll`, `gdi32.dll`, `kernel32.dll`, `msimg32.dll`) sem dependências CGO externas.

### Compilar o Cliente de Console
```bash
cd console
go build -ldflags="-s -w" -o bin/NightAuth-Console.exe .
```

### Compilar o Cliente Desktop (Sem Janela de Console)
```bash
cd desktop
go build -ldflags="-s -w -H=windowsgui" -o bin/NightAuth-Desktop.exe .
```

---

## ⚙️ Configuração & Inicialização

Configure as credenciais do seu aplicativo no painel de controle da NightAuth:

```go
package main

import (
	"fmt"
	"os"

	"nightauth"
)

// Configuração do Aplicativo
const (
	AppName    = ""
	AppApiKey  = ""
	AppVersion = "1.0"
	AppToken   = ""
)

func main() {
	// Inicializar cliente API
	auth := nightauth.NewApi(AppName, AppApiKey, AppVersion, AppToken, "")

	// Realizar handshake criptográfico (Curve25519 + Ed25519 + AES-256-GCM)
	auth.Init()

	if !auth.Response.Success {
		fmt.Printf("Falha na inicialização: %s\n", auth.Response.Message)
		os.Exit(1)
	}

	fmt.Println("Conectado com sucesso à nuvem NightAuth!")
}
```

---

## 🚀 Guia de Uso do SDK

### 1. Validação de Sessão & Verificação de Lista Negra

```go
if !auth.Check() {
	fmt.Printf("Sessão inválida: %s\n", auth.Response.Message)
}

hwid := nightauth.SecurityProtections{}.GetHardwareId()
if auth.CheckBlacklist(hwid) {
	fmt.Println("Acesso Negado: Máquina ou IP restrito.")
	os.Exit(1)
}
```

### 2. Login com Usuário e Senha (e 2FA)

```go
if auth.Login(username, password, totpCode) {
	fmt.Printf("Bem-vindo de volta, %s!\n", auth.User.Username)
} else {
	fmt.Printf("Falha no login: %s\n", auth.Response.Message)
}
```

### 3. Registro com Chave de Licença

```go
if auth.Register(username, password, licenseKey, email) {
	fmt.Printf("Registro realizado com sucesso: %s\n", auth.User.Username)
} else {
	fmt.Printf("Falha no registro: %s\n", auth.Response.Message)
}
```

### 4. Autenticação Apenas por Chave de Licença (e 2FA)

```go
// Iniciar sessão apenas com chave de licença
if auth.License(licenseKey) {
	fmt.Printf("Licença válida! Sessão iniciada como: %s\n", auth.User.Username)
} else {
	fmt.Printf("Falha na autenticação da licença: %s\n", auth.Response.Message)
}

// Iniciar sessão com chave de licença e código 2FA TOTP
if auth.LicenseWithTotp(licenseKey, "123456") {
	fmt.Printf("Licença verificada com TOTP! Sessão iniciada como: %s\n", auth.User.Username)
}
```

### 5. Atualizar Plano / Resgatar Nova Chave

```go
if auth.Upgrade(username, upgradeKey) {
	fmt.Println("Assinatura atualizada com sucesso!")
} else {
	fmt.Printf("Falha na atualização: %s\n", auth.Response.Message)
}
```

### 6. Redefinir Bloqueio de Hardware (HWID Lock)

```go
if auth.ResetHwid(username, password, totpCode) {
	fmt.Println("HWID redefinido. Você já pode fazer login neste dispositivo.")
} else {
	fmt.Printf("Falha ao redefinir HWID: %s\n", auth.Response.Message)
}
```

### 7. Recuperação de Conta / Esqueceu a Senha

```go
if auth.ForgotPassword(email) {
	fmt.Println("Link de recuperação enviado para o e-mail.")
}

if auth.ResetPassword(resetToken, newPassword) {
	fmt.Println("Senha atualizada com sucesso!")
}
```

### 8. Dados do Perfil de Usuário (`auth.User`)

```go
user := auth.User
fmt.Printf("Usuário:      %s\n", user.Username)
fmt.Printf("E-mail:       %s\n", user.Email)
fmt.Printf("HWID:         %s\n", user.Hwid)
fmt.Printf("IP Cliente:   %s\n", user.Ip)
fmt.Printf("Criado em:    %s\n", user.CreatedAt)
fmt.Printf("Último login: %s\n", user.LastLoginAt)

for _, sub := range user.Subscriptions {
	fmt.Printf("Plano: %s | Validade: %s | Chave: %s\n", sub.Subscription, sub.Expiry, sub.Key)
}
```

### 9. Variáveis na Nuvem (Globais)

```go
motd := auth.GetGlobalVar("server_motd")
fmt.Printf("Variável na Nuvem: %s\n", motd)
```

### 10. Variáveis do Usuário

```go
auth.SetUserVar("custom_theme", "dark_neon")

val := auth.GetUserVar("custom_theme")
fmt.Printf("Variável do Usuário: %s\n", val)

if auth.GetUserVariables() {
	rawVars, _ := auth.Response.Payload["variables"].([]interface{})
	for _, item := range rawVars {
		if m, ok := item.(map[string]interface{}); ok {
			fmt.Printf("Variável %s = %s\n", m["key"], m["value"])
		}
	}
}
```

### 11. Logs de Auditoria

```go
auth.Log("Tentativa de adulteração interceptada.")
```

### 12. Banir Usuário e Bloquear Máquina (Killswitch / Ban)

```go
auth.Ban("Falha de integridade Anti-Cheat: Modificação de memória detectada.")
```

### 13. Webhooks no Servidor & Upload de Arquivos

```go
payload := map[string]interface{}{
	"event": "alert_triggered",
	"level": "critical",
}
result := auth.TriggerWebhook("7kR0UedlVI", payload)
if result != nil && result.Success {
	fmt.Printf("Webhook entregue com sucesso (HTTP %d)\n", result.StatusCode)
}

fileBytes := []byte("conteúdo do relatório")
fileResult := auth.TriggerWebhookFile("7kR0UedlVI", fileBytes, "report.txt", payload)
if fileResult != nil && fileResult.Success {
	fmt.Printf("Webhook com arquivo entregue (HTTP %d)\n", fileResult.StatusCode)
}
```

### 14. Download Seguro na Memória RAM

```go
fileBytes := auth.DownloadFile("385624")
if len(fileBytes) > 0 {
	fmt.Printf("Baixados %d bytes com segurança na memória RAM.\n", len(fileBytes))
}
```

### 15. Canais de Chat em Tempo Real

```go
auth.SendChatMessage("Olá comunidade do Go!", "general")

messages := auth.GetChatMessages("general")
for _, msg := range messages {
	fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.Author, msg.Message)
}

paged := auth.GetChatMessagesPage("general", 0, true)
for _, msg := range paged {
	fmt.Printf("#%s [%s] %s: %s\n", msg.Id, msg.Timestamp, msg.Author, msg.Message)
}
```

### 16. Eventos em Tempo Real (Server-Sent Events)

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err := auth.SubscribeEvents(ctx, "", func(evt nightauth.EventFrame) {
	fmt.Printf("[EVENTO #%s] %s\n", evt.ID, evt.Data)
})
```

### 17. Autenticação de Dois Fatores (2FA)

```go
setupData := auth.Setup2FA()
if setupData != nil {
	fmt.Printf("Segredo: %s | URI: %s\n", setupData.Secret, setupData.Uri)
}

confirmData := auth.Confirm2FA("123456")
if confirmData != nil && confirmData.Success {
	fmt.Println("2FA Ativado com sucesso! Códigos de backup:")
	for _, code := range confirmData.BackupCodes {
		fmt.Printf(" - %s\n", code)
	}
}

if auth.Disable2FA(password, "123456") {
	fmt.Println("2FA desativado com sucesso.")
}
```

### 18. Ciclo de Vida da Sessão & Logout

```go
auth.Heartbeat()

auth.Logout()
auth.Close()
```

### 19. Alterar Nome de Utilizador

```go
if auth.ChangeUsername("novo_alias", "palavra_passe_atual") {
	fmt.Printf("Nome de utilizador atualizado com sucesso para: %s\n", auth.User.Username)
} else {
	fmt.Printf("Falha ao alterar nome de utilizador: %s\n", auth.Response.Message)
}
```

### 20. Estatísticas da Aplicação e Utilizadores Online

```go
// Obter estatísticas da aplicação
stats := auth.GetStats()
if stats != nil {
	fmt.Printf("Utilizadores: %d total, %d ativos, %d banidos\n", stats.TotalUsers, stats.ActiveUsers, stats.BannedUsers)
	fmt.Printf("Licenças: %d total, %d usadas, %d não usadas\n", stats.TotalLicenses, stats.UsedLicenses, stats.UnusedLicenses)
	fmt.Printf("Versão do servidor: %s\n", stats.Version)
}

// Obter lista de utilizadores atualmente online
onlineUsers := auth.GetOnlineUsers()
fmt.Printf("Utilizadores online (%d): %v\n", len(onlineUsers), onlineUsers)
```

### 21. Verificar Hash do Ficheiro

```go
// Verificar hash SHA-256 do binário contra os hashes aprovados pelo servidor
isValid := auth.VerifyFileHash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
if isValid {
	fmt.Println("Integridade do ficheiro verificada com sucesso.")
} else {
	fmt.Println("Falha na integridade: hash incorreto ou compilação não autorizada.")
}
```

### 22. Tickets de Acesso

```go
// Solicitar ticket de acesso temporário com âmbito
ticket := auth.AccessTicket("file_download", "file_abc123")
if ticket != "" {
	fmt.Printf("Ticket de acesso concedido: %s\n", ticket)
} else {
	fmt.Printf("Falha ao emitir ticket de acesso: %s\n", auth.Response.Message)
}
```

### 23. Callbacks Reativos & Listener em Segundo Plano

Configure manipuladores de eventos reativos e execute o despachante em segundo plano:

```go
auth.OnSessionInvalidated = func(reason string) {
	fmt.Printf("Sessão terminada: %s\n", reason)
	os.Exit(0)
}
auth.OnAppBanned = func(reason string) {
	fmt.Printf("Aplicação banida: %s\n", reason)
}
auth.OnAppStatusChanged = func(status string) {
	fmt.Printf("Estado da aplicação atualizado: %s\n", status)
}
auth.OnVersionUpdated = func(version string) {
	fmt.Printf("Nova versão disponível: %s\n", version)
}
auth.OnChatMessageReceived = func(msg nightauth.ChatMessage, channel string) {
	fmt.Printf("[%s] %s: %s\n", channel, msg.Author, msg.Message)
}

// Iniciar listener SSE em segundo plano
auth.StartRealtimeEventListener()

// Parar listener quando não for mais necessário
auth.StopRealtimeEventListener()
```

### 24. `SecurityProtections`

Utilitários de segurança estáticos e de instância contra adulteração:

```go
// Executar verificações automáticas (deteção de proxy HTTP do sistema)
nightauth.RunSecurityChecks()

// Detetar depuradores ativos (locais ou remotos Win32)
nightauth.CheckDebugger()

// Detetar ferramentas de análise na lista negra (x64dbg, Cheat Engine, IDA, Wireshark, Fiddler)
// Defina NIGHTAUTH_IGNORE_SECURITY=1 durante testes locais
nightauth.CheckBlacklistedProcesses()

// Obter identificador estável de hardware do dispositivo
hwid := nightauth.GetHardwareId()
fmt.Printf("HWID do dispositivo: %s\n", hwid)
```

### 25. Opções Avançadas e Chaves Personalizadas

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

### 26. Atualizações com Assinatura Criptográfica

Verifique e descarregue versões seguras contra adulteração assinadas com Ed25519:

```go
// Download direto verificado com proteção anti-rollback
targetPath, err := nightauth.SignedUpdateInstaller.Download(
	"https://releases.example.com/manifest.json",
	"CHAVE_PUBLICA_PINNED_ED25519_HEX",
	"my-app",
	"updates/my-app.exe",
	5*time.Minute,
)
if err != nil {
	fmt.Printf("Falha na atualização: %v\n", err)
} else {
	fmt.Printf("Atualização verificada instalada em: %s\n", targetPath)
}

// Ou controle o pipeline manualmente:
verified, err := nightauth.VerifySignedUpdate(manifestEnvelopeJSON, pinnedPublicKeyHex, "my-app", minimumSequence)
if err == nil {
	fmt.Printf("Versão verificada %s v%s (seq %d)\n", verified.Application, verified.Version, verified.Sequence)
	installed, err := nightauth.DownloadSignedUpdate(verified, "updates/my-app.exe", 5*time.Minute)
}
```

## 🖥️ Executar os Clientes de Console e Desktop

Execute qualquer um dos clientes de referência a partir do diretório `source`:

```bash
# Interactive Terminal Console Client
cd console
go run .

# Modern Pure Win32 Desktop Client (Windows)
cd desktop
go run .
```

---

## 📄 Licença

Este projeto está licenciado sob a **Elastic License 2.0**.

* Você não pode fornecer o software a terceiros como um serviço hospedado ou gerenciado.
* Você não pode remover, alterar ou burlar a funcionalidade de chave de licença.
* Você não pode alterar, remover ou ocultar avisos de direitos autorais ou licença.

Consulte o arquivo [LICENSE](../LICENSE) para obter o texto jurídico integral.
