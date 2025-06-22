package watcher

import (
	"os"
	"testing"
)

func TestNewWatcher_InvalidDir(t *testing.T) {
	_, err := NewWatcher("")
	if err == nil {
		t.Error("expected error for empty dir")
	}
}

func TestIsMediaFile(t *testing.T) {
	cases := map[string]bool{
		"song.mp3":  true,
		"audio.wav": true,
		"clip.m4a":  true,
		"voice.aac": true,
		"music.ogg": true,
		"track.flac": true,
		"video.mp4": true,
		"movie.mov": true,
		"film.avi":  true,
		"show.mkv":  true,
		"doc.txt":   false,
		"image.jpg": false,
	}
	for file, want := range cases {
		if got := isMediaFile(file); got != want {
			t.Errorf("isMediaFile(%q) = %v, want %v", file, got, want)
		}
	}
}

func TestScanAndProcess_SingleMediaFile(t *testing.T) {
	dir := t.TempDir()
	mediaFile := dir + "/test.mp3"
	os.WriteFile(mediaFile, []byte("dummy audio"), 0644)

	uploaded := ""
	mockUploader := func(filePath, folderID string) error {
		uploaded = filePath
		return nil
	}

	ScanAndProcess(dir, mockUploader)

	if uploaded != mediaFile {
		t.Errorf("expected %s to be uploaded, got %s", mediaFile, uploaded)
	}
}
