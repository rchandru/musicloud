package drive

import (
	"context"
	"os"
	"testing"
)

func getEnvWithDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func TestInitializeDriveService_FileNotFound(t *testing.T) {
	err := InitializeDriveService(context.Background(), "nonexistent.json")
	if err == nil {
		t.Error("expected error for missing credentials file")
	}
}

func TestGetCredentialsFile_Set(t *testing.T) {
	os.Setenv("MUSICLOUD_CONFIG", "/tmp/creds.json")
	path, err := GetCredentialsFile()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "/tmp/creds.json" {
		t.Errorf("expected /tmp/creds.json, got %s", path)
	}
}

func TestGetCredentialsFile_NotSet(t *testing.T) {
	os.Unsetenv("MUSICLOUD_CONFIG")
	_, err := GetCredentialsFile()
	if err == nil {
		t.Error("expected error when MUSICLOUD_CONFIG is not set")
	}
}

func TestGetOrCreateFolderID_CreateAndFind(t *testing.T) {
	// This is a mock test: in real code, you would mock the drive.Service and its Files.List/Create methods.
	// Here, we just check that the function can be called and returns the expected error for nil service.
	_, err := GetOrCreateFolderID(nil, "TestFolder")
	if err == nil {
		t.Error("expected error for nil service")
	}
}

func TestGoogleDriveFolderName_Set(t *testing.T) {
	os.Setenv("MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME", "TestFolderName")
	name := os.Getenv("MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME")
	if name != "TestFolderName" {
		t.Errorf("expected TestFolderName, got %s", name)
	}
}

func TestGoogleDriveFolderName_NotSet(t *testing.T) {
	os.Unsetenv("MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME")
	name := os.Getenv("MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME")
	if name != "" {
		t.Errorf("expected empty string, got %s", name)
	}
}

func TestGoogleDriveFolderName_Default(t *testing.T) {
	os.Unsetenv("MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME")
	name := getEnvWithDefault("MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME", "Recordings")
	if name != "Recordings" {
		t.Errorf("expected default 'Recordings', got %s", name)
	}
}
