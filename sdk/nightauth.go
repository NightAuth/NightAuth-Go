package nightauth

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

const (
	DefaultAPIURL             = "https://nightauth.org/api/v1"
	DefaultSigningPublicKey   = "cdb2285db481666ce40003aadb874c8751dbeb9a0672a37f740dfe881ecdc742"
	DefaultSealedBoxPublicKey = "b773fc4f6cf1a0938e78e714c3d84cb386584bea24728923008471a6cd203213"
	ServerSigningPublicKey    = DefaultSigningPublicKey
	ServerSealedBoxPublicKey  = DefaultSealedBoxPublicKey
)

const (
	maxFileBytes int64 = 64 * 1024 * 1024

	maxDownloadResponseBytes int64 = maxFileBytes*4/3 + 1024*1024
)

var (
	hkdfInfo = []byte("nightauth-v1-e2e")
	sealInfo = []byte("nightauth-v1-init-seal")
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`

	StatusCode int `json:"-"`

	RetryAfter float64 `json:"-"`

	Payload map[string]interface{} `json:"-"`
}

type SubscriptionData struct {
	Subscription string `json:"subscription"`
	Key          string `json:"key"`
	Expiry       string `json:"expiry"`
	TimeLeft     int64  `json:"timeLeft"`
}

type UserData struct {
	Username         string             `json:"username"`
	Email            string             `json:"email"`
	Hwid             string             `json:"hwid"`
	Ip               string             `json:"ip"`
	Status           string             `json:"status"`
	CreatedAt        string             `json:"createdAt"`
	LastLoginAt      string             `json:"lastLoginAt"`
	TwoFactorEnabled bool               `json:"twoFactorEnabled"`
	Subscriptions    []SubscriptionData `json:"subscriptions"`
}

type AppData struct {
	TotalUsers  int64  `json:"totalUsers"`
	OnlineUsers int64  `json:"onlineUsers"`
	TotalKeys   int64  `json:"totalKeys"`
	Version     string `json:"version"`
}

type AppStats struct {
	TotalUsers     int64  `json:"totalUsers"`
	ActiveUsers    int64  `json:"activeUsers"`
	BannedUsers    int64  `json:"bannedUsers"`
	TotalLicenses  int64  `json:"totalLicenses"`
	UsedLicenses   int64  `json:"usedLicenses"`
	UnusedLicenses int64  `json:"unusedLicenses"`
	Version        string `json:"version"`
}

type ChatMessage struct {
	Id        string `json:"id"`
	Author    string `json:"author"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type TwoFactorSetupData struct {
	Secret    string `json:"secret"`
	QrDataUrl string `json:"qrDataUrl"`
	Uri       string `json:"uri"`
}

type TwoFactorConfirmData struct {
	Success     bool     `json:"success"`
	BackupCodes []string `json:"backupCodes"`
}

type SecurityProtections struct{}

var blacklistedProcesses = []string{
	"fiddler", "wireshark", "httpdebugger", "httpdebuggerui", "x64dbg", "x32dbg",
	"ida", "ida64", "charles", "dnspy", "de4dot", "scylla", "processhacker",
	"cheatengine", "ollydbg",
}

func selfTerminate(reason string) {
	if runtime.GOOS == "windows" {
		cleanReason := strings.NewReplacer(`"`, `'`, `&`, ` `, `|`, ` `).Replace(reason)
		cmdArgs := fmt.Sprintf(`/c title NightAuth Security Alert & color 4f & echo ========================================================= & echo [NIGHTAUTH SECURITY PROTECTION TRIGGERED] & echo ========================================================= & echo. & echo Reason: %s & echo. & echo Process has been terminated immediately to prevent memory dumps. & echo This window will close automatically in 5 seconds... & timeout /t 5 >nul`, cleanReason)
		cmd := exec.Command("cmd.exe", "/c", cmdArgs)
		_ = cmd.Start()
	}
	os.Exit(1)
}

func (SecurityProtections) RunSecurityChecks() {
	RunSecurityChecks()
}

func RunSecurityChecks() {
	CheckSystemProxy()
}

func (SecurityProtections) CheckSystemProxy() {
	CheckSystemProxy()
}

func CheckSystemProxy() {
	if os.Getenv("HTTP_PROXY") != "" || os.Getenv("HTTPS_PROXY") != "" ||
		os.Getenv("http_proxy") != "" || os.Getenv("https_proxy") != "" ||
		os.Getenv("ALL_PROXY") != "" || os.Getenv("all_proxy") != "" {
		selfTerminate("Active HTTP proxy detected.")
		return
	}
	if runtime.GOOS == "windows" {
		advapi32 := syscall.NewLazyDLL("advapi32.dll")
		regOpenKeyExW := advapi32.NewProc("RegOpenKeyExW")
		regQueryValueExW := advapi32.NewProc("RegQueryValueExW")
		regCloseKey := advapi32.NewProc("RegCloseKey")

		subKey, _ := syscall.UTF16PtrFromString(`Software\Microsoft\Windows\CurrentVersion\Internet Settings`)
		var hKey uintptr
		// HKEY_CURRENT_USER = 0x80000001, KEY_READ = 0x20019
		ret, _, _ := regOpenKeyExW.Call(0x80000001, uintptr(unsafe.Pointer(subKey)), 0, 0x20019, uintptr(unsafe.Pointer(&hKey)))
		if ret == 0 {
			defer regCloseKey.Call(hKey)
			valName, _ := syscall.UTF16PtrFromString("ProxyEnable")
			var valType, valData, valSize uint32
			valSize = 4
			ret2, _, _ := regQueryValueExW.Call(hKey, uintptr(unsafe.Pointer(valName)), 0, uintptr(unsafe.Pointer(&valType)), uintptr(unsafe.Pointer(&valData)), uintptr(unsafe.Pointer(&valSize)))
			if ret2 == 0 && valData != 0 {
				selfTerminate("Active HTTP proxy detected.")
			}
		}
	}
}

func (SecurityProtections) CheckDebugger() {
	CheckDebugger()
}

func CheckDebugger() {
	if runtime.GOOS != "windows" {
		return
	}
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	pIsDebuggerPresent := kernel32.NewProc("IsDebuggerPresent")
	pCheckRemote := kernel32.NewProc("CheckRemoteDebuggerPresent")
	pGetCurrentProcess := kernel32.NewProc("GetCurrentProcess")

	r, _, _ := pIsDebuggerPresent.Call()
	if r != 0 {
		selfTerminate("Debugger detected.")
		return
	}

	hProc, _, _ := pGetCurrentProcess.Call()
	var isRemote int32
	r2, _, _ := pCheckRemote.Call(hProc, uintptr(unsafe.Pointer(&isRemote)))
	if r2 != 0 && isRemote != 0 {
		selfTerminate("Remote debugger detected.")
		return
	}
}

func (SecurityProtections) CheckBlacklistedProcesses() {
	CheckBlacklistedProcesses()
}

func CheckBlacklistedProcesses() {
	if os.Getenv("NIGHTAUTH_IGNORE_SECURITY") == "1" {
		return
	}
	if runtime.GOOS != "windows" {
		return
	}
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	pCreateToolhelp32Snapshot := kernel32.NewProc("CreateToolhelp32Snapshot")
	pProcess32FirstW := kernel32.NewProc("Process32FirstW")
	pProcess32NextW := kernel32.NewProc("Process32NextW")
	pCloseHandle := kernel32.NewProc("CloseHandle")

	// TH32CS_SNAPPROCESS = 0x00000002
	hSnap, _, _ := pCreateToolhelp32Snapshot.Call(0x00000002, 0)
	if hSnap == 0 || hSnap == ^uintptr(0) {
		return
	}
	defer pCloseHandle.Call(hSnap)

	type processEntry32W struct {
		DwSize              uint32
		CntUsage            uint32
		Th32ProcessID       uint32
		Th32DefaultHeapID   uintptr
		Th32ModuleID        uint32
		CntThreads          uint32
		Th32ParentProcessID uint32
		PcPriClassBase      int32
		DwFlags             uint32
		SzExeFile           [260]uint16
	}

	var pe processEntry32W
	pe.DwSize = uint32(unsafe.Sizeof(pe))

	ret, _, _ := pProcess32FirstW.Call(hSnap, uintptr(unsafe.Pointer(&pe)))
	for ret != 0 {
		exeName := syscall.UTF16ToString(pe.SzExeFile[:])
		lowerExe := strings.ToLower(exeName)
		for _, bad := range blacklistedProcesses {
			if strings.Contains(lowerExe, bad) {
				selfTerminate(fmt.Sprintf("Blacklisted analysis tool detected: %s", exeName))
				return
			}
		}
		ret, _, _ = pProcess32NextW.Call(hSnap, uintptr(unsafe.Pointer(&pe)))
	}
}

var (
	cachedHwid     string
	cachedHwidLock sync.Mutex
)

func (SecurityProtections) GetHardwareId() string {
	return GetHardwareId()
}

func GetHardwareId() string {
	cachedHwidLock.Lock()
	defer cachedHwidLock.Unlock()
	if cachedHwid != "" {
		return cachedHwid
	}
	hostname := os.Getenv("COMPUTERNAME")
	if hostname == "" {
		hostname = os.Getenv("HOSTNAME")
	}
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	if hostname == "" {
		hostname = "unknown_host"
	}
	hostname = strings.ToUpper(strings.TrimSpace(hostname))

	username := os.Getenv("USERNAME")
	if username == "" {
		username = os.Getenv("USER")
	}
	if username == "" {
		username = "unknown_user"
	}
	username = strings.ToLower(strings.TrimSpace(username))

	raw := fmt.Sprintf("%s-%s-nightauth-hwid", hostname, username)
	hash := sha256.Sum256([]byte(raw))
	cachedHwid = strings.ToLower(hex.EncodeToString(hash[:]))
	return cachedHwid
}

type Api struct {
	Name    string
	ApiKey  string
	Version string
	Token   string
	ApiUrl  string

	HardwareIdProvider func() string

	Response Response
	User     UserData
	App      AppData

	PendingDownloadUrl string
	RequiredVersion    string

	OnSessionInvalidated  func(reason string)
	OnAppBanned           func(reason string)
	OnAppStatusChanged    func(status string)
	OnVersionUpdated      func(version string)
	OnChatMessageReceived func(msg ChatMessage, room string)

	sessionKey            string
	sharedKey             []byte
	isInitialized         bool
	httpClient            *http.Client
	allowInsecureLoopback bool

	signingPin   string
	sealPin      string
	stopChan     chan struct{}
	cancelEvents context.CancelFunc
	mu           sync.Mutex
	disposed     bool
}

func validateApiUrl(rawUrl string, allowInsecureLoopback ...bool) error {
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return fmt.Errorf("invalid api url: %w", err)
	}
	host := parsed.Hostname()
	isLoopback := host == "localhost" || host == "127.0.0.1" || host == "::1"
	if parsed.Scheme != "https" && !(parsed.Scheme == "http" && isLoopback && len(allowInsecureLoopback) > 0 && allowInsecureLoopback[0]) {
		return errors.New("HTTPS is required for non-loopback hosts")
	}
	return nil
}

