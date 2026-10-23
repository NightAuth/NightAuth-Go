package main

import (
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

func init() {
	runtime.LockOSThread()
}

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	msimg32  = syscall.NewLazyDLL("msimg32.dll")

	pGetModuleHandleW    = kernel32.NewProc("GetModuleHandleW")
	pRegisterClassExW    = user32.NewProc("RegisterClassExW")
	pCreateWindowExW     = user32.NewProc("CreateWindowExW")
	pDefWindowProcW      = user32.NewProc("DefWindowProcW")
	pShowWindow          = user32.NewProc("ShowWindow")
	pUpdateWindow        = user32.NewProc("UpdateWindow")
	pSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	pGetMessageW         = user32.NewProc("GetMessageW")
	pTranslateMessage    = user32.NewProc("TranslateMessage")
	pDispatchMessageW    = user32.NewProc("DispatchMessageW")
	pPostQuitMessage     = user32.NewProc("PostQuitMessage")
	pBeginPaint          = user32.NewProc("BeginPaint")
	pEndPaint            = user32.NewProc("EndPaint")
	pInvalidateRect      = user32.NewProc("InvalidateRect")
	pGetClientRect       = user32.NewProc("GetClientRect")
	pSetWindowPos        = user32.NewProc("SetWindowPos")
	pGetWindowTextW      = user32.NewProc("GetWindowTextW")
	pSetWindowTextW      = user32.NewProc("SetWindowTextW")
	pSendMessageW        = user32.NewProc("SendMessageW")
	pTrackMouseEvent     = user32.NewProc("TrackMouseEvent")
	pGetSystemMetrics    = user32.NewProc("GetSystemMetrics")
	pDrawTextW           = user32.NewProc("DrawTextW")
	pFillRect            = user32.NewProc("FillRect")
	pGetFocus            = user32.NewProc("GetFocus")
	pGetDC               = user32.NewProc("GetDC")
	pReleaseDC           = user32.NewProc("ReleaseDC")
	pScreenToClient      = user32.NewProc("ScreenToClient")
	pLoadCursorW         = user32.NewProc("LoadCursorW")
	pMessageBoxW         = user32.NewProc("MessageBoxW")
	pSetWindowLongPtrW   = user32.NewProc("SetWindowLongPtrW")
	pCallWindowProcW     = user32.NewProc("CallWindowProcW")

	pCreateSolidBrush    = gdi32.NewProc("CreateSolidBrush")
	pCreatePen           = gdi32.NewProc("CreatePen")
	pSelectObject        = gdi32.NewProc("SelectObject")
	pDeleteObject        = gdi32.NewProc("DeleteObject")
	pRoundRect           = gdi32.NewProc("RoundRect")
	pEllipse             = gdi32.NewProc("Ellipse")
	pRectangle           = gdi32.NewProc("Rectangle")
	pCreateFontW         = gdi32.NewProc("CreateFontW")
	pSetTextColor        = gdi32.NewProc("SetTextColor")
	pSetBkColor          = gdi32.NewProc("SetBkColor")
	pSetBkMode           = gdi32.NewProc("SetBkMode")
	pCreateCompatibleDC  = gdi32.NewProc("CreateCompatibleDC")
	pCreateCompatibleBmp = gdi32.NewProc("CreateCompatibleBitmap")
	pCreateDIBSection    = gdi32.NewProc("CreateDIBSection")
	pGetStockObject      = gdi32.NewProc("GetStockObject")
	pGetTextExtentPointW = gdi32.NewProc("GetTextExtentPoint32W")
	pBitBlt              = gdi32.NewProc("BitBlt")
	pDeleteDC            = gdi32.NewProc("DeleteDC")
	pMoveToEx            = gdi32.NewProc("MoveToEx")
	pLineTo              = gdi32.NewProc("LineTo")

	pAlphaBlend = msimg32.NewProc("AlphaBlend")
)

const (
	WM_CREATE       = 0x0001
	WM_DESTROY      = 0x0002
	WM_PAINT        = 0x000F
	WM_SETFONT      = 0x0030
	WM_NCHITTEST    = 0x0084
	WM_COMMAND      = 0x0111
	WM_CTLCOLOREDIT = 0x0133
	WM_MOUSEMOVE    = 0x0200
	WM_LBUTTONDOWN  = 0x0201
	WM_MOUSELEAVE   = 0x02A3

	EM_SETMARGINS  = 0x00D3
	EC_LEFTMARGIN  = 0x0001
	EC_RIGHTMARGIN = 0x0002
	EN_SETFOCUS    = 0x0100
	EN_KILLFOCUS   = 0x0200

	WS_POPUP        = 0x80000000
	WS_VISIBLE      = 0x10000000
	WS_CHILD        = 0x40000000
	WS_EX_APPWINDOW = 0x00040000

	ES_PASSWORD    = 0x0020
	ES_AUTOHSCROLL = 0x0080

	SW_HIDE     = 0
	SW_SHOW     = 5
	SW_MINIMIZE = 6

	SWP_NOZORDER   = 0x0004
	SWP_SHOWWINDOW = 0x0040

	HTCAPTION   = 2
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1

	DT_LEFT         = 0x0000
	DT_CENTER       = 0x0001
	DT_VCENTER      = 0x0004
	DT_WORDBREAK    = 0x0010
	DT_SINGLELINE   = 0x0020
	DT_NOPREFIX     = 0x0800
	DT_END_ELLIPSIS = 0x8000

	SRCCOPY     = 0x00CC0020
	TRANSPARENT = 1
	NULL_BRUSH  = 5

	WIN_WIDTH   = 1020
	WIN_HEIGHT  = 680
	TITLEBAR_H  = 50
	STATUSBAR_H = 32
	SIDEBAR_W   = 228

	ID_BTN_MINIMIZE = 1001
	ID_BTN_CLOSE    = 1002
	ID_BTN_ACTION   = 1003
	ID_BTN_LOGOUT   = 1004

	ID_TAB_PRE_BASE  = 1010
	ID_DASH_TAB_BASE = 1020

	ID_BTN_HEARTBEAT       = 1040
	ID_BTN_REFRESH_SUBS    = 1041
	ID_BTN_GET_GLOBAL      = 1042
	ID_BTN_GET_USER_VAR    = 1043
	ID_BTN_SET_USER_VAR    = 1044
	ID_BTN_GET_ALL_VARS    = 1045
	ID_BTN_SEND_CHAT       = 1046
	ID_BTN_REFRESH_CHAT    = 1047
	ID_BTN_REFRESH_STATS   = 1048
	ID_BTN_SETUP_2FA       = 1049
	ID_BTN_CONFIRM_2FA     = 1050
	ID_BTN_DISABLE_2FA     = 1051
	ID_BTN_CHECK_BL        = 1052
	ID_BTN_VERIFY_HASH     = 1053
	ID_BTN_KILLSWITCH      = 1054
	ID_BTN_DOWNLOAD_FILE   = 1055
	ID_BTN_TRIGGER_HOOK    = 1056
	ID_BTN_CHANGE_USERNAME = 1057
)

func rgb(r, g, b byte) uint32 {
	return uint32(r) | (uint32(g) << 8) | (uint32(b) << 16)
}

var (
	colCanvasBlack   = rgb(9, 10, 15)
	colSidebarBg     = rgb(13, 15, 22)
	colHeaderBg      = rgb(11, 13, 19)
	colHeroBg        = rgb(13, 16, 26)
	colCardBg        = rgb(17, 20, 29)
	colCardHeaderBg  = rgb(21, 25, 36)
	colCardBorder    = rgb(36, 42, 58)
	colCardHighlight = rgb(52, 62, 86)
	colInputBg       = rgb(11, 13, 20)
	colInputBorder   = rgb(42, 49, 68)
	colInputFocus    = rgb(59, 130, 246)

	colBluePrimary = rgb(37, 99, 235)
	colBlueHover   = rgb(59, 130, 246)
	colBlueLight   = rgb(147, 197, 253)
	colBlueGlow    = rgb(96, 165, 250)
	colBlueSubtle  = rgb(23, 37, 68)

	colTextHead  = rgb(248, 250, 252)
	colTextBody  = rgb(203, 213, 225)
	colTextLabel = rgb(148, 163, 184)
	colTextMuted = rgb(100, 116, 139)

	colSuccess   = rgb(52, 211, 153)
	colSuccessBg = rgb(14, 46, 38)
	colWarning   = rgb(251, 191, 36)
	colWarningBg = rgb(54, 39, 14)
	colError     = rgb(248, 113, 113)
	colErrorBg   = rgb(56, 20, 26)
	colDanger    = rgb(220, 38, 38)
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

type POINT struct {
	X, Y int32
}

type SIZE struct {
	Cx, Cy int32
}

type MSG struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type PAINTSTRUCT struct {
	Hdc         uintptr
	FErase      int32
	RcPaint     RECT
	FRestore    int32
	FIncUpdate  int32
	RgbReserved [32]byte
}

type TRACKMOUSEEVENT struct {
	CbSize      uint32
	DwFlags     uint32
	HwndTrack   uintptr
	DwHoverTime uint32
}

type WNDCLASSEXW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}

type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors [1]uint32
}

type ViewMode int

const (
	ViewLogin ViewMode = iota
	ViewRegister
	ViewLicense
	ViewReset
	ViewUpgrade
	ViewForgot
	ViewTokenReset
	ViewDashboard
)

type DashboardTab int

const (
	DashOverview DashboardTab = iota
	DashSubscriptions
	DashVariables
	DashChat
	DashStats
	DashSecurity
	DashFilesHooks
	DashCount
)

type AppState struct {
	api           *nightauth.Api
	currentView   ViewMode
	activeDashTab DashboardTab
	authenticated bool

	hbrCanvas     uintptr
	hbrSidebar    uintptr
	hbrHeader     uintptr
	hbrCard       uintptr
	hbrCardHeader uintptr
	hbrInput      uintptr

	hFontBrand   uintptr
	hFontHero    uintptr
	hFontTitle   uintptr
	hFontSection uintptr
	hFontBody    uintptr
	hFontInput   uintptr
	hFontLabel   uintptr
	hFontBadge   uintptr
	hFontStatus  uintptr
	hFontMono    uintptr
	hFontKpi     uintptr
	hbmLogo      uintptr

	hWnd uintptr

	hEditUser       uintptr
	hEditPass       uintptr
	hEditKey        uintptr
	hEdit2FA        uintptr
	hEditEmail      uintptr
	hEditResetToken uintptr
	hEditNewPass    uintptr

	hEditChangeNewUser uintptr
	hEditChangePass    uintptr

	hEditGlobalKey   uintptr
	hEditVarKey      uintptr
	hEditVarVal      uintptr
	hEditChatRoom    uintptr
	hEditChatMsg     uintptr
	hEdit2FACode     uintptr
	hEditDisablePass uintptr
	hEditDisableCode uintptr
	hEditCheckHwid   uintptr
	hEditVerifyHash  uintptr
	hEditBanReason   uintptr
	hEditFileId      uintptr
	hEditHookId      uintptr
	hEditHookPayload uintptr

	statusText  string
	statusColor uint32
	isBusy      bool
	hoveredBtn  int

	lastGlobalVal     string
	lastUserVal       string
	userVarsList      []map[string]string
	chatMsgs          []nightauth.ChatMessage
	appStats          *nightauth.AppStats
	onlineUsers       []string
	tfaSetup          *nightauth.TwoFactorSetupData
	downloadStatus    string
	downloadStatusOk  bool
	webhookResponse   string
	webhookResponseOk bool
}

const (
	AppName    = ""
	AppApiKey  = ""
	AppVersion = "1.0"
	AppToken   = ""
)

var (
	gApp                 AppState
	origEditProc         uintptr
	editSubclassCallback uintptr
)

func editSubclassProc(hWnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	if msg == 0x0100 && wParam == 0x0D { // WM_KEYDOWN, VK_RETURN
		handleActionClick()
		return 0
	}
	ret, _, _ := pCallWindowProcW.Call(origEditProc, hWnd, uintptr(msg), wParam, lParam)
	return ret
}

func formatTimestamp(raw string, fallback string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed == "0" {
		return fallback
	}
	if sec, err := strconv.ParseInt(trimmed, 10, 64); err == nil && sec > 0 {
		if sec > 1000000000000 {
			sec /= 1000
		}
		return time.Unix(sec, 0).UTC().Format("2006-01-02 15:04 UTC")
	}
	if t, err := time.Parse(time.RFC3339, trimmed); err == nil {
		return t.UTC().Format("2006-01-02 15:04 UTC")
	}
	return trimmed
}

func setStatus(text string, col uint32) {
	gApp.statusText = text
	gApp.statusColor = col
	if gApp.hWnd != 0 {
		pInvalidateRect.Call(gApp.hWnd, 0, 0)
		pUpdateWindow.Call(gApp.hWnd)
	}
}

func getEditText(h uintptr) string {
	if h == 0 {
		return ""
	}
	buf := make([]uint16, 512)
	pGetWindowTextW.Call(h, uintptr(unsafe.Pointer(&buf[0])), 512)
	return syscall.UTF16ToString(buf)
}

func setEditText(h uintptr, text string) {
	if h == 0 {
		return
	}
	pSetWindowTextW.Call(h, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))))
}

func drawText(hdc uintptr, text string, rc *RECT, flags uint32) {
	if text == "" {
		return
	}
	u16, err := syscall.UTF16FromString(text)
	if err != nil || len(u16) == 0 {
		return
	}
	pDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&u16[0])), uintptr(len(u16)-1), uintptr(unsafe.Pointer(rc)), uintptr(flags|DT_NOPREFIX))
}

func drawLine(hdc uintptr, x1, y1, x2, y2 int32, col uint32) {
	hPen, _, _ := pCreatePen.Call(0, 1, uintptr(col))
	oldP, _, _ := pSelectObject.Call(hdc, hPen)
	pMoveToEx.Call(hdc, uintptr(x1), uintptr(y1), 0)
	pLineTo.Call(hdc, uintptr(x2), uintptr(y2))
	pSelectObject.Call(hdc, oldP)
	pDeleteObject.Call(hPen)
}

func drawFilledRoundRect(hdc uintptr, rc RECT, radius int32, fillCol, borderCol uint32) {
	hbr, _, _ := pCreateSolidBrush.Call(uintptr(fillCol))
	hPen, _, _ := pCreatePen.Call(0, 1, uintptr(borderCol))
	oldBr, _, _ := pSelectObject.Call(hdc, hbr)
	oldPen, _, _ := pSelectObject.Call(hdc, hPen)
	pRoundRect.Call(hdc, uintptr(rc.Left), uintptr(rc.Top), uintptr(rc.Right), uintptr(rc.Bottom), uintptr(radius), uintptr(radius))
	pSelectObject.Call(hdc, oldBr)
	pSelectObject.Call(hdc, oldPen)
	pDeleteObject.Call(hbr)
	pDeleteObject.Call(hPen)
}

func drawCardRect(hdc uintptr, rc RECT, radius int32, hbrFill uintptr, borderCol uint32) {
	hPen, _, _ := pCreatePen.Call(0, 1, uintptr(borderCol))
	oldPen, _, _ := pSelectObject.Call(hdc, hPen)
	oldBr, _, _ := pSelectObject.Call(hdc, hbrFill)
	pRoundRect.Call(hdc, uintptr(rc.Left), uintptr(rc.Top), uintptr(rc.Right), uintptr(rc.Bottom), uintptr(radius), uintptr(radius))
	pSelectObject.Call(hdc, oldBr)
	pSelectObject.Call(hdc, oldPen)
	pDeleteObject.Call(hPen)

	if rc.Right-rc.Left > 40 && borderCol == colCardBorder {
		drawLine(hdc, rc.Left+radius/2+2, rc.Top+1, rc.Right-radius/2-2, rc.Top+1, colCardHighlight)
	}
}

