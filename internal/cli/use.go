package cli

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/rexqwer911/jvt/internal/config"
	"github.com/rexqwer911/jvt/internal/install"
	"github.com/rexqwer911/jvt/internal/version"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <version>",
	Short: "Switch to a specific Java version",
	Long:  "Switch the active Java version definitively (updates both current session and system defaults).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		versionStr := args[0]

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		// Find installed version matching the input
		installer := install.NewInstaller(cfg.InstallDir)
		versions, err := installer.ListInstalled()
		if err != nil {
			return fmt.Errorf("failed to list installed versions: %w", err)
		}

		if len(versions) == 0 {
			return fmt.Errorf("no Java versions installed. Use 'jvt install <version>' first")
		}

		// Try to find matching version
		var matchedVersion string
		if _, err := strconv.Atoi(versionStr); err == nil {
			// Search by major version
			for _, v := range versions {
				if len(v) >= len(versionStr) && v[0:len(versionStr)] == versionStr {
					matchedVersion = v
					break
				}
			}
		} else {
			// Exact match
			for _, v := range versions {
				if v == versionStr {
					matchedVersion = v
					break
				}
			}
		}

		if matchedVersion == "" {
			return fmt.Errorf("version %s is not installed", versionStr)
		}

		// Manage environment
		mgr := version.NewManager(cfg.InstallDir)

		// 1. Set Persistent Environment (Registry)
		// First try User environment (always should succeed)
		if err := mgr.SetUserEnvironment(matchedVersion); err != nil {
			return fmt.Errorf("failed to set user environment: %w", err)
		}

		// Then try System environment (Windows only)
		if runtime.GOOS == "windows" {
			if err := mgr.SetSystemEnvironment(matchedVersion); err != nil {
				// Check if it's likely a permission error
				// On Windows, syscall.ERROR_ACCESS_DENIED is 5
				fmt.Printf("\nNote: Could not update System environment variables (requires Administrator).\n")
				fmt.Printf("   Reason: %v\n", err)
				fmt.Println("   Only User environment variables were updated.")
			} else {
				fmt.Println("✓ System environment variables updated.")
			}
		}

		// 2. Set Current Session Environment
		if err := mgr.SetEnvironment(matchedVersion); err != nil {
			// Warn but don't fail if session update fails (e.g. maybe restricted)
			// But usually it should work if registry worked?
			// Actually failure here is annoying for the user.
			fmt.Printf("Warning: failed to set current session environment: %v\n", err)
		}

		fmt.Printf("✓ Now using Java %s\n", matchedVersion)
		return nil
	},
}

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the currently active Java version",
	Long:  "Display which Java version is currently active.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		mgr := version.NewManager(cfg.InstallDir)
		currentVersion, err := mgr.GetCurrentVersion()
		if err != nil {
			fmt.Println("No jvt-managed Java version is currently active.")
			fmt.Println("Use 'jvt use <version>' to activate a version.")
			return nil
		}

		fmt.Printf("Current Java version: %s\n", currentVersion)
		return nil
	},
}
