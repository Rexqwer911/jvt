package version

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
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

// SetEnvironment sets JAVA_HOME and updates PATH for the current session
func (m *Manager) SetEnvironment(version string) error {
	javaHome := filepath.Join(m.installDir, version)

	// Verify installation exists
	if _, err := os.Stat(javaHome); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed", version)
	}

	javaBin := filepath.Join(javaHome, "bin")

	// Set JAVA_HOME for current process
	if err := os.Setenv("JAVA_HOME", javaHome); err != nil {
		return fmt.Errorf("failed to set JAVA_HOME: %w", err)
	}

	// Update PATH for current process
	currentPath := os.Getenv("PATH")

	// Remove any existing Java paths
	pathParts := strings.Split(currentPath, ";")
	var newPathParts []string
	for _, part := range pathParts {
		if !strings.Contains(strings.ToLower(part), "java") || strings.Contains(part, m.installDir) {
			if !strings.Contains(part, m.installDir) {
				newPathParts = append(newPathParts, part)
			}
		}
	}

	// Add new Java bin to the front
	newPathParts = append([]string{javaBin}, newPathParts...)
	newPath := strings.Join(newPathParts, ";")

	if err := os.Setenv("PATH", newPath); err != nil {
		return fmt.Errorf("failed to set PATH: %w", err)
	}

	return nil
}

// SetUserEnvironment sets JAVA_HOME and PATH in user environment variables (persistent)
func (m *Manager) SetUserEnvironment(version string) error {
	javaHome := filepath.Join(m.installDir, version)

	// Verify installation exists
	if _, err := os.Stat(javaHome); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed", version)
	}

	javaBin := filepath.Join(javaHome, "bin")

	// Open user environment registry key
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to open registry: %w", err)
	}
	defer key.Close()

	// Set JAVA_HOME
	if err := key.SetStringValue("JAVA_HOME", javaHome); err != nil {
		return fmt.Errorf("failed to set JAVA_HOME: %w", err)
	}

	// Update PATH
	currentPath, _, err := key.GetStringValue("Path")
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("failed to read PATH: %w", err)
	}

	// Remove ALL Java-related paths (not just jvt)
	pathParts := strings.Split(currentPath, ";")
	var newPathParts []string
	var removedPaths []string

	for _, part := range pathParts {
		partLower := strings.ToLower(part)
		// Remove if contains java, jdk, jre, or adoptium
		if strings.Contains(partLower, "java") ||
			strings.Contains(partLower, "jdk") ||
			strings.Contains(partLower, "jre") ||
			strings.Contains(partLower, "adoptium") ||
			strings.Contains(partLower, "temurin") {
			removedPaths = append(removedPaths, part)
		} else if part != "" {
			newPathParts = append(newPathParts, part)
		}
	}

	// Add new Java bin to the BEGINNING
	newPathParts = append([]string{javaBin}, newPathParts...)
	newPath := strings.Join(newPathParts, ";")

	if err := key.SetStringValue("Path", newPath); err != nil {
		return fmt.Errorf("failed to set PATH: %w", err)
	}

	// Show what was removed
	if len(removedPaths) > 0 {
		fmt.Println("\nRemoved the following Java paths from user PATH:")
		for _, p := range removedPaths {
			fmt.Printf("  - %s\n", p)
		}
	}

	// Check for system-level Java
	m.checkSystemJava()

	// Broadcast environment change
	m.broadcastEnvironmentChange()

	return nil
}

// checkSystemJava checks if there's a system-level Java installation
func (m *Manager) checkSystemJava() {
	// Try to open system environment key (read-only)
	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
		registry.QUERY_VALUE)
	if err != nil {
		return // Can't read system registry, skip check
	}
	defer key.Close()

	// Check system PATH
	systemPath, _, err := key.GetStringValue("Path")
	if err != nil {
		return
	}

	// Look for Java in system PATH
	pathParts := strings.Split(systemPath, ";")
	var systemJavaPaths []string

	for _, part := range pathParts {
		partLower := strings.ToLower(part)
		if strings.Contains(partLower, "java") ||
			strings.Contains(partLower, "jdk") ||
			strings.Contains(partLower, "jre") ||
			strings.Contains(partLower, "adoptium") ||
			strings.Contains(partLower, "temurin") {
			systemJavaPaths = append(systemJavaPaths, part)
		}
	}

	if len(systemJavaPaths) > 0 {
		fmt.Println("\n⚠️  WARNING: System-level Java installation detected!")
		fmt.Println("The following Java paths are in SYSTEM PATH (requires admin to remove):")
		for _, p := range systemJavaPaths {
			fmt.Printf("  - %s\n", p)
		}
		fmt.Println("\nThese may override your jvt-managed Java version.")
		fmt.Println("To fix this, you need to:")
		fmt.Println("  1. Open 'System Properties' → 'Environment Variables' (as Administrator)")
		fmt.Println("  2. In 'System variables', edit 'Path'")
		fmt.Println("  3. Remove the Java-related entries listed above")
		fmt.Println("  4. Restart your terminal")
	}
}

// broadcastEnvironmentChange notifies Windows of environment variable changes
func (m *Manager) broadcastEnvironmentChange() {
	// This would require syscall to SendMessageTimeout
	// For now, we'll just inform the user to restart their terminal
	fmt.Println("\nNote: Please restart your terminal or command prompt for changes to take effect.")
}

// GetCurrentVersion returns the currently active Java version
func (m *Manager) GetCurrentVersion() (string, error) {
	javaHome := os.Getenv("JAVA_HOME")

	if javaHome == "" {
		return "", fmt.Errorf("JAVA_HOME is not set")
	}

	// Check if it's a jvt-managed version
	if !strings.HasPrefix(javaHome, m.installDir) {
		return "", fmt.Errorf("current Java is not managed by jvt")
	}

	version := filepath.Base(javaHome)
	return version, nil
}