func NewApi(name, apiKey string, version, token, apiUrl string) *Api {
	return NewApiWithOptions(name, apiKey, version, token, apiUrl, false)
}

func NewApiWithOptions(name, apiKey string, version, token, apiUrl string, allowInsecureLoopback bool) *Api {
	SecurityProtections{}.RunSecurityChecks()

	v := version
	if v == "" {
		v = "1.0"
	}
	u := strings.TrimRight(apiUrl, "/")
	if u == "" {
		u = DefaultAPIURL
	}
	if err := validateApiUrl(u, allowInsecureLoopback); err != nil {
		panic(err)
	}

	return &Api{
		Name:                  name,
		ApiKey:                apiKey,
		Version:               v,
		Token:                 token,
		ApiUrl:                u,
		allowInsecureLoopback: allowInsecureLoopback,
		signingPin:            ServerSigningPublicKey,
		sealPin:               ServerSealedBoxPublicKey,
		httpClient:            &http.Client{Timeout: 20 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }},
		stopChan:              make(chan struct{}),
	}
}

func (a *Api) SetRequestTimeout(d time.Duration) error {
	if d <= 0 {
		return errors.New("request timeout must be positive")
	}
	a.httpClient.Timeout = d
	return nil
}

func (a *Api) SetServerKeys(signingKeyHex, sealedBoxKeyHex string) error {
	for _, k := range []string{signingKeyHex, sealedBoxKeyHex} {
		if b, err := hex.DecodeString(k); err != nil || len(b) != 32 {
			return errors.New("server keys must be 64 hex characters")
		}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.signingPin = strings.ToLower(signingKeyHex)
	a.sealPin = strings.ToLower(sealedBoxKeyHex)
	return nil
}

func (a *Api) pinnedSigningKey() string {
	if a.signingPin == "" {
		return ServerSigningPublicKey
	}
	return a.signingPin
}

func (a *Api) pinnedSealKey() string {
	if a.sealPin == "" {
		return ServerSealedBoxPublicKey
	}
	return a.sealPin
}

func (a *Api) Init() {
	SecurityProtections{}.RunSecurityChecks()

	var clientPriv [32]byte
	_, _ = rand.Read(clientPriv[:])
	var clientPub [32]byte
	curve25519.ScalarBaseMult(&clientPub, &clientPriv)

	k1Hex := strings.ToLower(hex.EncodeToString(clientPub[:]))

	a.sessionKey = ""
	a.sharedKey = nil
	a.isInitialized = false
	sessionIdBytes := make([]byte, 16)
	_, _ = rand.Read(sessionIdBytes)
	newSessionKey := hex.EncodeToString(sessionIdBytes)

	payload := map[string]interface{}{
		"apikey":  a.ApiKey,
		"name":    a.Name,
		"version": a.Version,
		"id":      newSessionKey,
		"k1":      k1Hex,
	}

	payloadBytes, _ := json.Marshal(payload)
	sealed, err := cryptoBoxSeal(string(payloadBytes), a.pinnedSealKey())
	if err != nil {
		a.setError(fmt.Sprintf("Init encryption failed: %v", err), "init_encryption_failed")
		return
	}

	reqBody, _ := json.Marshal(map[string]string{"p": sealed})
	status, body, retry, err := a.postRaw("/init", reqBody)
	if err != nil {
		code := classifyTransport(err)
		if errors.Is(err, errResponseTooLarge) {
			code = "response_too_large"
		}
		a.apply(localFailure(status, code, err.Error()))
		return
	}

	var respJson map[string]interface{}
	if err := json.Unmarshal(body, &respJson); err != nil || respJson == nil {
		a.apply(localFailure(status, "invalid_response", "Response body is not a JSON object."))
		return
	}
	respJson["http_status"] = status
	if retry > 0 {
		respJson["retry_after"] = retry
	}
	if !a.apply(respJson) {
		return
	}
	if status != http.StatusOK {
		a.apply(localFailure(status, "invalid_response", "Initialization answered without success."))
		return
	}

	k2Hex, _ := respJson["k2"].(string)
	vk, _ := respJson["vk"].(string)
	tsNum, _ := respJson["ts"].(float64)
	ts := int64(tsNum)
	sigB64, _ := respJson["sig"].(string)
	nowMs := time.Now().UnixMilli()
	if len(k2Hex) != 64 || !strings.EqualFold(vk, a.pinnedSigningKey()) || ts < nowMs-90000 || ts > nowMs+90000 {
		a.setError("Untrusted or stale initialization response.", "untrusted_response")
		return
	}
	if !verifyServerSignature(k2Hex, a.pinnedSigningKey(), ts, sigB64) {
		a.setError("Security integrity breach: Invalid server signature.", "untrusted_response")
		return
	}

	k2Bytes, err := hex.DecodeString(k2Hex)
	if err != nil || len(k2Bytes) != 32 {
		a.setError("Invalid server public key", "untrusted_response")
		return
	}

	var serverPub [32]byte
	copy(serverPub[:], k2Bytes)

	var rawSecret [32]byte
	curve25519.ScalarMult(&rawSecret, &clientPriv, &serverPub)

	hk := hkdf.New(sha256.New, rawSecret[:], nil, hkdfInfo)
	okm := make([]byte, 32)
	if _, err := io.ReadFull(hk, okm); err != nil {
		a.setError("Key derivation failed", "")
		return
	}

	a.sharedKey = okm
	a.sessionKey = newSessionKey
	a.isInitialized = true
	a.setSuccess("Initialized successfully.")
}

func (a *Api) Login(username, password string, totpCode string) bool {
	SecurityProtections{}.RunSecurityChecks()

	body := a.createBasePayload()
	body["version"] = a.Version
	body["username"] = username
	body["password"] = password
	body["hwid"] = a.hardwareID()
	if totpCode != "" {
		body["totpCode"] = totpCode
	}

	resp := a.sendEncrypted("/login", body)
	if !a.apply(resp) {
		return false
	}

	a.populateUserData(resp)
	a.setSuccess("Logged in successfully.")
	return true
}

func (a *Api) Register(username, password, key, email string) bool {
	SecurityProtections{}.RunSecurityChecks()

	body := a.createBasePayload()
	body["version"] = a.Version
	body["username"] = username
	body["password"] = password
	body["key"] = key
	if email != "" {
		body["email"] = email
	}
	body["hwid"] = a.hardwareID()

	resp := a.sendEncrypted("/register", body)
	if !a.apply(resp) {
		return false
	}

	a.populateUserData(resp)
	a.setSuccess("Registered successfully.")
	return true
}

func (a *Api) License(key string) bool {
	return a.LicenseWithTotp(key, "")
}

func (a *Api) LicenseWithTotp(key, totpCode string) bool {
	SecurityProtections{}.RunSecurityChecks()

	body := a.createBasePayload()
	body["version"] = a.Version
	body["key"] = key
	body["hwid"] = a.hardwareID()
	if totpCode != "" {
		body["totpCode"] = totpCode
	}

	resp := a.sendEncrypted("/license", body)
	if !a.apply(resp) {
		return false
	}

	a.populateUserData(resp)
	a.setSuccess("License verified successfully.")
	return true
}

func (a *Api) Upgrade(username, key string) bool {
	body := a.createBasePayload()

	if strings.TrimSpace(username) != "" {
		body["username"] = username
	}
	body["key"] = key

	resp := a.sendEncrypted("/upgrade", body)
	if !a.apply(resp) {
		return false
	}

	a.setSuccess("Subscription upgraded successfully.")
	return true
}

func (a *Api) Heartbeat() bool {
	body := a.createBasePayload()
	body["version"] = a.Version

	resp := a.sendEncrypted("/heartbeat", body)
	if !a.apply(resp) {
		a.invalidateRejectedAccess()
		return false
	}

	a.setSuccess("Session alive.")
	return true
}

func (a *Api) Check() bool {
	body := a.createBasePayload()
	body["version"] = a.Version

	resp := a.sendEncrypted("/check", body)
	if !a.apply(resp) {
		a.invalidateRejectedAccess()
		return false
	}

	a.populateUserData(resp)

	if subs, ok := resp["subscriptions"].([]interface{}); ok && len(a.User.Subscriptions) == 0 {
		for _, s := range subs {
			if smap, ok := s.(map[string]interface{}); ok {
				subName := coerceString(smap["subscription"])
				k := coerceString(smap["key"])
				exp := coerceString(smap["expiry"])
				if exp == "" {
					exp = coerceString(smap["expiresAt"])
				}
				tlNum, _ := smap["timeleft"].(float64)
				if tlNum == 0 {
					tlNum, _ = smap["timeLeft"].(float64)
				}
				a.User.Subscriptions = append(a.User.Subscriptions, SubscriptionData{
					Subscription: subName,
					Key:          k,
					Expiry:       exp,
					TimeLeft:     int64(tlNum),
				})
			}
		}
	}

	if b, ok := resp["twoFactorEnabled"].(bool); ok {
		a.User.TwoFactorEnabled = b
	}

	a.setSuccess("Session valid.")
	return true
}

func (a *Api) Logout() bool {
	body := a.createBasePayload()

	result := a.sendEncrypted("/logout", body)
	confirmed := a.apply(result)
	clientFailure, _ := result["client_failure"].(bool)
	serverCode := a.Response.Error
	serverMsg := a.Response.Message

	a.User = UserData{}

	a.sessionKey = ""
	a.sharedKey = nil
	a.isInitialized = false
	if confirmed {
		return true
	}
	if clientFailure || serverCode == "handshake_required" {
		return false
	}
	a.Response = Response{Success: false, Message: serverMsg, Error: "logout_unconfirmed", StatusCode: a.Response.StatusCode, RetryAfter: a.Response.RetryAfter, Payload: map[string]interface{}{"serverError": serverCode}}
	return false
}

func (a *Api) ChangeUsername(newUsername, password string) bool {
	body := a.createBasePayload()
	body["newUsername"] = newUsername
	body["password"] = password

	resp := a.sendEncrypted("/change-username", body)
	if !a.apply(resp) {
		return false
	}

	a.User.Username = newUsername
	a.setSuccess("Username changed successfully.")
	return true
}

func (a *Api) ResetHwid(username, password, totpCode string) bool {
	body := a.createBasePayload()
	body["username"] = username
	body["password"] = password
	if totpCode != "" {
		body["totpCode"] = totpCode
	}

	resp := a.sendEncrypted("/hwid-reset", body)
	if !a.apply(resp) {
		return false
	}

	a.sessionKey = ""
	a.sharedKey = nil
	a.isInitialized = false

	msg, _ := resp["message"].(string)
	if msg == "" {
		msg = "HWID reset successfully."
	}
	a.setSuccess(msg)
	return true
}

func (a *Api) Ban(reason string) bool {
	body := a.createBasePayload()
	if reason != "" {
		body["reason"] = reason
	}

	resp := a.sendEncrypted("/ban", body)
	if !a.apply(resp) {
		return false
	}

	a.setSuccess("Account banned.")
	return true
}

func (a *Api) CheckBlacklist(hwid string) bool {
	body := a.createBasePayload()
	if hwid == "" {
		hwid = a.hardwareID()
	}
	body["hwid"] = hwid

	resp := a.sendEncrypted("/check-blacklist", body)
	if !a.apply(resp) {
		return false
	}
	blacklisted, _ := resp["blacklisted"].(bool)
	return blacklisted
}

func (a *Api) VerifyFileHash(hash string) bool {
	body := a.createBasePayload()
	body["hash"] = hash

	resp := a.sendEncrypted("/verify-hash", body)
	if !a.apply(resp) {
		return false
	}
	verified, _ := resp["verified"].(bool)
	return verified
}

func (a *Api) ForgotPassword(email string) bool {
	body := a.createBasePayload()
	body["email"] = email

	resp := a.sendEncrypted("/forgot-password", body)
	return a.apply(resp)
}

func (a *Api) ResetPassword(token, newPassword string) bool {
	body, _ := json.Marshal(map[string]string{
		"token":       token,
		"newPassword": newPassword,
	})

	status, resBytes, retry, err := a.postRaw("/reset-password", body)
	if err != nil {
		code := classifyTransport(err)
		if errors.Is(err, errResponseTooLarge) {
			code = "response_too_large"
		}
		a.apply(localFailure(status, code, err.Error()))
		return false
	}

	var respJson map[string]interface{}
	if err := json.Unmarshal(resBytes, &respJson); err != nil || respJson == nil {
		a.apply(localFailure(status, "invalid_response", "Response body is not a JSON object."))
		return false
	}
	respJson["http_status"] = status
	if retry > 0 {
		respJson["retry_after"] = retry
	}
	if !a.apply(respJson) {
		return false
	}
	a.setSuccess("Password has been reset.")
	return true
}

func (a *Api) Setup2FA() *TwoFactorSetupData {
	body := a.createBasePayload()

	resp := a.sendEncrypted("/2fa/setup", body)
	if !a.apply(resp) {
		return nil
	}

	secret, _ := resp["secret"].(string)
	qr, _ := resp["qr"].(string)
	uri, _ := resp["uri"].(string)

	return &TwoFactorSetupData{
		Secret:    secret,
		QrDataUrl: qr,
		Uri:       uri,
	}
}

func (a *Api) Confirm2FA(code string) *TwoFactorConfirmData {
	body := a.createBasePayload()
	body["code"] = code

	resp := a.sendEncrypted("/2fa/confirm", body)
	if !a.apply(resp) {
		return nil
	}

	var backupCodes []string
	if arr, ok := resp["backupCodes"].([]interface{}); ok {
		for _, b := range arr {
			if s, ok := b.(string); ok {
				backupCodes = append(backupCodes, s)
			}
		}
	}

	return &TwoFactorConfirmData{Success: true, BackupCodes: backupCodes}
}

func (a *Api) Disable2FA(password, code string) bool {
	body := a.createBasePayload()
	body["password"] = password
	body["code"] = code

	resp := a.sendEncrypted("/2fa/disable", body)
	if !a.apply(resp) {
		return false
	}

	a.setSuccess("2FA disabled.")
	return true
}

func (a *Api) GetGlobalVar(key string) string {
	body := a.createBasePayload()
	body["varName"] = key

	resp := a.sendEncrypted("/global-var", body)
	if !a.apply(resp) {
		return ""
	}
	val, _ := resp["value"].(string)
	return val
}

func (a *Api) GetUserVar(key string) string {
	body := a.createBasePayload()
	body["key"] = key

	resp := a.sendEncrypted("/user-vars/get", body)
	if !a.apply(resp) {
		return ""
	}
	if val, ok := resp["value"].(string); ok {
		return val
	}
	if vars, ok := resp["variables"].([]interface{}); ok {
		for _, v := range vars {
			if vmap, ok := v.(map[string]interface{}); ok {
				if vmap["key"] == key {
					val, _ := vmap["value"].(string)
					return val
				}
			}
		}
	}
	return ""
}

func (a *Api) SetUserVar(key, value string) bool {
	body := a.createBasePayload()
	body["key"] = key
	body["value"] = value

	resp := a.sendEncrypted("/user-vars/set", body)
	if !a.apply(resp) {
		return false
	}

	a.setSuccess("User variable saved.")
	return true
}

func (a *Api) AccessTicket(scope, resourceId string) string {
	body := a.createBasePayload()
	body["scope"] = scope
	if resourceId != "" {
		body["resourceId"] = resourceId
	}

	resp := a.sendEncrypted("/access-ticket", body)
	if !a.apply(resp) {
		return ""
	}
	ticket, _ := resp["ticket"].(string)
	if ticket == "" {
		a.setError("", "invalid_response")
		return ""
	}
	a.setSuccess("Access ticket issued.")
	return ticket
}

func (a *Api) DownloadFile(fileId string) []byte {
	ticket := a.AccessTicket("file", fileId)
	if ticket == "" {
		return nil
	}

	status, raw, location, err := a.getRaw("/files/download?ticket="+url.QueryEscape(ticket), maxDownloadResponseBytes)
	if err != nil {
		code := classifyTransport(err)
		if errors.Is(err, errResponseTooLarge) {
			code = "response_too_large"
		}
		a.apply(localFailure(status, code, err.Error()))
		return nil
	}
	if status >= 300 && status < 400 {

		refusal := map[string]interface{}{"success": false, "error": "external_download", "http_status": status}
		if location != "" {
			refusal["url"] = location
		}
		a.apply(refusal)
		return nil
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(raw, &resp); err != nil || resp == nil {
		a.apply(localFailure(status, "invalid_response", "Response body is not a JSON object."))
		return nil
	}
	resp["http_status"] = status
	if !a.apply(resp) {
		return nil
	}

	dataB64, _ := resp["data"].(string)
	if dataB64 == "" {
		a.apply(localFailure(status, "download_decode_failed", "File response is missing content"))
		return nil
	}
	gz, err := base64.StdEncoding.DecodeString(dataB64)
	if err != nil {
		a.apply(localFailure(status, "download_decode_failed", "File content could not be decoded"))
		return nil
	}
	zr, err := gzip.NewReader(bytes.NewReader(gz))
	if err != nil {
		a.apply(localFailure(status, "download_decode_failed", "File content could not be decompressed"))
		return nil
	}
	out, err := io.ReadAll(io.LimitReader(zr, maxFileBytes+1))
	if err != nil {
		a.apply(localFailure(status, "download_decode_failed", "File content could not be decompressed"))
		return nil
	}
	if int64(len(out)) > maxFileBytes {
		a.apply(localFailure(status, "download_too_large", "File content exceeds the size limit"))
		return nil
	}
	a.setSuccess("File downloaded.")
	a.Response.Payload["data"] = base64.StdEncoding.EncodeToString(out)
	return out
}

type WebhookResult struct {
	Success    bool `json:"success"`
	Status     int  `json:"status"`
	StatusCode int  `json:"statusCode"`
	Data       interface{}
	Truncated  bool
}

func (a *Api) TriggerWebhook(webId string, customPayload map[string]interface{}) *WebhookResult {
	body := a.createBasePayload()
	body["webId"] = webId
	if customPayload != nil {
		body["payload"] = customPayload
	}

	resp := a.sendEncrypted("/webhook", body)

	if !a.apply(resp) {
		if a.Response.Error == "" {
			a.Response.Error = "webhook_failed"
		}
		if status, ok := resp["status"].(float64); ok {
			if msg, _ := resp["message"].(string); msg == "" {
				a.Response.Message = fmt.Sprintf("Remote endpoint returned HTTP %d", int(status))
			}
		}
		return nil
	}

	status, _ := resp["status"].(float64)
	truncated, _ := resp["truncated"].(bool)
	a.setSuccess("Webhook triggered.")
	return &WebhookResult{Success: true, Status: int(status), StatusCode: int(status), Data: resp["data"], Truncated: truncated}
}

type EventFrame struct {
	ID   string
	Data string
}

type StreamError struct {
	Status int
	Code   string
	Msg    string
}

func (e *StreamError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "event stream refused: " + e.Code
}

func (a *Api) SubscribeEvents(ctx context.Context, lastEventID string, onEvent func(EventFrame)) error {
	ticket := a.AccessTicket("events", "")
	if ticket == "" {
		return &StreamError{Status: a.Response.StatusCode, Code: a.Response.Error, Msg: "events ticket denied: " + a.Response.Error}
	}
	if err := validateApiUrl(a.ApiUrl, a.allowInsecureLoopback); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.ApiUrl+"/events?ticket="+url.QueryEscape(ticket), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("User-Agent", "NightAuth-Golang-SDK/1.0")
	if lastEventID != "" {
		req.Header.Set("Last-Event-ID", lastEventID)
	}

	stream := &http.Client{
		Transport:     a.httpClient.Transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}
	resp, err := stream.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return &StreamError{Status: 0, Code: classifyTransport(err), Msg: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		var body map[string]interface{}
		_ = json.Unmarshal(raw, &body)
		errStr, _ := body["error"].(string)
		if errStr == "" {
			errStr = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
		a.setError("", errStr)
		a.Response.StatusCode = resp.StatusCode
		return &StreamError{Status: resp.StatusCode, Code: errStr, Msg: "events stream rejected: " + errStr}
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var frame EventFrame
	hasData := false
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case line == "":
			if hasData {
				onEvent(frame)
			}
			frame = EventFrame{}
			hasData = false
		case strings.HasPrefix(line, ":"):

		case strings.HasPrefix(line, "id:"):
			frame.ID = strings.TrimPrefix(strings.TrimPrefix(line, "id:"), " ")
		case strings.HasPrefix(line, "data:"):
			value := strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " ")
			if hasData {
				frame.Data += "\n"
			}
			frame.Data += value
			hasData = true
		}
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return scanner.Err()
}

func (a *Api) StartRealtimeEventListener() {
	a.mu.Lock()
	if a.stopChan != nil {
		a.mu.Unlock()
		return
	}
	a.stopChan = make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelEvents = cancel
	stopCh := a.stopChan
	a.mu.Unlock()

	go func() {
		defer cancel()
		go func() {
			select {
			case <-stopCh:
				cancel()
			case <-ctx.Done():
			}
		}()

		_ = a.SubscribeEvents(ctx, "", func(ev EventFrame) {
			a.dispatchEvent(ev.Data)
		})
	}()
}

func (a *Api) StopRealtimeEventListener() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cancelEvents != nil {
		a.cancelEvents()
		a.cancelEvents = nil
	}
	if a.stopChan != nil {
		close(a.stopChan)
		a.stopChan = nil
	}
}

func (a *Api) dispatchEvent(jsonString string) {
	var node map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &node); err != nil {
		return
	}

	evType, _ := node["type"].(string)
	switch evType {
	case "session-invalid", "session-terminated":
		a.StopRealtimeEventListener()
		a.User = UserData{}
		if a.OnSessionInvalidated != nil {
			a.OnSessionInvalidated("Your session has been terminated by administrator or system.")
		}

	case "app-banned":
		a.StopRealtimeEventListener()
		a.User = UserData{}
		reason, _ := node["reason"].(string)
		if reason == "" {
			reason = "Application has been banned."
		}
		if a.OnAppBanned != nil {
			a.OnAppBanned(reason)
		}

	case "app-paused", "app-disabled":
		a.StopRealtimeEventListener()
		a.User = UserData{}
		if a.OnAppStatusChanged != nil {
			a.OnAppStatusChanged(fmt.Sprintf("Application status changed to: %s", evType))
		}

	case "version":
		newVer, _ := node["version"].(string)
		if newVer != "" {
			a.App.Version = newVer
			a.RequiredVersion = newVer
			if dl, ok := node["downloadUrl"].(string); ok && dl != "" {
				a.PendingDownloadUrl = dl
			}
			if a.OnVersionUpdated != nil {
				a.OnVersionUpdated(newVer)
			}
			if a.Version != "" && newVer != a.Version {
				a.StopRealtimeEventListener()
				a.User = UserData{}
				if a.OnSessionInvalidated != nil {
					a.OnSessionInvalidated(fmt.Sprintf("Your build is outdated (server now requires v%s). Session terminated.", newVer))
				}
			}
		}

	case "chat-message":
		room, _ := node["room"].(string)
		if room == "" {
			room = "general"
		}
		msgId, _ := node["id"].(string)
		author, _ := node["author"].(string)
		if author == "" {
			author = "unknown"
		}
		content, _ := node["content"].(string)
		if content == "" {
			content, _ = node["message"].(string)
		}
		ts, _ := node["createdAt"].(string)
		if ts == "" {
			ts = time.Now().UTC().Format(time.RFC3339)
		}
		msg := ChatMessage{
			Id:        msgId,
			Author:    author,
			Message:   content,
			Timestamp: ts,
		}
		if a.OnChatMessageReceived != nil {
			a.OnChatMessageReceived(msg, room)
		}
	}
}