func drawSectionCard(hdc uintptr, rc RECT, title, badgeText string) {
	drawCardRect(hdc, rc, 10, gApp.hbrCard, colCardBorder)

	accentRc := RECT{Left: rc.Left + 18, Top: rc.Top + 15, Right: rc.Left + 21, Bottom: rc.Top + 29}
	accBr, _, _ := pCreateSolidBrush.Call(uintptr(colBlueGlow))
	pFillRect.Call(hdc, uintptr(unsafe.Pointer(&accentRc)), accBr)
	pDeleteObject.Call(accBr)

	pSetBkMode.Call(hdc, TRANSPARENT)
	pSetTextColor.Call(hdc, uintptr(colTextHead))
	pSelectObject.Call(hdc, gApp.hFontSection)
	titleRc := RECT{Left: rc.Left + 29, Top: rc.Top + 12, Right: rc.Right - 140, Bottom: rc.Top + 32}
	drawText(hdc, title, &titleRc, DT_LEFT|DT_VCENTER|DT_SINGLELINE)

	if badgeText != "" {
		var sz SIZE
		pSelectObject.Call(hdc, gApp.hFontBadge)
		u16, _ := syscall.UTF16FromString(badgeText)
		if len(u16) > 1 {
			pGetTextExtentPointW.Call(hdc, uintptr(unsafe.Pointer(&u16[0])), uintptr(len(u16)-1), uintptr(unsafe.Pointer(&sz)))
		}
		bw := sz.Cx + 18
		bRc := RECT{Left: rc.Right - 18 - bw, Top: rc.Top + 12, Right: rc.Right - 18, Bottom: rc.Top + 32}
		drawFilledRoundRect(hdc, bRc, 6, colBlueSubtle, rgb(37, 62, 112))
		pSetTextColor.Call(hdc, uintptr(colBlueLight))
		drawText(hdc, badgeText, &bRc, DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	}

	drawLine(hdc, rc.Left+1, rc.Top+42, rc.Right-1, rc.Top+42, rgb(28, 33, 46))
}

func drawInputBox(hdc uintptr, rc RECT, hEdit uintptr) {
	focusWnd, _, _ := pGetFocus.Call()
	focused := focusWnd == hEdit && hEdit != 0
	border := colInputBorder
	if focused {
		border = colInputFocus
	}
	drawCardRect(hdc, rc, 6, gApp.hbrInput, border)
	if focused {
		glowRc := RECT{Left: rc.Left - 1, Top: rc.Top - 1, Right: rc.Right + 1, Bottom: rc.Bottom + 1}
		hPen, _, _ := pCreatePen.Call(0, 1, uintptr(rgb(30, 64, 138)))
		oldPen, _, _ := pSelectObject.Call(hdc, hPen)
		nullBr, _, _ := pGetStockObject.Call(NULL_BRUSH)
		oldBr, _, _ := pSelectObject.Call(hdc, nullBr)
		pRoundRect.Call(hdc, uintptr(glowRc.Left), uintptr(glowRc.Top), uintptr(glowRc.Right), uintptr(glowRc.Bottom), 8, 8)
		pSelectObject.Call(hdc, oldBr)
		pSelectObject.Call(hdc, oldPen)
		pDeleteObject.Call(hPen)
	}
}

func drawBadge(hdc uintptr, x, y int32, text string, fg, bg, border uint32) {
	pSelectObject.Call(hdc, gApp.hFontBadge)
	var sz SIZE
	u16, _ := syscall.UTF16FromString(text)
	if len(u16) > 1 {
		pGetTextExtentPointW.Call(hdc, uintptr(unsafe.Pointer(&u16[0])), uintptr(len(u16)-1), uintptr(unsafe.Pointer(&sz)))
	}
	rc := RECT{Left: x, Top: y, Right: x + sz.Cx + 18, Bottom: y + 22}
	drawFilledRoundRect(hdc, rc, 6, bg, border)
	pSetBkMode.Call(hdc, TRANSPARENT)
	pSetTextColor.Call(hdc, uintptr(fg))
	drawText(hdc, text, &rc, DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}

func drawStatusDot(hdc uintptr, cx, cy, r int32, col uint32) {
	hbr, _, _ := pCreateSolidBrush.Call(uintptr(col))
	hpen, _, _ := pCreatePen.Call(0, 1, uintptr(col))
	oldBr, _, _ := pSelectObject.Call(hdc, hbr)
	oldPen, _, _ := pSelectObject.Call(hdc, hpen)
	pEllipse.Call(hdc, uintptr(cx-r), uintptr(cy-r), uintptr(cx+r), uintptr(cy+r))
	pSelectObject.Call(hdc, oldBr)
	pSelectObject.Call(hdc, oldPen)
	pDeleteObject.Call(hbr)
	pDeleteObject.Call(hpen)
}

func drawTabPill(hdc uintptr, rc RECT, text string, hovered, active bool) {
	bgCol := rgb(13, 15, 22)
	borderCol := rgb(28, 33, 46)
	textCol := colTextLabel
	if active {
		bgCol = colBluePrimary
		borderCol = colBlueGlow
		textCol = rgb(255, 255, 255)
	} else if hovered {
		bgCol = rgb(25, 30, 44)
		borderCol = rgb(55, 68, 98)
		textCol = colTextHead
	}

	drawFilledRoundRect(hdc, rc, 6, bgCol, borderCol)

	pSetBkMode.Call(hdc, TRANSPARENT)
	pSetTextColor.Call(hdc, uintptr(textCol))
	if active {
		pSelectObject.Call(hdc, gApp.hFontBadge)
	} else {
		pSelectObject.Call(hdc, gApp.hFontBody)
	}
	drawText(hdc, text, &rc, DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}

func drawSidebarPill(hdc uintptr, rc RECT, code, text string, hovered, active bool) {
	bgCol := colSidebarBg
	borderCol := colSidebarBg
	textCol := colTextLabel
	if active {
		bgCol = rgb(22, 35, 64)
		borderCol = rgb(45, 82, 156)
		textCol = colTextHead
	} else if hovered {
		bgCol = rgb(20, 24, 35)
		borderCol = colCardBorder
		textCol = colTextBody
	}

	drawFilledRoundRect(hdc, rc, 8, bgCol, borderCol)

	if active {
		barRc := RECT{Left: rc.Left + 4, Top: rc.Top + 8, Right: rc.Left + 7, Bottom: rc.Bottom - 8}
		barBr, _, _ := pCreateSolidBrush.Call(uintptr(colBlueGlow))
		pFillRect.Call(hdc, uintptr(unsafe.Pointer(&barRc)), barBr)
		pDeleteObject.Call(barBr)
	}

	idxRc := RECT{Left: rc.Left + 13, Top: rc.Top + 8, Right: rc.Left + 35, Bottom: rc.Bottom - 8}
	idxBg := rgb(18, 21, 30)
	idxBorder := rgb(34, 40, 56)
	if active {
		idxBg = colBluePrimary
		idxBorder = colBlueGlow
	} else if hovered {
		idxBg = rgb(28, 33, 48)
	}
	drawFilledRoundRect(hdc, idxRc, 5, idxBg, idxBorder)

	pSetBkMode.Call(hdc, TRANSPARENT)
	if active {
		pSetTextColor.Call(hdc, uintptr(rgb(255, 255, 255)))
	} else {
		pSetTextColor.Call(hdc, uintptr(colTextMuted))
	}
	pSelectObject.Call(hdc, gApp.hFontBadge)
	drawText(hdc, code, &idxRc, DT_CENTER|DT_VCENTER|DT_SINGLELINE)

	pSetTextColor.Call(hdc, uintptr(textCol))
	if active {
		pSelectObject.Call(hdc, gApp.hFontSection)
	} else {
		pSelectObject.Call(hdc, gApp.hFontBody)
	}
	textRc := RECT{Left: rc.Left + 44, Top: rc.Top, Right: rc.Right - 10, Bottom: rc.Bottom}
	drawText(hdc, text, &textRc, DT_LEFT|DT_VCENTER|DT_SINGLELINE)
}

func drawPrimaryButton(hdc uintptr, rc RECT, text string, hovered bool) {
	bgCol := colBluePrimary
	borderCol := colBlueGlow
	topHighlight := rgb(96, 165, 250)
	if hovered {
		bgCol = colBlueHover
		borderCol = colBlueLight
		topHighlight = rgb(147, 197, 253)
	}

	drawFilledRoundRect(hdc, rc, 7, bgCol, borderCol)
	drawLine(hdc, rc.Left+4, rc.Top+1, rc.Right-4, rc.Top+1, topHighlight)

	pSetBkMode.Call(hdc, TRANSPARENT)
	pSetTextColor.Call(hdc, uintptr(rgb(255, 255, 255)))
	pSelectObject.Call(hdc, gApp.hFontSection)
	drawText(hdc, text, &rc, DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}

func drawSecondaryButton(hdc uintptr, rc RECT, text string, hovered bool) {
	bgCol := rgb(21, 25, 36)
	borderCol := rgb(45, 53, 74)
	textCol := colTextBody
	if hovered {
		bgCol = rgb(28, 34, 49)
		borderCol = colBlueGlow
		textCol = colTextHead
	}

	drawFilledRoundRect(hdc, rc, 7, bgCol, borderCol)

	pSetBkMode.Call(hdc, TRANSPARENT)
	pSetTextColor.Call(hdc, uintptr(textCol))
	pSelectObject.Call(hdc, gApp.hFontSection)
	drawText(hdc, text, &rc, DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}

func drawDangerButton(hdc uintptr, rc RECT, text string, hovered bool) {
	bgCol := colDanger
	borderCol := rgb(248, 113, 113)
	if hovered {
		bgCol = rgb(239, 68, 68)
		borderCol = rgb(254, 202, 202)
	}

	drawFilledRoundRect(hdc, rc, 7, bgCol, borderCol)

	pSetBkMode.Call(hdc, TRANSPARENT)
	pSetTextColor.Call(hdc, uintptr(rgb(255, 255, 255)))
	pSelectObject.Call(hdc, gApp.hFontSection)
	drawText(hdc, text, &rc, DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}

func drawBentoCard(hdc uintptr, rc RECT, label, value string, monoValue bool, valColor uint32) {
	drawCardRect(hdc, rc, 8, gApp.hbrCard, colCardBorder)
	pSetBkMode.Call(hdc, TRANSPARENT)

	pSetTextColor.Call(hdc, uintptr(colTextMuted))
	pSelectObject.Call(hdc, gApp.hFontLabel)
	lRc := RECT{Left: rc.Left + 16, Top: rc.Top + 14, Right: rc.Right - 14, Bottom: rc.Top + 30}
	drawText(hdc, label, &lRc, DT_LEFT|DT_SINGLELINE)

	pSetTextColor.Call(hdc, uintptr(valColor))
	if monoValue {
		pSelectObject.Call(hdc, gApp.hFontMono)
	} else {
		pSelectObject.Call(hdc, gApp.hFontSection)
	}
	vRc := RECT{Left: rc.Left + 16, Top: rc.Top + 38, Right: rc.Right - 14, Bottom: rc.Bottom - 12}
	val := value
	if strings.TrimSpace(val) == "" {
		val = "N/A"
	}
	drawText(hdc, val, &vRc, DT_LEFT|DT_VCENTER|DT_SINGLELINE|DT_END_ELLIPSIS)
}

func drawInfoTile(hdc uintptr, rc RECT, label, value string, monoValue bool, valColor uint32) {
	drawFilledRoundRect(hdc, rc, 7, colInputBg, rgb(30, 36, 50))
	pSetBkMode.Call(hdc, TRANSPARENT)

	pSetTextColor.Call(hdc, uintptr(colTextMuted))
	pSelectObject.Call(hdc, gApp.hFontLabel)
	lRc := RECT{Left: rc.Left + 14, Top: rc.Top + 10, Right: rc.Right - 12, Bottom: rc.Top + 26}
	drawText(hdc, label, &lRc, DT_LEFT|DT_SINGLELINE)

	pSetTextColor.Call(hdc, uintptr(valColor))
	if monoValue {
		pSelectObject.Call(hdc, gApp.hFontMono)
	} else {
		pSelectObject.Call(hdc, gApp.hFontSection)
	}
	vRc := RECT{Left: rc.Left + 14, Top: rc.Top + 30, Right: rc.Right - 12, Bottom: rc.Bottom - 8}
	val := value
	if strings.TrimSpace(val) == "" {
		val = "N/A"
	}
	drawText(hdc, val, &vRc, DT_LEFT|DT_VCENTER|DT_SINGLELINE|DT_END_ELLIPSIS)
}

func updateControlsVisibility() {
	allEdits := []uintptr{
		gApp.hEditUser, gApp.hEditPass, gApp.hEditKey, gApp.hEdit2FA,
		gApp.hEditEmail, gApp.hEditResetToken, gApp.hEditNewPass,
		gApp.hEditChangeNewUser, gApp.hEditChangePass,
		gApp.hEditGlobalKey, gApp.hEditVarKey, gApp.hEditVarVal,
		gApp.hEditChatRoom, gApp.hEditChatMsg, gApp.hEdit2FACode,
		gApp.hEditDisablePass, gApp.hEditDisableCode, gApp.hEditCheckHwid,
		gApp.hEditVerifyHash, gApp.hEditBanReason, gApp.hEditFileId,
		gApp.hEditHookId, gApp.hEditHookPayload,
	}
	for _, h := range allEdits {
		if h != 0 {
			pShowWindow.Call(h, SW_HIDE)
		}
	}

	flags := uintptr(SWP_NOZORDER | SWP_SHOWWINDOW)

	if gApp.currentView != ViewDashboard {
		cardW := int32(860)
		cardX := (int32(WIN_WIDTH) - cardW) / 2
		formX := cardX + 340
		formW := cardW - 340
		fieldX := formX + 36
		fieldW := formW - 72

		switch gApp.currentView {
		case ViewLogin:
			pSetWindowPos.Call(gApp.hEditUser, 0, uintptr(fieldX+4), 249, uintptr(fieldW-8), 24, flags)
			pSetWindowPos.Call(gApp.hEditPass, 0, uintptr(fieldX+4), 319, uintptr(fieldW-8), 24, flags)
			pSetWindowPos.Call(gApp.hEdit2FA, 0, uintptr(fieldX+4), 389, uintptr(fieldW-8), 24, flags)
		case ViewRegister:
			pSetWindowPos.Call(gApp.hEditUser, 0, uintptr(fieldX+4), 232, uintptr(fieldW-8), 22, flags)
			pSetWindowPos.Call(gApp.hEditPass, 0, uintptr(fieldX+4), 290, uintptr(fieldW-8), 22, flags)
			pSetWindowPos.Call(gApp.hEditKey, 0, uintptr(fieldX+4), 348, uintptr(fieldW-8), 22, flags)
			pSetWindowPos.Call(gApp.hEditEmail, 0, uintptr(fieldX+4), 406, uintptr(fieldW-8), 22, flags)
		case ViewLicense:
			pSetWindowPos.Call(gApp.hEditKey, 0, uintptr(fieldX+4), 269, uintptr(fieldW-8), 24, flags)
		case ViewReset:
			pSetWindowPos.Call(gApp.hEditUser, 0, uintptr(fieldX+4), 249, uintptr(fieldW-8), 24, flags)
			pSetWindowPos.Call(gApp.hEditPass, 0, uintptr(fieldX+4), 319, uintptr(fieldW-8), 24, flags)
			pSetWindowPos.Call(gApp.hEdit2FA, 0, uintptr(fieldX+4), 389, uintptr(fieldW-8), 24, flags)
		case ViewUpgrade:
			pSetWindowPos.Call(gApp.hEditUser, 0, uintptr(fieldX+4), 259, uintptr(fieldW-8), 24, flags)
			pSetWindowPos.Call(gApp.hEditKey, 0, uintptr(fieldX+4), 329, uintptr(fieldW-8), 24, flags)
		case ViewForgot:
			pSetWindowPos.Call(gApp.hEditEmail, 0, uintptr(fieldX+4), 269, uintptr(fieldW-8), 24, flags)
		case ViewTokenReset:
			pSetWindowPos.Call(gApp.hEditResetToken, 0, uintptr(fieldX+4), 259, uintptr(fieldW-8), 24, flags)
			pSetWindowPos.Call(gApp.hEditNewPass, 0, uintptr(fieldX+4), 329, uintptr(fieldW-8), 24, flags)
		}
	} else {
		contentX := int32(256)
		cw := int32(WIN_WIDTH) - contentX - 28

		switch gApp.activeDashTab {
		case DashVariables:
			pSetWindowPos.Call(gApp.hEditGlobalKey, 0, uintptr(contentX+24), 179, 252, 22, flags)
			pSetWindowPos.Call(gApp.hEditVarKey, 0, uintptr(contentX+24), 314, 212, 22, flags)
			pSetWindowPos.Call(gApp.hEditVarVal, 0, uintptr(contentX+254), 314, 212, 22, flags)
		case DashChat:
			pSetWindowPos.Call(gApp.hEditChatRoom, 0, uintptr(contentX+24), 169, 192, 22, flags)
			pSetWindowPos.Call(gApp.hEditChatMsg, 0, uintptr(contentX+24), 576, uintptr(cw-172), 22, flags)
		case DashSecurity:
			pSetWindowPos.Call(gApp.hEditChangeNewUser, 0, uintptr(contentX+24), 174, 242, 22, flags)
			pSetWindowPos.Call(gApp.hEditChangePass, 0, uintptr(contentX+284), 174, 242, 22, flags)
			pSetWindowPos.Call(gApp.hEdit2FACode, 0, uintptr(contentX+24), 343, 172, 22, flags)
			pSetWindowPos.Call(gApp.hEditDisablePass, 0, uintptr(contentX+379), 343, 132, 22, flags)
			pSetWindowPos.Call(gApp.hEditDisableCode, 0, uintptr(contentX+529), 343, 96, 22, flags)
			pSetWindowPos.Call(gApp.hEditCheckHwid, 0, uintptr(contentX+24), 485, 342, 22, flags)
			pSetWindowPos.Call(gApp.hEditVerifyHash, 0, uintptr(contentX+24), 536, 342, 22, flags)
			pSetWindowPos.Call(gApp.hEditBanReason, 0, uintptr(contentX+24), 587, 342, 22, flags)
		case DashFilesHooks:
			pSetWindowPos.Call(gApp.hEditFileId, 0, uintptr(contentX+24), 179, 352, 22, flags)
			pSetWindowPos.Call(gApp.hEditHookId, 0, uintptr(contentX+24), 369, 352, 22, flags)
			pSetWindowPos.Call(gApp.hEditHookPayload, 0, uintptr(contentX+24), 439, 352, 22, flags)
		}
	}
	pInvalidateRect.Call(gApp.hWnd, 0, 0)
}

func onPaint(hWnd uintptr) {
	var ps PAINTSTRUCT
	hdcMem, _, _ := pBeginPaint.Call(hWnd, uintptr(unsafe.Pointer(&ps)))

	var clientRc RECT
	pGetClientRect.Call(hWnd, uintptr(unsafe.Pointer(&clientRc)))
	width := clientRc.Right
	height := clientRc.Bottom

	hdc, _, _ := pCreateCompatibleDC.Call(hdcMem)
	hbm, _, _ := pCreateCompatibleBmp.Call(hdcMem, uintptr(width), uintptr(height))
	oldBmp, _, _ := pSelectObject.Call(hdc, hbm)

	pFillRect.Call(hdc, uintptr(unsafe.Pointer(&clientRc)), gApp.hbrCanvas)

	topBar := RECT{Left: 0, Top: 0, Right: width, Bottom: TITLEBAR_H}
	pFillRect.Call(hdc, uintptr(unsafe.Pointer(&topBar)), gApp.hbrHeader)
	drawLine(hdc, 0, TITLEBAR_H-1, width, TITLEBAR_H-1, colCardBorder)

	if gApp.hbmLogo != 0 {
		hdcLogo, _, _ := pCreateCompatibleDC.Call(hdc)
		oldL, _, _ := pSelectObject.Call(hdcLogo, gApp.hbmLogo)
		var bf uintptr = (0x01 << 24) | (0xFF << 16)
		pAlphaBlend.Call(hdc, 16, 11, LOGO_WIDTH, LOGO_HEIGHT, hdcLogo, 0, 0, LOGO_WIDTH, LOGO_HEIGHT, bf)
		pSelectObject.Call(hdcLogo, oldL)
		pDeleteDC.Call(hdcLogo)
	}

	pSetBkMode.Call(hdc, TRANSPARENT)
	pSetTextColor.Call(hdc, uintptr(colTextHead))
	pSelectObject.Call(hdc, gApp.hFontBrand)
	brandRc := RECT{Left: 54, Top: 10, Right: 260, Bottom: 28}
	drawText(hdc, "NIGHTAUTH", &brandRc, DT_LEFT|DT_SINGLELINE)

	pSetTextColor.Call(hdc, uintptr(colTextMuted))
	pSelectObject.Call(hdc, gApp.hFontLabel)
	subBrandRc := RECT{Left: 54, Top: 27, Right: 300, Bottom: 42}
	drawText(hdc, "NATIVE GO SECURITY CLIENT", &subBrandRc, DT_LEFT|DT_SINGLELINE)

	minRc := RECT{Left: width - 76, Top: 10, Right: width - 44, Bottom: 40}
	closeRc := RECT{Left: width - 40, Top: 10, Right: width - 8, Bottom: 40}
	if gApp.hoveredBtn == ID_BTN_MINIMIZE {
		drawFilledRoundRect(hdc, minRc, 6, rgb(26, 31, 44), colCardBorder)
	}
	if gApp.hoveredBtn == ID_BTN_CLOSE {
		drawFilledRoundRect(hdc, closeRc, 6, colDanger, rgb(248, 113, 113))
	}
	if gApp.hoveredBtn == ID_BTN_MINIMIZE {
		pSetTextColor.Call(hdc, uintptr(colTextHead))
	} else {
		pSetTextColor.Call(hdc, uintptr(colTextLabel))
	}
	pSelectObject.Call(hdc, gApp.hFontSection)
	drawText(hdc, "-", &minRc, DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	if gApp.hoveredBtn == ID_BTN_CLOSE {
		pSetTextColor.Call(hdc, uintptr(rgb(255, 255, 255)))
	} else {
		pSetTextColor.Call(hdc, uintptr(colTextLabel))
	}
	drawText(hdc, "X", &closeRc, DT_CENTER|DT_VCENTER|DT_SINGLELINE)

	if gApp.currentView != ViewDashboard {
		cardW := int32(860)
		cardH := int32(540)
		cardX := (width - cardW) / 2
		cardY := int32(68)
		heroW := int32(340)
		formX := cardX + heroW
		formW := cardW - heroW

		mainCard := RECT{Left: cardX, Top: cardY, Right: cardX + cardW, Bottom: cardY + cardH}
		drawCardRect(hdc, mainCard, 12, gApp.hbrCard, colCardBorder)

		// Left Hero Branding Panel
		heroPanel := RECT{Left: cardX + 1, Top: cardY + 1, Right: formX, Bottom: cardY + cardH - 1}
		heroBr, _, _ := pCreateSolidBrush.Call(uintptr(rgb(13, 16, 25)))
		pFillRect.Call(hdc, uintptr(unsafe.Pointer(&heroPanel)), heroBr)
		pDeleteObject.Call(heroBr)
		drawLine(hdc, formX, cardY+1, formX, cardY+cardH-1, colCardBorder)

		// Top Cobalt Accent Bar on Hero
		heroAccent := RECT{Left: cardX + 28, Top: cardY + 32, Right: cardX + 76, Bottom: cardY + 36}
		accBr, _, _ := pCreateSolidBrush.Call(uintptr(colBluePrimary))
		pFillRect.Call(hdc, uintptr(unsafe.Pointer(&heroAccent)), accBr)
		pDeleteObject.Call(accBr)

		drawBadge(hdc, cardX+28, cardY+50, "●  ZERO-TRUST TELEMETRY", colBlueLight, colBlueSubtle, rgb(37, 62, 112))

		pSetTextColor.Call(hdc, uintptr(colTextHead))
		pSelectObject.Call(hdc, gApp.hFontHero)
		hTitle1 := RECT{Left: cardX + 28, Top: cardY + 86, Right: formX - 24, Bottom: cardY + 116}
		drawText(hdc, "Enterprise Cloud", &hTitle1, DT_LEFT|DT_SINGLELINE)
		hTitle2 := RECT{Left: cardX + 28, Top: cardY + 114, Right: formX - 24, Bottom: cardY + 144}
		pSetTextColor.Call(hdc, uintptr(colBlueGlow))
		drawText(hdc, "Licensing & Identity", &hTitle2, DT_LEFT|DT_SINGLELINE)

		pSetTextColor.Call(hdc, uintptr(colTextLabel))
		pSelectObject.Call(hdc, gApp.hFontBody)
		hDesc := RECT{Left: cardX + 28, Top: cardY + 154, Right: formX - 28, Bottom: cardY + 218}
		drawText(hdc, "Cryptographically bound sessions, hardware fingerprinting, real-time cloud variables, and tamper-evident delivery.", &hDesc, DT_WORDBREAK)

		// Feature Checkmarks
		heroFeatures := [4]string{
			"Instant License Key & Account Activation",
			"Automatic Device HWID Binding & Reset",
			"Cloud Variables, Live Chat & File Delivery",
			"Two-Factor (2FA) Account Protection",
		}
		for i := 0; i < 4; i++ {
			fy := cardY + 232 + int32(i)*34
			dotBox := RECT{Left: cardX + 28, Top: fy + 2, Right: cardX + 46, Bottom: fy + 20}
			drawFilledRoundRect(hdc, dotBox, 5, colSuccessBg, rgb(24, 92, 70))
			pSetTextColor.Call(hdc, uintptr(colSuccess))
			pSelectObject.Call(hdc, gApp.hFontBadge)
			drawText(hdc, "+", &dotBox, DT_CENTER|DT_VCENTER|DT_SINGLELINE)

			pSetTextColor.Call(hdc, uintptr(colTextBody))
			pSelectObject.Call(hdc, gApp.hFontBody)
			fTxt := RECT{Left: cardX + 56, Top: fy, Right: formX - 20, Bottom: fy + 22}
			drawText(hdc, heroFeatures[i], &fTxt, DT_LEFT|DT_VCENTER|DT_SINGLELINE)
		}

		// Bound Application & Device Card at bottom of Hero
		hwidCard := RECT{Left: cardX + 24, Top: cardY + cardH - 124, Right: formX - 24, Bottom: cardY + cardH - 24}
		drawFilledRoundRect(hdc, hwidCard, 8, colInputBg, rgb(30, 36, 52))

		pSetTextColor.Call(hdc, uintptr(colBlueGlow))
		pSelectObject.Call(hdc, gApp.hFontLabel)
		hwLbl := RECT{Left: hwidCard.Left + 14, Top: hwidCard.Top + 12, Right: hwidCard.Right - 14, Bottom: hwidCard.Top + 26}
		drawText(hdc, "BOUND APPLICATION & DEVICE", &hwLbl, DT_LEFT|DT_SINGLELINE)

		appNameDisp := gApp.api.Name
		if appNameDisp == "" {
			appNameDisp = "Unconfigured"
		}
		pSetTextColor.Call(hdc, uintptr(colTextHead))
		pSelectObject.Call(hdc, gApp.hFontSection)
		appLbl := RECT{Left: hwidCard.Left + 14, Top: hwidCard.Top + 32, Right: hwidCard.Right - 14, Bottom: hwidCard.Top + 52}
		drawText(hdc, appNameDisp, &appLbl, DT_LEFT|DT_SINGLELINE|DT_END_ELLIPSIS)

		hwidPreview := nightauth.SecurityProtections{}.GetHardwareId()
		if hwidPreview == "" {
			hwidPreview = "Detected on Auth"
		}
		pSetTextColor.Call(hdc, uintptr(colTextMuted))
		pSelectObject.Call(hdc, gApp.hFontMono)
		hwVal := RECT{Left: hwidCard.Left + 14, Top: hwidCard.Top + 58, Right: hwidCard.Right - 14, Bottom: hwidCard.Top + 82}
		drawText(hdc, "HWID: "+hwidPreview, &hwVal, DT_LEFT|DT_SINGLELINE|DT_END_ELLIPSIS)

		// Right Side: Navigation Pills + Active Auth Form
		fieldX := formX + 36
		fieldW := formW - 72

		tabGap := int32(5)
		tabW := (fieldW - tabGap*6) / 7
		tabY := cardY + 26
		tabH := int32(30)

		preTabs := []string{"Login", "Register", "License", "HWID", "Upgrade", "Forgot", "Token"}
		for i, name := range preTabs {
			tabRc := RECT{Left: fieldX + int32(i)*(tabW+tabGap), Top: tabY, Right: fieldX + int32(i)*(tabW+tabGap) + tabW, Bottom: tabY + tabH}
			drawTabPill(hdc, tabRc, name, gApp.hoveredBtn == (ID_TAB_PRE_BASE+i), gApp.currentView == ViewMode(i))
		}

		viewTitles := []string{
			"Sign in to your account",
			"Create a new account",
			"Instant license key access",
			"Reset hardware binding",
			"Upgrade account subscription",
			"Recover account password",
			"Complete password reset",
		}
		viewSubs := []string{
			"Enter your credentials to establish an encrypted cloud session.",
			"Activate your license key and bind your credentials.",
			"Authenticate directly using a valid license key without a password.",
			"Unbind your previous machine fingerprint to authorize this device.",
			"Apply a new license key tier to an existing registered user.",
			"Request a time-limited recovery token to your registered email.",
			"Provide the recovery token and choose a new account password.",
		}

		pSetTextColor.Call(hdc, uintptr(colTextHead))
		pSelectObject.Call(hdc, gApp.hFontTitle)
		vtRc := RECT{Left: fieldX, Top: cardY + 74, Right: fieldX + fieldW, Bottom: cardY + 102}
		drawText(hdc, viewTitles[int(gApp.currentView)], &vtRc, DT_LEFT|DT_SINGLELINE)

		pSetTextColor.Call(hdc, uintptr(colTextLabel))
		pSelectObject.Call(hdc, gApp.hFontBody)
		vsRc := RECT{Left: fieldX, Top: cardY + 104, Right: fieldX + fieldW, Bottom: cardY + 126}
		drawText(hdc, viewSubs[int(gApp.currentView)], &vsRc, DT_LEFT|DT_SINGLELINE)

		drawLine(hdc, fieldX, cardY+136, fieldX+fieldW, cardY+136, rgb(28, 33, 46))

		btnRc := RECT{Left: fieldX, Top: 456, Right: fieldX + fieldW, Bottom: 500}

		switch gApp.currentView {
		case ViewLogin:
			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l1 := RECT{Left: fieldX, Top: 224, Right: fieldX + fieldW, Bottom: 240}
			drawText(hdc, "USERNAME", &l1, DT_LEFT)
			f1 := RECT{Left: fieldX, Top: 243, Right: fieldX + fieldW, Bottom: 279}
			drawInputBox(hdc, f1, gApp.hEditUser)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l2 := RECT{Left: fieldX, Top: 294, Right: fieldX + fieldW, Bottom: 310}
			drawText(hdc, "PASSWORD", &l2, DT_LEFT)
			f2 := RECT{Left: fieldX, Top: 313, Right: fieldX + fieldW, Bottom: 349}
			drawInputBox(hdc, f2, gApp.hEditPass)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l3 := RECT{Left: fieldX, Top: 364, Right: fieldX + fieldW, Bottom: 380}
			drawText(hdc, "TWO-FACTOR CODE (OPTIONAL)", &l3, DT_LEFT)
			f3 := RECT{Left: fieldX, Top: 383, Right: fieldX + fieldW, Bottom: 419}
			drawInputBox(hdc, f3, gApp.hEdit2FA)

			btnLabel := "Sign In to Dashboard"
			if gApp.isBusy {
				btnLabel = "Authenticating Session..."
			}
			drawPrimaryButton(hdc, btnRc, btnLabel, gApp.hoveredBtn == ID_BTN_ACTION)

		case ViewRegister:
			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l1 := RECT{Left: fieldX, Top: 209, Right: fieldX + fieldW, Bottom: 224}
			drawText(hdc, "USERNAME", &l1, DT_LEFT)
			f1 := RECT{Left: fieldX, Top: 226, Right: fieldX + fieldW, Bottom: 260}
			drawInputBox(hdc, f1, gApp.hEditUser)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l2 := RECT{Left: fieldX, Top: 267, Right: fieldX + fieldW, Bottom: 282}
			drawText(hdc, "PASSWORD", &l2, DT_LEFT)
			f2 := RECT{Left: fieldX, Top: 284, Right: fieldX + fieldW, Bottom: 318}
			drawInputBox(hdc, f2, gApp.hEditPass)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l3 := RECT{Left: fieldX, Top: 325, Right: fieldX + fieldW, Bottom: 340}
			drawText(hdc, "ACTIVATION LICENSE KEY", &l3, DT_LEFT)
			f3 := RECT{Left: fieldX, Top: 342, Right: fieldX + fieldW, Bottom: 376}
			drawInputBox(hdc, f3, gApp.hEditKey)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l4 := RECT{Left: fieldX, Top: 383, Right: fieldX + fieldW, Bottom: 398}
			drawText(hdc, "EMAIL ADDRESS (OPTIONAL)", &l4, DT_LEFT)
			f4 := RECT{Left: fieldX, Top: 400, Right: fieldX + fieldW, Bottom: 434}
			drawInputBox(hdc, f4, gApp.hEditEmail)

			btnLabel := "Create Account"
			if gApp.isBusy {
				btnLabel = "Creating Account..."
			}
			drawPrimaryButton(hdc, btnRc, btnLabel, gApp.hoveredBtn == ID_BTN_ACTION)

		case ViewLicense:
			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l1 := RECT{Left: fieldX, Top: 244, Right: fieldX + fieldW, Bottom: 260}
			drawText(hdc, "LICENSE KEY", &l1, DT_LEFT)
			f1 := RECT{Left: fieldX, Top: 263, Right: fieldX + fieldW, Bottom: 299}
			drawInputBox(hdc, f1, gApp.hEditKey)

			infoCard := RECT{Left: fieldX, Top: 322, Right: fieldX + fieldW, Bottom: 398}
			drawFilledRoundRect(hdc, infoCard, 8, colBlueSubtle, rgb(37, 62, 112))
			pSetTextColor.Call(hdc, uintptr(colBlueLight))
			pSelectObject.Call(hdc, gApp.hFontBody)
			noteRc := RECT{Left: fieldX + 16, Top: 336, Right: fieldX + fieldW - 16, Bottom: 386}
			drawText(hdc, "License-only mode binds your key directly to this machine's hardware ID for instant passwordless authentication.", &noteRc, DT_WORDBREAK)

			btnLabel := "Authenticate with Key"
			if gApp.isBusy {
				btnLabel = "Validating License..."
			}
			drawPrimaryButton(hdc, btnRc, btnLabel, gApp.hoveredBtn == ID_BTN_ACTION)

		case ViewReset:
			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l1 := RECT{Left: fieldX, Top: 224, Right: fieldX + fieldW, Bottom: 240}
			drawText(hdc, "ACCOUNT USERNAME", &l1, DT_LEFT)
			f1 := RECT{Left: fieldX, Top: 243, Right: fieldX + fieldW, Bottom: 279}
			drawInputBox(hdc, f1, gApp.hEditUser)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l2 := RECT{Left: fieldX, Top: 294, Right: fieldX + fieldW, Bottom: 310}
			drawText(hdc, "ACCOUNT PASSWORD", &l2, DT_LEFT)
			f2 := RECT{Left: fieldX, Top: 313, Right: fieldX + fieldW, Bottom: 349}
			drawInputBox(hdc, f2, gApp.hEditPass)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l3 := RECT{Left: fieldX, Top: 364, Right: fieldX + fieldW, Bottom: 380}
			drawText(hdc, "2FA CODE (IF ENABLED)", &l3, DT_LEFT)
			f3 := RECT{Left: fieldX, Top: 383, Right: fieldX + fieldW, Bottom: 419}
			drawInputBox(hdc, f3, gApp.hEdit2FA)

			btnLabel := "Reset Hardware Fingerprint"
			if gApp.isBusy {
				btnLabel = "Resetting Binding..."
			}
			drawPrimaryButton(hdc, btnRc, btnLabel, gApp.hoveredBtn == ID_BTN_ACTION)

		case ViewUpgrade:
			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l1 := RECT{Left: fieldX, Top: 234, Right: fieldX + fieldW, Bottom: 250}
			drawText(hdc, "EXISTING USERNAME", &l1, DT_LEFT)
			f1 := RECT{Left: fieldX, Top: 253, Right: fieldX + fieldW, Bottom: 289}
			drawInputBox(hdc, f1, gApp.hEditUser)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l2 := RECT{Left: fieldX, Top: 304, Right: fieldX + fieldW, Bottom: 320}
			drawText(hdc, "NEW UPGRADE LICENSE KEY", &l2, DT_LEFT)
			f2 := RECT{Left: fieldX, Top: 323, Right: fieldX + fieldW, Bottom: 359}
			drawInputBox(hdc, f2, gApp.hEditKey)

			btnLabel := "Apply Subscription Upgrade"
			if gApp.isBusy {
				btnLabel = "Applying Upgrade..."
			}
			drawPrimaryButton(hdc, btnRc, btnLabel, gApp.hoveredBtn == ID_BTN_ACTION)

		case ViewForgot:
			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l1 := RECT{Left: fieldX, Top: 244, Right: fieldX + fieldW, Bottom: 260}
			drawText(hdc, "REGISTERED EMAIL ADDRESS", &l1, DT_LEFT)
			f1 := RECT{Left: fieldX, Top: 263, Right: fieldX + fieldW, Bottom: 299}
			drawInputBox(hdc, f1, gApp.hEditEmail)

			infoCard := RECT{Left: fieldX, Top: 322, Right: fieldX + fieldW, Bottom: 398}
			drawFilledRoundRect(hdc, infoCard, 8, colBlueSubtle, rgb(37, 62, 112))
			pSetTextColor.Call(hdc, uintptr(colBlueLight))
			pSelectObject.Call(hdc, gApp.hFontBody)
			noteRc := RECT{Left: fieldX + 16, Top: 336, Right: fieldX + fieldW - 16, Bottom: 386}
			drawText(hdc, "We will dispatch a one-time cryptographic reset token to your verified email address. Then switch to the Token tab.", &noteRc, DT_WORDBREAK)

			btnLabel := "Send Recovery Token"
			if gApp.isBusy {
				btnLabel = "Dispatching Email..."
			}
			drawPrimaryButton(hdc, btnRc, btnLabel, gApp.hoveredBtn == ID_BTN_ACTION)

		case ViewTokenReset:
			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l1 := RECT{Left: fieldX, Top: 234, Right: fieldX + fieldW, Bottom: 250}
			drawText(hdc, "RECOVERY RESET TOKEN", &l1, DT_LEFT)
			f1 := RECT{Left: fieldX, Top: 253, Right: fieldX + fieldW, Bottom: 289}
			drawInputBox(hdc, f1, gApp.hEditResetToken)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			l2 := RECT{Left: fieldX, Top: 304, Right: fieldX + fieldW, Bottom: 320}
			drawText(hdc, "NEW ACCOUNT PASSWORD", &l2, DT_LEFT)
			f2 := RECT{Left: fieldX, Top: 323, Right: fieldX + fieldW, Bottom: 359}
			drawInputBox(hdc, f2, gApp.hEditNewPass)

			btnLabel := "Confirm New Password"
			if gApp.isBusy {
				btnLabel = "Updating Password..."
			}
			drawPrimaryButton(hdc, btnRc, btnLabel, gApp.hoveredBtn == ID_BTN_ACTION)
		}
	} else {
		contentBottom := height - STATUSBAR_H
		sideRc := RECT{Left: 0, Top: TITLEBAR_H, Right: SIDEBAR_W, Bottom: contentBottom}
		pFillRect.Call(hdc, uintptr(unsafe.Pointer(&sideRc)), gApp.hbrSidebar)
		drawLine(hdc, SIDEBAR_W, TITLEBAR_H, SIDEBAR_W, contentBottom, colCardBorder)

		pSetTextColor.Call(hdc, uintptr(colTextMuted))
		pSelectObject.Call(hdc, gApp.hFontLabel)
		navHdr := RECT{Left: 20, Top: TITLEBAR_H + 16, Right: SIDEBAR_W - 16, Bottom: TITLEBAR_H + 32}
		drawText(hdc, "WORKSPACE MODULES", &navHdr, DT_LEFT|DT_SINGLELINE)

		dashTabs := []string{
			"Overview", "Subscriptions", "Cloud Variables",
			"Live Chat", "App Statistics", "Security & 2FA", "Files & Hooks",
		}
		dashCodes := []string{"01", "02", "03", "04", "05", "06", "07"}

		sPadY := int32(TITLEBAR_H + 38)
		sTabH := int32(38)
		sGap := int32(5)
		for i, name := range dashTabs {
			pillRc := RECT{Left: 12, Top: sPadY + int32(i)*(sTabH+sGap), Right: SIDEBAR_W - 12, Bottom: sPadY + int32(i)*(sTabH+sGap) + sTabH}
			drawSidebarPill(hdc, pillRc, dashCodes[i], name, gApp.hoveredBtn == (ID_DASH_TAB_BASE+i), gApp.activeDashTab == DashboardTab(i))
		}

		profCard := RECT{Left: 12, Top: contentBottom - 108, Right: 216, Bottom: contentBottom - 12}
		drawFilledRoundRect(hdc, profCard, 8, colCardBg, colCardBorder)

		uname := gApp.api.User.Username
		if uname == "" {
			uname = "Authenticated"
		}

		pSetTextColor.Call(hdc, uintptr(colTextHead))
		pSelectObject.Call(hdc, gApp.hFontSection)
		pNameRc := RECT{Left: profCard.Left + 14, Top: profCard.Top + 10, Right: profCard.Right - 12, Bottom: profCard.Top + 28}
		drawText(hdc, uname, &pNameRc, DT_LEFT|DT_SINGLELINE|DT_END_ELLIPSIS)

		pSetTextColor.Call(hdc, uintptr(colSuccess))
		pSelectObject.Call(hdc, gApp.hFontBadge)
		pStatRc := RECT{Left: profCard.Left + 14, Top: profCard.Top + 28, Right: profCard.Right - 12, Bottom: profCard.Top + 44}
		drawText(hdc, "● Authenticated Session", &pStatRc, DT_LEFT|DT_SINGLELINE)

		logoutRc := RECT{Left: profCard.Left + 12, Top: profCard.Bottom - 40, Right: profCard.Right - 12, Bottom: profCard.Bottom - 10}
		drawSecondaryButton(hdc, logoutRc, "Sign Out", gApp.hoveredBtn == ID_BTN_LOGOUT)

		cx := int32(256)
		cw := width - cx - 28

		dashSubtitles := []string{
			"Session identity, hardware fingerprint, and live application telemetry.",
			"Active license entitlements, expiration windows, and tier assignments.",
			"Encrypted server-side global constants and per-user key-value storage.",
			"Real-time channel messaging with authenticated workspace members.",
			"Application-wide user metrics and active concurrent sessions.",
			"TOTP two-factor authentication, HWID blacklist check, and emergency lockout.",
			"Protected asset delivery via scoped access tickets and webhook dispatch.",
		}

		pSetTextColor.Call(hdc, uintptr(colTextHead))
		pSelectObject.Call(hdc, gApp.hFontTitle)
		headRc := RECT{Left: cx, Top: 58, Right: width - 200, Bottom: 84}
		drawText(hdc, dashTabs[int(gApp.activeDashTab)], &headRc, DT_LEFT|DT_SINGLELINE)

		pSetTextColor.Call(hdc, uintptr(colTextLabel))
		pSelectObject.Call(hdc, gApp.hFontBody)
		subHeadRc := RECT{Left: cx, Top: 84, Right: width - 200, Bottom: 104}
		drawText(hdc, dashSubtitles[int(gApp.activeDashTab)], &subHeadRc, DT_LEFT|DT_SINGLELINE)

		switch gApp.activeDashTab {
		case DashOverview:
			hbBtn := RECT{Left: width - 190, Top: 62, Right: width - 28, Bottom: 96}
			drawPrimaryButton(hdc, hbBtn, "Send Heartbeat", gApp.hoveredBtn == ID_BTN_HEARTBEAT)

			colW := (cw - 24) / 3
			cardH := int32(88)
			row1Y := int32(112)
			row2Y := int32(212)

			c1 := RECT{Left: cx, Top: row1Y, Right: cx + colW, Bottom: row1Y + cardH}
			drawBentoCard(hdc, c1, "ACCOUNT USERNAME", uname, false, colTextHead)

			email := gApp.api.User.Email
			if email == "" {
				email = "Not Bound"
			}
			c2 := RECT{Left: cx + colW + 12, Top: row1Y, Right: cx + colW*2 + 12, Bottom: row1Y + cardH}
			drawBentoCard(hdc, c2, "EMAIL ADDRESS", email, false, colTextBody)

			statusStr := "● ACTIVE"
			if gApp.api.User.Status != "" {
				statusStr = "● " + strings.ToUpper(gApp.api.User.Status)
			}
			c3 := RECT{Left: cx + colW*2 + 24, Top: row1Y, Right: cx + cw, Bottom: row1Y + cardH}
			drawBentoCard(hdc, c3, "ACCOUNT STATUS", statusStr, false, colSuccess)

			ip := gApp.api.User.Ip
			if ip == "" {
				ip = "Protected / TLS"
			}
			c4 := RECT{Left: cx, Top: row2Y, Right: cx + colW, Bottom: row2Y + cardH}
			drawBentoCard(hdc, c4, "CLIENT IP ADDRESS", ip, true, colBlueLight)

			createdStr := formatTimestamp(gApp.api.User.CreatedAt, "—")
			c5 := RECT{Left: cx + colW + 12, Top: row2Y, Right: cx + colW*2 + 12, Bottom: row2Y + cardH}
			drawBentoCard(hdc, c5, "CREATED DATE", createdStr, true, colTextBody)

			lastLoginStr := formatTimestamp(gApp.api.User.LastLoginAt, "—")
			c6 := RECT{Left: cx + colW*2 + 24, Top: row2Y, Right: cx + cw, Bottom: row2Y + cardH}
			drawBentoCard(hdc, c6, "LAST LOGIN", lastLoginStr, true, colTextBody)

			hwidBox := RECT{Left: cx, Top: 314, Right: cx + cw, Bottom: 388}
			drawCardRect(hdc, hwidBox, 8, gApp.hbrCard, colCardBorder)
			pSetTextColor.Call(hdc, uintptr(colBlueGlow))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			hwL := RECT{Left: cx + 20, Top: 328, Right: cx + cw - 20, Bottom: 344}
			drawText(hdc, "HARDWARE FINGERPRINT (HWID)", &hwL, DT_LEFT|DT_SINGLELINE)

			hwid := gApp.api.User.Hwid
			if hwid == "" {
				hwid = nightauth.SecurityProtections{}.GetHardwareId()
			}
			pSetTextColor.Call(hdc, uintptr(colTextHead))
			pSelectObject.Call(hdc, gApp.hFontMono)
			hwV := RECT{Left: cx + 20, Top: 350, Right: cx + cw - 20, Bottom: 376}
			drawText(hdc, hwid, &hwV, DT_LEFT|DT_VCENTER|DT_SINGLELINE|DT_END_ELLIPSIS)

			appBox := RECT{Left: cx, Top: 402, Right: cx + cw, Bottom: 590}
			drawSectionCard(hdc, appBox, "APPLICATION TELEMETRY & CLOUD SESSION", "REALTIME")

			thirdW := (cw - 64) / 3
			appName := gApp.api.Name
			if appName == "" {
				appName = "NightAuth App"
			}
			a1 := RECT{Left: cx + 20, Top: 458, Right: cx + 20 + thirdW, Bottom: 520}
			drawInfoTile(hdc, a1, "APPLICATION NAME", appName, false, colTextHead)

			ver := gApp.api.App.Version
			if ver == "" {
				ver = gApp.api.Version
			}
			a2 := RECT{Left: cx + 32 + thirdW, Top: 458, Right: cx + 32 + thirdW*2, Bottom: 520}
			drawInfoTile(hdc, a2, "TARGET VERSION", ver, true, colBlueGlow)

			onlineBuf := fmt.Sprintf("%d Online", gApp.api.App.OnlineUsers)
			a3 := RECT{Left: cx + 44 + thirdW*2, Top: 458, Right: cx + cw - 20, Bottom: 520}
			drawInfoTile(hdc, a3, "ACTIVE ONLINE USERS", onlineBuf, false, colSuccess)

			bannerRc := RECT{Left: cx + 20, Top: 532, Right: cx + cw - 20, Bottom: 576}
			drawFilledRoundRect(hdc, bannerRc, 6, colBlueSubtle, rgb(37, 62, 112))
			pSetTextColor.Call(hdc, uintptr(colBlueLight))
			pSelectObject.Call(hdc, gApp.hFontBody)
			bDesc := RECT{Left: bannerRc.Left + 14, Top: bannerRc.Top, Right: bannerRc.Right - 14, Bottom: bannerRc.Bottom}
			drawText(hdc, "All API payloads are verified against Ed25519 server signatures with anti-replay and HWID binding.", &bDesc, DT_LEFT|DT_VCENTER|DT_SINGLELINE|DT_END_ELLIPSIS)

		case DashSubscriptions:
			subBox := RECT{Left: cx, Top: 114, Right: cx + cw, Bottom: 618}
			subCount := len(gApp.api.User.Subscriptions)
			badgeStr := fmt.Sprintf("%d ACTIVE", subCount)
			drawSectionCard(hdc, subBox, "ACTIVE SUBSCRIPTIONS & ENTITLEMENTS", badgeStr)

			if subCount == 0 {
				emptyCard := RECT{Left: cx + 20, Top: 176, Right: cx + cw - 20, Bottom: 266}
				drawFilledRoundRect(hdc, emptyCard, 8, colInputBg, colCardBorder)
				pSetTextColor.Call(hdc, uintptr(colTextHead))
				pSelectObject.Call(hdc, gApp.hFontSection)
				e1 := RECT{Left: emptyCard.Left + 20, Top: emptyCard.Top + 22, Right: emptyCard.Right - 20, Bottom: emptyCard.Top + 44}
				drawText(hdc, "No Active Subscription Tiers", &e1, DT_LEFT|DT_SINGLELINE)

				pSetTextColor.Call(hdc, uintptr(colTextMuted))
				pSelectObject.Call(hdc, gApp.hFontBody)
				e2 := RECT{Left: emptyCard.Left + 20, Top: emptyCard.Top + 46, Right: emptyCard.Right - 20, Bottom: emptyCard.Top + 74}
				drawText(hdc, "This account does not currently have any recurring or lifetime subscription entitlements attached.", &e2, DT_LEFT|DT_SINGLELINE)
			} else {
				for i, sub := range gApp.api.User.Subscriptions {
					if i >= 7 {
						break
					}
					sy := int32(172 + i*58)
					itemRc := RECT{Left: cx + 20, Top: sy, Right: cx + cw - 20, Bottom: sy + 48}
					drawFilledRoundRect(hdc, itemRc, 8, colInputBg, colInputBorder)

					stripe := RECT{Left: itemRc.Left + 4, Top: itemRc.Top + 8, Right: itemRc.Left + 7, Bottom: itemRc.Bottom - 8}
					sBr, _, _ := pCreateSolidBrush.Call(uintptr(colSuccess))
					pFillRect.Call(hdc, uintptr(unsafe.Pointer(&stripe)), sBr)
					pDeleteObject.Call(sBr)

					pSetTextColor.Call(hdc, uintptr(colTextHead))
					pSelectObject.Call(hdc, gApp.hFontSection)
					tierRc := RECT{Left: itemRc.Left + 20, Top: sy + 8, Right: itemRc.Left + 260, Bottom: sy + 28}
					drawText(hdc, sub.Subscription, &tierRc, DT_LEFT|DT_SINGLELINE)

					keyVal := sub.Key
					if keyVal == "" {
						keyVal = "Active Account"
					}
					keyStr := "Key: " + keyVal
					pSetTextColor.Call(hdc, uintptr(colTextMuted))
					pSelectObject.Call(hdc, gApp.hFontMono)
					keyRc := RECT{Left: itemRc.Left + 20, Top: sy + 27, Right: itemRc.Right - 240, Bottom: sy + 44}
					drawText(hdc, keyStr, &keyRc, DT_LEFT|DT_SINGLELINE)

					expFormatted := formatTimestamp(sub.Expiry, "LIFETIME")
					expStr := "EXPIRES: " + expFormatted
					drawBadge(hdc, itemRc.Right-225, sy+13, expStr, colSuccess, colSuccessBg, rgb(24, 92, 70))
				}
			}

		case DashVariables:
			gBox := RECT{Left: cx, Top: 114, Right: cx + cw, Bottom: 224}
			drawSectionCard(hdc, gBox, "GLOBAL CLOUD VARIABLE", "READ-ONLY")

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lgk := RECT{Left: cx + 20, Top: 158, Right: cx + 280, Bottom: 173}
			drawText(hdc, "VARIABLE KEY", &lgk, DT_LEFT)

			f1 := RECT{Left: cx + 20, Top: 174, Right: cx + 280, Bottom: 208}
			drawInputBox(hdc, f1, gApp.hEditGlobalKey)
			b1 := RECT{Left: cx + 295, Top: 174, Right: cx + 425, Bottom: 208}
			drawPrimaryButton(hdc, b1, "Fetch Global", gApp.hoveredBtn == ID_BTN_GET_GLOBAL)

			gResPill := RECT{Left: cx + 440, Top: 174, Right: cx + cw - 20, Bottom: 208}
			drawFilledRoundRect(hdc, gResPill, 6, colInputBg, rgb(32, 38, 54))
			if gApp.lastGlobalVal != "" {
				pSetTextColor.Call(hdc, uintptr(colSuccess))
			} else {
				pSetTextColor.Call(hdc, uintptr(colTextMuted))
			}
			pSelectObject.Call(hdc, gApp.hFontMono)
			rVal := RECT{Left: gResPill.Left + 12, Top: gResPill.Top, Right: gResPill.Right - 12, Bottom: gResPill.Bottom}
			gText := gApp.lastGlobalVal
			if gText == "" {
				gText = "Result will appear here..."
			}
			drawText(hdc, gText, &rVal, DT_LEFT|DT_VCENTER|DT_SINGLELINE|DT_END_ELLIPSIS)

			uBox := RECT{Left: cx, Top: 238, Right: cx + cw, Bottom: 384}
			drawSectionCard(hdc, uBox, "USER CLOUD VARIABLE (GET / SET)", "READ & WRITE")

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			luk := RECT{Left: cx + 20, Top: 293, Right: cx + 240, Bottom: 308}
			drawText(hdc, "VARIABLE KEY", &luk, DT_LEFT)

			luv := RECT{Left: cx + 250, Top: 293, Right: cx + 470, Bottom: 308}
			drawText(hdc, "VALUE (FOR SET)", &luv, DT_LEFT)

			f2 := RECT{Left: cx + 20, Top: 309, Right: cx + 240, Bottom: 343}
			drawInputBox(hdc, f2, gApp.hEditVarKey)
			f3 := RECT{Left: cx + 250, Top: 309, Right: cx + 470, Bottom: 343}
			drawInputBox(hdc, f3, gApp.hEditVarVal)

			b2 := RECT{Left: cx + 485, Top: 309, Right: cx + 590, Bottom: 343}
			drawSecondaryButton(hdc, b2, "Get Value", gApp.hoveredBtn == ID_BTN_GET_USER_VAR)
			b3 := RECT{Left: cx + 600, Top: 309, Right: cx + 705, Bottom: 343}
			drawPrimaryButton(hdc, b3, "Save Value", gApp.hoveredBtn == ID_BTN_SET_USER_VAR)

			if gApp.lastUserVal != "" {
				ruVal := RECT{Left: cx + 20, Top: 350, Right: cx + cw - 20, Bottom: 374}
				pSetTextColor.Call(hdc, uintptr(colSuccess))
				pSelectObject.Call(hdc, gApp.hFontMono)
				drawText(hdc, "Retrieved Value: "+gApp.lastUserVal, &ruVal, DT_LEFT|DT_VCENTER|DT_SINGLELINE)
			}

			allBox := RECT{Left: cx, Top: 398, Right: cx + cw, Bottom: 618}
			drawSectionCard(hdc, allBox, "ALL USER VARIABLES", "")

			bAll := RECT{Left: cx + cw - 156, Top: 405, Right: cx + cw - 18, Bottom: 434}
			drawSecondaryButton(hdc, bAll, "Sync All Vars", gApp.hoveredBtn == ID_BTN_GET_ALL_VARS)

			if len(gApp.userVarsList) == 0 {
				pSetTextColor.Call(hdc, uintptr(colTextMuted))
				pSelectObject.Call(hdc, gApp.hFontBody)
				emptyV := RECT{Left: cx + 20, Top: 462, Right: cx + cw - 20, Bottom: 490}
				drawText(hdc, "Click 'Sync All Vars' to inspect all key-value pairs stored for this account.", &emptyV, DT_LEFT|DT_SINGLELINE)
			} else {
				for i, v := range gApp.userVarsList {
					if i >= 4 {
						break
					}
					vy := int32(452 + i*38)
					rowRc := RECT{Left: cx + 20, Top: vy, Right: cx + cw - 20, Bottom: vy + 32}
					drawFilledRoundRect(hdc, rowRc, 6, colInputBg, rgb(30, 36, 50))

					vstr := fmt.Sprintf("%s  =  %s", v["key"], v["value"])
					vr := RECT{Left: rowRc.Left + 14, Top: rowRc.Top, Right: rowRc.Right - 14, Bottom: rowRc.Bottom}
					pSetTextColor.Call(hdc, uintptr(colBlueLight))
					pSelectObject.Call(hdc, gApp.hFontMono)
					drawText(hdc, vstr, &vr, DT_LEFT|DT_VCENTER|DT_SINGLELINE)
				}
			}

		case DashChat:
			cBox := RECT{Left: cx, Top: 114, Right: cx + cw, Bottom: 618}
			drawSectionCard(hdc, cBox, "LIVE CHANNEL STREAM", "REALTIME")

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lr := RECT{Left: cx + 20, Top: 148, Right: cx + 220, Bottom: 163}
			drawText(hdc, "CHANNEL / ROOM", &lr, DT_LEFT)

			fRoom := RECT{Left: cx + 20, Top: 164, Right: cx + 220, Bottom: 198}
			drawInputBox(hdc, fRoom, gApp.hEditChatRoom)
			bRefChat := RECT{Left: cx + 235, Top: 164, Right: cx + 365, Bottom: 198}
			drawSecondaryButton(hdc, bRefChat, "Refresh Channel", gApp.hoveredBtn == ID_BTN_REFRESH_CHAT)

			msgBox := RECT{Left: cx + 20, Top: 210, Right: cx + cw - 20, Bottom: 542}
			drawFilledRoundRect(hdc, msgBox, 8, colCanvasBlack, colCardBorder)

			if len(gApp.chatMsgs) == 0 {
				emptyChat := RECT{Left: msgBox.Left + 20, Top: msgBox.Top + 145, Right: msgBox.Right - 20, Bottom: msgBox.Top + 185}
				pSetTextColor.Call(hdc, uintptr(colTextMuted))
				pSelectObject.Call(hdc, gApp.hFontBody)
				drawText(hdc, "No messages loaded in this channel yet. Click 'Refresh Channel' or send a message below.", &emptyChat, DT_CENTER|DT_VCENTER|DT_SINGLELINE)
			} else {
				for i, m := range gApp.chatMsgs {
					if i >= 8 {
						break
					}
					my := int32(222 + i*38)
					bubble := RECT{Left: msgBox.Left + 12, Top: my, Right: msgBox.Right - 12, Bottom: my + 32}
					drawFilledRoundRect(hdc, bubble, 6, colCardBg, rgb(30, 35, 48))

					pSetTextColor.Call(hdc, uintptr(colBlueGlow))
					pSelectObject.Call(hdc, gApp.hFontSection)
					authRc := RECT{Left: bubble.Left + 12, Top: bubble.Top, Right: bubble.Left + 150, Bottom: bubble.Bottom}
					drawText(hdc, m.Author, &authRc, DT_LEFT|DT_VCENTER|DT_SINGLELINE|DT_END_ELLIPSIS)

					pSetTextColor.Call(hdc, uintptr(colTextBody))
					pSelectObject.Call(hdc, gApp.hFontBody)
					mTxtRc := RECT{Left: bubble.Left + 156, Top: bubble.Top, Right: bubble.Right - 12, Bottom: bubble.Bottom}
					drawText(hdc, m.Message, &mTxtRc, DT_LEFT|DT_VCENTER|DT_SINGLELINE|DT_END_ELLIPSIS)
				}
			}

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lm := RECT{Left: cx + 20, Top: 554, Right: cx + cw - 140, Bottom: 569}
			drawText(hdc, "COMPOSE MESSAGE", &lm, DT_LEFT)

			fMsg := RECT{Left: cx + 20, Top: 571, Right: cx + cw - 144, Bottom: 605}
			drawInputBox(hdc, fMsg, gApp.hEditChatMsg)
			bSend := RECT{Left: cx + cw - 132, Top: 571, Right: cx + cw - 20, Bottom: 605}
			drawPrimaryButton(hdc, bSend, "Send", gApp.hoveredBtn == ID_BTN_SEND_CHAT)

		case DashStats:
			cardW := (cw - 36) / 4
			statLabels := []string{"TOTAL USERS", "ACTIVE USERS", "BANNED USERS", "TOTAL LICENSES"}
			statAccents := []uint32{colBlueGlow, colSuccess, colError, colWarning}
			var statVals [4]int64
			if gApp.appStats != nil {
				statVals[0] = gApp.appStats.TotalUsers
				statVals[1] = gApp.appStats.ActiveUsers
				statVals[2] = gApp.appStats.BannedUsers
				statVals[3] = gApp.appStats.TotalLicenses
			}

			for i := 0; i < 4; i++ {
				kpi := RECT{Left: cx + int32(i)*(cardW+12), Top: 114, Right: cx + int32(i)*(cardW+12) + cardW, Bottom: 212}
				drawCardRect(hdc, kpi, 10, gApp.hbrCard, colCardBorder)

				topStrip := RECT{Left: kpi.Left + 16, Top: kpi.Top + 14, Right: kpi.Left + 44, Bottom: kpi.Top + 17}
				stBr, _, _ := pCreateSolidBrush.Call(uintptr(statAccents[i]))
				pFillRect.Call(hdc, uintptr(unsafe.Pointer(&topStrip)), stBr)
				pDeleteObject.Call(stBr)

				pSetTextColor.Call(hdc, uintptr(colTextLabel))
				pSelectObject.Call(hdc, gApp.hFontLabel)
				kl := RECT{Left: kpi.Left + 16, Top: kpi.Top + 26, Right: kpi.Right - 14, Bottom: kpi.Top + 42}
				drawText(hdc, statLabels[i], &kl, DT_LEFT)

				vbuf := fmt.Sprintf("%d", statVals[i])
				pSetTextColor.Call(hdc, uintptr(colTextHead))
				pSelectObject.Call(hdc, gApp.hFontKpi)
				kv := RECT{Left: kpi.Left + 16, Top: kpi.Top + 46, Right: kpi.Right - 14, Bottom: kpi.Bottom - 10}
				drawText(hdc, vbuf, &kv, DT_LEFT|DT_VCENTER|DT_SINGLELINE)
			}

			onlineBox := RECT{Left: cx, Top: 228, Right: cx + cw, Bottom: 618}
			drawSectionCard(hdc, onlineBox, "ONLINE ACTIVE SESSIONS", "")

			bRefStats := RECT{Left: cx + cw - 156, Top: 235, Right: cx + cw - 18, Bottom: 264}
			drawSecondaryButton(hdc, bRefStats, "Refresh Metrics", gApp.hoveredBtn == ID_BTN_REFRESH_STATS)

			if len(gApp.onlineUsers) == 0 {
				pSetTextColor.Call(hdc, uintptr(colTextMuted))
				pSelectObject.Call(hdc, gApp.hFontBody)
				emptyO := RECT{Left: cx + 20, Top: 294, Right: cx + cw - 20, Bottom: 322}
				drawText(hdc, "Click 'Refresh Metrics' to query live application telemetry and active users.", &emptyO, DT_LEFT|DT_SINGLELINE)
			} else {
				for i, u := range gApp.onlineUsers {
					if i >= 8 {
						break
					}
					oy := int32(284 + i*38)
					orc := RECT{Left: cx + 20, Top: oy, Right: cx + cw - 20, Bottom: oy + 32}
					drawFilledRoundRect(hdc, orc, 6, colInputBg, rgb(30, 36, 50))
					drawStatusDot(hdc, orc.Left+16, oy+16, 4, colSuccess)

					pSetTextColor.Call(hdc, uintptr(colTextHead))
					pSelectObject.Call(hdc, gApp.hFontSection)
					uTxt := RECT{Left: orc.Left + 30, Top: oy, Right: orc.Right - 150, Bottom: oy + 32}
					drawText(hdc, u, &uTxt, DT_LEFT|DT_VCENTER|DT_SINGLELINE)

					drawBadge(hdc, orc.Right-120, oy+5, "ACTIVE SESSION", colSuccess, colSuccessBg, rgb(24, 92, 70))
				}
			}

		case DashSecurity:
			userBox := RECT{Left: cx, Top: 114, Right: cx + cw, Bottom: 222}
			drawSectionCard(hdc, userBox, "ACCOUNT IDENTITY & USERNAME CHANGE", "CREDENTIALS")

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lNewU := RECT{Left: cx + 20, Top: 158, Right: cx + 265, Bottom: 173}
			drawText(hdc, "NEW ACCOUNT USERNAME", &lNewU, DT_LEFT)
			fNewU := RECT{Left: cx + 20, Top: 169, Right: cx + 265, Bottom: 203}
			drawInputBox(hdc, fNewU, gApp.hEditChangeNewUser)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lPassU := RECT{Left: cx + 280, Top: 158, Right: cx + 525, Bottom: 173}
			drawText(hdc, "CURRENT PASSWORD (VERIFY)", &lPassU, DT_LEFT)
			fPassU := RECT{Left: cx + 280, Top: 169, Right: cx + 525, Bottom: 203}
			drawInputBox(hdc, fPassU, gApp.hEditChangePass)

			bChangeUser := RECT{Left: cx + 550, Top: 169, Right: cx + 720, Bottom: 203}
			drawPrimaryButton(hdc, bChangeUser, "Update Username", gApp.hoveredBtn == ID_BTN_CHANGE_USERNAME)

			tfaBox := RECT{Left: cx, Top: 234, Right: cx + cw, Bottom: 420}
			drawSectionCard(hdc, tfaBox, "TWO-FACTOR AUTHENTICATION (TOTP)", "RFC 6238")

			bSetup := RECT{Left: cx + 20, Top: 278, Right: cx + 185, Bottom: 310}
			drawPrimaryButton(hdc, bSetup, "Generate Secret", gApp.hoveredBtn == ID_BTN_SETUP_2FA)

			secPill := RECT{Left: cx + 200, Top: 278, Right: cx + cw - 20, Bottom: 310}
			drawFilledRoundRect(hdc, secPill, 6, colInputBg, rgb(32, 38, 54))
			sbuf := "Click 'Generate Secret' to enroll an authenticator app."
			if gApp.tfaSetup != nil && gApp.tfaSetup.Secret != "" {
				pSetTextColor.Call(hdc, uintptr(colSuccess))
				sbuf = "TOTP Secret: " + gApp.tfaSetup.Secret
			} else {
				pSetTextColor.Call(hdc, uintptr(colTextMuted))
			}
			pSelectObject.Call(hdc, gApp.hFontMono)
			sr := RECT{Left: secPill.Left + 12, Top: secPill.Top, Right: secPill.Right - 12, Bottom: secPill.Bottom}
			drawText(hdc, sbuf, &sr, DT_LEFT|DT_VCENTER|DT_SINGLELINE)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lCode := RECT{Left: cx + 20, Top: 322, Right: cx + 200, Bottom: 337}
			drawText(hdc, "CONFIRM TOTP CODE", &lCode, DT_LEFT)

			fCode := RECT{Left: cx + 20, Top: 338, Right: cx + 200, Bottom: 372}
			drawInputBox(hdc, fCode, gApp.hEdit2FACode)
			bConf := RECT{Left: cx + 212, Top: 338, Right: cx + 336, Bottom: 372}
			drawPrimaryButton(hdc, bConf, "Confirm 2FA", gApp.hoveredBtn == ID_BTN_CONFIRM_2FA)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lDisPass := RECT{Left: cx + 375, Top: 322, Right: cx + 515, Bottom: 337}
			drawText(hdc, "PASSWORD", &lDisPass, DT_LEFT)
			fDisPass := RECT{Left: cx + 375, Top: 338, Right: cx + 515, Bottom: 372}
			drawInputBox(hdc, fDisPass, gApp.hEditDisablePass)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lDisCode := RECT{Left: cx + 525, Top: 322, Right: cx + 629, Bottom: 337}
			drawText(hdc, "2FA CODE", &lDisCode, DT_LEFT)
			fDisCode := RECT{Left: cx + 525, Top: 338, Right: cx + 629, Bottom: 372}
			drawInputBox(hdc, fDisCode, gApp.hEditDisableCode)

			bDis := RECT{Left: cx + 640, Top: 338, Right: cx + cw - 20, Bottom: 372}
			drawSecondaryButton(hdc, bDis, "Disable 2FA", gApp.hoveredBtn == ID_BTN_DISABLE_2FA)

			diagBox := RECT{Left: cx, Top: 432, Right: cx + cw, Bottom: 618}
			drawSectionCard(hdc, diagBox, "SECURITY INTEGRITY & EMERGENCY KILLSWITCH", "PROTECTION")

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lh := RECT{Left: cx + 20, Top: 466, Right: cx + 370, Bottom: 480}
			drawText(hdc, "HARDWARE ID (BLACKLIST CHECK)", &lh, DT_LEFT)
			fHwid := RECT{Left: cx + 20, Top: 480, Right: cx + 370, Bottom: 514}
			drawInputBox(hdc, fHwid, gApp.hEditCheckHwid)
			bCheckBl := RECT{Left: cx + 385, Top: 480, Right: cx + 535, Bottom: 514}
			drawSecondaryButton(hdc, bCheckBl, "Check Blacklist", gApp.hoveredBtn == ID_BTN_CHECK_BL)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lha := RECT{Left: cx + 20, Top: 517, Right: cx + 370, Bottom: 531}
			drawText(hdc, "FILE INTEGRITY HASH (SHA-256)", &lha, DT_LEFT)
			fHash := RECT{Left: cx + 20, Top: 531, Right: cx + 370, Bottom: 565}
			drawInputBox(hdc, fHash, gApp.hEditVerifyHash)
			bHash := RECT{Left: cx + 385, Top: 531, Right: cx + 535, Bottom: 565}
			drawSecondaryButton(hdc, bHash, "Verify Hash", gApp.hoveredBtn == ID_BTN_VERIFY_HASH)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lkr := RECT{Left: cx + 20, Top: 568, Right: cx + 370, Bottom: 582}
			drawText(hdc, "KILLSWITCH REASON (EMERGENCY LOCKOUT)", &lkr, DT_LEFT)
			fReason := RECT{Left: cx + 20, Top: 582, Right: cx + 370, Bottom: 612}
			drawInputBox(hdc, fReason, gApp.hEditBanReason)
			bKill := RECT{Left: cx + 385, Top: 582, Right: cx + 535, Bottom: 612}
			drawDangerButton(hdc, bKill, "Killswitch / Ban", gApp.hoveredBtn == ID_BTN_KILLSWITCH)

		case DashFilesHooks:
			fileBox := RECT{Left: cx, Top: 114, Right: cx + cw, Bottom: 286}
			drawSectionCard(hdc, fileBox, "PROTECTED REMOTE FILE DOWNLOADER", "SCOPED TICKET")

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lfid := RECT{Left: cx + 20, Top: 158, Right: cx + 380, Bottom: 173}
			drawText(hdc, "FILE ID / ASSET IDENTIFIER", &lfid, DT_LEFT)

			fFile := RECT{Left: cx + 20, Top: 174, Right: cx + 380, Bottom: 208}
			drawInputBox(hdc, fFile, gApp.hEditFileId)
			bDown := RECT{Left: cx + 395, Top: 174, Right: cx + 555, Bottom: 208}
			drawPrimaryButton(hdc, bDown, "Download Asset", gApp.hoveredBtn == ID_BTN_DOWNLOAD_FILE)

			dStatusBox := RECT{Left: cx + 20, Top: 224, Right: cx + cw - 20, Bottom: 266}
			drawFilledRoundRect(hdc, dStatusBox, 6, colInputBg, rgb(30, 36, 50))
			dTxt := "Ready to request scoped download ticket."
			if gApp.downloadStatus != "" {
				dTxt = gApp.downloadStatus
				if gApp.downloadStatusOk {
					pSetTextColor.Call(hdc, uintptr(colSuccess))
				} else {
					pSetTextColor.Call(hdc, uintptr(colError))
				}
			} else {
				pSetTextColor.Call(hdc, uintptr(colTextMuted))
			}
			pSelectObject.Call(hdc, gApp.hFontMono)
			rDown := RECT{Left: dStatusBox.Left + 12, Top: dStatusBox.Top, Right: dStatusBox.Right - 12, Bottom: dStatusBox.Bottom}
			drawText(hdc, dTxt, &rDown, DT_LEFT|DT_VCENTER|DT_SINGLELINE)

			hookBox := RECT{Left: cx, Top: 304, Right: cx + cw, Bottom: 618}
			drawSectionCard(hdc, hookBox, "AUTHENTICATED WEBHOOK TRIGGER", "SERVER-SIDE")

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lhid := RECT{Left: cx + 20, Top: 348, Right: cx + 380, Bottom: 363}
			drawText(hdc, "WEBHOOK IDENTIFIER / ID", &lhid, DT_LEFT)

			fHookId := RECT{Left: cx + 20, Top: 364, Right: cx + 380, Bottom: 398}
			drawInputBox(hdc, fHookId, gApp.hEditHookId)

			pSetTextColor.Call(hdc, uintptr(colTextLabel))
			pSelectObject.Call(hdc, gApp.hFontLabel)
			lhpl := RECT{Left: cx + 20, Top: 418, Right: cx + 380, Bottom: 433}
			drawText(hdc, "JSON PAYLOAD (DATA)", &lhpl, DT_LEFT)

			fHookPayload := RECT{Left: cx + 20, Top: 434, Right: cx + 380, Bottom: 468}
			drawInputBox(hdc, fHookPayload, gApp.hEditHookPayload)
			bHook := RECT{Left: cx + 395, Top: 434, Right: cx + 555, Bottom: 468}
			drawPrimaryButton(hdc, bHook, "Trigger Webhook", gApp.hoveredBtn == ID_BTN_TRIGGER_HOOK)

			hRespBox := RECT{Left: cx + 20, Top: 486, Right: cx + cw - 20, Bottom: 598}
			drawFilledRoundRect(hdc, hRespBox, 6, colInputBg, rgb(30, 36, 50))
			hTxt := "Webhook response payload will be displayed here..."
			if gApp.webhookResponse != "" {
				hTxt = gApp.webhookResponse
				if gApp.webhookResponseOk {
					pSetTextColor.Call(hdc, uintptr(colSuccess))
				} else {
					pSetTextColor.Call(hdc, uintptr(colError))
				}
			} else {
				pSetTextColor.Call(hdc, uintptr(colTextMuted))
			}
			pSelectObject.Call(hdc, gApp.hFontMono)
			rHook := RECT{Left: hRespBox.Left + 12, Top: hRespBox.Top + 10, Right: hRespBox.Right - 12, Bottom: hRespBox.Bottom - 10}
			drawText(hdc, hTxt, &rHook, DT_WORDBREAK)
		}
	}

	statusBar := RECT{Left: 0, Top: height - STATUSBAR_H, Right: width, Bottom: height}
	pFillRect.Call(hdc, uintptr(unsafe.Pointer(&statusBar)), gApp.hbrHeader)
	drawLine(hdc, 0, height-STATUSBAR_H, width, height-STATUSBAR_H, colCardBorder)

	dotCol := gApp.statusColor
	if dotCol == 0 {
		dotCol = colBlueGlow
	}
	drawStatusDot(hdc, 20, height-STATUSBAR_H/2, 4, dotCol)

	statusRc := RECT{Left: 32, Top: height - STATUSBAR_H, Right: width - 20, Bottom: height}
	sTxt := "Ready."
	if gApp.statusText != "" {
		sTxt = gApp.statusText
		pSetTextColor.Call(hdc, uintptr(gApp.statusColor))
	} else {
		pSetTextColor.Call(hdc, uintptr(colTextMuted))
	}
	pSelectObject.Call(hdc, gApp.hFontStatus)
	drawText(hdc, sTxt, &statusRc, DT_LEFT|DT_VCENTER|DT_SINGLELINE)

	framePen, _, _ := pCreatePen.Call(0, 1, uintptr(rgb(42, 49, 68)))
	oldFP, _, _ := pSelectObject.Call(hdc, framePen)
	nullBr, _, _ := pGetStockObject.Call(NULL_BRUSH)
	oldFB, _, _ := pSelectObject.Call(hdc, nullBr)
	pRectangle.Call(hdc, 0, 0, uintptr(width), uintptr(height))
	pSelectObject.Call(hdc, oldFB)
	pSelectObject.Call(hdc, oldFP)
	pDeleteObject.Call(framePen)

	pBitBlt.Call(hdcMem, 0, 0, uintptr(width), uintptr(height), hdc, 0, 0, SRCCOPY)

	pSelectObject.Call(hdc, oldBmp)
	pDeleteObject.Call(hbm)
	pDeleteDC.Call(hdc)
	pEndPaint.Call(hWnd, uintptr(unsafe.Pointer(&ps)))
}

func handleActionClick() {
	if gApp.isBusy {
		return
	}
	gApp.isBusy = true
	setStatus("Connecting to secure cloud...", colBlueGlow)

	user := strings.TrimSpace(getEditText(gApp.hEditUser))
	pass := strings.TrimSpace(getEditText(gApp.hEditPass))
	key := strings.TrimSpace(getEditText(gApp.hEditKey))
	code := strings.TrimSpace(getEditText(gApp.hEdit2FA))
	email := strings.TrimSpace(getEditText(gApp.hEditEmail))
	tok := strings.TrimSpace(getEditText(gApp.hEditResetToken))
	newPass := strings.TrimSpace(getEditText(gApp.hEditNewPass))

	go func() {
		defer func() {
			gApp.isBusy = false
			if gApp.hWnd != 0 {
				pInvalidateRect.Call(gApp.hWnd, 0, 0)
			}
		}()

		var ok bool
		switch gApp.currentView {
		case ViewLogin:
			if user == "" || pass == "" {
				setStatus("Please enter username and password.", colError)
				return
			}
			ok = gApp.api.Login(user, pass, code)
		case ViewRegister:
			if user == "" || pass == "" || key == "" {
				setStatus("Please provide username, password and license key.", colError)
				return
			}
			ok = gApp.api.Register(user, pass, key, email)
		case ViewLicense:
			if key == "" {
				setStatus("Please enter a valid license key.", colError)
				return
			}
			ok = gApp.api.License(key)
		case ViewReset:
			if user == "" || pass == "" {
				setStatus("Please enter username and password.", colError)
				return
			}
			ok = gApp.api.ResetHwid(user, pass, code)
		case ViewUpgrade:
			if user == "" || key == "" {
				setStatus("Please enter username and upgrade key.", colError)
				return
			}
			ok = gApp.api.Upgrade(user, key)
		case ViewForgot:
			if email == "" {
				setStatus("Please enter registered email address.", colError)
				return
			}
			ok = gApp.api.ForgotPassword(email)
		case ViewTokenReset:
			if tok == "" || newPass == "" {
				setStatus("Please enter reset token and new password.", colError)
				return
			}
			ok = gApp.api.ResetPassword(tok, newPass)
		}

		if ok {
			if gApp.currentView == ViewLogin || gApp.currentView == ViewRegister || gApp.currentView == ViewLicense {
				gApp.authenticated = true
				gApp.currentView = ViewDashboard
				gApp.activeDashTab = DashOverview
				setStatus("Session authenticated securely.", colSuccess)
				updateControlsVisibility()
			} else {
				msg := gApp.api.Response.Message
				if msg == "" {
					msg = "Operation succeeded."
				}
				setStatus(msg, colSuccess)
			}
		} else {
			msg := gApp.api.Response.Message
			if msg == "" {
				msg = "Operation failed."
			}
			setStatus(msg, colError)
		}
	}()
}

func wndProc(hWnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_CREATE:
		gApp.hWnd = hWnd

		gApp.hbrCanvas, _, _ = pCreateSolidBrush.Call(uintptr(colCanvasBlack))
		gApp.hbrSidebar, _, _ = pCreateSolidBrush.Call(uintptr(colSidebarBg))
		gApp.hbrHeader, _, _ = pCreateSolidBrush.Call(uintptr(colHeaderBg))
		gApp.hbrCard, _, _ = pCreateSolidBrush.Call(uintptr(colCardBg))
		gApp.hbrCardHeader, _, _ = pCreateSolidBrush.Call(uintptr(colCardHeaderBg))
		gApp.hbrInput, _, _ = pCreateSolidBrush.Call(uintptr(colInputBg))

		fName := syscall.StringToUTF16Ptr("Segoe UI")
		fMono := syscall.StringToUTF16Ptr("Consolas")

		neg := func(v int32) uintptr {
			return uintptr(uint32(v))
		}

		gApp.hFontBrand, _, _ = pCreateFontW.Call(neg(-15), 0, 0, 0, 700, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(fName)))
		gApp.hFontHero, _, _ = pCreateFontW.Call(neg(-22), 0, 0, 0, 700, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(fName)))
		gApp.hFontTitle, _, _ = pCreateFontW.Call(neg(-19), 0, 0, 0, 700, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(fName)))
		gApp.hFontSection, _, _ = pCreateFontW.Call(neg(-13), 0, 0, 0, 600, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(fName)))
		gApp.hFontBody, _, _ = pCreateFontW.Call(neg(-12), 0, 0, 0, 400, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(fName)))
		gApp.hFontInput, _, _ = pCreateFontW.Call(neg(-13), 0, 0, 0, 400, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(fName)))
		gApp.hFontLabel, _, _ = pCreateFontW.Call(neg(-10), 0, 0, 0, 700, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(fName)))
		gApp.hFontBadge, _, _ = pCreateFontW.Call(neg(-10), 0, 0, 0, 600, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(fName)))
		gApp.hFontStatus, _, _ = pCreateFontW.Call(neg(-11), 0, 0, 0, 400, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(fName)))
		gApp.hFontMono, _, _ = pCreateFontW.Call(neg(-12), 0, 0, 0, 400, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(fMono)))
		gApp.hFontKpi, _, _ = pCreateFontW.Call(neg(-24), 0, 0, 0, 700, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(fName)))

		var bmi BITMAPINFO
		bmi.BmiHeader.BiSize = uint32(unsafe.Sizeof(bmi.BmiHeader))
		bmi.BmiHeader.BiWidth = LOGO_WIDTH
		bmi.BmiHeader.BiHeight = -LOGO_HEIGHT
		bmi.BmiHeader.BiPlanes = 1
		bmi.BmiHeader.BiBitCount = 32
		bmi.BmiHeader.BiCompression = 0

		var pvBits unsafe.Pointer
		hdcScreen, _, _ := pGetDC.Call(0)
		hbmLogo, _, _ := pCreateDIBSection.Call(hdcScreen, uintptr(unsafe.Pointer(&bmi)), 0, uintptr(unsafe.Pointer(&pvBits)), 0, 0)
		pReleaseDC.Call(0, hdcScreen)
		if hbmLogo != 0 && pvBits != nil {
			dst := unsafe.Slice((*byte)(pvBits), len(gLogoBgra))
			copy(dst, gLogoBgra[:])
			gApp.hbmLogo = hbmLogo
		}

		hInst, _, _ := pGetModuleHandleW.Call(0)
		editClass := syscall.StringToUTF16Ptr("EDIT")

		createEdit := func(pwd bool) uintptr {
			flags := uintptr(WS_CHILD | ES_AUTOHSCROLL)
			if pwd {
				flags |= ES_PASSWORD
			}
			h, _, _ := pCreateWindowExW.Call(0, uintptr(unsafe.Pointer(editClass)), 0, flags, 0, 0, 0, 0, hWnd, 0, hInst, 0)
			pSendMessageW.Call(h, WM_SETFONT, gApp.hFontInput, 1)
			pSendMessageW.Call(h, EM_SETMARGINS, EC_LEFTMARGIN|EC_RIGHTMARGIN, uintptr((8<<16)|8))
			return h
		}

		gApp.hEditUser = createEdit(false)
		gApp.hEditPass = createEdit(true)
		gApp.hEditKey = createEdit(false)
		gApp.hEdit2FA = createEdit(false)
		gApp.hEditEmail = createEdit(false)
		gApp.hEditResetToken = createEdit(false)
		gApp.hEditNewPass = createEdit(true)

		editSubclassCallback = syscall.NewCallback(editSubclassProc)
		subclassEdit := func(h uintptr) {
			GWLP_WNDPROC := uintptr(^uintptr(3))
			prev, _, _ := pSetWindowLongPtrW.Call(h, GWLP_WNDPROC, editSubclassCallback)
			if origEditProc == 0 && prev != 0 {
				origEditProc = prev
			}
		}
		subclassEdit(gApp.hEditUser)
		subclassEdit(gApp.hEditPass)
		subclassEdit(gApp.hEditKey)
		subclassEdit(gApp.hEdit2FA)
		subclassEdit(gApp.hEditEmail)
		subclassEdit(gApp.hEditResetToken)
		subclassEdit(gApp.hEditNewPass)

		gApp.hEditChangeNewUser = createEdit(false)
		gApp.hEditChangePass = createEdit(true)

		gApp.hEditGlobalKey = createEdit(false)
		setEditText(gApp.hEditGlobalKey, "version")
		gApp.hEditVarKey = createEdit(false)
		gApp.hEditVarVal = createEdit(false)
		gApp.hEditChatRoom = createEdit(false)
		setEditText(gApp.hEditChatRoom, "general")
		gApp.hEditChatMsg = createEdit(false)
		gApp.hEdit2FACode = createEdit(false)
		gApp.hEditDisablePass = createEdit(true)
		gApp.hEditDisableCode = createEdit(false)
		gApp.hEditCheckHwid = createEdit(false)
		gApp.hEditVerifyHash = createEdit(false)
		gApp.hEditBanReason = createEdit(false)
		setEditText(gApp.hEditBanReason, "Security violation")
		gApp.hEditFileId = createEdit(false)
		setEditText(gApp.hEditFileId, "1")
		gApp.hEditHookId = createEdit(false)
		setEditText(gApp.hEditHookId, "1")
		gApp.hEditHookPayload = createEdit(false)
		setEditText(gApp.hEditHookPayload, "{\"alert\":\"critical\"}")

		updateControlsVisibility()
		return 0

	case WM_COMMAND:
		hi := uint16((wParam >> 16) & 0xFFFF)
		if hi == EN_SETFOCUS || hi == EN_KILLFOCUS {
			pInvalidateRect.Call(hWnd, 0, 0)
		}
		return 0

	case WM_CTLCOLOREDIT:
		hdc := wParam
		pSetTextColor.Call(hdc, uintptr(colTextHead))
		pSetBkColor.Call(hdc, uintptr(colInputBg))
		return gApp.hbrInput

	case WM_NCHITTEST:
		pt := POINT{X: int32(int16(lParam & 0xFFFF)), Y: int32(int16((lParam >> 16) & 0xFFFF))}
		pScreenToClient.Call(hWnd, uintptr(unsafe.Pointer(&pt)))
		if pt.Y >= 0 && pt.Y <= TITLEBAR_H && pt.X < (WIN_WIDTH-84) {
			return HTCAPTION
		}
		ret, _, _ := pDefWindowProcW.Call(hWnd, uintptr(msg), wParam, lParam)
		return ret

	case WM_MOUSEMOVE:
		x := int32(int16(lParam & 0xFFFF))
		y := int32(int16((lParam >> 16) & 0xFFFF))
		newHovered := 0

		if y <= TITLEBAR_H {
			if x >= WIN_WIDTH-76 && x < WIN_WIDTH-44 {
				newHovered = ID_BTN_MINIMIZE
			} else if x >= WIN_WIDTH-40 && x <= WIN_WIDTH-8 {
				newHovered = ID_BTN_CLOSE
			}
		} else if gApp.currentView != ViewDashboard {
			cardW := int32(860)
			cardX := (int32(WIN_WIDTH) - cardW) / 2
			cardY := int32(68)
			formX := cardX + 340
			formW := cardW - 340
			fieldX := formX + 36
			fieldW := formW - 72
			tabGap := int32(5)
			tabW := (fieldW - tabGap*6) / 7
			tabY := cardY + 26

			if y >= tabY && y <= tabY+30 {
				for i := 0; i < 7; i++ {
					tx := fieldX + int32(i)*(tabW+tabGap)
					if x >= tx && x <= tx+tabW {
						newHovered = ID_TAB_PRE_BASE + i
						break
					}
				}
			} else if y >= 456 && y <= 500 && x >= fieldX && x <= fieldX+fieldW {
				newHovered = ID_BTN_ACTION
			}
		} else {
			contentBottom := int32(WIN_HEIGHT) - STATUSBAR_H
			if x <= SIDEBAR_W {
				sPadY := int32(TITLEBAR_H + 38)
				sTabH := int32(38)
				sGap := int32(5)
				for i := 0; i < int(DashCount); i++ {
					ty := sPadY + int32(i)*(sTabH+sGap)
					if y >= ty && y <= ty+sTabH && x >= 12 && x <= SIDEBAR_W-12 {
						newHovered = ID_DASH_TAB_BASE + i
						break
					}
				}
				if y >= contentBottom-52 && y <= contentBottom-22 && x >= 24 && x <= 204 {
					newHovered = ID_BTN_LOGOUT
				}
			} else {
				cx := int32(256)
				cw := int32(WIN_WIDTH) - cx - 28
				switch gApp.activeDashTab {
				case DashOverview:
					if x >= WIN_WIDTH-190 && x <= WIN_WIDTH-28 && y >= 62 && y <= 96 {
						newHovered = ID_BTN_HEARTBEAT
					}
				case DashVariables:
					if x >= cx+295 && x <= cx+425 && y >= 174 && y <= 208 {
						newHovered = ID_BTN_GET_GLOBAL
					} else if x >= cx+485 && x <= cx+590 && y >= 309 && y <= 343 {
						newHovered = ID_BTN_GET_USER_VAR
					} else if x >= cx+600 && x <= cx+705 && y >= 309 && y <= 343 {
						newHovered = ID_BTN_SET_USER_VAR
					} else if x >= cx+cw-156 && x <= cx+cw-18 && y >= 405 && y <= 434 {
						newHovered = ID_BTN_GET_ALL_VARS
					}
				case DashChat:
					if x >= cx+235 && x <= cx+365 && y >= 164 && y <= 198 {
						newHovered = ID_BTN_REFRESH_CHAT
					} else if x >= cx+cw-132 && x <= cx+cw-20 && y >= 571 && y <= 605 {
						newHovered = ID_BTN_SEND_CHAT
					}
				case DashStats:
					if x >= cx+cw-156 && x <= cx+cw-18 && y >= 235 && y <= 264 {
						newHovered = ID_BTN_REFRESH_STATS
					}
				case DashSecurity:
					if x >= cx+550 && x <= cx+720 && y >= 169 && y <= 203 {
						newHovered = ID_BTN_CHANGE_USERNAME
					} else if x >= cx+20 && x <= cx+185 && y >= 278 && y <= 310 {
						newHovered = ID_BTN_SETUP_2FA
					} else if x >= cx+212 && x <= cx+336 && y >= 338 && y <= 372 {
						newHovered = ID_BTN_CONFIRM_2FA
					} else if x >= cx+640 && x <= cx+cw-20 && y >= 338 && y <= 372 {
						newHovered = ID_BTN_DISABLE_2FA
					} else if x >= cx+385 && x <= cx+535 && y >= 480 && y <= 514 {
						newHovered = ID_BTN_CHECK_BL
					} else if x >= cx+385 && x <= cx+535 && y >= 531 && y <= 565 {
						newHovered = ID_BTN_VERIFY_HASH
					} else if x >= cx+385 && x <= cx+535 && y >= 582 && y <= 612 {
						newHovered = ID_BTN_KILLSWITCH
					}
				case DashFilesHooks:
					if x >= cx+395 && x <= cx+555 && y >= 174 && y <= 208 {
						newHovered = ID_BTN_DOWNLOAD_FILE
					} else if x >= cx+395 && x <= cx+555 && y >= 434 && y <= 468 {
						newHovered = ID_BTN_TRIGGER_HOOK
					}
				}
			}
		}

		if newHovered != gApp.hoveredBtn {
			gApp.hoveredBtn = newHovered
			pInvalidateRect.Call(hWnd, 0, 0)
		}

		var tme TRACKMOUSEEVENT
		tme.CbSize = uint32(unsafe.Sizeof(tme))
		tme.DwFlags = 0x00000002
		tme.HwndTrack = hWnd
		pTrackMouseEvent.Call(uintptr(unsafe.Pointer(&tme)))
		return 0

	case WM_MOUSELEAVE:
		if gApp.hoveredBtn != 0 {
			gApp.hoveredBtn = 0
			pInvalidateRect.Call(hWnd, 0, 0)
		}
		return 0

	case WM_LBUTTONDOWN:
		x := int32(int16(lParam & 0xFFFF))
		y := int32(int16((lParam >> 16) & 0xFFFF))

		if y <= TITLEBAR_H {
			if x >= WIN_WIDTH-76 && x < WIN_WIDTH-44 {
				pShowWindow.Call(hWnd, SW_MINIMIZE)
				return 0
			} else if x >= WIN_WIDTH-40 && x <= WIN_WIDTH-8 {
				pPostQuitMessage.Call(0)
				return 0
			}
		}

		if gApp.currentView != ViewDashboard {
			cardW := int32(860)
			cardX := (int32(WIN_WIDTH) - cardW) / 2
			cardY := int32(68)
			formX := cardX + 340
			formW := cardW - 340
			fieldX := formX + 36
			fieldW := formW - 72
			tabGap := int32(5)
			tabW := (fieldW - tabGap*6) / 7
			tabY := cardY + 26

			if y >= tabY && y <= tabY+30 {
				for i := 0; i < 7; i++ {
					tx := fieldX + int32(i)*(tabW+tabGap)
					if x >= tx && x <= tx+tabW {
						gApp.currentView = ViewMode(i)
						updateControlsVisibility()
						setStatus("Ready.", colTextMuted)
						return 0
					}
				}
			}

			if y >= 456 && y <= 500 && x >= fieldX && x <= fieldX+fieldW {
				handleActionClick()
				return 0
			}
		} else {
			contentBottom := int32(WIN_HEIGHT) - STATUSBAR_H
			if x <= SIDEBAR_W {
				sPadY := int32(TITLEBAR_H + 38)
				sTabH := int32(38)
				sGap := int32(5)
				for i := 0; i < int(DashCount); i++ {
					ty := sPadY + int32(i)*(sTabH+sGap)
					if y >= ty && y <= ty+sTabH && x >= 12 && x <= SIDEBAR_W-12 {
						gApp.activeDashTab = DashboardTab(i)
						updateControlsVisibility()
						setStatus("Module loaded.", colTextMuted)
						return 0
					}
				}
				if y >= contentBottom-52 && y <= contentBottom-22 && x >= 24 && x <= 204 {
					gApp.api.Logout()
					gApp.authenticated = false
					gApp.currentView = ViewLogin
					updateControlsVisibility()
					setStatus("Signed out cleanly.", colTextMuted)
					return 0
				}
			} else {
				cx := int32(256)
				cw := int32(WIN_WIDTH) - cx - 28
				switch gApp.activeDashTab {
				case DashOverview:
					if x >= WIN_WIDTH-190 && x <= WIN_WIDTH-28 && y >= 62 && y <= 96 {
						ok := gApp.api.Heartbeat()
						if ok {
							setStatus("Heartbeat acknowledged by cloud.", colSuccess)
						} else {
							setStatus("Heartbeat failed.", colError)
						}
						return 0
					}
				case DashVariables:
					if x >= cx+295 && x <= cx+425 && y >= 174 && y <= 208 {
						key := strings.TrimSpace(getEditText(gApp.hEditGlobalKey))
						gApp.lastGlobalVal = gApp.api.GetGlobalVar(key)
						if gApp.lastGlobalVal != "" {
							setStatus("Global variable fetched.", colSuccess)
						} else {
							setStatus("Failed to fetch global variable.", colError)
						}
						pInvalidateRect.Call(hWnd, 0, 0)
						return 0
					} else if x >= cx+485 && x <= cx+590 && y >= 309 && y <= 343 {
						key := strings.TrimSpace(getEditText(gApp.hEditVarKey))
						gApp.lastUserVal = gApp.api.GetUserVar(key)
						if gApp.lastUserVal != "" {
							setStatus("User variable retrieved.", colSuccess)
						} else {
							setStatus("User variable not found.", colError)
						}
						pInvalidateRect.Call(hWnd, 0, 0)
						return 0
					} else if x >= cx+600 && x <= cx+705 && y >= 309 && y <= 343 {
						key := strings.TrimSpace(getEditText(gApp.hEditVarKey))
						val := strings.TrimSpace(getEditText(gApp.hEditVarVal))
						ok := gApp.api.SetUserVar(key, val)
						if ok {
							setStatus("User variable saved to cloud.", colSuccess)
						} else {
							setStatus("Failed to save user variable.", colError)
						}
						return 0
					} else if x >= cx+cw-156 && x <= cx+cw-18 && y >= 405 && y <= 434 {
						ok := gApp.api.GetUserVariables()
						gApp.userVarsList = nil
						if ok && gApp.api.Response.Payload != nil {
							if vars, ok := gApp.api.Response.Payload["variables"].([]interface{}); ok {
								for _, v := range vars {
									if vm, ok := v.(map[string]interface{}); ok {
										k, _ := vm["key"].(string)
										val, _ := vm["value"].(string)
										gApp.userVarsList = append(gApp.userVarsList, map[string]string{"key": k, "value": val})
									}
								}
							}
						}
						if ok {
							setStatus("User variables refreshed.", colSuccess)
						} else {
							setStatus("Failed to load variables.", colError)
						}
						pInvalidateRect.Call(hWnd, 0, 0)
						return 0
					}
				case DashChat:
					if x >= cx+235 && x <= cx+365 && y >= 164 && y <= 198 {
						room := strings.TrimSpace(getEditText(gApp.hEditChatRoom))
						gApp.chatMsgs = gApp.api.GetChatMessages(room)
						setStatus("Chat channel synchronized.", colSuccess)
						pInvalidateRect.Call(hWnd, 0, 0)
						return 0
					} else if x >= cx+cw-132 && x <= cx+cw-20 && y >= 571 && y <= 605 {
						msgText := strings.TrimSpace(getEditText(gApp.hEditChatMsg))
						room := strings.TrimSpace(getEditText(gApp.hEditChatRoom))
						if msgText != "" {
							if gApp.api.SendChatMessage(msgText, room) {
								setEditText(gApp.hEditChatMsg, "")
								gApp.chatMsgs = gApp.api.GetChatMessages(room)
								setStatus("Message sent.", colSuccess)
							} else {
								setStatus("Failed to send chat message.", colError)
							}
							pInvalidateRect.Call(hWnd, 0, 0)
						}
						return 0
					}
				case DashStats:
					if x >= cx+cw-156 && x <= cx+cw-18 && y >= 235 && y <= 264 {
						gApp.appStats = gApp.api.GetStats()
						gApp.onlineUsers = gApp.api.GetOnlineUsers()
						setStatus("Metrics and active sessions refreshed.", colSuccess)
						pInvalidateRect.Call(hWnd, 0, 0)
						return 0
					}
				case DashSecurity:
					if x >= cx+550 && x <= cx+720 && y >= 169 && y <= 203 {
						newU := strings.TrimSpace(getEditText(gApp.hEditChangeNewUser))
						pass := strings.TrimSpace(getEditText(gApp.hEditChangePass))
						if newU == "" || pass == "" {
							setStatus("Please enter new username and current password.", colError)
							return 0
						}
						setStatus("Updating username with cloud...", colBlueGlow)
						pUpdateWindow.Call(hWnd)
						if gApp.api.ChangeUsername(newU, pass) {
							pMsg, _ := syscall.UTF16PtrFromString("Username updated successfully! For security, your active session has been terminated. Please sign in again with your new username.")
							pTitle, _ := syscall.UTF16PtrFromString("Username Changed")
							pMessageBoxW.Call(hWnd, uintptr(unsafe.Pointer(pMsg)), uintptr(unsafe.Pointer(pTitle)), 0x40)

							gApp.authenticated = false
							gApp.activeDashTab = DashOverview
							gApp.currentView = ViewLogin
							setEditText(gApp.hEditUser, newU)
							setEditText(gApp.hEditPass, "")
							setEditText(gApp.hEditChangeNewUser, "")
							setEditText(gApp.hEditChangePass, "")
							setStatus("Session closed. Please sign in with your new username.", colSuccess)
							updateControlsVisibility()
						} else {
							msg := gApp.api.Response.Message
							if msg == "" {
								msg = "Failed to change username."
							}
							setStatus(msg, colError)
						}
						return 0
					} else if x >= cx+20 && x <= cx+185 && y >= 278 && y <= 310 {
						gApp.tfaSetup = gApp.api.Setup2FA()
						if gApp.tfaSetup != nil {
							setStatus("2FA Setup secret generated.", colSuccess)
						} else {
							setStatus("Failed to setup 2FA.", colError)
						}
						pInvalidateRect.Call(hWnd, 0, 0)
						return 0
					} else if x >= cx+212 && x <= cx+336 && y >= 338 && y <= 372 {
						code := strings.TrimSpace(getEditText(gApp.hEdit2FACode))
						conf := gApp.api.Confirm2FA(code)
						if conf != nil {
							setStatus("2FA confirmed and activated.", colSuccess)
						} else {
							setStatus("Invalid 2FA code.", colError)
						}
						return 0
					} else if x >= cx+640 && x <= cx+cw-20 && y >= 338 && y <= 372 {
						pwd := strings.TrimSpace(getEditText(gApp.hEditDisablePass))
						code := strings.TrimSpace(getEditText(gApp.hEditDisableCode))
						ok := gApp.api.Disable2FA(pwd, code)
						if ok {
							setStatus("2FA disabled successfully.", colSuccess)
						} else {
							setStatus("Failed to disable 2FA.", colError)
						}
						return 0
					} else if x >= cx+385 && x <= cx+535 && y >= 480 && y <= 514 {
						hwid := strings.TrimSpace(getEditText(gApp.hEditCheckHwid))
						isBl := gApp.api.CheckBlacklist(hwid)
						if isBl {
							setStatus("HWID is BLACKLISTED.", colError)
						} else {
							setStatus("HWID is clean (not blacklisted).", colSuccess)
						}
						return 0
					} else if x >= cx+385 && x <= cx+535 && y >= 531 && y <= 565 {
						hash := strings.TrimSpace(getEditText(gApp.hEditVerifyHash))
						ok := gApp.api.VerifyFileHash(hash)
						if ok {
							setStatus("Hash verified with server signature.", colSuccess)
						} else {
							setStatus("Hash verification failed / mismatched.", colError)
						}
						return 0
					} else if x >= cx+385 && x <= cx+535 && y >= 582 && y <= 612 {
						reason := strings.TrimSpace(getEditText(gApp.hEditBanReason))
						ok := gApp.api.Ban(reason)
						if ok {
							setStatus("Account banned / killswitch executed.", colSuccess)
						} else {
							setStatus("Failed to execute ban.", colError)
						}
						return 0
					}
				case DashFilesHooks:
					if x >= cx+395 && x <= cx+555 && y >= 174 && y <= 208 {
						fileId := strings.TrimSpace(getEditText(gApp.hEditFileId))
						data := gApp.api.DownloadFile(fileId)
						if data != nil {
							gApp.downloadStatusOk = true
							gApp.downloadStatus = fmt.Sprintf("Downloaded %d bytes -> downloaded.bin", len(data))
							setStatus("File downloaded successfully.", colSuccess)
						} else {
							gApp.downloadStatusOk = false
							gApp.downloadStatus = "Download failed for file ID: " + fileId
							setStatus("File download failed.", colError)
						}
						pInvalidateRect.Call(hWnd, 0, 0)
						return 0
					} else if x >= cx+395 && x <= cx+555 && y >= 434 && y <= 468 {
						hookId := strings.TrimSpace(getEditText(gApp.hEditHookId))
						payload := strings.TrimSpace(getEditText(gApp.hEditHookPayload))
						customMap := map[string]interface{}{"data": payload}
						res := gApp.api.TriggerWebhook(hookId, customMap)
						if res != nil {
							gApp.webhookResponseOk = true
							gApp.webhookResponse = fmt.Sprintf("HTTP %d  |  %s", res.Status, res.Data)
							setStatus("Webhook executed successfully.", colSuccess)
						} else {
							gApp.webhookResponseOk = false
							gApp.webhookResponse = "Webhook dispatch failed: " + gApp.api.Response.Message
							setStatus("Webhook failed.", colError)
						}
						pInvalidateRect.Call(hWnd, 0, 0)
						return 0
					}
				}
			}
		}
		return 0

	case WM_PAINT:
		onPaint(hWnd)
		return 0

	case WM_DESTROY:
		pDeleteObject.Call(gApp.hbrCanvas)
		pDeleteObject.Call(gApp.hbrSidebar)
		pDeleteObject.Call(gApp.hbrHeader)
		pDeleteObject.Call(gApp.hbrCard)
		pDeleteObject.Call(gApp.hbrCardHeader)
		pDeleteObject.Call(gApp.hbrInput)
		pDeleteObject.Call(gApp.hFontBrand)
		pDeleteObject.Call(gApp.hFontHero)
		pDeleteObject.Call(gApp.hFontTitle)
		pDeleteObject.Call(gApp.hFontSection)
		pDeleteObject.Call(gApp.hFontBody)
		pDeleteObject.Call(gApp.hFontInput)
		pDeleteObject.Call(gApp.hFontLabel)
		pDeleteObject.Call(gApp.hFontBadge)
		pDeleteObject.Call(gApp.hFontStatus)
		pDeleteObject.Call(gApp.hFontMono)
		pDeleteObject.Call(gApp.hFontKpi)
		if gApp.hbmLogo != 0 {
			pDeleteObject.Call(gApp.hbmLogo)
		}
		pPostQuitMessage.Call(0)
		return 0
	}

	ret, _, _ := pDefWindowProcW.Call(hWnd, uintptr(msg), wParam, lParam)
	return ret
}

