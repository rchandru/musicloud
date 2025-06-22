package organizer

import (
	"fmt"
	"time"

	"google.golang.org/api/drive/v3"
)

type Metadata struct {
	GroupName   string
	Teacher     string
	SessionType string
	SongsTaught []string
	Ragas       []string
	Talas       []string
	Composers   []string
}

func OrganizeFiles(service *drive.Service, fileID string, metadata Metadata) error {
	if service == nil {
		return fmt.Errorf("service is nil")
	}

	// Create a folder in Google Drive based on the recording date
	recordingDate := time.Now().Format("2006-01-02")
	folderName := fmt.Sprintf("%s - %s", recordingDate, metadata.GroupName)

	folderID, err := createFolder(service, folderName)
	if err != nil {
		return err
	}

	// Move the uploaded file to the newly created folder
	err = moveFileToFolder(service, fileID, folderID)
	if err != nil {
		return err
	}

	// Optionally, save metadata to a file or database
	err = saveMetadata(metadata, folderID)
	if err != nil {
		return err
	}

	return nil
}

func createFolder(service *drive.Service, folderName string) (string, error) {
	if service == nil {
		return "", fmt.Errorf("service is nil")
	}

	folder := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
	}

	createdFolder, err := service.Files.Create(folder).Do()
	if err != nil {
		return "", err
	}

	return createdFolder.Id, nil
}

func moveFileToFolder(service *drive.Service, fileID string, folderID string) error {
	// Retrieve the existing file
	file, err := service.Files.Get(fileID).Do()
	if err != nil {
		return err
	}

	// Update the file's parents to include the new folder
	file.Parents = append(file.Parents, folderID)

	_, err = service.Files.Update(fileID, file).Do()
	if err != nil {
		return err
	}

	return nil
}

func saveMetadata(metadata Metadata, folderID string) error {
	// Implement saving metadata logic here (e.g., to a database or a file)
	// This is a placeholder for the actual implementation
	return nil
}