package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRejectsMalformedDocuments(t *testing.T) {
	for _, s := range []string{"", "unknown: true", "version: x\nversion: y", "version: x\n---\nversion: y"} {
		p := filepath.Join(t.TempDir(), "grimley.yaml")
		if e := os.WriteFile(p, []byte(s), 0600); e != nil {
			t.Fatal(e)
		}
		if _, e := Load(p); e == nil {
			t.Fatalf("accepted %q", s)
		}
	}
}
