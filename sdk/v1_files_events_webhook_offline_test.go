package nightauth

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// ftEncryptedHandler receives the decrypted request body and returns the plaintext reply.
type ftEncryptedHandler func(path string, body map[string]interface{}) map[string]interface{}

func ftNewApi(t *testing.T, encrypted ftEncryptedHandler, plain http.HandlerFunc) *Api {
	t.Helper()
	var api *Api
	api = newTestApi(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if plain == nil {
				http.NotFound(w, r)
				return
			}
			plain(w, r)
			return
		}
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
		plainReq, err := aesGcmDecrypt(api.sharedKey, ct, n, reqAad)
		if err != nil {
			t.Errorf("decrypt request %s: %v", r.URL.Path, err)
			return
		}
		var body map[string]interface{}
		if err := json.Unmarshal([]byte(plainReq), &body); err != nil {
			t.Errorf("request plaintext is not JSON: %v", err)
			return
		}
		reply, _ := json.Marshal(encrypted(r.URL.Path, body))
		resAad := fmt.Sprintf("v1:res:%s:%s:%s", api.sessionKey, rid, ts)
		resCt, resN, err := aesGcmEncrypt(api.sharedKey, string(reply), resAad)
		if err != nil {
			t.Errorf("encrypt reply: %v", err)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id": api.sessionKey, "ct": resCt, "n": resN, "rid": rid, "ts": env["ts"],
		})
	})
	return api
}

func ftGzipB64(t *testing.T, content string) string {
	t.Helper()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

const ftTicket = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" // 43 base64url chars, like the server

func TestAccessTicketSendsScopeAndResourceId(t *testing.T) {
	var got map[string]interface{}
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		if path != "/access-ticket" {
			t.Errorf("unexpected path %s", path)
		}
		got = body
		return map[string]interface{}{"success": true, "ticket": ftTicket, "expiresIn": 60}
	}, nil)

	if ticket := api.AccessTicket("file", "42"); ticket != ftTicket {
		t.Fatalf("AccessTicket returned %q, error %+v", ticket, api.Response)
	}
	if got["scope"] != "file" || got["resourceId"] != "42" {
		t.Fatalf("file ticket body wrong: %+v", got)
	}

	if ticket := api.AccessTicket("events", ""); ticket != ftTicket {
		t.Fatalf("events AccessTicket returned %q", ticket)
	}
	if _, present := got["resourceId"]; present || got["scope"] != "events" {
		t.Fatalf("events ticket must not carry resourceId: %+v", got)
	}
}

func TestAccessTicketDeniedSetsError(t *testing.T) {
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"success": false, "error": "feature_disabled"}
	}, nil)
	if ticket := api.AccessTicket("file", "42"); ticket != "" {
		t.Fatalf("denied ticket must be empty, got %q", ticket)
	}
	if api.Response.Error != "feature_disabled" {
		t.Fatalf("expected feature_disabled, got %+v", api.Response)
	}
}

func TestDownloadFileUsesTicketWithGetAndDecompresses(t *testing.T) {
	const content = "hello from the file store"
	var downloadPath, downloadQuery string
	var ticketScope string
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		ticketScope, _ = body["scope"].(string)
		return map[string]interface{}{"success": true, "ticket": ftTicket, "expiresIn": 60}
	}, func(w http.ResponseWriter, r *http.Request) {
		downloadPath = r.URL.Path
		downloadQuery = r.URL.RawQuery
		if r.Method != http.MethodGet {
			t.Errorf("download must use GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true, "name": "a.txt", "data": ftGzipB64(t, content),
		})
	})

	data := api.DownloadFile("42")
	if string(data) != content {
		t.Fatalf("DownloadFile returned %q, error %+v", data, api.Response)
	}
	if !api.Response.Success || api.Response.Error != "" {
		t.Fatalf("successful download must clear the previous error: %+v", api.Response)
	}
	if ticketScope != "file" {
		t.Fatalf("download must request a file ticket, got scope %q", ticketScope)
	}
	if downloadPath != "/files/download" {
		t.Fatalf("download path = %s, want /files/download", downloadPath)
	}
	if downloadQuery != "ticket="+ftTicket {
		t.Fatalf("download query = %q, want ticket=<ticket>", downloadQuery)
	}
}

