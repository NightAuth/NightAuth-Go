package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"nightauth"
)

const (
	AppName    = ""
	AppApiKey  = ""
	AppVersion = "1.0"
	AppToken   = ""

	CReset   = "\x1b[0m"
	CBold    = "\x1b[1m"
	CDim     = "\x1b[38;2;100;116;139m"
	CMuted   = "\x1b[38;2;148;163;184m"
	CText    = "\x1b[38;2;226;232;240m"
	CWhite   = "\x1b[38;2;248;250;252m"
	CBlue    = "\x1b[38;2;59;130;246m"
	CGlow    = "\x1b[38;2;96;165;250m"
	CCyan    = "\x1b[38;2;56;189;248m"
	CGreen   = "\x1b[38;2;52;211;153m"
	CAmber   = "\x1b[38;2;251;191;36m"
	CRed     = "\x1b[38;2;248;113;113m"
	CBorder  = "\x1b[38;2;42;50;70m"
	CBadgeBg = "\x1b[48;2;23;37;68m"
	CGreenBg = "\x1b[48;2;14;46;38m"
	CRedBg   = "\x1b[48;2;56;20;26m"
)

var nightAuth *nightauth.Api
var reader = bufio.NewReader(os.Stdin)

func enableVirtualTerminal() {
	if runtime.GOOS != "windows" {
		return
	}
	k32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleOutputCP := k32.NewProc("SetConsoleOutputCP")
	setConsoleCP := k32.NewProc("SetConsoleCP")
	getConsoleMode := k32.NewProc("GetConsoleMode")
	setConsoleMode := k32.NewProc("SetConsoleMode")
	setConsoleTitleW := k32.NewProc("SetConsoleTitleW")

	_, _, _ = setConsoleOutputCP.Call(65001)
	_, _, _ = setConsoleCP.Call(65001)

	hOut := uintptr(syscall.Stdout)
	var mode uint32
	if r, _, _ := getConsoleMode.Call(hOut, uintptr(unsafe.Pointer(&mode))); r != 0 {
		_, _, _ = setConsoleMode.Call(hOut, uintptr(mode|0x0004))
	}

	if titlePtr, err := syscall.UTF16PtrFromString("NightAuth — Native Go Security Console"); err == nil {
		_, _, _ = setConsoleTitleW.Call(uintptr(unsafe.Pointer(titlePtr)))
	}
}

func clearScreen() {
	fmt.Print("\x1b[2J\x1b[H")
}

func printHr() {
	fmt.Printf("  %s────────────────────────────────────────────────────────────────────────────%s\n", CBorder, CReset)
}

func printSectionHeader(title, subtitle string) {
	fmt.Printf("\n  %s▌%s %s%s%s%s", CBlue, CReset, CWhite, CBold, title, CReset)
	if subtitle != "" {
		fmt.Printf("  %s·  %s%s", CDim, subtitle, CReset)
	}
	fmt.Println()
	printHr()
}

func printOk(title, detail string) {
	fmt.Printf("\n  %s%s%s ✓ SUCCESS %s %s%s%s%s", CGreenBg, CGreen, CBold, CReset, CWhite, CBold, title, CReset)
	if detail != "" {
		fmt.Printf(" %s— %s%s", CMuted, detail, CReset)
	}
	fmt.Print("\n\n")
}

func printErr(title, detail string) {
	fmt.Printf("\n  %s%s%s ✕ ERROR %s %s%s%s%s", CRedBg, CRed, CBold, CReset, CWhite, CBold, title, CReset)
	if detail != "" {
		fmt.Printf(" %s— %s%s", CRed, detail, CReset)
	}
	fmt.Print("\n\n")
}

func printInfo(label, value string) {
	fmt.Printf("  %s•%s %s%-18s%s %s%s%s\n", CBlue, CReset, CMuted, label, CReset, CCyan, value, CReset)
}

func ask(prompt string) string {
	fmt.Printf("  %s❯%s %s%s%s", CGlow, CReset, CText, prompt, CReset)
	text, err := reader.ReadString('\n')
	if err != nil && len(text) == 0 {
		time.Sleep(100 * time.Millisecond)
		return ""
	}
	return strings.TrimRight(text, "\r\n")
}

var sessionCancel context.CancelFunc