func (a *Api) getRaw(pathAndQuery string, maxBytes int64) (int, []byte, string, error) {
	if err := validateApiUrl(a.ApiUrl, a.allowInsecureLoopback); err != nil {
		return 0, nil, "", err
	}
	req, err := http.NewRequest(http.MethodGet, a.ApiUrl+pathAndQuery, nil)
	if err != nil {
		return 0, nil, "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NightAuth-Golang-SDK/1.0")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return 0, nil, "", err
	}
	defer resp.Body.Close()
	location := resp.Header.Get("Location")

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return resp.StatusCode, nil, location, err
	}
	if int64(len(body)) > maxBytes {
		return resp.StatusCode, nil, location, fmt.Errorf("response exceeds size limit")
	}
	return resp.StatusCode, body, location, nil
}

func (a *Api) GetChatMessages(channel string) []ChatMessage {
	body := a.createBasePayload()
	if channel == "" {
		channel = "general"
	}
	body["room"] = channel

	resp := a.sendEncrypted("/chat/get", body)
	if !a.apply(resp) {
		return nil
	}

	var list []ChatMessage
	if arr, ok := resp["messages"].([]interface{}); ok {
		for _, item := range arr {
			if imap, ok := item.(map[string]interface{}); ok {
				id := ""
				switch v := imap["id"].(type) {
				case string:
					id = v
				case float64:
					id = strconv.FormatFloat(v, 'f', -1, 64)
				}
				author, _ := imap["author"].(string)
				content, _ := imap["content"].(string)
				if content == "" {
					content, _ = imap["message"].(string)
				}
				ts, _ := imap["timestamp"].(string)
				if ts == "" {
					ts, _ = imap["createdAt"].(string)
				}
				list = append(list, ChatMessage{
					Id:        id,
					Author:    author,
					Message:   content,
					Timestamp: ts,
				})
			}
		}
	}
	return list
}

