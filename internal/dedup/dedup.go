package dedup

import (
	"os"
	"path/filepath"
)

// CheckDuplicate checks if a file already exists in the specified directory.
func CheckDuplicate(filePath string, directory string) (bool, error) {
	// Get the base name of the file
	fileName := filepath.Base(filePath)

	// Iterate through the files in the directory
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the file name matches and if it's not the same file
		if info.Name() == fileName && path != filePath {
			return filepath.SkipDir // Duplicate found, skip further processing
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	// If we reach here, no duplicate was found
	return false, nil
}