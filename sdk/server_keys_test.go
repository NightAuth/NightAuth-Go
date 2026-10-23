package nightauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestServerKeysPerClient(t *testing.T) {
	other := strings.Repeat("1", 64)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true, "k2": other, "vk": other, "ts": time.Now().UnixMilli(), "sig": "AAAA",
		})
	}))
	defer server.Close()
	api := NewApiWithOptions("test-app", "test-key", "1.0", "", strings.Replace(server.URL, "127.0.0.1", "localhost", 1), true)

	if err := api.SetServerKeys("abc", other); err == nil {
		t.Fatal("a short signing key must be refused")
	}
	if err := api.SetServerKeys(other, strings.Repeat("1", 62)+"zz"); err == nil {
		t.Fatal("a non-hex sealed-box key must be refused")
	}
	if api.pinnedSigningKey() != ServerSigningPublicKey || api.pinnedSealKey() != ServerSealedBoxPublicKey {
		t.Fatal("a refused change must keep the built-in keys")
	}

	api.Init()
	if api.isInitialized || !strings.Contains(api.Response.Message, "Untrusted") {
		t.Fatalf("with the built-in keys a response for another key must be untrusted, got %q", api.Response.Message)
	}
	if api.Response.Error != "untrusted_response" {
		t.Fatalf("an untrusted handshake must carry untrusted_response, got %q", api.Response.Error)
	}

	if err := api.SetServerKeys(other, other); err != nil {
		t.Fatalf("two 64-hex keys must be accepted: %v", err)
	}
	api.Init()
	if api.isInitialized || !strings.Contains(api.Response.Message, "signature") || api.Response.Error != "untrusted_response" {
		t.Fatalf("after the change the key check passes and the forged signature fails, got %q %q", api.Response.Message, api.Response.Error)
	}
}

func TestSetRequestTimeoutBoundsACall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	api := NewApiWithOptions("test-app", "test-key", "1.0", "", strings.Replace(server.URL, "127.0.0.1", "localhost", 1), true)
	if err := api.SetRequestTimeout(0); err == nil {
		t.Fatal("a non-positive timeout must be refused")
	}
	if err := api.SetRequestTimeout(300 * time.Millisecond); err != nil {
		t.Fatal(err)
	}
	api.sessionKey = "session-under-test"
	api.sharedKey = make([]byte, 32)
	api.isInitialized = true
	started := time.Now()
	if api.Heartbeat() {
		t.Fatal("a call past its deadline must fail")
	}
	if api.Response.Error != "timeout" {
		t.Fatalf("expected timeout, got %+v", api.Response)
	}
	if time.Since(started) > 1500*time.Millisecond {
		t.Fatalf("the client timeout did not end the call (%v)", time.Since(started))
	}
}
