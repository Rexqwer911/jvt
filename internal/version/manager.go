package version

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

// IsVersionActive checks if the specified version is currently active
func (m *Manager) IsVersionActive(version string) (bool, error) {
	currentVersion, err := m.GetCurrentVersion()
	if err != nil {
		return false, nil // If no version is set, it's not active
	}
	return currentVersion == version, nil
}

// CompareVersions compares two Java version strings.
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2.
// Format: major.minor.patch+build (e.g., "17.0.10+7"). Missing parts are treated as 0,
// so shorter forms like "17", "17.0", or "17.0.1" are valid and equivalent to "17.0.0+0", "17.0.0+0", and "17.0.1+0".
func CompareVersions(v1, v2 string) (int, error) {
	parts1, err := parseVersion(v1)
	if err != nil {
		return 0, fmt.Errorf("invalid version v1: %w", err)
	}

	parts2, err := parseVersion(v2)
	if err != nil {
		return 0, fmt.Errorf("invalid version v2: %w", err)
	}

	// Compare each part: major, minor, patch, build
	for i := 0; i < 4; i++ {
		if parts1[i] < parts2[i] {
			return -1, nil
		}
		if parts1[i] > parts2[i] {
			return 1, nil
		}
	}

	return 0, nil
}

// parseVersion parses a version string into [major, minor, patch, build]
func parseVersion(version string) ([4]int, error) {
	var parts [4]int

	// Split by '+' to separate version from build
	mainAndBuild := strings.Split(version, "+")
	if mainAndBuild[0] == "" {
		return parts, fmt.Errorf("invalid version format: %s", version)
	}

	// Parse major.minor.patch
	versionParts := strings.Split(mainAndBuild[0], ".")
	for i := 0; i < len(versionParts) && i < 3; i++ {
		val, err := strconv.Atoi(versionParts[i])
		if err != nil {
			return parts, fmt.Errorf("invalid version number: %s", versionParts[i])
		}
		parts[i] = val
	}

	// Parse build number if present
	if len(mainAndBuild) > 1 {
		build, err := strconv.Atoi(mainAndBuild[1])
		if err != nil {
			return parts, fmt.Errorf("invalid build number: %s", mainAndBuild[1])
		}
		parts[3] = build
	}

	return parts, nil
}
