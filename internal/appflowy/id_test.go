package appflowy

import "testing"

func TestParseDepth(t *testing.T) {
	d, err := ParseDepth("10")
	if err != nil || d != "10" {
		t.Fatalf("got %q %v", d, err)
	}
	if _, err := ParseDepth("0"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := ParseDepth("11"); err == nil {
		t.Fatal("expected error")
	}
}