func askPassword(prompt string) string {
	fmt.Printf("  %s❯%s %s%s%s", CGlow, CReset, CText, prompt, CReset)
	if runtime.GOOS == "windows" {
		k32 := syscall.NewLazyDLL("kernel32.dll")
		getStdHandle := k32.NewProc("GetStdHandle")
		getConsoleMode := k32.NewProc("GetConsoleMode")
		setConsoleMode := k32.NewProc("SetConsoleMode")
		readConsoleW := k32.NewProc("ReadConsoleW")

		hIn, _, _ := getStdHandle.Call(^uintptr(9))
		var mode uint32
		if r, _, _ := getConsoleMode.Call(hIn, uintptr(unsafe.Pointer(&mode))); r != 0 {
			rawMode := mode &^ (0x0004 | 0x0002)
			setConsoleMode.Call(hIn, uintptr(rawMode))
			defer setConsoleMode.Call(hIn, uintptr(mode))

			var pass []rune
			for {
				var ch uint16
				var read uint32
				r, _, _ := readConsoleW.Call(hIn, uintptr(unsafe.Pointer(&ch)), 1, uintptr(unsafe.Pointer(&read)), 0)
				if r == 0 || read == 0 {
					break
				}
				if ch == '\r' || ch == '\n' {
					fmt.Println()
					break
				}
				if ch == '\b' || ch == 0x7f {
					if len(pass) > 0 {
						pass = pass[:len(pass)-1]
						fmt.Print("\b \b")
					}
				} else if ch >= 32 {
					pass = append(pass, rune(ch))
					fmt.Print("*")
				}
			}
			return string(pass)
		}
	}
	text, _ := reader.ReadString('\n')
	return strings.TrimRight(text, "\r\n")
}

func startRealtimeEventListener() {
	var ctx context.Context
	ctx, sessionCancel = context.WithCancel(context.Background())
	go func() {
		_ = nightAuth.SubscribeEvents(ctx, "", func(ev nightauth.EventFrame) {
			var parsed map[string]interface{}
			if json.Unmarshal([]byte(ev.Data), &parsed) == nil {
				evType, _ := parsed["type"].(string)
				if evType == "session-terminated" || evType == "session-invalid" {
					fmt.Printf("\n  %s%s%s SECURITY ALERT %s %sSession was terminated by administrator.%s\n",
						CRedBg, CRed, CBold, CReset, CWhite, CReset)
					os.Exit(0)
				}
			}
		})
	}()
}

func startHeartbeatThread() {
	go func() {
		ticker := time.NewTicker(45 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if sessionCancel == nil {
				return
			}
			nightAuth.Heartbeat()
		}
	}()
}

func padOrTrunc(s string, width int) string {
	runes := []rune(s)
	if len(runes) > width {
		return string(runes[:width])
	}
	return s + strings.Repeat(" ", width-len(runes))
}

func formatTimestamp(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" || s == "0" {
		return "-"
	}
	if unix, err := strconv.ParseInt(s, 10, 64); err == nil && unix > 0 {
		if unix > 100_000_000_000 {
			unix /= 1000
		}
		return time.Unix(unix, 0).UTC().Format("2006-01-02 15:04 UTC")
	}
	for _, layout := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05", "2006-01-02T15:04:05"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC().Format("2006-01-02 15:04 UTC")
		}
	}
	return s
}

func renderBanner(appName, hwid string, connected bool) {
	fmt.Println()
	fmt.Printf("  %s%s  _   _ _       _     _      _         _   _     \n%s", CBlue, CBold, CReset)
	fmt.Printf("  %s%s | \\ | (_) __ _| |__ | |_   / \\  _   _| |_| |__  \n%s", CBlue, CBold, CReset)
	fmt.Printf("  %s%s |  \\| | |/ _` | '_ \\| __| / _ \\| | | | __| '_ \\ \n%s", CGlow, CBold, CReset)
	fmt.Printf("  %s%s | |\\  | | (_| | | | | |_ / ___ \\ |_| | |_| | | |\n%s", CGlow, CBold, CReset)
	fmt.Printf("  %s%s |_| \\_|_|\\__, |_| |_|\\__/_/   \\_\\__,_|\\__|_| |_|\n%s", CCyan, CBold, CReset)
	fmt.Printf("  %s%s          |___/                                  %s  %s%s%s NATIVE GO CLI %s\n\n", CCyan, CBold, CReset, CBadgeBg, CGlow, CBold, CReset)

	dispName := appName
	if dispName == "" {
		dispName = "Unconfigured"
	}
	statusDot := CAmber + "◌ "
	statusText := "Initializing"
	if connected {
		statusDot = CGreen + "● "
		statusText = "Connected"
	}

	fmt.Printf("  %s╭──────────────────────────────────────────────────────────────────────────╮%s\n", CBorder, CReset)
	fmt.Printf("  %s│%s  %sAPPLICATION%s  %s%s%s%s  %sSTATUS%s  %s%-16s%s %s│%s\n",
		CBorder, CReset,
		CDim, CReset,
		CWhite, CBold, padOrTrunc(dispName, 26), CReset,
		CDim, CReset,
		statusDot, statusText, CReset,
		CBorder, CReset)
	if hwid != "" {
		fmt.Printf("  %s│%s  %sDEVICE HWID%s  %s%s%s %s│%s\n",
			CBorder, CReset,
			CDim, CReset,
			CCyan, padOrTrunc(hwid, 58), CReset,
			CBorder, CReset)
	}
	fmt.Printf("  %s╰──────────────────────────────────────────────────────────────────────────╯%s\n", CBorder, CReset)
}

