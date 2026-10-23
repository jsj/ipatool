package emulator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAuthDiagnosticRequiresPrivateDirectory(t *testing.T) {
	directory := t.TempDir()
	if err := os.Chmod(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("IPATOOL_AUTH_RESPONSE_DIR", directory)
	if err := captureAuthResponse([]byte("fixture")); err == nil {
		t.Fatal("expected public directory to be rejected")
	}
	entries, err := os.ReadDir(directory)
	if err != nil || len(entries) != 0 {
		t.Fatalf("diagnostic wrote into a public directory: %v", err)
	}
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(directory, link); err != nil {
		t.Fatal(err)
	}
	t.Setenv("IPATOOL_AUTH_RESPONSE_DIR", link)
	if err := captureAuthResponse([]byte("fixture")); err == nil {
		t.Fatal("expected symlink directory to be rejected")
	}
}

func TestAuthDiagnosticDisabledByDefault(t *testing.T) {
	t.Setenv("IPATOOL_AUTH_RESPONSE_DIR", "")
	if err := captureAuthResponse([]byte("fixture")); err != nil {
		t.Fatal(err)
	}
}