func (a *Api) SendChatMessage(message, channel string) bool {
	body := a.createBasePayload()
	if channel == "" {
		channel = "general"
	}
	body["room"] = channel
	body["content"] = message

	resp := a.sendEncrypted("/chat/send", body)
	if !a.apply(resp) {
		return false
	}
	a.setSuccess("Message sent.")
	return true
}

func (a *Api) Log(message string) bool {
	body := a.createBasePayload()
	body["message"] = message

	resp := a.sendEncrypted("/log", body)
	return a.apply(resp)
}

func (a *Api) GetStats() *AppStats {
	body := a.createBasePayload()

	resp := a.sendEncrypted("/stats", body)
	if !a.apply(resp) {
		return nil
	}

	users, _ := resp["users"].(map[string]interface{})
	licenses, _ := resp["licenses"].(map[string]interface{})

	uTotal, _ := users["total"].(float64)
	uActive, _ := users["active"].(float64)
	uBanned, _ := users["banned"].(float64)

	lTotal, _ := licenses["total"].(float64)
	lUsed, _ := licenses["used"].(float64)
	lNotUsed, _ := licenses["notUsed"].(float64)

	ver, _ := resp["version"].(string)
	if ver == "" {
		ver = a.Version
	}

	return &AppStats{
		TotalUsers:     int64(uTotal),
		ActiveUsers:    int64(uActive),
		BannedUsers:    int64(uBanned),
		TotalLicenses:  int64(lTotal),
		UsedLicenses:   int64(lUsed),
		UnusedLicenses: int64(lNotUsed),
		Version:        ver,
	}
}

