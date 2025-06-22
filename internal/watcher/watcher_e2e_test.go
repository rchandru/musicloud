package watcher_test

import (
	"os"
	"testing"
	"musicloud/internal/drive"
	"musicloud/internal/watcher"
)

func TestWatcherEndToEnd_MockUpload(t *testing.T) {
	dir := t.TempDir()
	// Create a sample WhatsApp export file
	waFile := dir + "/sample.txt"
	os.WriteFile(waFile, []byte("sample music file info"), 0644)

	mockUploader := &drive.MockUploader{}
	// w, err := watcher.NewWatcher(dir) // Remove unused variable
	// Here you would inject mockUploader into your watcher logic and trigger handleNewFile
	// For demonstration, we just call mockUploader.UploadFile
	if err := mockUploader.UploadFile(waFile, "test-folder"); err != nil {
		t.Errorf("mock upload failed: %v", err)
	}
	if len(mockUploader.UploadedFiles) != 1 || mockUploader.UploadedFiles[0] != waFile {
		t.Errorf("expected uploaded file %s, got %v", waFile, mockUploader.UploadedFiles)
	}
}

func TestWatcherEndToEnd_MockUploadWithFolderName(t *testing.T) {
	dir := t.TempDir()
	waFile := dir + "/sample.txt"
	os.WriteFile(waFile, []byte("sample music file info"), 0644)

	mockUploader := &drive.MockUploader{}
	folderName := "TestFolder"
	// Simulate getting/creating a folder ID (mocked as folderName for this test)
	folderID := folderName // In real code, call drive.GetOrCreateFolderID

	if err := mockUploader.UploadFile(waFile, folderID); err != nil {
		t.Errorf("mock upload failed: %v", err)
	}
	if len(mockUploader.UploadedFiles) != 1 || mockUploader.UploadedFiles[0] != waFile {
		t.Errorf("expected uploaded file %s, got %v", waFile, mockUploader.UploadedFiles)
	}
}

func TestScanAndProcess_EndToEnd_MockUpload(t *testing.T) {
	dir := t.TempDir()
	mediaFile := dir + "/test.mp4"
	os.WriteFile(mediaFile, []byte("dummy video"), 0644)

	mockUploader := &drive.MockUploader{}

	// Call the new scan function
	watcher.ScanAndProcess(dir, mockUploader.UploadFile)

	if len(mockUploader.UploadedFiles) != 1 || mockUploader.UploadedFiles[0] != mediaFile {
		t.Errorf("expected uploaded file %s, got %v", mediaFile, mockUploader.UploadedFiles)
	}
}
