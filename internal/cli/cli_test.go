package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testEnv(key string) string {
	switch key {
	case "APPFLOWY_TOKEN":
		return "test-token-secret"
	case "APPFLOWY_URL":
		return ""
	default:
		return ""
	}
}

func runCLI(t *testing.T, srvURL string, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	getenv := testEnv
	if srvURL != "" {
		prev := getenv
		getenv = func(k string) string {
			if k == "APPFLOWY_URL" {
				return srvURL
			}
			return prev(k)
		}
	}
	var out, errOut bytes.Buffer
	code = Execute(args, strings.NewReader(""), &out, &errOut, getenv)
	return out.String(), errOut.String(), code
}

func TestWorkspacesList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token-secret" {
			t.Fatalf("bad auth: %q", r.Header.Get("Authorization"))
		}
		if r.URL.Path != "/api/workspace" {
			t.Fatalf("path %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"data":[]}`))
	}))
	defer srv.Close()

	out, errOut, code := runCLI(t, srv.URL, "workspaces", "list")
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errOut)
	}
	if !json.Valid([]byte(out)) {
		t.Fatalf("stdout: %s", out)
	}
}

func TestInvalidWorkspaceUUID(t *testing.T) {
	out, errOut, code := runCLI(t, "https://example.test", "folders", "list", "../bad")
	if code != exitUsage {
		t.Fatalf("code=%d", code)
	}
	if !strings.Contains(errOut, `"type":"usage"`) {
		t.Fatalf("stderr %s stdout %s", errOut, out)
	}
}

func TestAPIPathRejected(t *testing.T) {
	getenv := func(k string) string {
		switch k {
		case "APPFLOWY_URL":
			return "https://example.test"
		case "APPFLOWY_TOKEN":
			return "tok"
		default:
			return ""
		}
	}
	var out, errOut bytes.Buffer
	code := Execute([]string{"api", "/api/workspace?x=1"}, strings.NewReader(""), &out, &errOut, getenv)
	if code != exitUsage {
		t.Fatalf("code=%d", code)
	}
}

func TestWriteBlockedWithoutFlag(t *testing.T) {
	getenv := func(k string) string {
		if k == "APPFLOWY_URL" {
			return "https://example.test"
		}
		if k == "APPFLOWY_TOKEN" {
			return "tok"
		}
		return ""
	}
	var out, errOut bytes.Buffer
	code := Execute([]string{"api", "/api/workspace", "--method", "POST"}, strings.NewReader(""), &out, &errOut, getenv)
	if code != exitPolicy {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
}

func TestRedactInError(t *testing.T) {
	secret := "leaky-token-value-12345"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid token "+secret, http.StatusUnauthorized)
	}))
	defer srv.Close()
	getenv := func(k string) string {
		if k == "APPFLOWY_URL" {
			return srv.URL
		}
		if k == "APPFLOWY_TOKEN" {
			return secret
		}
		return ""
	}
	var out, errOut bytes.Buffer
	code := Execute([]string{"workspaces", "list"}, strings.NewReader(""), &out, &errOut, getenv)
	if code != exitAuth {
		t.Fatalf("code=%d", code)
	}
	if strings.Contains(errOut.String(), secret) {
		t.Fatalf("token leaked in stderr: %s", errOut.String())
	}
}

func TestDryRun(t *testing.T) {
	getenv := func(k string) string {
		if k == "APPFLOWY_URL" {
			return "https://af.test"
		}
		return ""
	}
	var out, errOut bytes.Buffer
	code := Execute([]string{"workspaces", "list", "--dry-run"}, strings.NewReader(""), &out, &errOut, getenv)
	if code != 0 {
		t.Fatalf("code=%d", code)
	}
	var v map[string]any
	if err := json.Unmarshal([]byte(out.String()), &v); err != nil {
		t.Fatal(err)
	}
	hdr := v["header"].(map[string]any)
	if hdr["Authorization"] != "Bearer <TOKEN>" {
		t.Fatalf("hdr %v", hdr)
	}
}