func (a *Api) GetOnlineUsers() []string {
	body := a.createBasePayload()

	resp := a.sendEncrypted("/online-users", body)
	if !a.apply(resp) {
		return nil
	}

	var users []string
	if arr, ok := resp["users"].([]interface{}); ok {
		for _, u := range arr {
			if s, ok := u.(string); ok {
				users = append(users, s)
			}
		}
	}
	return users
}

func (a *Api) checkInit() {
	if !a.isInitialized || a.sharedKey == nil || a.sessionKey == "" {
		panic("NightAuth must be initialized before calling API methods. Call app.Init() first.")
	}
}

func (a *Api) createBasePayload() map[string]interface{} {
	obj := map[string]interface{}{
		"apikey": a.ApiKey,
		"name":   a.Name,
	}
	if a.Token != "" {
		obj["token"] = a.Token
	}
	return obj
}

func (a *Api) sendEncrypted(path string, payload map[string]interface{}) map[string]interface{} {
	return a.sendEncryptedWith(path, payload, nil)
}

func (a *Api) sendEncryptedWith(path string, payload map[string]interface{}, mk envelopeTransport) map[string]interface{} {
	if !a.isInitialized || a.sharedKey == nil || a.sessionKey == "" {
		a.Init()
		if !a.isInitialized {
			return localFailure(a.Response.StatusCode, a.Response.Error, a.Response.Message)
		}
	}
	result := a.sendEncryptedOnce(path, payload, mk)
	if success, _ := result["success"].(bool); !success && isHandshakeExpired(result) {
		a.Init()
		if a.isInitialized {
			result = a.sendEncryptedOnce(path, payload, mk)
		}
	}
	return result
}

func isHandshakeExpired(result map[string]interface{}) bool {
	if _, local := result["client_failure"]; local {
		return false
	}
	msg, _ := result["message"].(string)
	msg = strings.ToLower(msg)
	errCode, _ := result["error"].(string)
	return strings.Contains(msg, "handshake") || strings.Contains(msg, "initialize app") ||
		errCode == "handshake_required" || errCode == "session_invalid" || errCode == "unauthenticated"
}

func (a *Api) sendEncryptedOnce(path string, payload map[string]interface{}, mk envelopeTransport) map[string]interface{} {
	ridBytes := make([]byte, 16)
	_, _ = rand.Read(ridBytes)
	rid := hex.EncodeToString(ridBytes)

	ts := time.Now().UnixMilli()

	fullURL := a.ApiUrl + path
	targetPath := path
	if u, err := url.Parse(fullURL); err == nil && u.Path != "" {
		targetPath = u.Path
	}
	aad := fmt.Sprintf("POST:%s:%s:%s:%d", targetPath, a.sessionKey, rid, ts)
	jsonPlain, _ := json.Marshal(payload)

	ctB64, nonceB64, err := aesGcmEncrypt(a.sharedKey, string(jsonPlain), aad)
	if err != nil {
		return localFailure(0, "encryption_failed", err.Error())
	}

	envelope := map[string]interface{}{
		"id":  a.sessionKey,
		"ct":  ctB64,
		"n":   nonceB64,
		"rid": rid,
		"ts":  ts,
	}

	envBytes, _ := json.Marshal(envelope)
	ctype, reqBody := "application/json", envBytes
	if mk != nil {
		var err error
		if ctype, reqBody, err = mk(envBytes); err != nil {
			return localFailure(0, "invalid_request", err.Error())
		}
	}
	status, respRaw, retry, err := a.postRawTyped(path, reqBody, ctype)
	if err != nil {
		if errors.Is(err, errResponseTooLarge) {
			return localFailure(status, "response_too_large", "Response exceeds size limit.")
		}
		return localFailure(0, classifyTransport(err), err.Error())
	}

	var encResp map[string]interface{}
	if err := json.Unmarshal(respRaw, &encResp); err != nil {
		return localFailure(status, "invalid_response", "Response body is not a JSON object.")
	}

	ct, hasCt := encResp["ct"].(string)
	n, hasN := encResp["n"].(string)
	resRid, hasRid := encResp["rid"].(string)
	resTsNum, hasTs := encResp["ts"].(float64)
	if hasCt && hasN && hasRid && hasTs {
		resTs := int64(resTsNum)
		if resRid != rid || resTs != ts {
			return localFailure(status, "response_binding_mismatch", "Response is not bound to this request.")
		}

		resAad := fmt.Sprintf("v1:res:%s:%s:%d", a.sessionKey, resRid, resTs)
		decrypted, err := aesGcmDecrypt(a.sharedKey, ct, n, resAad)
		if err != nil {
			return localFailure(status, "unauthenticated_response", "Response is not authenticated to this request.")
		}
		var decJson map[string]interface{}
		if err := json.Unmarshal([]byte(decrypted), &decJson); err != nil || decJson == nil {
			return localFailure(status, "invalid_response", "Decrypted response is not a JSON object.")
		}
		decJson["http_status"] = status
		if retry > 0 {
			decJson["retry_after"] = retry
		}
		return decJson
	}

	if status >= 200 && status < 300 {
		return localFailure(status, "unauthenticated_response", "Response is not authenticated to this request.")
	}

	refusal := make(map[string]interface{}, len(encResp)+3)
	for k, v := range encResp {
		refusal[k] = v
	}
	refusal["success"] = false
	refusal["http_status"] = status
	if retry > 0 {
		refusal["retry_after"] = retry
	}
	if status == http.StatusPreconditionRequired || refusal["error"] == "handshake_required" {
		a.sessionKey = ""
		a.sharedKey = nil
		a.isInitialized = false
	}
	return refusal
}

