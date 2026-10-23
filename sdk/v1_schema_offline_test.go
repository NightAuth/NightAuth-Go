package nightauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func captureV1Body(t *testing.T, call func(a *Api)) map[string]interface{} {
	t.Helper()
	var captured map[string]interface{}
	var api *Api
	api = newTestApi(t, func(w http.ResponseWriter, r *http.Request) {
		var env map[string]interface{}
		dec := json.NewDecoder(r.Body)
		dec.UseNumber()
		if err := dec.Decode(&env); err != nil {
			t.Errorf("decode envelope: %v", err)
			return
		}
		rid, _ := env["rid"].(string)
		ts := fmt.Sprint(env["ts"])
		ct, _ := env["ct"].(string)
		n, _ := env["n"].(string)
		reqAad := fmt.Sprintf("POST:%s:%s:%s:%s", r.URL.Path, api.sessionKey, rid, ts)
		plain, err := aesGcmDecrypt(api.sharedKey, ct, n, reqAad)
		if err != nil {
			t.Errorf("decrypt request: %v", err)
			return
		}
		if err := json.Unmarshal([]byte(plain), &captured); err != nil {
			t.Errorf("request plaintext is not JSON: %v", err)
			return
		}
		resAad := fmt.Sprintf("v1:res:%s:%s:%s", api.sessionKey, rid, ts)
		resCt, resN, err := aesGcmEncrypt(api.sharedKey, `{"success":true}`, resAad)
		if err != nil {
			t.Errorf("encrypt response: %v", err)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id": api.sessionKey, "ct": resCt, "n": resN, "rid": rid, "ts": env["ts"],
		})
	})
	api.Version = "2.5"
	call(api)
	if captured == nil {
		t.Fatal("no request reached the server")
	}
	return captured
}

func TestHeartbeatSendsVersion(t *testing.T) {
	body := captureV1Body(t, func(a *Api) { a.Heartbeat() })
	if got, _ := body["version"].(string); got != "2.5" {
		t.Fatalf("Heartbeat body must carry version=2.5, got keys %v", keysOf(body))
	}
}

func TestCheckSendsVersion(t *testing.T) {
	body := captureV1Body(t, func(a *Api) { a.Check() })
	if got, _ := body["version"].(string); got != "2.5" {
		t.Fatalf("Check body must carry version=2.5, got keys %v", keysOf(body))
	}
}

func TestRegisterOmitsEmptyEmail(t *testing.T) {
	body := captureV1Body(t, func(a *Api) { a.Register("user1", "pw", "KEY", "") })
	if _, present := body["email"]; present {
		t.Fatalf("Register must not send an empty email field, got keys %v", keysOf(body))
	}
}

func TestRegisterSendsNonEmptyEmail(t *testing.T) {
	body := captureV1Body(t, func(a *Api) { a.Register("user1", "pw", "KEY", "a@example.com") })
	if got, _ := body["email"].(string); got != "a@example.com" {
		t.Fatalf("Register must send the non-empty email, got %q", got)
	}
}

func keysOf(m map[string]interface{}) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ",")
}
