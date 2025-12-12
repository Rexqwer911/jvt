package cli

import (
	"fmt"

	"github.com/rexqwer911/jvt/internal/config"
	"github.com/rexqwer911/jvt/internal/install"
	"github.com/rexqwer911/jvt/internal/registry"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List installed Java versions",
	Long:    "Display all Java versions currently installed on your system.",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		installer := install.NewInstaller(cfg.InstallDir)
		versions, err := installer.ListInstalled()
		if err != nil {
			return fmt.Errorf("failed to list installed versions: %w", err)
		}

		if len(versions) == 0 {
			fmt.Println("No Java versions installed.")
			fmt.Println("Use 'jvt install <version>' to install a version.")
			return nil
		}

		fmt.Println("Installed Java versions:")
		for _, v := range versions {
			fmt.Printf("  * %s\n", v)
		}

		return nil
	},
}

var listRemoteCmd = &cobra.Command{
	Use:     "list-remote",
	Short:   "List available Java versions for download",
	Long:    "Display all Java versions available for download from configured sources.",
	Aliases: []string{"ls-remote"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Fetching available Java versions from Adoptium...")

		reg := registry.NewRegistry()
		if err := reg.FetchAvailableVersions(); err != nil {
			return fmt.Errorf("failed to fetch versions: %w", err)
		}

		versions := reg.GetVersions()
		if len(versions) == 0 {
			fmt.Println("No versions found.")
			return nil
		}

		fmt.Println("\nAvailable Java versions (Temurin/OpenJDK):")
		fmt.Println("Major | Full Version")
		fmt.Println("------|-------------")

		currentMajor := -1
		for _, v := range versions {
			if v.MajorVersion != currentMajor {
				fmt.Printf("  %2d  | %s\n", v.MajorVersion, v.Version)
				currentMajor = v.MajorVersion
			}
		}

		fmt.Println("\nUse 'jvt install <major-version>' to install (e.g., 'jvt install 17')")

		return nil
	},
}
