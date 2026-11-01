package appflowy

import "testing"

func TestLoadConfigStripsURLQuery(t *testing.T) {
	getenv := func(k string) string {
		switch k {
		case "APPFLOWY_URL":
			return "https://cloud.example/?q=1#frag"
		case "APPFLOWY_TOKEN":
			return "tok"
		default:
			return ""
		}
	}
	cfg, err := LoadConfig(getenv)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.URL != "https://cloud.example" {
		t.Fatalf("got %q", cfg.URL)
	}
}
