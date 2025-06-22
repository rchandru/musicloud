package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"musicloud/internal/drive"
	"musicloud/internal/ffmpeg"
	"musicloud/internal/metadata"
	"musicloud/internal/organizer"
)

type Watcher struct {
	watcher  *fsnotify.Watcher
	dir      string
	metadata metadata.Metadata
}

func NewWatcher(dir string) (*Watcher, error) {
	if dir == "" {
		return nil, fmt.Errorf("directory cannot be empty")
	}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher: w,
		dir:     dir,
	}, nil
}

// New creates a new Watcher. This is an alias for NewWatcher for external use.
func New(dir string) (*Watcher, error) {
	return NewWatcher(dir)
}

func (w *Watcher) Start() {
	err := w.watcher.Add(w.dir)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					w.handleNewFile(event.Name)
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
}

func isMediaFile(filePath string) bool {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".mp3", ".wav", ".m4a", ".aac", ".ogg", ".flac", ".mp4", ".mov", ".avi", ".mkv":
		return true
	default:
		return false
	}
}

func (w *Watcher) handleNewFile(filePath string) {
	if !isMediaFile(filePath) {
		return
	}

	log.Printf("New media file detected: %s\n", filePath)

	ffmpegAvailable, _ := ffmpeg.IsFFmpegInstalled()
	if !ffmpegAvailable {
		log.Printf("FFmpeg not found in environment. Skipping audio conversion step for this file.")
	}

	inputFile := filePath
	outputFile := inputFile
	if ffmpegAvailable && filepath.Ext(inputFile) != ".mp4" {
		outputFile = ffmpeg.GetOutputFilePath(inputFile)
		err := ffmpeg.ConvertToMP4(inputFile, outputFile)
		if err != nil {
			log.Printf("Error converting file to MP4: %s\n", err)
			return
		}
	}

	err := drive.UploadFile(outputFile, "")
	if err != nil {
		log.Printf("Error uploading file to Google Drive: %s\n", err)
		return
	}

	err = organizer.OrganizeFiles(nil, outputFile, organizer.Metadata(w.metadata))
	if err != nil {
		log.Printf("Error organizing file: %s\n", err)
		return
	}

	log.Printf("Processed and uploaded: %s\n", outputFile)
}

func (w *Watcher) Close() {
	w.watcher.Close()
}

// UploaderFunc defines the signature for uploading a file
// (filePath, folderID string) error
// This allows for dependency injection in tests.
type UploaderFunc func(filePath, folderID string) error

// ScanAndProcess scans the directory for media files and processes them using the provided uploader.
func ScanAndProcess(dir string, uploader UploaderFunc) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		filePath := filepath.Join(dir, entry.Name())
		if isMediaFile(filePath) {
			log.Printf("Found media file: %s\n", filePath)
			processMediaFile(filePath, uploader)
		}
	}
}

// processMediaFile processes a single media file using the provided uploader.
func processMediaFile(filePath string, uploader UploaderFunc) {
	ffmpegAvailable, _ := ffmpeg.IsFFmpegInstalled()
	if !ffmpegAvailable {
		log.Printf("FFmpeg not found in environment. Skipping audio conversion step for this file.")
	}

	inputFile := filePath
	outputFile := inputFile
	if ffmpegAvailable && filepath.Ext(inputFile) != ".mp4" {
		outputFile = ffmpeg.GetOutputFilePath(inputFile)
		err := ffmpeg.ConvertToMP4(inputFile, outputFile)
		if err != nil {
			log.Printf("Error converting file to MP4: %s\n", err)
			return
		}
	}

	err := uploader(outputFile, "")
	if err != nil {
		log.Printf("Error uploading file to Google Drive: %s\n", err)
		return
	}

	// Organizer and metadata can be added here if needed
	log.Printf("Processed and uploaded: %s\n", outputFile)
}

// For production use, call ScanAndProcess with drive.UploadFile as the uploader.
// watcher.ScanAndProcess(dir, drive.UploadFile)