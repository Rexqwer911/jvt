//go:build linux || darwin

package version

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SetEnvironment sets JAVA_HOME and updates PATH for the current session
func (m *Manager) SetEnvironment(version string) error {
	// On Unix, the shell wrapper function defined in .bashrc/.zshrc handles the immediate update.
	// If the user hasn't restarted their shell since installing, they might need to source manually once.
	return nil
}

// SetUserEnvironment sets JAVA_HOME and PATH in user environment variables (persistent)
func (m *Manager) SetUserEnvironment(version string) error {
	javaHome := filepath.Join(m.installDir, version)

	// Update shell config files
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	rcFiles := []string{".zshrc", ".bashrc", ".profile", ".bash_profile"}
	updated := false

	// Ensure ~/.jvt/jvt.sh exists and is sourced
	jvtScriptPath := filepath.Join(home, ".jvt", "jvt.sh")
	jvtScriptContent := fmt.Sprintf("export JAVA_HOME=\"%s\"\nexport PATH=\"$JAVA_HOME/bin:$PATH\"\n", javaHome)

	if err := os.MkdirAll(filepath.Dir(jvtScriptPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(jvtScriptPath, []byte(jvtScriptContent), 0644); err != nil {
		return err
	}

	for _, rc := range rcFiles {
		rcPath := filepath.Join(home, rc)
		if _, err := os.Stat(rcPath); err == nil {
			// File exists, check if it sources jvt.sh AND has the wrapper
			content, err := os.ReadFile(rcPath)
			if err == nil {
				strContent := string(content)
				// Check if we need to update:
				// 1. Missing the basic source line?
				// 2. Missing the function wrapper?
				needsUpdate := !strings.Contains(strContent, ".jvt/jvt.sh") || !strings.Contains(strContent, "jvt() {")

				if needsUpdate {
					// Append source command and wrapper
					f, err := os.OpenFile(rcPath, os.O_APPEND|os.O_WRONLY, 0644)
					if err == nil {
						if !strings.HasSuffix(strContent, "\n") {
							f.WriteString("\n")
						}

						// If the file already has the source line but not the wrapper, we might be duplicating the source line
						// inside the script block. That's acceptable (idempotent-ish).
						// To be cleaner, we could check, but appending the whole block is safer ensuring consistency.

						// Inject shell function wrapper to allow auto-sourcing
						script := `
# JVT Java Version Tool
export JVT_HOME="$HOME/.jvt"
[ -s "$JVT_HOME/jvt.sh" ] && . "$JVT_HOME/jvt.sh"

jvt() {
    command jvt "$@"
    local exit_code=$?
    if [ $exit_code -eq 0 ] && [ "$1" = "use" ]; then
        [ -s "$HOME/.jvt/jvt.sh" ] && . "$HOME/.jvt/jvt.sh"
    fi
    return $exit_code
}
`
						f.WriteString(script)
						f.Close()
						fmt.Printf("Updated %s with shell wrapper\n", rc)
						updated = true
					}
				} else {
					updated = true // Already configured
				}
			}
		}
	}

	if !updated {
		fmt.Println("Could not find common shell configuration files (.zshrc, .bashrc, etc.).")
		fmt.Println("Please add the following execution to your shell startup script:")
		fmt.Printf(`
export JVT_HOME="$HOME/.jvt"
[ -s "$JVT_HOME/jvt.sh" ] && . "$JVT_HOME/jvt.sh"

jvt() {
    command jvt "$@"
    local exit_code=$?
    if [ $exit_code -eq 0 ] && [ "$1" = "use" ]; then
        [ -s "$HOME/.jvt/jvt.sh" ] && . "$HOME/.jvt/jvt.sh"
    fi
    return $exit_code
}
`)
	}

	fmt.Printf("Java %s configured.\n", version)
	return nil
}

// SetSystemEnvironment sets JAVA_HOME and PATH in SYSTEM environment variables (persistent)
func (m *Manager) SetSystemEnvironment(version string) error {
	return fmt.Errorf("system-wide configuration not supported on Unix yet (requires sudo)")
}

// checkSystemJava checks if there's a system-level Java installation
func (m *Manager) checkSystemJava() {
	// Not implemented for Unix yet
}

func (m *Manager) updatePathString(currentPath, javaBin string) string {
	// Not used in Unix implementation as we use sourcing
	return ""
}
