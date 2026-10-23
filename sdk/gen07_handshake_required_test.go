package nightauth

import (
	"net/http"
	"sync/atomic"
	"testing"
)

func TestHandshakeRequiredTriggersOneReinitAndNoRetryWhenReinitFails(t *testing.T) {
	var calls, inits int32
	api := newTestApi(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/init" {
			atomic.AddInt32(&inits, 1)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"success":false,"error":"internal_error"}`))
			return
		}
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusPreconditionRequired)
		_, _ = w.Write([]byte(`{"success":false,"error":"handshake_required","message":"Initialize app first"}`))
	})

	resp := api.sendEncrypted("/check", map[string]interface{}{})
	if success, _ := resp["success"].(bool); success {
		t.Fatalf("a 428 must never be a success: %+v", resp)
	}
	if resp["error"] != "handshake_required" {
		t.Fatalf("the first answer must be returned when the re-init fails, got %+v", resp)
	}
	if got := atomic.LoadInt32(&inits); got != 1 {
		t.Fatalf("expected exactly one re-init, got %d", got)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("no retry is allowed after a failed re-init, got %d calls", got)
	}
	if api.isInitialized || api.sharedKey != nil || api.sessionKey != "" {
		t.Fatal("a failed re-init must not leave a mixed session state")
	}
}

func TestPlainRefusalWithoutHandshakeIsNotRetried(t *testing.T) {
	var calls, inits int32
	api := newTestApi(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/init" {
			atomic.AddInt32(&inits, 1)
		} else {
			atomic.AddInt32(&calls, 1)
		}
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		_, _ = w.Write([]byte(`{"success":false,"error":"payload_too_large"}`))
	})
	resp := api.sendEncrypted("/check", map[string]interface{}{})
	if resp["error"] != "payload_too_large" || resp["http_status"] != 413 {
		t.Fatalf("the server refusal must be surfaced with its status, got %+v", resp)
	}
	if atomic.LoadInt32(&inits) != 0 || atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("a 413 must not re-init or retry (inits=%d calls=%d)", inits, calls)
	}
	if !api.isInitialized {
		t.Fatal("a 413 must keep the session")
	}
}
