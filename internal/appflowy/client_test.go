package appflowy

import "testing"

func TestValidateBaseURL(t *testing.T) {
	n, insecure, err := ValidateBaseURL("https://cloud.example/")
	if err != nil || insecure || n != "https://cloud.example" {
		t.Fatalf("got %q %v %v", n, insecure, err)
	}
	n, _, err = ValidateBaseURL("https://cloud.example/app?x=1#frag")
	if err != nil || n != "https://cloud.example/app" {
		t.Fatalf("strip query/frag: %q err=%v", n, err)
	}
	if _, _, err := ValidateBaseURL("https://user:pass@cloud.example"); err == nil {
		t.Fatal("expected userinfo error")
	}
	if _, _, err := ValidateBaseURL("ftp://x"); err == nil {
		t.Fatal("expected scheme error")
	}
}

func TestValidateAPIPath(t *testing.T) {
	p, err := ValidateAPIPath("/api/workspace")
	if err != nil || p != "/api/workspace" {
		t.Fatalf("got %q %v", p, err)
	}
	if _, err := ValidateAPIPath("/api/../admin"); err == nil {
		t.Fatal("expected .. error")
	}
	if _, err := ValidateAPIPath("/api/x%2e%2e/y"); err == nil {
		t.Fatal("expected decoded .. error")
	}
	if _, err := ValidateAPIPath("/api/x?y=1"); err == nil {
		t.Fatal("expected query error")
	}
	if _, err := ValidateAPIPath("/other"); err == nil {
		t.Fatal("expected /api prefix error")
	}
}

func TestParseUUID(t *testing.T) {
	got, err := ParseUUID("550E8400-E29B-41D4-A716-446655440000")
	if err != nil || got != "550e8400-e29b-41d4-a716-446655440000" {
		t.Fatalf("got %q err=%v", got, err)
	}
	if _, err := ParseUUID("not-a-uuid"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := ParseUUID("550e8400-e29b-41d4-a716"); err == nil {
		t.Fatal("expected error")
	}
}

func TestClientRedact(t *testing.T) {
	secret := "super-secret-jwt-value"
	c := &Client{Token: secret}
	msg := "Authorization: Bearer " + secret + " failed"
	out := c.Redact(msg)
	if out == msg || contains(out, secret) {
		t.Fatalf("token not redacted: %q", out)
	}
	if !contains(out, "<TOKEN>") {
		t.Fatalf("expected placeholder: %q", out)
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
