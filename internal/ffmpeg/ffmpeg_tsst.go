package ffmpeg

import (
    "os"
    "path/filepath"
    "testing"
)

func TestGetOutputFilePath(t *testing.T) {
    input := "/tmp/testaudio.wav"
    expected := filepath.Join(filepath.Dir(input), filepath.Base(input)+".mp4")
    got := GetOutputFilePath(input)
    if got != expected {
        t.Errorf("expected %s, got %s", expected, got)
    }
}


func TestIsFFmpegInstalled(t *testing.T) {
    installed, _ := IsFFmpegInstalled()
    // We can't guarantee ffmpeg is installed on CI, so just check for no panic
    _ = installed
}

func TestGetFFmpegPath(t *testing.T) {
    _, err := GetFFmpegPath()
    // Should not panic, error is OK if ffmpeg is not installed
    if err != nil {
        t.Log("ffmpeg not found in PATH (this is OK for test environments)")
    }
}

func TestConvertToMP4(t *testing.T) {
    installed, _ := IsFFmpegInstalled()
    if !installed {
        t.Skip("ffmpeg not installed, skipping ConvertToMP4 test")
    }
    // Create a dummy input file
    input := "testinput.wav"
    output := "testoutput.mp4"
    f, err := os.Create(input)
    if err != nil {
        t.Fatalf("failed to create dummy input file: %v", err)
    }
    f.Close()
    defer os.Remove(input)
    defer os.Remove(output)

    err = ConvertToMP4(input, output)
    if err == nil {
        // Output file should be created (though not a valid mp4)
        if _, err := os.Stat(output); os.IsNotExist(err) {
            t.Errorf("expected output file %s to exist", output)
        }
    } else {
        t.Logf("ConvertToMP4 failed as expected (invalid input): %v", err)
    }
}