func (a *Api) postRaw(path string, jsonBody []byte) (int, []byte, float64, error) {
	return a.postRawTyped(path, jsonBody, "application/json")
}

func (a *Api) postRawTyped(path string, jsonBody []byte, contentType string) (int, []byte, float64, error) {
	if err := validateApiUrl(a.ApiUrl, a.allowInsecureLoopback); err != nil {
		return 0, nil, 0, err
	}
	reqUrl := a.ApiUrl + path
	req, err := http.NewRequest("POST", reqUrl, bytes.NewReader(jsonBody))
	if err != nil {
		return 0, nil, 0, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "NightAuth-Golang-SDK/1.0")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return 0, nil, 0, err
	}
	defer resp.Body.Close()

	retry := parseRetryAfter(resp.Header.Get("Retry-After"))
	const maxResponseBytes = 5 * 1024 * 1024
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err == nil && len(body) > maxResponseBytes {
		return resp.StatusCode, nil, retry, errResponseTooLarge
	}
	return resp.StatusCode, body, retry, err
}

func aesGcmEncrypt(key []byte, plaintext, aad string) (string, string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", "", err
	}

	sealed := gcm.Seal(nil, nonce, []byte(plaintext), []byte(aad))
	return base64.StdEncoding.EncodeToString(sealed), base64.StdEncoding.EncodeToString(nonce), nil
}

func aesGcmDecrypt(key []byte, ctB64, nB64, aad string) (string, error) {
	nonce, err := base64.StdEncoding.DecodeString(nB64)
	if err != nil {
		return "", err
	}
	ct, err := base64.StdEncoding.DecodeString(ctB64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plain, err := gcm.Open(nil, nonce, ct, []byte(aad))
	if err != nil {
		return "", err
	}

	return string(plain), nil
}

func verifyServerSignature(k2Hex, vkHex string, ts int64, sigB64 string) bool {
	vkBytes, err := hex.DecodeString(vkHex)
	if err != nil || len(vkBytes) != 32 {
		return false
	}

	sigBytes, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil || len(sigBytes) != 64 {
		return false
	}

	data := fmt.Sprintf("%s.%s.%d", k2Hex, vkHex, ts)
	return ed25519.Verify(ed25519.PublicKey(vkBytes), []byte(data), sigBytes)
}

func cryptoBoxSeal(plaintext, recipientPubKeyHex string) (string, error) {
	recPubBytes, err := hex.DecodeString(recipientPubKeyHex)
	if err != nil || len(recPubBytes) != 32 {
		return "", errors.New("invalid recipient public key")
	}

	var recPub [32]byte
	copy(recPub[:], recPubBytes)

	var ephemPriv [32]byte
	_, _ = rand.Read(ephemPriv[:])
	var ephemPub [32]byte
	curve25519.ScalarBaseMult(&ephemPub, &ephemPriv)

	var rawSecret [32]byte
	curve25519.ScalarMult(&rawSecret, &ephemPriv, &recPub)

	hk := hkdf.New(sha256.New, rawSecret[:], nil, sealInfo)
	encKey := make([]byte, 32)
	if _, err := io.ReadFull(hk, encKey); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	_, _ = rand.Read(nonce)

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), ephemPub[:])
	if len(ciphertext) < 16 {
		return "", errors.New("ciphertext too short")
	}

	tag := ciphertext[len(ciphertext)-16:]
	cipherBytes := ciphertext[:len(ciphertext)-16]

	var result []byte
	result = append(result, ephemPub[:]...)
	result = append(result, nonce...)
	result = append(result, tag...)
	result = append(result, cipherBytes...)

	return base64.StdEncoding.EncodeToString(result), nil
}

func coerceString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatInt(int64(val), 10)
	case float32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case int:
		return strconv.Itoa(val)
	case json.Number:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

func (a *Api) populateUserData(data map[string]interface{}) {
	u, ok := data["user"].(map[string]interface{})
	if !ok {
		u, _ = data["info"].(map[string]interface{})
	}
	if u == nil {
		return
	}

	a.User.Username = coerceString(u["username"])
	a.User.Email = coerceString(u["email"])
	a.User.Hwid = coerceString(u["hwid"])
	a.User.Ip = coerceString(u["ip"])
	a.User.Status = coerceString(u["status"])
	if a.User.Status == "" {
		a.User.Status = "active"
	}
	a.User.CreatedAt = coerceString(u["createdAt"])
	if a.User.CreatedAt == "" {
		a.User.CreatedAt = coerceString(u["createdate"])
	}
	a.User.LastLoginAt = coerceString(u["lastLoginAt"])
	if a.User.LastLoginAt == "" {
		a.User.LastLoginAt = coerceString(u["lastlogin"])
	}
	if b, ok := u["twoFactorEnabled"].(bool); ok {
		a.User.TwoFactorEnabled = b
	}

	a.User.Subscriptions = nil
	if subs, ok := u["subscriptions"].([]interface{}); ok {
		for _, s := range subs {
			if smap, ok := s.(map[string]interface{}); ok {
				subName := coerceString(smap["subscription"])
				k := coerceString(smap["key"])
				exp := coerceString(smap["expiry"])
				if exp == "" {
					exp = coerceString(smap["expiresAt"])
				}
				tlNum, _ := smap["timeleft"].(float64)
				if tlNum == 0 {
					tlNum, _ = smap["timeLeft"].(float64)
				}
				a.User.Subscriptions = append(a.User.Subscriptions, SubscriptionData{
					Subscription: subName,
					Key:          k,
					Expiry:       exp,
					TimeLeft:     int64(tlNum),
				})
			}
		}
	}

	if appObj, ok := data["app"].(map[string]interface{}); ok {
		tu, _ := appObj["totalUsers"].(float64)
		ou, _ := appObj["onlineUsers"].(float64)
		tk, _ := appObj["totalKeys"].(float64)
		ver := coerceString(appObj["version"])
		if ver == "" {
			ver = a.Version
		}
		a.App.TotalUsers = int64(tu)
		a.App.OnlineUsers = int64(ou)
		a.App.TotalKeys = int64(tk)
		a.App.Version = ver
	}
}

func formatHumanErrorMessage(msg, err string) string {
	candidate := strings.TrimSpace(msg)
	if candidate == "" {
		candidate = err
	}
	if candidate == "" {
		candidate = "Operation failed"
	}

	mapping := map[string]string{
		"rate_limited":           "Too many requests. Please wait a moment and try again.",
		"invalid_credentials":    "The credentials provided do not match our records.",
		"invalid_apikey":         "The API key provided is invalid or unrecognized.",
		"app_not_found":          "Application was not found. Please check your credentials.",
		"app_disabled":           "This application is currently disabled by administrator policy.",
		"app_paused":             "Application is currently paused for maintenance.",
		"invalid_token":          "Authentication token is invalid or expired.",
		"feature_disabled":       "This feature is currently disabled for this application.",
		"session_invalid":        "Session is no longer valid. Please sign in again.",
		"session_not_validated":  "Session is no longer valid. Please sign in again.",
		"invalid_hwid":           "Device signature is invalid or unrecognized.",
		"hwid_mismatch":          "Device signature mismatch. Please reset your HWID.",
		"vpn_blocked":            "Anonymous proxies and VPN connections are not permitted.",
		"blacklisted":            "Your hardware identifier or IP address has been restricted.",
		"username_taken":         "This username is already registered. Please choose another.",
		"email_taken":            "This email address is already registered.",
		"invalid_license":        "The license key provided is invalid or unrecognized.",
		"license_used":           "This license key has already been activated.",
		"no_active_subscription": "No active subscription found for this account.",
		"expired":                "No active subscription found for this account.",
		"subscription_paused":    "Your subscription is temporarily suspended.",
		"paused":                 "Your subscription is temporarily suspended.",
	}

	if val, found := mapping[candidate]; found {
		return val
	}
	return candidate
}

func (a *Api) setError(msg, err string) {
	a.Response = Response{
		Success: false,
		Message: formatHumanErrorMessage(msg, err),
		Error:   err,
	}
}

