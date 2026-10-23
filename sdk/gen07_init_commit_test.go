package nightauth

import (
	"net/http"
	"testing"
)

func TestFailedInitLeavesNoHalfBuiltSession(t *testing.T) {
	cases := map[string]func(w http.ResponseWriter, r *http.Request){
		"http 500": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"success":false,"error":"internal_error"}`))
		},
		"unparseable": func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("<html>proxy</html>")) },
		"success=false": func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"success":false,"error":"app_paused"}`))
		},
		"bad signature": func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"success":true,"k2":"` + string(make([]byte, 0)) + `0000000000000000000000000000000000000000000000000000000000000000","ts":1,"sig":"AAAA"}`))
		},
	}
	for name, handler := range cases {
		t.Run(name, func(t *testing.T) {
			api := newTestApi(t, handler)
			previousKey := api.sharedKey
			api.Init()
			if api.isInitialized {
				t.Fatal("a failed handshake must not leave the client initialized")
			}
			if api.sharedKey != nil || api.sessionKey != "" {
				t.Fatalf("a failed handshake must not keep a session (sessionKey=%q, sharedKey kept=%v)", api.sessionKey, api.sharedKey != nil)
			}
			_ = previousKey
		})
	}
}