func TestDownloadFileRefusesRedirect(t *testing.T) {
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"success": true, "ticket": ftTicket, "expiresIn": 60}
	}, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/files/download" {
			http.Redirect(w, r, "/elsewhere", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"success":true,"data":""}`)
	})
	api.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }

	if data := api.DownloadFile("42"); data != nil {
		t.Fatalf("redirect must not be followed, got %q", data)
	}
	if api.Response.Error != "external_download" {
		t.Fatalf("expected external_download, got %+v", api.Response)
	}
	if api.Response.Payload["url"] != "/elsewhere" {
		t.Fatalf("expected the redirect target in the payload, got %+v", api.Response.Payload)
	}
}

func TestDownloadFileReportsServerError(t *testing.T) {
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"success": true, "ticket": ftTicket, "expiresIn": 60}
	}, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"success":false,"error":"invalid_or_expired_ticket"}`)
	})
	if data := api.DownloadFile("42"); data != nil {
		t.Fatalf("expected nil, got %q", data)
	}
	if api.Response.Error != "invalid_or_expired_ticket" {
		t.Fatalf("expected invalid_or_expired_ticket, got %+v", api.Response)
	}
}

func TestDownloadFileRejectsCorruptGzip(t *testing.T) {
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"success": true, "ticket": ftTicket, "expiresIn": 60}
	}, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true, "name": "a.txt", "data": base64.StdEncoding.EncodeToString([]byte("not gzip")),
		})
	})
	if data := api.DownloadFile("42"); data != nil {
		t.Fatalf("corrupt payload must not be returned, got %q", data)
	}
	if api.Response.Error != "download_decode_failed" {
		t.Fatalf("expected download_decode_failed, got %+v", api.Response)
	}
}

func TestTriggerWebhookSendsWebIdAndPayloadEncrypted(t *testing.T) {
	var got map[string]interface{}
	var gotPath string
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		gotPath = path
		got = body
		return map[string]interface{}{"success": true, "status": 200, "data": map[string]interface{}{"ok": true}}
	}, nil)

	res := api.TriggerWebhook("hook-1", map[string]interface{}{"probe": "yes"})
	if res == nil {
		t.Fatalf("TriggerWebhook failed: %+v", api.Response)
	}
	if gotPath != "/webhook" {
		t.Fatalf("webhook path = %s", gotPath)
	}
	if got["webId"] != "hook-1" {
		t.Fatalf("webId must travel in the encrypted body: %+v", got)
	}
	payload, _ := got["payload"].(map[string]interface{})
	if payload["probe"] != "yes" {
		t.Fatalf("payload missing: %+v", got)
	}
	if res.Status != 200 {
		t.Fatalf("status = %d, want 200", res.Status)
	}
	if data, _ := res.Data.(map[string]interface{}); data["ok"] != true {
		t.Fatalf("remote data not returned: %+v", res.Data)
	}
}

func TestTriggerWebhookRemoteFailureSetsError(t *testing.T) {
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"success": false, "status": 502, "data": "bad gateway"}
	}, nil)
	if res := api.TriggerWebhook("hook-1", nil); res != nil {
		t.Fatalf("remote failure must not report success: %+v", res)
	}
	if api.Response.Error != "webhook_failed" {
		t.Fatalf("expected webhook_failed, got %+v", api.Response)
	}
}

func TestTriggerWebhookUnknownIdReportsServerError(t *testing.T) {
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"success": false, "error": "webhook_not_found"}
	}, nil)
	if res := api.TriggerWebhook("missing", nil); res != nil {
		t.Fatalf("expected nil, got %+v", res)
	}
	if api.Response.Error != "webhook_not_found" {
		t.Fatalf("expected webhook_not_found, got %+v", api.Response)
	}
	if api.Response.StatusCode != http.StatusOK {
		t.Fatalf("expected the answer's status 200, got %+v", api.Response)
	}
}

