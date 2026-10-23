package nightauth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApi(t *testing.T, handler http.HandlerFunc) *Api {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	return &Api{
		Name: "test", ApiKey: "na_test", Version: "1.0", ApiUrl: server.URL,
		httpClient:            server.Client(),
		allowInsecureLoopback: true,
		sessionKey:            "test-session",
		sharedKey:             key,
		isInitialized:         true,
	}
}

func TestSendEncryptedRejectsMismatchedResponseBinding(t *testing.T) {
	var api *Api
	api = newTestApi(t, func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		aad := fmt.Sprintf("v1:res:%s:%s:%d", api.sessionKey, "not-the-real-rid", req["ts"])
		ct, n, err := aesGcmEncrypt(api.sharedKey, `{"success":true,"stolen":"data"}`, aad)
		if err != nil {
			t.Fatal(err)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id": req["id"], "ct": ct, "n": n, "rid": "not-the-real-rid", "ts": req["ts"],
		})
	})

	resp := api.sendEncrypted("/test", map[string]interface{}{})
	if success, _ := resp["success"].(bool); success {
		t.Fatalf("expected rejection of a response bound to a different rid, got success: %+v", resp)
	}
	if resp["error"] != "response_binding_mismatch" {
		t.Fatalf("expected response_binding_mismatch, got %+v", resp)
	}
}

func TestSendEncryptedRejectsPlaintextSuccess(t *testing.T) {
	api := newTestApi(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "not actually encrypted"})
	})

	resp := api.sendEncrypted("/test", map[string]interface{}{})
	if success, _ := resp["success"].(bool); success {
		t.Fatalf("an unencrypted 2xx body must never be treated as an authenticated success, got: %+v", resp)
	}
}

func TestSendEncryptedAcceptsCorrectlyBoundResponse(t *testing.T) {
	var api *Api
	api = newTestApi(t, func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		reqTs := int64(req["ts"].(float64))
		aad := fmt.Sprintf("v1:res:%s:%s:%d", api.sessionKey, req["rid"], reqTs)
		ct, n, err := aesGcmEncrypt(api.sharedKey, `{"success":true}`, aad)
		if err != nil {
			t.Fatal(err)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id": req["id"], "ct": ct, "n": n, "rid": req["rid"], "ts": req["ts"],
		})
	})

	resp := api.sendEncrypted("/test", map[string]interface{}{})
	if success, _ := resp["success"].(bool); !success {
		t.Fatalf("expected success for a response correctly bound to the request's rid/ts, got: %+v", resp)
	}
}

func TestTransportRequiresExplicitLoopbackOptIn(t *testing.T) {
	for _, raw := range []string{"http://localhost:3000", "http://example.com", "http://localhost.evil"} {
		if validateApiUrl(raw) == nil {
			t.Fatalf("accepted insecure URL %s", raw)
		}
	}
	if validateApiUrl("http://127.0.0.1:3000", true) != nil {
		t.Fatal("explicit loopback rejected")
	}
	if validateApiUrl("http://localhost.evil", true) == nil {
		t.Fatal("lookalike accepted")
	}
}

func TestPostRawRejectsOversizedResponse(t *testing.T) {
	api := newTestApi(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(make([]byte, 5*1024*1024+1))
	})
	_, body, _, err := api.postRaw("/test", []byte(`{}`))
	if err == nil || body != nil {
		t.Fatalf("expected oversized response rejection, body bytes=%d err=%v", len(body), err)
	}
}