func main() {
	appName := os.Getenv("NIGHTAUTH_APP_NAME")
	if appName == "" {
		appName = AppName
	}
	apiKey := os.Getenv("NIGHTAUTH_API_KEY")
	if apiKey == "" {
		apiKey = AppApiKey
	}
	appVer := os.Getenv("NIGHTAUTH_APP_VERSION")
	if appVer == "" {
		appVer = AppVersion
	}
	appToken := os.Getenv("NIGHTAUTH_APP_TOKEN")
	if appToken == "" {
		appToken = AppToken
	}
	apiUrl := os.Getenv("NIGHTAUTH_API_URL")
	if apiUrl == "" {
		apiUrl = nightauth.DefaultAPIURL
	}
	allowLoopback := os.Getenv("NIGHTAUTH_ALLOW_INSECURE_LOOPBACK") == "1"
	gApp.api = nightauth.NewApiWithOptions(appName, apiKey, appVer, appToken, apiUrl, allowLoopback)

	hInst, _, _ := pGetModuleHandleW.Call(0)
	className := syscall.StringToUTF16Ptr("NightAuthGoDesktopUI")

	hCur, _, _ := pLoadCursorW.Call(0, 32512)
	var wc WNDCLASSEXW
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = syscall.NewCallback(wndProc)
	wc.HInstance = hInst
	wc.HCursor = hCur
	wc.HbrBackground = 0
	wc.LpszClassName = className

	pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	screenW, _, _ := pGetSystemMetrics.Call(SM_CXSCREEN)
	screenH, _, _ := pGetSystemMetrics.Call(SM_CYSCREEN)
	posX := (int32(screenW) - WIN_WIDTH) / 2
	posY := (int32(screenH) - WIN_HEIGHT) / 2

	hWnd, _, _ := pCreateWindowExW.Call(
		WS_EX_APPWINDOW,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("NightAuth Enterprise - Desktop"))),
		WS_POPUP|WS_VISIBLE,
		uintptr(posX), uintptr(posY), WIN_WIDTH, WIN_HEIGHT,
		0, 0, hInst, 0,
	)

	if hWnd == 0 {
		return
	}

	pShowWindow.Call(hWnd, SW_SHOW)
	pUpdateWindow.Call(hWnd)
	pSetForegroundWindow.Call(hWnd)

	setStatus("Establishing secure session with cloud...", colBlueGlow)
	go func() {
		gApp.api.Init()
		if gApp.api.Response.Success {
			setStatus("Connected securely. Ready to authenticate.", colSuccess)
		} else {
			msg := gApp.api.Response.Message
			if msg == "" {
				msg = "Handshake failed."
			}
			setStatus(msg, colError)
		}
	}()

	var msg MSG
	for {
		ret, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 || int32(ret) == -1 {
			break
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}