func (a *Api) setSuccess(msg string) {
	a.Response.Success = true
	a.Response.Message = msg
	a.Response.Error = ""
}

func (a *Api) Close() {
	a.mu.Lock()
	if a.disposed {
		a.mu.Unlock()
		return
	}
	a.disposed = true
	a.mu.Unlock()

	a.StopRealtimeEventListener()

	a.mu.Lock()
	a.sessionKey = ""
	a.sharedKey = nil
	a.isInitialized = false
	a.User = UserData{}
	a.mu.Unlock()
}

var errResponseTooLarge = errors.New("response exceeds size limit")

func localFailure(status int, code, message string) map[string]interface{} {
	return map[string]interface{}{
		"success":        false,
		"error":          code,
		"message":        message,
		"http_status":    status,
		"client_failure": true,
	}
}

func classifyTransport(err error) string {
	var netErr net.Error
	if errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &netErr) && netErr.Timeout()) {
		return "timeout"
	}
	return "network_error"
}

func parseRetryAfter(v string) float64 {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	if secs, err := strconv.ParseFloat(v, 64); err == nil {
		if secs < 0 {
			return 0
		}
		return secs
	}
	if t, err := http.ParseTime(v); err == nil {
		if wait := time.Until(t).Seconds(); wait > 0 {
			return wait
		}
	}
	return 0
}

func (a *Api) apply(resp map[string]interface{}) bool {
	status, _ := resp["http_status"].(int)
	retry, _ := resp["retry_after"].(float64)
	clientFailure, _ := resp["client_failure"].(bool)
	success, _ := resp["success"].(bool)
	msg, _ := resp["message"].(string)
	code, _ := resp["error"].(string)

	payload := make(map[string]interface{}, len(resp))
	for k, v := range resp {
		if k == "http_status" || k == "client_failure" {
			continue
		}
		payload[k] = v
	}

	if dl, ok := resp["downloadUrl"].(string); ok && dl != "" {
		a.PendingDownloadUrl = dl
	}
	if reqVer, ok := resp["requiredVersion"].(string); ok && reqVer != "" {
		a.RequiredVersion = reqVer
	}

	if success && !clientFailure && status >= 200 && status < 300 {
		a.Response = Response{Success: true, Message: msg, StatusCode: status, RetryAfter: retry, Payload: payload}
		return true
	}
	a.Response = Response{
		Success:    false,
		Message:    formatHumanErrorMessage(msg, code),
		Error:      code,
		StatusCode: status,
		RetryAfter: retry,
		Payload:    payload,
	}
	return false
}

func (a *Api) hardwareID() string {
	if a.HardwareIdProvider != nil {
		return a.HardwareIdProvider()
	}
	return SecurityProtections{}.GetHardwareId()
}

func (a *Api) invalidateRejectedAccess() {
	switch a.Response.Error {
	case "session_invalid", "session_not_validated", "unauthenticated", "version_mismatch", "app_disabled", "app_paused", "blacklisted":
		a.sessionKey = ""
		a.sharedKey = nil
		a.isInitialized = false
	}
}

type envelopeTransport func(envelope []byte) (string, []byte, error)

func (a *Api) GetUserVariables() bool {
	body := a.createBasePayload()
	resp := a.sendEncrypted("/user-vars/get", body)
	return a.apply(resp)
}

func (a *Api) GetChatMessagesPage(channel string, before int64, markRead bool) []ChatMessage {
	body := a.createBasePayload()
	if channel == "" {
		channel = "general"
	}
	body["room"] = channel
	if before > 0 {
		body["before"] = before
	}
	if markRead {
		body["markRead"] = true
	}
	resp := a.sendEncrypted("/chat/get", body)
	if !a.apply(resp) {
		return nil
	}
	return chatMessagesFrom(resp)
}

func chatMessagesFrom(resp map[string]interface{}) []ChatMessage {
	var list []ChatMessage
	arr, _ := resp["messages"].([]interface{})
	for _, item := range arr {
		imap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		id := ""
		switch v := imap["id"].(type) {
		case string:
			id = v
		case float64:
			id = strconv.FormatFloat(v, 'f', -1, 64)
		}
		author, _ := imap["author"].(string)
		content, _ := imap["content"].(string)
		if content == "" {
			content, _ = imap["message"].(string)
		}
		ts, _ := imap["timestamp"].(string)
		if ts == "" {
			ts, _ = imap["createdAt"].(string)
		}
		list = append(list, ChatMessage{Id: id, Author: author, Message: content, Timestamp: ts})
	}
	return list
}

func (a *Api) TriggerWebhookFile(webId string, file []byte, fileName string, customPayload map[string]interface{}) *WebhookResult {
	if strings.TrimSpace(fileName) == "" {
		a.setError("file name is required", "invalid_request")
		return nil
	}
	body := a.createBasePayload()
	body["webId"] = webId
	if customPayload != nil {
		body["payload"] = customPayload
	}
	resp := a.sendEncryptedWith("/webhook", body, func(envelope []byte) (string, []byte, error) {
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		ep, err := w.CreatePart(textproto.MIMEHeader{"Content-Disposition": {`form-data; name="envelope"`}})
		if err != nil {
			return "", nil, err
		}
		if _, err := ep.Write(envelope); err != nil {
			return "", nil, err
		}
		safeName := strings.NewReplacer(`"`, "%22", "\r", "%0D", "\n", "%0A").Replace(fileName)
		fp, err := w.CreatePart(textproto.MIMEHeader{
			"Content-Disposition": {`form-data; name="file"; filename="` + safeName + `"`},
			"Content-Type":        {"application/octet-stream"},
		})
		if err != nil {
			return "", nil, err
		}
		if _, err := fp.Write(file); err != nil {
			return "", nil, err
		}
		if err := w.Close(); err != nil {
			return "", nil, err
		}
		return w.FormDataContentType(), buf.Bytes(), nil
	})
	if !a.apply(resp) {
		return nil
	}
	status, _ := resp["status"].(float64)
	truncated, _ := resp["truncated"].(bool)
	a.setSuccess("Webhook triggered.")
	return &WebhookResult{Success: true, Status: int(status), StatusCode: int(status), Data: resp["data"], Truncated: truncated}
}

func (a *Api) ApplyUpdate(downloadUrl string) bool {
	a.setError("Unsigned auto-update is disabled. Verify a signed release manifest with SignedUpdateInstaller.", "unsigned_update_disabled")
	return false
}

func (a *Api) DownloadSignedUpdate(manifestURL, pinnedReleaseKeyHex, destination string, timeout time.Duration) (string, error) {
	return DownloadSignedUpdateWithInstaller(manifestURL, pinnedReleaseKeyHex, a.Name, destination, timeout)
}

const (
	MaxSignedUpdateArtifactBytes int64 = 256 * 1024 * 1024
	MaxSignedUpdateManifestBytes int64 = 64 * 1024
)

type VerifiedUpdate struct {
	Application string
	Version     string
	Sequence    int64
	URL         string
	Size        int64
	ExpiresAt   time.Time
	Digest      []byte
}

