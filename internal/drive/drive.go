package drive

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

var (
	driveService *drive.Service
)

// InitializeDriveService initializes the Google Drive API service
func InitializeDriveService(ctx context.Context, credentialsFile string) error {
	b, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)
	driveService, err = drive.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve drive client: %v", err)
	}

	return nil
}

// getClient retrieves a token, saves the token, and returns an authenticated HTTP client
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// tokenFromFile retrieves a token from a local file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// getTokenFromWeb requests a token from the web, then returns the retrieved token
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n", url)

	var code string
	fmt.Print("Enter the authorization code: ")
	fmt.Scan(&code)

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// saveToken saves a token to a file
func saveToken(file string, token *oauth2.Token) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// UploadFile uploads a file to Google Drive
func UploadFile(filePath string, folderID string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	fileMetadata := &drive.File{
		Name:    filepath.Base(filePath),
		Parents: []string{folderID},
	}

	f, err := driveService.Files.Create(fileMetadata).Media(file).Do()
	if err != nil {
		return fmt.Errorf("unable to upload file: %v", err)
	}

	fmt.Printf("File uploaded successfully: %s\n", f.WebViewLink)
	return nil
}

// GetCredentialsFile returns the path to the credentials file from the MUSICLOUD_CONFIG environment variable, or an error if not set.
func GetCredentialsFile() (string, error) {
	path := os.Getenv("MUSICLOUD_CONFIG")
	if path == "" {
		return "", fmt.Errorf("MUSICLOUD_CONFIG environment variable not set")
	}
	return path, nil
}

// GetOrCreateFolderID finds a folder by name in the user's Drive, or creates it if it doesn't exist.
func GetOrCreateFolderID(service *drive.Service, folderName string) (string, error) {
	if service == nil {
		return "", fmt.Errorf("drive service is nil")
	}
	// Search for the folder by name
	query := fmt.Sprintf("mimeType='application/vnd.google-apps.folder' and name='%s' and trashed=false", folderName)
	fileList, err := service.Files.List().Q(query).Fields("files(id, name)").Do()
	if err != nil {
		return "", fmt.Errorf("error searching for folder: %v", err)
	}
	if len(fileList.Files) > 0 {
		return fileList.Files[0].Id, nil
	}
	// Not found, create it
	folder := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
	}
	created, err := service.Files.Create(folder).Do()
	if err != nil {
		return "", fmt.Errorf("error creating folder: %v", err)
	}
	return created.Id, nil
}

// GetDriveService returns the initialized Google Drive service instance.
func GetDriveService() *drive.Service {
	return driveService
}