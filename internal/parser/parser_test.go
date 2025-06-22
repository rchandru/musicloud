package parser

import (
	"testing"
)

func TestParseWhatsAppExport_FileNotFound(t *testing.T) {
	_, err := ParseWhatsAppExport("nonexistent.txt")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

// Add more tests for isMusicRelated and extractMusicFileInfo if exported
