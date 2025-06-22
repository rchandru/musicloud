# WhatsApp Music Uploader

## Project Summary
The WhatsApp Music Uploader is a Go application that scans a folder for new media files (audio/video), optionally converts audio to MP4 using FFmpeg, and uploads the results to Google Drive. It ensures no duplicates are uploaded and organizes files by recording date and user-supplied metadata.

**Key Capabilities:**
- Scans a specified folder for new media files (audio/video) and uploads them in a single batch when run.
- Optionally converts audio files to MP4 format using FFmpeg (if installed).
- Uploads processed files to Google Drive using OAuth 2.0 authentication.
- Prevents duplicate uploads.
- Organizes files in Google Drive by recording date and metadata (group, teacher, session type, songs, ragas, talas, composers).
- Allows user-supplied metadata for each recording.
- Configuration is managed via environment variables, with secure handling of Google credentials using the `MUSICLOUD_CONFIG` variable.
- Includes clear setup instructions and environment variable documentation in the README.

---

## Technical Details

### Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/whatsapp-music-uploader.git
   ```
2. Navigate to the project directory:
   ```sh
   cd whatsapp-music-uploader
   ```
3. Install the necessary dependencies:
   ```sh
   go mod tidy
   ```

### Configuration
Before running the application, ensure that you have configured the necessary settings in `config/config.go`, including paths and API keys for Google Drive.

### Usage
1. Place your media files in the folder specified by `MUSICLOUD_WATCH_FOLDER` (or use `-dir`).
2. Run the application:
   ```sh
   go run cmd/main.go
   ```
3. The application will scan the folder for supported media files and process/upload them to Google Drive in a single batch.

### Supported Media File Types
The application will process files with the following extensions:
- .mp3, .wav, .m4a, .aac, .ogg, .flac, .mp4, .mov, .avi, .mkv

### Metadata Input
The application allows users to input metadata for each recording, including:
- Group Name
- Teacher
- Session Type (In-person or Virtual)
- Songs Taught
- Ragas
- Talas
- Composers

### Google Drive API Setup
To use Google Drive upload features, you must set up OAuth credentials in Google Cloud Console:

1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project (or select an existing one).
3. Enable the Google Drive API:
   - Go to "APIs & Services" > "Library".
   - Search for "Google Drive API" and click "Enable".
4. Create OAuth 2.0 credentials:
   - Go to "APIs & Services" > "Credentials".
   - Click "+ CREATE CREDENTIALS" > "OAuth client ID".
   - Choose "Desktop app" as the application type.
   - Name it (e.g., "Musicloud Desktop").
   - Click "Create".
   - Download the JSON file (usually named `credentials.json`).
5. **Do NOT check this file into version control.**
6. Set the environment variable before running the app:
   ```sh
   export MUSICLOUD_CONFIG=/path/to/your/credentials.json
   ```
7. On first run, the app will prompt you to authorize access in your browser. A `token.json` will be created for future use.

For more details, see the Google [OAuth 2.0 for Desktop Apps documentation](https://developers.google.com/identity/protocols/oauth2/native-app).

### FFmpeg Optional Usage

If FFmpeg is not installed or not found in your environment, the application will skip the audio conversion step and upload the original file as-is. You will see a log message indicating that FFmpeg was not found and conversion was skipped. All other processing and uploads will continue as normal.

- To enable audio conversion, ensure FFmpeg is installed and available in your PATH, or set the `MUSICLOUD_FFMPEG_PATH` environment variable to the correct binary location.

### Environment Variables

| Variable                          | Default Value         | Description                                                    |
|-----------------------------------|----------------------|----------------------------------------------------------------|
| MUSICLOUD_WATCH_FOLDER            | ./watched            | Folder to watch for new WhatsApp exports                       |
| MUSICLOUD_GOOGLE_DRIVE_ID         | (empty)              | Google Drive folder ID (if used, overrides folder name)         |
| MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME| Recordings            | Google Drive folder name (auto-creates/uses folder by name)     |
| MUSICLOUD_FFMPEG_PATH             | ffmpeg               | Path to ffmpeg binary                                          |
| MUSICLOUD_OAUTH_TOKEN             | (empty)              | OAuth token (not used directly, see Drive setup)               |
| MUSICLOUD_CONFIG                  | (none, must be set)  | Path to Google API credentials JSON file                       |

- `MUSICLOUD_CONFIG` must be set to use Google Drive features.
- If both `MUSICLOUD_GOOGLE_DRIVE_ID` and `MUSICLOUD_GOOGLE_DRIVE_FOLDER_NAME` are set, the ID takes precedence.
- Other variables have defaults and can be overridden as needed.

---

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.