package install

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Installer handles Java installation
type Installer struct {
	installDir string
}

// NewInstaller creates a new installer instance
func NewInstaller(installDir string) *Installer {
	return &Installer{
		installDir: installDir,
	}
}

// Install extracts a Java archive to the installation directory
func (i *Installer) Install(archivePath, version string) error {
	// Ensure install directory exists
	if err := os.MkdirAll(i.installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	versionDir := filepath.Join(i.installDir, version)

	// Check if version already installed
	if _, err := os.Stat(versionDir); err == nil {
		return fmt.Errorf("version %s is already installed", version)
	}

	// Create version directory
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		return fmt.Errorf("failed to create version directory: %w", err)
	}

	// Extract archive
	fmt.Printf("Extracting to %s...\n", versionDir)

	if strings.HasSuffix(archivePath, ".zip") {
		if err := i.extractZip(archivePath, versionDir); err != nil {
			os.RemoveAll(versionDir) // Clean up on error
			return fmt.Errorf("failed to extract zip archive: %w", err)
		}
	} else if strings.HasSuffix(archivePath, ".tar.gz") {
		if err := i.extractTarGz(archivePath, versionDir); err != nil {
			os.RemoveAll(versionDir) // Clean up on error
			return fmt.Errorf("failed to extract tar.gz archive: %w", err)
		}
	} else {
		os.RemoveAll(versionDir)
		return fmt.Errorf("unsupported archive format: %s", archivePath)
	}

	fmt.Printf("Java %s installed successfully!\n", version)
	return nil
}

// extractZip extracts a zip archive to the destination directory
func (i *Installer) extractZip(archivePath, destDir string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Find the root directory in the archive (usually jdk-x.x.x+x)
	var rootDir string
	if len(r.File) > 0 {
		firstFile := r.File[0].Name
		parts := strings.Split(firstFile, "/")
		if len(parts) > 0 {
			rootDir = parts[0]
		}
	}

	for _, f := range r.File {
		// Skip the root directory itself
		relativePath := f.Name
		if rootDir != "" && strings.HasPrefix(relativePath, rootDir+"/") {
			relativePath = strings.TrimPrefix(relativePath, rootDir+"/")
		}

		if relativePath == "" {
			continue
		}

		fpath := filepath.Join(destDir, relativePath)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Extract file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// extractTarGz extracts a tar.gz archive to the destination directory
func (i *Installer) extractTarGz(archivePath, destDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	uncompressedStream, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)

	// Find the root directory
	// In tar.gz, we invoke Next() to get headers. We can't easily peek first like in zip.
	// So we'll have to detect rootDir on the fly or assuming the first entry determines it if it's a dir.
	// However, tar ordering isn't guaranteed to be depth-first strictly, but usually is.

	// A simpler approach for stripping root component:
	// Determine root component from the first entry.
	var rootDir string

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if rootDir == "" {
			parts := strings.Split(header.Name, "/")
			if len(parts) > 0 {
				rootDir = parts[0]
			}
		}

		relativePath := header.Name
		if rootDir != "" && strings.HasPrefix(relativePath, rootDir+"/") {
			relativePath = strings.TrimPrefix(relativePath, rootDir+"/")
		}

		if relativePath == "" {
			continue
		}

		fpath := filepath.Join(destDir, relativePath)

		// Check for ZipSlip
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(fpath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
				return err
			}
			outFile, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// Uninstall removes an installed Java version
func (i *Installer) Uninstall(version string) error {
	versionDir := filepath.Join(i.installDir, version)

	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed", version)
	}

	if err := os.RemoveAll(versionDir); err != nil {
		return fmt.Errorf("failed to remove version: %w", err)
	}

	fmt.Printf("Java %s uninstalled successfully!\n", version)
	return nil
}

// ListInstalled returns a list of installed Java versions
func (i *Installer) ListInstalled() ([]string, error) {
	entries, err := os.ReadDir(i.installDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read install directory: %w", err)
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}

	return versions, nil
}

// IsInstalled checks if a version is installed
func (i *Installer) IsInstalled(version string) bool {
	versionDir := filepath.Join(i.installDir, version)
	_, err := os.Stat(versionDir)
	return err == nil
}

// GetJavaHome returns the JAVA_HOME path for a version
func (i *Installer) GetJavaHome(version string) string {
	return filepath.Join(i.installDir, version)
}
