package nightauth

import (
	"testing"
	"time"
)

func TestSecurityProtectionsHardwareId(t *testing.T) {
	hwid1 := SecurityProtections{}.GetHardwareId()
	if hwid1 == "" || len(hwid1) != 64 {
		t.Fatalf("expected 64-hex hardware ID, got %q", hwid1)
	}

	hwid2 := GetHardwareId()
	if hwid1 != hwid2 {
		t.Fatalf("cached HWID mismatch: %q vs %q", hwid1, hwid2)
	}
}

func TestTwoFactorConfirmDataFields(t *testing.T) {
	data := TwoFactorConfirmData{
		Success:     true,
		BackupCodes: []string{"abc-123", "def-456"},
	}
	if !data.Success || len(data.BackupCodes) != 2 {
		t.Fatalf("unexpected 2FA confirm data: %+v", data)
	}
}

func TestWebhookResultFields(t *testing.T) {
	res := WebhookResult{
		Success:    true,
		Status:     200,
		StatusCode: 200,
		Data:       "ok",
		Truncated:  false,
	}
	if !res.Success || res.Status != 200 || res.StatusCode != 200 {
		t.Fatalf("unexpected webhook result: %+v", res)
	}
}

func TestRealtimeEventListenerStartStop(t *testing.T) {
	api := NewApiWithOptions("test-app", "test-key", "1.0", "", "http://localhost:13100", true)
	// Starting listener when not connected should fail gracefully without panic
	api.StartRealtimeEventListener()
	time.Sleep(50 * time.Millisecond)
	api.StopRealtimeEventListener()

	// Calling Stop again should be safe and idempotent
	api.StopRealtimeEventListener()

	// Close should be safe and idempotent
	api.Close()
	api.Close()
}

func TestEventDispatching(t *testing.T) {
	api := NewApiWithOptions("test-app", "test-key", "1.0", "", "http://localhost:13100", true)

	var sessionInvalidReason string
	api.OnSessionInvalidated = func(reason string) {
		sessionInvalidReason = reason
	}

	var chatMsg ChatMessage
	var chatRoom string
	api.OnChatMessageReceived = func(msg ChatMessage, room string) {
		chatMsg = msg
		chatRoom = room
	}

	api.dispatchEvent(`{"type":"chat-message","room":"general","id":"msg1","author":"alice","content":"hello world"}`)
	if chatRoom != "general" || chatMsg.Author != "alice" || chatMsg.Message != "hello world" {
		t.Fatalf("unexpected chat message dispatch: %+v in %q", chatMsg, chatRoom)
	}

	api.dispatchEvent(`{"type":"session-invalid"}`)
	if sessionInvalidReason == "" {
		t.Fatalf("expected session invalidation reason to be set")
	}
}

func TestSignedUpdateHelperInterface(t *testing.T) {
	// Verify that helper methods exist and reject invalid envelopes cleanly
	_, err := SignedUpdate.Verify("{}", "invalid-key", "test-app", 0)
	if err == nil {
		t.Fatal("expected error on invalid envelope")
	}

	_, err = SignedUpdateInstaller.Download("invalid-url", "key", "app", "dest", 1*time.Second)
	if err == nil {
		t.Fatal("expected error on invalid URL")
	}
}

func TestPopulateUserDataNumericTimestampsAndTwoFactor(t *testing.T) {
	api := &Api{}
	payload := map[string]interface{}{
		"info": map[string]interface{}{
			"username":         "alice",
			"email":            "alice@example.com",
			"hwid":             "HWID-12345",
			"ip":               "1.2.3.4",
			"status":           "active",
			"createdate":       float64(1768473045),
			"lastlogin":        float64(1787731200),
			"twoFactorEnabled": true,
			"subscriptions": []interface{}{
				map[string]interface{}{
					"subscription": "VIP",
					"key":          "KEY-999",
					"expiry":       float64(1800000000),
					"timeleft":     float64(86400),
				},
			},
		},
	}

	api.populateUserData(payload)

	if api.User.Username != "alice" {
		t.Fatalf("expected username alice, got %q", api.User.Username)
	}
	if api.User.CreatedAt != "1768473045" {
		t.Fatalf("expected CreatedAt 1768473045, got %q", api.User.CreatedAt)
	}
	if api.User.LastLoginAt != "1787731200" {
		t.Fatalf("expected LastLoginAt 1787731200, got %q", api.User.LastLoginAt)
	}
	if !api.User.TwoFactorEnabled {
		t.Fatalf("expected TwoFactorEnabled true, got false")
	}
	if len(api.User.Subscriptions) != 1 {
		t.Fatalf("expected 1 subscription, got %d", len(api.User.Subscriptions))
	}
	if api.User.Subscriptions[0].Expiry != "1800000000" {
		t.Fatalf("expected subscription Expiry 1800000000, got %q", api.User.Subscriptions[0].Expiry)
	}
	if api.User.Subscriptions[0].TimeLeft != 86400 {
		t.Fatalf("expected TimeLeft 86400, got %d", api.User.Subscriptions[0].TimeLeft)
	}
}