func performHandshake(appName string) {
	fmt.Printf("\n  %s◌ Establishing encrypted session with NightAuth Cloud...%s\r", CMuted, CReset)
	nightAuth.Init()
	fmt.Print("                                                                      \r")

	if !nightAuth.Response.Success {
		renderBanner(appName, "", false)
		printErr("Handshake Failed", nightAuth.Response.Message)
		fmt.Printf("  %sPlease verify your App Name, API Key, and Server URL.%s\n\n", CDim, CReset)
		ask("Press Enter to exit...")
		os.Exit(1)
	}

	hwid := nightauth.SecurityProtections{}.GetHardwareId()
	clearScreen()
	renderBanner(appName, hwid, true)
}

func handleLogin() bool {
	printSectionHeader("ACCOUNT SIGN IN", "Authenticate with your username and password")
	user := ask("Username: ")
	pass := askPassword("Password: ")
	totp := ask("2FA Code " + CDim + "(Optional, Enter to skip)" + CText + ": ")

	if nightAuth.Login(user, pass, totp) {
		printOk("Session Authenticated", "Welcome back to NightAuth.")
		return true
	}

	printErr("Login Failed", nightAuth.Response.Message)
	return false
}

func handleRegister() bool {
	printSectionHeader("CREATE NEW ACCOUNT", "Activate a license key and bind your credentials")
	user := ask("Username: ")
	pass := askPassword("Password: ")
	key := ask("License Key: ")
	email := ask("Email Address " + CDim + "(Optional, Enter to skip)" + CText + ": ")

	if nightAuth.Register(user, pass, key, email) {
		printOk("Account Registered & Activated", "Session started automatically.")
		return true
	}

	detail := nightAuth.Response.Message
	if nightAuth.Response.Error != "" {
		detail += " [" + nightAuth.Response.Error + "]"
	}
	printErr("Registration Failed", detail)
	return false
}

func handleLicenseOnly() bool {
	printSectionHeader("LICENSE KEY AUTHENTICATION", "Instant passwordless login bound to your HWID")
	key := ask("Enter License Key: ")

	if nightAuth.License(key) {
		printOk("License Validated", "Direct hardware-bound session active.")
		return true
	}

	printErr("License Validation Failed", nightAuth.Response.Message)
	return false
}

func handleUpgrade(knownUser string) bool {
	printSectionHeader("UPGRADE ACCOUNT SUBSCRIPTION", "Attach a new license key tier to an account")
	user := knownUser
	if user != "" {
		printInfo("Target Account:", user)
	} else {
		user = ask("Target Username: ")
	}
	key := ask("Upgrade License Key: ")

	if nightAuth.Upgrade(user, key) {
		printOk("Upgrade Complete", nightAuth.Response.Message)
		return true
	}
	printErr("Upgrade Failed", nightAuth.Response.Message)
	return false
}

func handleForgotPassword() {
	printSectionHeader("PASSWORD RECOVERY", "Request a reset link to your registered email")
	email := ask("Registered Email: ")

	if nightAuth.ForgotPassword(email) {
		printOk("Recovery Dispatched", "If an account matches, a password reset link was sent to your email.")
	} else {
		printErr("Request Failed", nightAuth.Response.Message)
	}
}

func handleResetPassword() {
	printSectionHeader("COMPLETE PASSWORD RESET", "Confirm recovery token and set a new password")
	token := ask("Reset Token: ")
	newPass := askPassword("New Password: ")

	if nightAuth.ResetPassword(token, newPass) {
		printOk("Password Updated", "Your password has been reset successfully.")
	} else {
		printErr("Reset Failed", nightAuth.Response.Message)
	}
}

func handleResetHwid() {
	printSectionHeader("RESET HARDWARE FINGERPRINT (HWID)", "Unbind previous device to authorize a new machine")
	user := ask("Username: ")
	pass := askPassword("Password: ")
	totp := ask("2FA Code " + CDim + "(Optional, Enter to skip)" + CText + ": ")

	if nightAuth.ResetHwid(user, pass, totp) {
		printOk("HWID Reset", nightAuth.Response.Message)
	} else {
		printErr("HWID Reset Failed", nightAuth.Response.Message)
	}
}

