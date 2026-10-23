# NightAuth Go SDK — 在 GitHub 上为我们点赞 ⭐

官方 **NightAuth** Go SDK 与参考客户端套件（基于**原生 Win32 GDI** 的现代桌面 UI 与交互式终端 CLI）。几分钟内为您的 Go 软件提供防逆向工程、防篡改与防盗版保护，内置云端授权、硬件 HWID 绑定、端到端加密实时聊天以及内存安全载荷执行。

官方门户: [https://nightauth.org](https://nightauth.org)

---

**Languages / Idiomas / Langues / Sprachen:**  
[English](../README.md) • [Español](es.md) • [Français](fr.md) • [Deutsch](de.md) • [Italiano](it.md) • [Português](pt.md) • [Русский](ru.md) • [Türkçe](tr.md) • [Polski](pl.md) • [中文](zh.md)

---

> [!NOTE]
> ### 📺 视频教程与系统设置指南
> 详细说明如何在 NightAuth 控制面板中配置应用程序、订阅计划、云变量、Webhook 和许可证的官方视频教程和分步指南**即将发布**在我们的官方网站和 YouTube 频道。

---

## 📸 屏幕截图

### 桌面客户端 (原生 Win32 GDI 现代深色主题)

| 启动 / 身份验证界面 | 控制台与订阅管理 |
| :---: | :---: |
| ![NightAuth 桌面登录](../image/desktop.png) | ![NightAuth 桌面控制台](../image/desktop-dashboard.png) |

### 控制台客户端 (交互式终端 CLI)

| 终端登录界面 | 交互式终端控制台 |
| :---: | :---: |
| ![NightAuth 终端登录](../image/console.png) | ![NightAuth 终端控制台](../image/console-dashboard.png) |

---

## 🛡️ 推荐的安全实践

* **高级 Go 混淆:** 标准 Go 二进制文件包含运行时反射信息和符号名称。发布生产版本时，请使用 **Garble** (`garble -literals -tiny build`) 剥离类型元数据、混淆控制流并加密字符串常量。
* **商业二进制加壳:** 对编译后的可执行文件使用商业加壳保护工具（如 **VMProtect** 或 **Themida**），防止 IDA Pro、Ghidra 和 x64dbg 静态分析与动态调试。
* **内存完整性与反调试:** NightAuth 内置安全防护，可检测常见调试器（`x64dbg`、`CheatEngine`、`Wireshark`、`HTTPDebugger`）及系统代理。在启动时调用 `nightauth.CheckDebugger()` 和 `nightauth.CheckBlacklistedProcesses()` 以获得最高保护级别（参见 [SecurityProtections](#24-securityprotections)）。
* **内存中安全载荷执行:** 使用 `auth.DownloadFile()` 检索受保护的二进制模块或 Shellcode 时，**切勿将字节写入物理磁盘**。直接在虚拟内存中执行载荷，并在执行后立即将内存缓冲区清零。

---

## 📦 环境要求与构建

需要 **Go 1.20** 或更高版本。

原生 Win32 桌面应用程序使用纯 Windows 系统 DLL（`user32.dll`、`gdi32.dll`、`kernel32.dll`、`msimg32.dll`），零 CGO 依赖。

### 构建控制台客户端
```bash
cd console
go build -ldflags="-s -w" -o bin/NightAuth-Console.exe .
```

### 构建桌面客户端 (无控制台窗口)
```bash
cd desktop
go build -ldflags="-s -w -H=windowsgui" -o bin/NightAuth-Desktop.exe .
```

---

## ⚙️ 配置与初始化

在 NightAuth 控制台的 **Manage Apps** 页面获取应用程序凭证：

```go
package main

import (
	"fmt"
	"os"

	"nightauth"
)

// 应用程序配置 (来自 NightAuth 控制台)
const (
	AppName    = ""
	AppApiKey  = ""
	AppVersion = "1.0"
	AppToken   = ""
)

func main() {
	// 初始化 API 客户端
	auth := nightauth.NewApi(AppName, AppApiKey, AppVersion, AppToken, "")

	// 执行加密握手 (Curve25519 + Ed25519 + AES-256-GCM)
	auth.Init()

	if !auth.Response.Success {
		fmt.Printf("初始化失败: %s\n", auth.Response.Message)
		os.Exit(1)
	}

	fmt.Println("成功连接到 NightAuth 云端服务！")
}
```

---

## 🚀 SDK 使用指南

### 1. 验证会话与黑名单状态

```go
// 验证与服务端的当前会话状态
if !auth.Check() {
	fmt.Printf("会话无效: %s\n", auth.Response.Message)
}

// 检查当前设备 HWID 或 IP 是否被列入黑名单
hwid := nightauth.SecurityProtections{}.GetHardwareId()
if auth.CheckBlacklist(hwid) {
	fmt.Println("拒绝访问: 当前机器或 IP 已被列入黑名单。")
	os.Exit(1)
}
```

### 2. 用户名和密码登录 (及 2FA)

```go
if auth.Login(username, password, totpCode) {
	fmt.Printf("欢迎回来，%s！\n", auth.User.Username)
} else {
	fmt.Printf("登录失败: %s\n", auth.Response.Message)
}
```

### 3. 使用许可证密钥注册

```go
if auth.Register(username, password, licenseKey, email) {
	fmt.Printf("注册成功: %s\n", auth.User.Username)
} else {
	fmt.Printf("注册失败: %s\n", auth.Response.Message)
}
```

### 4. 仅使用许可证密钥验证（及 2FA）

```go
// 仅使用许可证密钥登录
if auth.License(licenseKey) {
	fmt.Printf("许可证有效！登录身份：%s\n", auth.User.Username)
} else {
	fmt.Printf("许可证认证失败：%s\n", auth.Response.Message)
}

// 使用许可证密钥及 2FA TOTP 动态码登录
if auth.LicenseWithTotp(licenseKey, "123456") {
	fmt.Printf("许可证通过 TOTP 验证！登录身份：%s\n", auth.User.Username)
}
```

### 5. 升级账户 / 兑换许可证密钥

```go
if auth.Upgrade(username, upgradeKey) {
	fmt.Println("订阅升级成功！")
} else {
	fmt.Printf("升级失败: %s\n", auth.Response.Message)
}
```

### 6. 重置硬件 ID (HWID 锁定)

```go
if auth.ResetHwid(username, password, totpCode) {
	fmt.Println("HWID 重置成功。您现在可以在此设备上登录。")
} else {
	fmt.Printf("HWID 重置失败: %s\n", auth.Response.Message)
}
```

### 7. 忘记密码 / 账户找回

```go
// 步骤 1: 请求密码重置邮件
if auth.ForgotPassword(email) {
	fmt.Println("密码重置链接已发送至您的邮箱。")
}

// 步骤 2: 使用令牌确认新密码
if auth.ResetPassword(resetToken, newPassword) {
	fmt.Println("密码更新成功！")
}
```

### 8. 用户个人资料 (`auth.User`)

```go
user := auth.User
fmt.Printf("用户名:     %s\n", user.Username)
fmt.Printf("邮箱:       %s\n", user.Email)
fmt.Printf("硬件 ID:    %s\n", user.Hwid)
fmt.Printf("IP 地址:    %s\n", user.Ip)
fmt.Printf("注册时间:   %s\n", user.CreatedAt)
fmt.Printf("最后登录:   %s\n", user.LastLoginAt)

for _, sub := range user.Subscriptions {
	fmt.Printf("方案: %s | 到期时间: %s | 密钥: %s\n", sub.Subscription, sub.Expiry, sub.Key)
}
```

### 9. 全局云变量

```go
motd := auth.GetGlobalVar("server_motd")
fmt.Printf("云变量 (全局): %s\n", motd)
```

### 10. 用户专属变量

```go
// 设置自定义用户变量
auth.SetUserVar("custom_theme", "dark_neon")

// 获取存储的变量
val := auth.GetUserVar("custom_theme")
fmt.Printf("用户自定义变量: %s\n", val)

// 列出所有用户变量
if auth.GetUserVariables() {
	rawVars, _ := auth.Response.Payload["variables"].([]interface{})
	for _, item := range rawVars {
		if m, ok := item.(map[string]interface{}); ok {
			fmt.Printf("变量 %s = %s\n", m["key"], m["value"])
		}
	}
}
```

### 11. 应用程序审计日志

```go
auth.Log("捕获到可疑的完整性异常行为。")
```

### 12. 封禁用户并列入机器黑名单

```go
auth.Ban("反作弊完整性校验失败: 检测到内存修改尝试。")
```

### 13. 服务端 Webhook 与文件上传

```go
// 标准 JSON Webhook
payload := map[string]interface{}{
	"event": "alert_triggered",
	"level": "critical",
}
result := auth.TriggerWebhook("7kR0UedlVI", payload)
if result != nil && result.Success {
	fmt.Printf("Webhook 发送成功 (HTTP %d)\n", result.StatusCode)
}

// 带文件附件的多部分 (Multipart) Webhook
fileBytes := []byte("报告示例数据")
fileResult := auth.TriggerWebhookFile("7kR0UedlVI", fileBytes, "report.txt", payload)
if fileResult != nil && fileResult.Success {
	fmt.Printf("文件 Webhook 发送成功 (HTTP %d)\n", fileResult.StatusCode)
}
```

### 14. 内存中安全文件下载

```go
fileBytes := auth.DownloadFile("385624")
if len(fileBytes) > 0 {
	fmt.Printf("已安全下载 %d 字节至内存 (RAM)。\n", len(fileBytes))
	// 无需写入物理硬盘，直接在内存中执行或解密处理
}
```

### 15. 实时聊天频道

```go
// 向频道发送消息
auth.SendChatMessage("来自 Go SDK 的社区问候！", "general")

// 获取频道的最新消息
messages := auth.GetChatMessages("general")
for _, msg := range messages {
	fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.Author, msg.Message)
}

// 分页获取指定 ID 之前的历史消息
paged := auth.GetChatMessagesPage("general", 0, true)
for _, msg := range paged {
	fmt.Printf("#%s [%s] %s: %s\n", msg.Id, msg.Timestamp, msg.Author, msg.Message)
}
```

### 16. 实时事件流 (Server-Sent Events)

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err := auth.SubscribeEvents(ctx, "", func(evt nightauth.EventFrame) {
	fmt.Printf("[事件 #%s] %s\n", evt.ID, evt.Data)
})
if err != nil {
	fmt.Printf("事件流已终止: %v\n", err)
}
```

### 17. 双因素身份验证 (2FA / TOTP)

```go
// 发起 TOTP 2FA 设置
setupData := auth.Setup2FA()
if setupData != nil {
	fmt.Printf("密钥: %s | URI: %s\n", setupData.Secret, setupData.Uri)
}

// 使用验证器代码确认 2FA
confirmData := auth.Confirm2FA("123456")
if confirmData != nil && confirmData.Success {
	fmt.Println("2FA 启用成功！备用恢复代码:")
	for _, code := range confirmData.BackupCodes {
		fmt.Printf(" - %s\n", code)
	}
}

// 停用 2FA
if auth.Disable2FA(password, "123456") {
	fmt.Println("2FA 已成功停用。")
}
```

### 18. 会话生命周期与注销

```go
// 发送周期性心跳以保持会话活跃
auth.Heartbeat()

// 在云端注销会话并断开事件监听
auth.Logout()
auth.Close()
```

### 19. 修改用户名

```go
if auth.ChangeUsername("new_alias", "current_password") {
	fmt.Printf("用户名已成功更新为：%s\n", auth.User.Username)
} else {
	fmt.Printf("修改用户名失败：%s\n", auth.Response.Message)
}
```

### 20. 应用统计与在线用户

```go
// 获取应用聚合统计数据
stats := auth.GetStats()
if stats != nil {
	fmt.Printf("用户数：%d 总计, %d 活跃, %d 已封禁\n", stats.TotalUsers, stats.ActiveUsers, stats.BannedUsers)
	fmt.Printf("许可证：%d 总计, %d 已使用, %d 未使用\n", stats.TotalLicenses, stats.UsedLicenses, stats.UnusedLicenses)
	fmt.Printf("服务端版本：%s\n", stats.Version)
}

// 获取当前在线用户列表
onlineUsers := auth.GetOnlineUsers()
fmt.Printf("当前在线用户 (%d)：%v\n", len(onlineUsers), onlineUsers)
```

### 21. 校验文件哈希

```go
// 校验二进制文件的 SHA-256 哈希是否匹配服务端授权版本
isValid := auth.VerifyFileHash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
if isValid {
	fmt.Println("文件完整性校验成功。")
} else {
	fmt.Println("完整性校验失败：文件哈希不匹配或版本未授权。")
}
```

### 22. 访问票据

```go
// 请求短期作用域访问票据（如用于安全数据下载或流式传输）
ticket := auth.AccessTicket("file_download", "file_abc123")
if ticket != "" {
	fmt.Printf("访问票据已签发：%s\n", ticket)
} else {
	fmt.Printf("签发访问票据失败：%s\n", auth.Response.Message)
}
```

### 23. 实时事件回调与后台监听器

配置响应式事件处理器并启动后台分发器：

```go
auth.OnSessionInvalidated = func(reason string) {
	fmt.Printf("会话已终止：%s\n", reason)
	os.Exit(0)
}
auth.OnAppBanned = func(reason string) {
	fmt.Printf("应用已被封禁：%s\n", reason)
}
auth.OnAppStatusChanged = func(status string) {
	fmt.Printf("应用状态已更新：%s\n", status)
}
auth.OnVersionUpdated = func(version string) {
	fmt.Printf("发现新版本：%s\n", version)
}
auth.OnChatMessageReceived = func(msg nightauth.ChatMessage, channel string) {
	fmt.Printf("[%s] %s: %s\n", channel, msg.Author, msg.Message)
}

// 启动后台 SSE 事件监听协程
auth.StartRealtimeEventListener()

// 不需要时停止监听
auth.StopRealtimeEventListener()
```

### 24. `SecurityProtections`

静态防篡改与反调试安全组件：

```go
// 执行自动安全检查（系统 HTTP 代理检测）
nightauth.RunSecurityChecks()

// 检测活跃调试器（本地或远程 Win32 调试器）
nightauth.CheckDebugger()

// 检测黑名单分析工具（x64dbg, Cheat Engine, IDA, Wireshark, Fiddler）
// 本地测试期间可设置环境变量 NIGHTAUTH_IGNORE_SECURITY=1 跳过检测
nightauth.CheckBlacklistedProcesses()

// 获取稳定的硬件指纹标识符
hwid := nightauth.GetHardwareId()
fmt.Printf("设备 HWID：%s\n", hwid)
```

### 25. 独立配置与自定义服务端公钥

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

### 26. 签名自更新

校验并下载经 Ed25519 签名的高安全性防篡改发布包：

```go
// 单次调用完成带防重放与防回滚保护的校验更新下载
targetPath, err := nightauth.SignedUpdateInstaller.Download(
	"https://releases.example.com/manifest.json",
	"YOUR_PINNED_ED25519_PUBLIC_KEY_HEX",
	"my-app",
	"updates/my-app.exe",
	5*time.Minute,
)
if err != nil {
	fmt.Printf("更新失败：%v\n", err)
} else {
	fmt.Printf("已验证更新安装至：%s\n", targetPath)
}

// 或手动驱动完整更新管线：
verified, err := nightauth.VerifySignedUpdate(manifestEnvelopeJSON, pinnedPublicKeyHex, "my-app", minimumSequence)
if err == nil {
	fmt.Printf("已验证发布包 %s v%s (序号 %d)\n", verified.Application, verified.Version, verified.Sequence)
	installed, err := nightauth.DownloadSignedUpdate(verified, "updates/my-app.exe", 5*time.Minute)
}
```

## 🖥️ 运行控制台与桌面客户端

从 `source` 目录运行任意参考客户端：

```bash
# Interactive Terminal Console Client
cd console
go run .

# Modern Pure Win32 Desktop Client (Windows)
cd desktop
go run .
```

---

## 📄 许可证

本项目遵循 **Elastic License 2.0**。

* 您不得将本软件作为托管或云端服务提供给第三方。
* 您不得移动、修改、禁用或绕过许可证密钥功能。
* 您不得篡改、移除或遮挡任何许可与版权声明。

完整法律条款请参阅 [LICENSE](../LICENSE) 文件。
