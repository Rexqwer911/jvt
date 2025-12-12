package download

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

// Downloader handles downloading files with progress tracking
type Downloader struct {
	cacheDir string
}

// NewDownloader creates a new downloader instance
func NewDownloader(cacheDir string) *Downloader {
	return &Downloader{
		cacheDir: cacheDir,
	}
}

// Download downloads a file from URL to the cache directory
func (d *Downloader) Download(url, filename string, showProgress bool) (string, error) {
	// Ensure cache directory exists
	if err := os.MkdirAll(d.cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	destPath := filepath.Join(d.cacheDir, filename)

	// Check if file already exists
	if _, err := os.Stat(destPath); err == nil {
		fmt.Printf("File already exists in cache: %s\n", filename)
		return destPath, nil
	}

	// Create HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create destination file
	out, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Download with progress bar
	if showProgress {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			fmt.Sprintf("Downloading %s", filename),
		)
		_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	} else {
		_, err = io.Copy(out, resp.Body)
	}

	if err != nil {
		os.Remove(destPath) // Clean up partial download
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return destPath, nil
}

// VerifyChecksum verifies the SHA256 checksum of a file
func (d *Downloader) VerifyChecksum(filePath, expectedChecksum string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	actualChecksum := hex.EncodeToString(hash.Sum(nil))

	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

// DownloadAndVerify downloads a file and verifies its checksum
func (d *Downloader) DownloadAndVerify(url, filename, checksum string) (string, error) {
	filePath, err := d.Download(url, filename, true)
	if err != nil {
		return "", err
	}

	if checksum != "" {
		fmt.Println("Verifying checksum...")
		if err := d.VerifyChecksum(filePath, checksum); err != nil {
			return "", fmt.Errorf("checksum verification failed: %w", err)
		}
		fmt.Println("Checksum verified successfully!")
	}

	return filePath, nil
}