func handleCheckBlacklist() {
	printSectionHeader("HARDWARE BLACKLIST CHECK", "Verify if a machine fingerprint is restricted")
	targetHwid := ask("Hardware ID " + CDim + "(Press Enter for local machine HWID)" + CText + ": ")
	if targetHwid == "" {
		targetHwid = nightauth.SecurityProtections{}.GetHardwareId()
	}

	isBlacklisted := nightAuth.CheckBlacklist(targetHwid)
	if nightAuth.Response.Success {
		if isBlacklisted {
			printErr("Blacklisted", "This Hardware ID is currently restricted.")
		} else {
			printOk("Clean Fingerprint", "Hardware ID is not on the blacklist.")
		}
	} else {
		printErr("Blacklist Query Failed", nightAuth.Response.Message)
	}
}

func handleVerifyHash() {
	printSectionHeader("BINARY INTEGRITY CHECK", "Validate SHA-256 checksum against server registry")
	hash := ask("SHA-256 File Hash: ")
	verified := nightAuth.VerifyFileHash(hash)
	if nightAuth.Response.Success {
		if verified {
			printOk("Integrity Verified", "Checksum matches authorized server release.")
		} else {
			printErr("Verification Failed", "Checksum does not match any approved file release.")
		}
	} else {
		printErr("Verification Failed", nightAuth.Response.Message)
	}
}

