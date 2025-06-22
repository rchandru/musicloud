package main

import (
	"flag"
	"fmt"
	"log"
	"musicloud/internal/drive"
	"musicloud/internal/watcher"
	"context"
	"os"
)

func printHelp() {
	fmt.Println(`WhatsApp Music Uploader
Usage:
  musicloud [options]

Options:
  -dir string
        Path to the folder to monitor for WhatsApp exports (default: ./watched or $MUSICLOUD_WATCH_FOLDER)
  -help
        Show this help message and exit

Environment variables:
  MUSICLOUD_WATCH_FOLDER              Folder to watch for new WhatsApp exports
  MUSICLOUD_GOOGLE_DRIVE_ID           Google Drive folder ID (takes precedence if set)
  MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME  Google Drive folder name (used if ID is not set; will be created if missing)
  MUSICLOUD_FFMPEG_PATH               Path to ffmpeg binary
  MUSICLOUD_CONFIG                    Path to Google API credentials JSON file (required)
  MUSICLOUD_OAUTH_TOKEN               OAuth token (managed automatically; not required)`)
	fmt.Println("\nEnvironment variable summary:")
	fmt.Printf("  %-32s %-20q %-20q %-20q\n", "Variable", "Current Value", "Default", "Effective (used)")
	fmt.Printf("  %-32s %-20q %-20q %-20q\n", "MUSICLOUD_WATCH_FOLDER", os.Getenv("MUSICLOUD_WATCH_FOLDER"), "./watched", getEnvWithDefault("MUSICLOUD_WATCH_FOLDER", "./watched"))
	fmt.Printf("  %-32s %-20q %-20q %-20q\n", "MUSICLOUD_GOOGLE_DRIVE_ID", os.Getenv("MUSICLOUD_GOOGLE_DRIVE_ID"), "", getEnvWithDefault("MUSICLOUD_GOOGLE_DRIVE_ID", ""))
	fmt.Printf("  %-32s %-20q %-20q %-20q\n", "MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME", os.Getenv("MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME"), "Recordings", getEnvWithDefault("MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME", "Recordings"))
	fmt.Printf("  %-32s %-20q %-20q %-20q\n", "MUSICLOUD_FFMPEG_PATH", os.Getenv("MUSICLOUD_FFMPEG_PATH"), "ffmpeg", getEnvWithDefault("MUSICLOUD_FFMPEG_PATH", "ffmpeg"))
	fmt.Printf("  %-32s %-20q %-20q %-20q\n", "MUSICLOUD_OAUTH_TOKEN", os.Getenv("MUSICLOUD_OAUTH_TOKEN"), "", getEnvWithDefault("MUSICLOUD_OAUTH_TOKEN", ""))
	fmt.Printf("  %-32s %-20q %-20q %-20q\n", "MUSICLOUD_CONFIG", os.Getenv("MUSICLOUD_CONFIG"), "(required)", os.Getenv("MUSICLOUD_CONFIG"))
}

func printConfig() {
	fmt.Println("Current Musicloud Configuration:")
	fmt.Printf("  %-30s %s (default: %s)\n", "MUSICLOUD_WATCH_FOLDER:", getEnvWithDefault("MUSICLOUD_WATCH_FOLDER", "./watched"), "./watched")
	fmt.Printf("  %-30s %s (default: %s)\n", "MUSICLOUD_GOOGLE_DRIVE_ID:", getEnvWithDefault("MUSICLOUD_GOOGLE_DRIVE_ID", ""), "empty")
	fmt.Printf("  %-30s %s (default: %s)\n", "MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME:", getEnvWithDefault("MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME", "Recordings"), "Recordings")
	fmt.Printf("  %-30s %s (default: %s)\n", "MUSICLOUD_FFMPEG_PATH:", getEnvWithDefault("MUSICLOUD_FFMPEG_PATH", "ffmpeg"), "ffmpeg")
	fmt.Printf("  %-30s %s (default: %s)\n", "MUSICLOUD_OAUTH_TOKEN:", getEnvWithDefault("MUSICLOUD_OAUTH_TOKEN", ""), "empty")
	fmt.Printf("  %-30s %s (required)\n", "MUSICLOUD_CONFIG:", os.Getenv("MUSICLOUD_CONFIG"))
	fmt.Println()
}

func getEnvWithDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func main() {
	help := flag.Bool("help", false, "Show help")
	dir := flag.String("dir", os.Getenv("MUSICLOUD_WATCH_FOLDER"), "Path to folder to scan")
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	printConfig()

	if *dir == "" {
		*dir = "./watched"
	}

	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		log.Fatalf("The folder to scan ('%s') does not exist. Please create it or specify a valid path using -dir or MUSICLOUD_WATCH_FOLDER.", *dir)
	}

	// Initialize Google Drive service
	creds, err := drive.GetCredentialsFile()
	if err != nil {
		log.Fatalf("Google Drive credentials error: %v", err)
	}
	if err := drive.InitializeDriveService(context.Background(), creds); err != nil {
		log.Fatalf("Failed to initialize Google Drive service: %v", err)
	}

	// Determine Google Drive folder ID
	folderID := os.Getenv("MUSICLOUD_GOOGLE_DRIVE_ID")
	if folderID == "" {
		folderName := os.Getenv("MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME")
		if folderName != "" {
			id, err := drive.GetOrCreateFolderID(drive.GetDriveService(), folderName)
			if err != nil {
				log.Fatalf("Failed to get or create Google Drive folder: %v", err)
			}
			folderID = id
		} else {
			folderID = "root"
		}
	}

	// Run the scan-and-upload batch process, always passing a valid folderID
	uploader := func(filePath, _ string) error {
		return drive.UploadFile(filePath, folderID)
	}
	watcher.ScanAndProcess(*dir, uploader)
}