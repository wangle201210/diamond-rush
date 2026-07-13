package originalgame

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePathFallsBackToExecutableDirectory(t *testing.T) {
	executableDir := t.TempDir()
	relativePath := filepath.Join("testdata", "assets", "manifest.json")
	want := filepath.Join(executableDir, relativePath)
	if err := os.MkdirAll(filepath.Dir(want), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(want, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := resolvePathFromExecutable(relativePath, filepath.Join(executableDir, "originalrush"))
	if got != want {
		t.Fatalf("resolvePathFromExecutable() = %q, want %q", got, want)
	}
}
