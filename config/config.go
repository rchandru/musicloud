package config

import (
	"os"
)

type Config struct {
	WatchFolder   string
	GoogleDriveID string
	FFmpegPath    string
	OAuthToken    string
}

func LoadConfig() (*Config, error) {
	return &Config{
		WatchFolder:   getEnv("MUSICLOUD_WATCH_FOLDER", "./watched"),
		GoogleDriveID: getEnv("MUSICLOUD_GOOGLE_DRIVE_ID", ""),
		FFmpegPath:    getEnv("MUSICLOUD_FFMPEG_PATH", "ffmpeg"),
		OAuthToken:    getEnv("MUSICLOUD_OAUTH_TOKEN", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}