package dedup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckDuplicate_NoDuplicate(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file1.txt")
	os.WriteFile(file, []byte("test"), 0644)
	found, err := CheckDuplicate(file, dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Errorf("expected no duplicate, got duplicate")
	}
}
