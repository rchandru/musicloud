package drive

import "fmt"

// Uploader defines the interface for uploading files.
type Uploader interface {
	UploadFile(filePath string, folderID string) error
}

// MockUploader is a mock implementation of Uploader for testing.
type MockUploader struct {
	UploadedFiles []string
	ShouldFail    bool
}

func (m *MockUploader) UploadFile(filePath string, folderID string) error {
	if m.ShouldFail {
		return fmt.Errorf("mock upload failed")
	}
	m.UploadedFiles = append(m.UploadedFiles, filePath)
	return nil
}