func handleAuthenticationFlow() bool {
	for {
		printSectionHeader("AUTHENTICATION & PRE-SESSION GATEWAY", "Select an operation to continue")

		fmt.Printf("  %s%s%s 01 %s %sSign In (Username & Password)%s          %s%s%s 06 %s %sComplete Password Reset (Token)%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset,
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("  %s%s%s 02 %s %sRegister New Account + License%s         %s%s%s 07 %s %sReset Hardware Binding (HWID)%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset,
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("  %s%s%s 03 %s %sInstant License Key Login%s              %s%s%s 08 %s %sCheck HWID Blacklist Status%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset,
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("  %s%s%s 04 %s %sUpgrade Existing Account Plan%s          %s%s%s 09 %s %sVerify Binary SHA-256 Hash%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset,
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("  %s%s%s 05 %s %sForgot Password (Email Link)%s           %s%s%s 10 %s %sExit Application%s\n\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset,
			CRedBg, CRed, CBold, CReset, CMuted, CReset)

		sel := strings.TrimSpace(ask("Select option [1-10]: "))

		switch sel {
		case "1", "01":
			if handleLogin() {
				return true
			}
		case "2", "02":
			if handleRegister() {
				return true
			}
		case "3", "03":
			if handleLicenseOnly() {
				return true
			}
		case "4", "04":
			handleUpgrade("")
		case "5", "05":
			handleForgotPassword()
		case "6", "06":
			handleResetPassword()
		case "7", "07":
			handleResetHwid()
		case "8", "08":
			handleCheckBlacklist()
		case "9", "09":
			handleVerifyHash()
		case "10", "0":
			return false
		}
	}
}

func renderUserProfile() {
	u := nightAuth.User
	username := u.Username
	if username == "" {
		username = "Authenticated"
	}
	status := u.Status
	if status == "" {
		status = "ACTIVE"
	}
	email := u.Email
	if email == "" {
		email = "Not linked"
	}
	ip := u.Ip
	if ip == "" {
		ip = "Protected / TLS"
	}
	created := formatTimestamp(u.CreatedAt)
	if created == "" {
		created = "-"
	}
	lastLogin := formatTimestamp(u.LastLoginAt)
	if lastLogin == "" {
		lastLogin = "-"
	}

	fmt.Println()
	fmt.Printf("  %s╭─ %s%sAUTHENTICATED IDENTITY & ENTITLEMENTS%s ──────────────────────────────────╮%s\n",
		CBorder, CGlow, CBold, CBorder, CReset)
	fmt.Printf("  %s│%s  %sUsername :%s %s%s%s%s  %sStatus     :%s %s%s%s%s %s│%s\n",
		CBorder, CReset,
		CDim, CReset, CWhite, CBold, padOrTrunc(username, 24), CReset,
		CDim, CReset, CGreen, CBold, padOrTrunc(status, 20), CReset,
		CBorder, CReset)
	fmt.Printf("  %s│%s  %sEmail    :%s %s%s%s  %sClient IP  :%s %s%s%s %s│%s\n",
		CBorder, CReset,
		CDim, CReset, CText, padOrTrunc(email, 24), CReset,
		CDim, CReset, CCyan, padOrTrunc(ip, 20), CReset,
		CBorder, CReset)
	fmt.Printf("  %s│%s  %sCreated  :%s %s%s%s  %sLast Login :%s %s%s%s %s│%s\n",
		CBorder, CReset,
		CDim, CReset, CMuted, padOrTrunc(created, 24), CReset,
		CDim, CReset, CMuted, padOrTrunc(lastLogin, 20), CReset,
		CBorder, CReset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────┤%s\n", CBorder, CReset)

	if len(u.Subscriptions) > 0 {
		for _, sub := range u.Subscriptions {
			key := sub.Key
			if key == "" {
				key = "Active Account"
			}
			exp := formatTimestamp(sub.Expiry)
			if exp == "" {
				exp = "-"
			}
			fmt.Printf("  %s│%s  %s●%s %s%s%s%s %sKey:%s %s%s%s %sExp:%s %s%s%s %s│%s\n",
				CBorder, CReset,
				CGreen, CReset,
				CWhite, CBold, padOrTrunc(sub.Subscription, 16), CReset,
				CDim, CReset, CCyan, padOrTrunc(key, 22), CReset,
				CDim, CReset, CAmber, padOrTrunc(exp, 16), CReset,
				CBorder, CReset)
		}
	} else {
		fmt.Printf("  %s│%s  %sNo active subscription tiers attached to this account.%s                   %s│%s\n",
			CBorder, CReset, CDim, CReset, CBorder, CReset)
	}
	fmt.Printf("  %s╰──────────────────────────────────────────────────────────────────────────╯%s\n", CBorder, CReset)
}

func handleChatRoom() {
	printSectionHeader("LIVE CHANNEL CHAT", "Type /exit at any time to leave the channel")
	room := ask("Channel name (e.g. general): ")
	if room == "" {
		room = "general"
	}

	msgs := nightAuth.GetChatMessagesPage(room, 0, true)
	fmt.Printf("\n  %s╭─ %s%s#%s%s ──────────────────────────────────────────────────╮%s\n",
		CBorder, CGlow, CBold, padOrTrunc(room, 20), CBorder, CReset)
	if len(msgs) == 0 {
		fmt.Printf("  %s│%s  %sNo messages in this channel yet. Start the conversation below.%s          %s│%s\n",
			CBorder, CReset, CDim, CReset, CBorder, CReset)
	} else {
		for _, m := range msgs {
			fmt.Printf("  %s│%s %s[%s]%s %s%s%s%s: %s%s%s\n",
				CBorder, CReset, CDim, m.Timestamp, CReset, CCyan, CBold, m.Author, CReset, CText, m.Message, CReset)
		}
	}
	fmt.Printf("  %s╰──────────────────────────────────────────────────────────────────────────╯%s\n\n", CBorder, CReset)

	for {
		input := ask(CCyan + "#" + room + CReset + " message: ")
		if strings.EqualFold(input, "/exit") {
			fmt.Printf("  %sLeft channel #%s.%s\n\n", CDim, room, CReset)
			return
		}
		if input == "" {
			continue
		}
		if !nightAuth.SendChatMessage(input, room) {
			printErr("Send Failed", nightAuth.Response.Message)
		}
	}
}

func renderStatsTable() {
	stats := nightAuth.GetStats()
	if stats != nil {
		ver := stats.Version
		if ver == "" {
			ver = AppVersion
		}
		fmt.Println()
		fmt.Printf("  %s╭─ %s%sAPPLICATION LIVE TELEMETRY%s ─────────────────────────────────────────────╮%s\n",
			CBorder, CGlow, CBold, CBorder, CReset)
		fmt.Printf("  %s│%s  %sTotal Users     :%s %s%s%-14d%s  %sTotal Licenses  :%s %s%s%-18d%s %s│%s\n",
			CBorder, CReset, CDim, CReset, CWhite, CBold, stats.TotalUsers, CReset, CDim, CReset, CWhite, CBold, stats.TotalLicenses, CReset, CBorder, CReset)
		fmt.Printf("  %s│%s  %sActive Users    :%s %s%s%-14d%s  %sUsed Licenses   :%s %s%s%-18d%s %s│%s\n",
			CBorder, CReset, CDim, CReset, CGreen, CBold, stats.ActiveUsers, CReset, CDim, CReset, CCyan, CBold, stats.UsedLicenses, CReset, CBorder, CReset)
		fmt.Printf("  %s│%s  %sBanned Users    :%s %s%s%-14d%s  %sUnused Licenses :%s %s%s%-18d%s %s│%s\n",
			CBorder, CReset, CDim, CReset, CRed, CBold, stats.BannedUsers, CReset, CDim, CReset, CAmber, CBold, stats.UnusedLicenses, CReset, CBorder, CReset)
		fmt.Printf("  %s│%s  %sTarget Version  :%s %s%s%s %s│%s\n",
			CBorder, CReset, CDim, CReset, CGlow, padOrTrunc(ver, 52), CReset, CBorder, CReset)
		fmt.Printf("  %s╰──────────────────────────────────────────────────────────────────────────╯%s\n\n", CBorder, CReset)
	} else {
		printErr("Telemetry Error", nightAuth.Response.Message)
	}
}

func handleSetup2FA() {
	printSectionHeader("TWO-FACTOR AUTHENTICATION (TOTP)", "Enroll an authenticator app (RFC 6238)")
	setup := nightAuth.Setup2FA()
	if setup != nil {
		printInfo("TOTP Secret Key:", setup.Secret)
		printInfo("Provisioning URI:", setup.Uri)
		code := ask("Enter 6-digit Authenticator code: ")
		conf := nightAuth.Confirm2FA(code)
		if conf != nil {
			printOk("2FA Activated", "Store these recovery backup codes safely:")
			for _, bc := range conf.BackupCodes {
				fmt.Printf("    %s• %s%s\n", CAmber, bc, CReset)
			}
			fmt.Println()
		} else {
			printErr("Verification Failed", nightAuth.Response.Message)
		}
	} else {
		printErr("2FA Setup Failed", nightAuth.Response.Message)
	}
}

func handleDownloadFile() {
	printSectionHeader("PROTECTED ASSET DOWNLOADER", "Download encrypted files via scoped access ticket")
	fileId := ask("Enter File ID or Asset Name: ")
	content := nightAuth.DownloadFile(fileId)
	if content != nil {
		printOk("Download Complete", fmt.Sprintf("Retrieved %d bytes securely.", len(content)))
		for i := range content {
			content[i] = 0
		}
	} else {
		printErr("Download Failed", nightAuth.Response.Message)
	}
}

func handleListAllUserVars() {
	if !nightAuth.GetUserVariables() {
		printErr("Fetch Failed", nightAuth.Response.Message)
		return
	}
	printSectionHeader("USER CLOUD VARIABLES", "Per-user encrypted key-value pairs")
	rawVars, _ := nightAuth.Response.Payload["variables"].([]interface{})
	if len(rawVars) == 0 {
		fmt.Printf("  %sNo user variables stored for this account yet.%s\n\n", CDim, CReset)
		return
	}
	for _, item := range rawVars {
		if m, ok := item.(map[string]interface{}); ok {
			k, _ := m["key"].(string)
			v, _ := m["value"].(string)
			fmt.Printf("  %s•%s %s%s%-22s%s = %s%s%s\n", CBlue, CReset, CCyan, CBold, k, CReset, CWhite, v, CReset)
		}
	}
	fmt.Println()
}

func handleTriggerWebhook() {
	printSectionHeader("SERVER-SIDE WEBHOOK TRIGGER", "Execute an authenticated webhook endpoint")
	webId := ask("Webhook ID: ")
	payloadStr := ask("JSON Payload " + CDim + "(Enter for {})" + CText + ": ")
	if payloadStr == "" {
		payloadStr = "{}"
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		printErr("Invalid Webhook Payload", err.Error())
		return
	}

	res := nightAuth.TriggerWebhook(webId, payload)
	if res != nil {
		printOk("Webhook Dispatched", fmt.Sprintf("HTTP %d OK", res.Status))
	} else {
		printErr("Webhook Failed", nightAuth.Response.Message)
	}
}

func handleEventsFeed() {
	printSectionHeader("REAL-TIME APPLICATION EVENTS (SSE)", "Press Enter at any time to stop listening")
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		_, _ = reader.ReadString('\n')
		cancel()
	}()

	err := nightAuth.SubscribeEvents(ctx, "", func(ev nightauth.EventFrame) {
		var parsed map[string]interface{}
		evType := "event"
		if json.Unmarshal([]byte(ev.Data), &parsed) == nil {
			if t, ok := parsed["type"].(string); ok && t != "" {
				evType = t
			}
		}
		fmt.Printf("  %s%s EVENT %s %s[%s]%s %s%s%s%s: %s\n",
			CBadgeBg, CGlow, CReset, CDim, ev.ID, CReset, CCyan, CBold, evType, CReset, ev.Data)
	})

	cancel()
	if err == nil || ctx.Err() != nil {
		printOk("Stream Closed", "Disconnected from real-time event feed.")
	} else {
		printErr("Stream Interrupted", nightAuth.Response.Message)
	}
}

func handleSelfBan() {
	printSectionHeader("EMERGENCY KILLSWITCH / BAN", "Report security violation and lock this session")
	reason := ask("Ban Reason " + CDim + "(Optional)" + CText + ": ")

	if nightAuth.Ban(reason) {
		printOk("Killswitch Executed", nightAuth.Response.Message)
	} else {
		printErr("Ban Report Failed", nightAuth.Response.Message)
	}
}

func handleHeartbeat() {
	if nightAuth.Heartbeat() {
		printOk("Heartbeat Acknowledged", nightAuth.Response.Message)
	} else {
		printErr("Heartbeat Failed", nightAuth.Response.Message)
	}
}

func handleKillSessionDemo() {
	printSectionHeader("SESSION REVOCATION & INVALIDATION VERIFICATION", "3-Step Live Security Audit")
	fmt.Printf("  %s1. Verifying active session before revocation...%s\n", CMuted, CReset)
	if nightAuth.Check() {
		printOk("Step 1 Passed", "Session is active and authenticated on server.")
	} else {
		printErr("Step 1 Failed", nightAuth.Response.Message)
		return
	}

	fmt.Printf("  %s2. Revoking and terminating session on server (/logout)...%s\n", CMuted, CReset)
	if nightAuth.Logout() {
		printOk("Step 2 Passed", "Server confirmed session revocation.")
	} else {
		printErr("Step 2 Failed", nightAuth.Response.Message)
	}

	fmt.Printf("  %s3. Verifying that subsequent session check is rejected...%s\n", CMuted, CReset)
	if nightAuth.Check() {
		printErr("Step 3 Unexpected", "Session was still accepted after revocation.")
	} else {
		printOk("Step 3 Passed", "Server rejected revoked session as expected.")
	}
}

func runApplicationFeaturesMenu() {
	for {
		printSectionHeader("WORKSPACE MODULES & ACTIONS", "Select a command [1-21]")

		fmt.Printf("  %sCLOUD VARIABLES & CHAT%s                    %sSECURITY & TELEMETRY%s\n", CDim, CReset, CDim, CReset)
		fmt.Printf("  %s%s%s 01 %s %sRead Global Variable%s                 %s%s%s 10 %s %sSetup Two-Factor Auth (2FA)%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset, CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("  %s%s%s 02 %s %sGet User Cloud Variable%s              %s%s%s 11 %s %sDisable Two-Factor Auth (2FA)%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset, CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("  %s%s%s 03 %s %sSet User Cloud Variable%s              %s%s%s 12 %s %sDownload Protected Asset%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset, CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("  %s%s%s 04 %s %sList All User Variables%s              %s%s%s 13 %s %sTrigger Server Webhook%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset, CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("  %s%s%s 05 %s %sEnter Live Chat Channel%s              %s%s%s 14 %s %sSend Audit Log Event%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset, CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("  %s%s%s 19 %s %sFetch Paginated Chat History%s         %s%s%s 15 %s %sVerify Session Entitlement%s\n\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset, CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)

		fmt.Printf("  %sACCOUNT & INSIGHTS%s                        %sSESSION LIFECYCLE%s\n", CDim, CReset, CDim, CReset)
		fmt.Printf("  %s%s%s 06 %s %sApplication Live Statistics%s          %s%s%s 16 %s %sSend Heartbeat Ping%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset, CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("  %s%s%s 07 %s %sView Online Active Users%s             %s%s%s 17 %s %sEmergency Killswitch / Ban%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset, CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("  %s%s%s 08 %s %sChange Account Username%s              %s%s%s 18 %s %sStream Live Events (SSE)%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset, CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("  %s%s%s 09 %s %sUpgrade Plan with License Key%s        %s%s%s 20 %s %sSession Revocation Audit Demo%s\n",
			CBadgeBg, CGlow, CBold, CReset, CWhite, CReset, CBadgeBg, CGlow, CBold, CReset, CWhite, CReset)
		fmt.Printf("                                            %s%s%s 21 %s %sSign Out & Terminate Session%s\n\n",
			CRedBg, CRed, CBold, CReset, CMuted, CReset)

		choice := strings.TrimSpace(ask("Action [1-21]: "))
		num, _ := strconv.Atoi(choice)

		switch num {
		case 1:
			gKey := ask("Global Variable Key: ")
			gVal := nightAuth.GetGlobalVar(gKey)
			if gVal != "" {
				printOk(gKey, gVal)
			} else {
				printErr("Fetch Failed", nightAuth.Response.Message)
			}
		case 2:
			uKey := ask("User Variable Key: ")
			uVal := nightAuth.GetUserVar(uKey)
			if uVal != "" {
				printOk(uKey, uVal)
			} else {
				printErr("Fetch Failed", nightAuth.Response.Message)
			}
		case 3:
			sKey := ask("User Variable Key: ")
			sVal := ask("Value to Store: ")
			if nightAuth.SetUserVar(sKey, sVal) {
				printOk("Variable Saved", sKey)
			} else {
				printErr("Save Failed", nightAuth.Response.Message)
			}
		case 4:
			handleListAllUserVars()
		case 5:
			handleChatRoom()
		case 6:
			renderStatsTable()
		case 7:
			users := nightAuth.GetOnlineUsers()
			hdr := fmt.Sprintf("%d active session(s)", len(users))
			printSectionHeader("ONLINE ACTIVE SESSIONS", hdr)
			for _, u := range users {
				fmt.Printf("  %s●%s %s%s%s\n", CGreen, CReset, CWhite, u, CReset)
			}
			fmt.Println()
		case 8:
			newU := ask("New Username: ")
			p1 := askPassword("Confirm Password: ")
			if nightAuth.ChangeUsername(newU, p1) {
				printOk("Username Updated", nightAuth.User.Username)
				fmt.Printf("  %sFor security, your session has been closed. Please sign in with your new username.%s\n\n", CAmber, CReset)
				return
			}
			printErr("Update Failed", nightAuth.Response.Message)
		case 9:
			if handleUpgrade(nightAuth.User.Username) {
				nightAuth.Check()
				renderUserProfile()
			}
		case 10:
			handleSetup2FA()
		case 11:
			p3 := askPassword("Account Password: ")
			code3 := ask("6-digit Authenticator Code: ")
			if nightAuth.Disable2FA(p3, code3) {
				printOk("2FA Disabled", "Two-factor protection has been removed from this account.")
			} else {
				printErr("Disable Failed", nightAuth.Response.Message)
			}
		case 12:
			handleDownloadFile()
		case 13:
			handleTriggerWebhook()
		case 14:
			logMsg := ask("Audit Log Message: ")
			if nightAuth.Log(logMsg) {
				printOk("Audit Log Recorded", "Delivered to application telemetry feed.")
			} else {
				printErr("Log Failed", nightAuth.Response.Message)
			}
		case 15:
			if nightAuth.Check() {
				printOk("Session Verified", "Active and cryptographically valid.")
				renderUserProfile()
			} else {
				printErr("Session Invalid", nightAuth.Response.Message)
			}
		case 16:
			handleHeartbeat()
		case 17:
			handleSelfBan()
		case 18:
			handleEventsFeed()
		case 19:
			histRoom := ask("Channel Name: ")
			histBefore := ask("Before Message ID " + CDim + "(Optional, Enter for latest)" + CText + ": ")
			histMr := ask("Mark messages as read? (y/n): ")
			var beforeId int64
			if histBefore != "" {
				beforeId, _ = strconv.ParseInt(histBefore, 10, 64)
			}
			markRead := strings.HasPrefix(strings.ToLower(histMr), "y")

			hist := nightAuth.GetChatMessagesPage(histRoom, beforeId, markRead)
			if !nightAuth.Response.Success {
				printErr("History Query Failed", nightAuth.Response.Message)
			} else {
				sub := fmt.Sprintf("%d message(s) retrieved", len(hist))
				printSectionHeader(histRoom, sub)
				for _, m := range hist {
					fmt.Printf("  %s[%s]%s %s%s%s%s: %s%s%s %s(#%s · -)%s\n",
						CDim, m.Timestamp, CReset, CCyan, CBold, m.Author, CReset, CText, m.Message, CReset, CDim, m.Id, CReset)
				}
				fmt.Println()
			}
		case 20:
			handleKillSessionDemo()
			return
		case 21:
			nightAuth.Logout()
			printOk("Signed Out", "Session revoked cleanly. Goodbye!")
			return
		}
	}
}

func main() {
	enableVirtualTerminal()

	appName := strings.TrimSpace(os.Getenv("NIGHTAUTH_APP_NAME"))
	if appName == "" {
		appName = AppName
	}
	apiKey := strings.TrimSpace(os.Getenv("NIGHTAUTH_API_KEY"))
	if apiKey == "" {
		apiKey = AppApiKey
	}
	version := strings.TrimSpace(os.Getenv("NIGHTAUTH_APP_VERSION"))
	if version == "" {
		version = AppVersion
	}
	apiUrl := strings.TrimSpace(os.Getenv("NIGHTAUTH_API_URL"))
	if apiUrl == "" {
		apiUrl = nightauth.DefaultAPIURL
	}
	allowLoopback := os.Getenv("NIGHTAUTH_ALLOW_INSECURE_LOOPBACK") == "1"

	token := strings.TrimSpace(os.Getenv("NIGHTAUTH_APP_TOKEN"))
	if token == "" {
		token = AppToken
	}

	hwid := nightauth.SecurityProtections{}.GetHardwareId()
	clearScreen()
	renderBanner(appName, hwid, false)

	nightAuth = nightauth.NewApiWithOptions(appName, apiKey, version, token, apiUrl, allowLoopback)

	performHandshake(appName)

	if !handleAuthenticationFlow() {
		nightAuth.Close()
		return
	}

	hwid = nightauth.SecurityProtections{}.GetHardwareId()
	clearScreen()
	renderBanner(appName, hwid, true)
	renderUserProfile()

	startRealtimeEventListener()
	startHeartbeatThread()
	runApplicationFeaturesMenu()
	if sessionCancel != nil {
		sessionCancel()
	}
	nightAuth.Close()
}