type signedUpdateEnvelope struct {
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

type signedUpdatePayload struct {
	Schema      int    `json:"schema"`
	Application string `json:"application"`
	Version     string `json:"version"`
	Sequence    int64  `json:"sequence"`
	URL         string `json:"url"`
	Size        int64  `json:"size"`
	ExpiresAt   int64  `json:"expiresAt"`
	SHA256      string `json:"sha256"`
}

func VerifySignedUpdate(envelopeJSON, pinnedPublicKeyHex, expectedApplication string, minimumSequence int64) (*VerifiedUpdate, error) {
	if int64(len(envelopeJSON)) > MaxSignedUpdateManifestBytes || minimumSequence < 0 || strings.TrimSpace(expectedApplication) == "" {
		return nil, errors.New("invalid manifest parameters")
	}

	var env signedUpdateEnvelope
	if err := json.Unmarshal([]byte(envelopeJSON), &env); err != nil {
		return nil, fmt.Errorf("manifest envelope is not valid JSON: %w", err)
	}
	if env.Payload == "" || env.Signature == "" {
		return nil, errors.New("manifest envelope missing payload or signature")
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(env.Payload)
	if err != nil {
		return nil, fmt.Errorf("invalid payload base64: %w", err)
	}
	signatureBytes, err := base64.StdEncoding.DecodeString(env.Signature)
	if err != nil {
		return nil, fmt.Errorf("invalid signature base64: %w", err)
	}
	keyBytes, err := hex.DecodeString(pinnedPublicKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid pinned public key hex: %w", err)
	}

	if len(payloadBytes) > 32768 || len(signatureBytes) != ed25519.SignatureSize || len(keyBytes) != ed25519.PublicKeySize {
		return nil, errors.New("invalid manifest size, signature or key length")
	}

	if !ed25519.Verify(ed25519.PublicKey(keyBytes), payloadBytes, signatureBytes) {
		return nil, errors.New("release signature rejected")
	}

	var data signedUpdatePayload
	if err := json.Unmarshal(payloadBytes, &data); err != nil {
		return nil, fmt.Errorf("manifest payload is not valid JSON: %w", err)
	}

	if data.Schema != 1 || data.Application != expectedApplication || strings.TrimSpace(data.Version) == "" || len(data.Version) > 100 {
		return nil, errors.New("release identity or schema rejected")
	}
	if data.Sequence <= minimumSequence || data.Sequence > 9007199254740991 {
		return nil, errors.New("release sequence rejected")
	}
	if data.Size < 1 || data.Size > MaxSignedUpdateArtifactBytes {
		return nil, errors.New("release artifact size rejected")
	}
	now := time.Now().Unix()
	if data.ExpiresAt <= now {
		return nil, errors.New("release expired")
	}

	parsedURL, err := url.Parse(data.URL)
	if err != nil || parsedURL.Scheme != "https" || parsedURL.User != nil || parsedURL.Fragment != "" {
		return nil, errors.New("release URL must be HTTPS without userinfo or fragment")
	}

	if len(data.SHA256) != 64 {
		return nil, errors.New("invalid sha256 hex length")
	}
	digest, err := hex.DecodeString(data.SHA256)
	if err != nil || len(digest) != 32 {
		return nil, errors.New("invalid sha256 digest")
	}

	return &VerifiedUpdate{
		Application: data.Application,
		Version:     data.Version,
		Sequence:    data.Sequence,
		URL:         data.URL,
		Size:        data.Size,
		ExpiresAt:   time.Unix(data.ExpiresAt, 0).UTC(),
		Digest:      digest,
	}, nil
}

func DownloadSignedUpdate(update *VerifiedUpdate, destination string, timeout time.Duration) (string, error) {
	if update == nil {
		return "", errors.New("update is nil")
	}
	if time.Now().UTC().After(update.ExpiresAt) {
		return "", errors.New("release expired")
	}

	target, err := filepath.Abs(destination)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(target); err == nil {
		return "", errors.New("destination file already exists")
	}

	dir := filepath.Dir(target)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	randBytes := make([]byte, 8)
	_, _ = rand.Read(randBytes)
	temporary := fmt.Sprintf("%s.%s.partial", target, hex.EncodeToString(randBytes))

	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("redirects are not permitted")
		},
	}

	req, err := http.NewRequest(http.MethodGet, update.URL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "NightAuth-Updater/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("artifact must return HTTP 200 without redirects (got %d)", resp.StatusCode)
	}
	if resp.ContentLength > 0 && resp.ContentLength != update.Size {
		return "", errors.New("artifact length differs from signed manifest")
	}

	out, err := os.OpenFile(temporary, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer func() {
		out.Close()
		_ = os.Remove(temporary)
	}()

	hasher := sha256.New()
	buf := make([]byte, 64*1024)
	var total int64
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			total += int64(n)
			if total > update.Size {
				return "", errors.New("artifact exceeds signed size")
			}
			hasher.Write(buf[:n])
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				return "", writeErr
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return "", readErr
		}
	}

	if total != update.Size {
		return "", errors.New("artifact size differs from signed manifest")
	}
	if !bytes.Equal(hasher.Sum(nil), update.Digest) {
		return "", errors.New("artifact hash rejected")
	}
	if time.Now().UTC().After(update.ExpiresAt) {
		return "", errors.New("release expired during download")
	}

	if err := out.Sync(); err != nil {
		return "", err
	}
	out.Close()

	if err := os.Rename(temporary, target); err != nil {
		return "", err
	}
	return target, nil
}

func getReleaseSequencePath(application string) (string, error) {
	if strings.TrimSpace(application) == "" || len(application) > 100 || strings.ContainsAny(application, "<>:\"/\\|?*") {
		return "", errors.New("invalid application identifier")
	}
	var root string
	if runtime.GOOS == "windows" {
		root = os.Getenv("LOCALAPPDATA")
		if root == "" {
			home, _ := os.UserHomeDir()
			root = filepath.Join(home, "AppData", "Local")
		}
	} else {
		root = os.Getenv("XDG_DATA_HOME")
		if root == "" {
			home, _ := os.UserHomeDir()
			root = filepath.Join(home, ".local", "share")
		}
	}
	dir := filepath.Join(root, "NightAuth", "release-sequences")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	hash := sha256.Sum256([]byte(application))
	return filepath.Join(dir, hex.EncodeToString(hash[:])+".seq"), nil
}

func readReleaseSequence(seqPath string) (int64, error) {
	data, err := os.ReadFile(seqPath)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, errors.New("stored release sequence is corrupt; refusing update")
	}
	trimmed := strings.TrimSpace(string(data))
	seq, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil || seq < 0 {
		return 0, errors.New("stored release sequence is corrupt; refusing update")
	}
	return seq, nil
}

func writeReleaseSequenceAtomically(seqPath string, sequence int64) error {
	randBytes := make([]byte, 8)
	_, _ = rand.Read(randBytes)
	temporary := fmt.Sprintf("%s.%s.tmp", seqPath, hex.EncodeToString(randBytes))

	f, err := os.OpenFile(temporary, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(f, "%d\n", sequence); err != nil {
		f.Close()
		_ = os.Remove(temporary)
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		_ = os.Remove(temporary)
		return err
	}
	f.Close()
	return os.Rename(temporary, seqPath)
}

var updateInstallerMutex sync.Mutex

func DownloadSignedUpdateWithInstaller(manifestURL, pinnedReleaseKeyHex, application, destination string, timeout time.Duration) (string, error) {
	parsedURL, err := url.Parse(manifestURL)
	if err != nil || parsedURL.Scheme != "https" || parsedURL.User != nil || parsedURL.Fragment != "" {
		return "", errors.New("release manifest must use an HTTPS URL without user info or fragment")
	}

	statePath, err := getReleaseSequencePath(application)
	if err != nil {
		return "", err
	}
	lockPath := statePath + ".lock"

	updateInstallerMutex.Lock()
	defer updateInstallerMutex.Unlock()

	start := time.Now()
	var lockFile *os.File
	for {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
		if err == nil {
			lockFile = f
			break
		}
		if time.Since(start) > 10*time.Second {
			return "", errors.New("could not acquire update lock")
		}
		time.Sleep(100 * time.Millisecond)
	}
	defer func() {
		if lockFile != nil {
			lockFile.Close()
			_ = os.Remove(lockPath)
		}
	}()

	currentSequence, err := readReleaseSequence(statePath)
	if err != nil {
		return "", err
	}

	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("redirects are not permitted")
		},
	}
	req, err := http.NewRequest(http.MethodGet, manifestURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "NightAuth-Updater/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("release manifest must return HTTP 200 without redirects (got %d)", resp.StatusCode)
	}

	manifestBytes, err := io.ReadAll(io.LimitReader(resp.Body, MaxSignedUpdateManifestBytes+1))
	if err != nil {
		return "", err
	}
	if int64(len(manifestBytes)) > MaxSignedUpdateManifestBytes {
		return "", errors.New("release manifest exceeds the maximum allowed size")
	}

	verified, err := VerifySignedUpdate(string(manifestBytes), pinnedReleaseKeyHex, application, currentSequence)
	if err != nil {
		return "", err
	}

	targetPath, err := DownloadSignedUpdate(verified, destination, 5*time.Minute)
	if err != nil {
		return "", err
	}

	if err := writeReleaseSequenceAtomically(statePath, verified.Sequence); err != nil {
		return "", err
	}

	return targetPath, nil
}

type signedUpdateHelper struct{}

var SignedUpdate signedUpdateHelper

func (signedUpdateHelper) Verify(envelopeJSON, pinnedPublicKeyHex, expectedApplication string, minimumSequence int64) (*VerifiedUpdate, error) {
	return VerifySignedUpdate(envelopeJSON, pinnedPublicKeyHex, expectedApplication, minimumSequence)
}

func (signedUpdateHelper) Download(update *VerifiedUpdate, destination string, timeout time.Duration) (string, error) {
	return DownloadSignedUpdate(update, destination, timeout)
}

type signedUpdateInstallerHelper struct{}

var SignedUpdateInstaller signedUpdateInstallerHelper

func (signedUpdateInstallerHelper) Download(manifestURL, pinnedReleaseKeyHex, application, destination string, timeout time.Duration) (string, error) {
	return DownloadSignedUpdateWithInstaller(manifestURL, pinnedReleaseKeyHex, application, destination, timeout)
}
