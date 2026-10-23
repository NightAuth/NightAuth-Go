package nightauth

import (
	"net/http"
	"testing"
)

func fakeInitializedApi() *Api {
	api := NewApi("test-app", "test-key", "1.0", "", "https://localhost:3000")
	// Simulate a completed handshake without hitting a real server.
	api.sessionKey = "session-under-test"
	api.sharedKey = make([]byte, 32)
	for i := range api.sharedKey {
		api.sharedKey[i] = 7
	}
	api.isInitialized = true
	return api
}

func TestLogoutClearsLocalSessionStateEvenWithoutNetwork(t *testing.T) {
	api := fakeInitializedApi()

	// The server-side /logout call will fail (no server is running
	// in this test), but Logout() must still wipe local state.
	if api.Logout() {
		t.Fatal("Logout() must report unconfirmed server revocation")
	}

	if api.sessionKey != "" {
		t.Error("sessionKey must be cleared after logout")
	}
	if api.sharedKey != nil {
		t.Error("sharedKey must be cleared after logout")
	}
	if api.isInitialized {
		t.Error("isInitialized must be false after logout")
	}
	if api.User.Username != "" {
		t.Error("user profile must be cleared after logout")
	}
}

func TestUsingSessionAfterLogoutRequiresReinit(t *testing.T) {
	var paths []string
	api := newTestApi(t, func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"success":false,"error":"rate_limited"}`))
	})
	api.sessionKey = "session-under-test"
	api.isInitialized = true
	if api.Logout() {
		t.Fatal("Logout() must report unconfirmed server revocation")
	}
	paths = nil
	_ = api.sendEncrypted("/check", map[string]interface{}{})
	if len(paths) == 0 || paths[0] != "/init" {
		t.Fatalf("first request after logout must be a fresh handshake, got %v", paths)
	}
	if api.sessionKey == "session-under-test" {
		t.Fatal("the old session key must not survive a failed re-init")
	}
}