func TestSubscribeEventsDeliversFramesUntilCancelled(t *testing.T) {
	var gotTicket, gotPath, gotLastID string
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"success": true, "ticket": ftTicket, "expiresIn": 60}
	}, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotTicket = r.URL.Query().Get("ticket")
		gotLastID = r.Header.Get("Last-Event-ID")
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		fl, _ := w.(http.Flusher)
		_, _ = io.WriteString(w, ": ping\n\n")
		_, _ = io.WriteString(w, "id: 1\ndata: {\"type\":\"chat-message\",\"id\":1}\n\n")
		_, _ = io.WriteString(w, "id: 2\ndata: {\"type\":\"global-var-updated\"}\n\n")
		if fl != nil {
			fl.Flush()
		}
		<-r.Context().Done()
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var frames []EventFrame
	done := make(chan error, 1)
	go func() {
		done <- api.SubscribeEvents(ctx, "0", func(f EventFrame) {
			frames = append(frames, f)
			if len(frames) == 2 {
				cancel()
			}
		})
	}()
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("SubscribeEvents returned %v, want context.Canceled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("SubscribeEvents did not return after cancel")
	}

	if gotPath != "/events" || gotTicket != ftTicket {
		t.Fatalf("stream request path=%s ticket=%q", gotPath, gotTicket)
	}
	if gotLastID != "0" {
		t.Fatalf("Last-Event-ID = %q, want 0", gotLastID)
	}
	if len(frames) != 2 {
		t.Fatalf("expected 2 frames, got %d", len(frames))
	}
	if frames[0].ID != "1" || !strings.Contains(frames[0].Data, `"chat-message"`) {
		t.Fatalf("first frame wrong: %+v", frames[0])
	}
	if frames[1].ID != "2" || !strings.Contains(frames[1].Data, `"global-var-updated"`) {
		t.Fatalf("second frame wrong: %+v", frames[1])
	}
}

func TestSubscribeEventsRejectedTicketReturnsError(t *testing.T) {
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"success": true, "ticket": ftTicket, "expiresIn": 60}
	}, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"success":false,"error":"invalid_or_expired_ticket"}`)
	})
	err := api.SubscribeEvents(context.Background(), "", func(EventFrame) {
		t.Fatal("no frame expected from a rejected stream")
	})
	if err == nil || !strings.Contains(err.Error(), "invalid_or_expired_ticket") {
		t.Fatalf("expected invalid_or_expired_ticket error, got %v", err)
	}
}

func TestUpgradeOmitsBlankUsername(t *testing.T) {
	var bodies []map[string]interface{}
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		bodies = append(bodies, body)
		return map[string]interface{}{"success": true}
	}, nil)
	api.Upgrade("", "KEY-1")
	api.Upgrade("  ", "KEY-1")
	api.Upgrade("bob", "KEY-1")
	if len(bodies) != 3 {
		t.Fatalf("expected 3 requests, got %d", len(bodies))
	}
	if _, ok := bodies[0]["username"]; ok {
		t.Fatalf("a blank username must not be sent: %+v", bodies[0])
	}
	if _, ok := bodies[1]["username"]; ok {
		t.Fatalf("a whitespace username must not be sent: %+v", bodies[1])
	}
	if bodies[2]["username"] != "bob" {
		t.Fatalf("the username must reach the server: %+v", bodies[2])
	}
}

func TestLicenseWithTotpAndTypedStreamRefusal(t *testing.T) {
	var bodies []map[string]interface{}
	api := ftNewApi(t, func(path string, body map[string]interface{}) map[string]interface{} {
		bodies = append(bodies, body)
		return map[string]interface{}{"success": false, "error": "feature_disabled"}
	}, nil)
	api.License("KEY-1")
	api.LicenseWithTotp("KEY-1", "123456")
	if _, ok := bodies[0]["totpCode"]; ok || bodies[1]["totpCode"] != "123456" {
		t.Fatalf("totpCode must be sent only when given: %+v", bodies)
	}
	err := api.SubscribeEvents(context.Background(), "", func(EventFrame) {})
	var se *StreamError
	if !errors.As(err, &se) || se.Code != "feature_disabled" {
		t.Fatalf("expected a StreamError with feature_disabled, got %v", err)
	}
}
