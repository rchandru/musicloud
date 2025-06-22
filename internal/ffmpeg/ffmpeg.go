package ffmpeg

import (
	"os/exec"
	"path/filepath"
)

// ConvertToMP4 converts the given audio file to MP4 format using FFmpeg.
func ConvertToMP4(inputFile string, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-codec:a", "aac", "-b:a", "192k", outputFile)
	return cmd.Run()
}

// GetFFmpegPath returns the path to the FFmpeg executable.
func GetFFmpegPath() (string, error) {
	return exec.LookPath("ffmpeg")
}

// IsFFmpegInstalled checks if FFmpeg is installed on the system.
func IsFFmpegInstalled() (bool, error) {
	_, err := GetFFmpegPath()
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetOutputFilePath generates the output file path for the converted MP4 file.
func GetOutputFilePath(inputFile string) string {
	return filepath.Join(filepath.Dir(inputFile), filepath.Base(inputFile)+".mp4")
}