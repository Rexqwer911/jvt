package version

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Manager handles Java version switching
type Manager struct {
	installDir string
}

// NewManager creates a new version manager
func NewManager(installDir string) *Manager {
	return &Manager{
		installDir: installDir,
	}
}

// GetCurrentVersion returns the currently active Java version
func (m *Manager) GetCurrentVersion() (string, error) {
	javaHome := os.Getenv("JAVA_HOME")

	if javaHome == "" {
		return "", fmt.Errorf("JAVA_HOME is not set")
	}

	// Check if it's a jvt-managed version
	// Only strict check if we are sure jvt solely manages it, but simple check is okay
	if !strings.HasPrefix(javaHome, m.installDir) {
		// On windows path separators might differ, but HasPrefix usually handles generic paths if normalized.
		// Let's rely on standard logic.
		return "", fmt.Errorf("current Java is not managed by jvt")
	}

	version := filepath.Base(javaHome)
	return version, nil
}
