package cli

import (
	"fmt"
	"strconv"

	"github.com/rexqwer911/jvt/internal/config"
	"github.com/rexqwer911/jvt/internal/install"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:     "uninstall <version>",
	Short:   "Uninstall a specific Java version",
	Long:    "Remove a specific Java version from your system.",
	Aliases: []string{"remove", "rm"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		versionStr := args[0]

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		installer := install.NewInstaller(cfg.InstallDir)
		versions, err := installer.ListInstalled()
		if err != nil {
			return fmt.Errorf("failed to list installed versions: %w", err)
		}

		// Find matching version
		var matchedVersion string
		if _, err := strconv.Atoi(versionStr); err == nil {
			for _, v := range versions {
				if len(v) >= len(versionStr) && v[0:len(versionStr)] == versionStr {
					matchedVersion = v
					break
				}
			}
		} else {
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

		// Confirm and uninstall
		fmt.Printf("Uninstalling Java %s...\n", matchedVersion)
		if err := installer.Uninstall(matchedVersion); err != nil {
			return fmt.Errorf("uninstall failed: %w", err)
		}

		fmt.Printf("✓ Java %s uninstalled successfully!\n", matchedVersion)

		return nil
	},
}
