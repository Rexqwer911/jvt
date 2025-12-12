package cli

import (
	"fmt"

	"github.com/rexqwer911/jvt/internal/config"
	"github.com/rexqwer911/jvt/internal/download"
	"github.com/rexqwer911/jvt/internal/install"
	"github.com/rexqwer911/jvt/internal/registry"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Install a specific Java version",
	Long:  "Download and install a specific Java version.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		versionStr := args[0]

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		if err := cfg.EnsureDirectories(); err != nil {
			return fmt.Errorf("failed to create directories: %w", err)
		}

		// Fetch available versions
		fmt.Println("Fetching available versions...")
		reg := registry.NewRegistry()
		if err := reg.FetchAvailableVersions(); err != nil {
			return fmt.Errorf("failed to fetch versions: %w", err)
		}

		// Find the requested version
		javaVersion, err := reg.FindVersion(versionStr)
		if err != nil {
			return fmt.Errorf("version not found: %w", err)
		}

		fmt.Printf("\nFound: Java %s (%s)\n", javaVersion.Version, javaVersion.Distribution)

		// Check if already installed
		installer := install.NewInstaller(cfg.InstallDir)
		if installer.IsInstalled(javaVersion.Version) {
			fmt.Printf("Java %s is already installed.\n", javaVersion.Version)
			return nil
		}

		// Download
		downloader := download.NewDownloader(cfg.CacheDir)
		fmt.Printf("\nDownloading from: %s\n", javaVersion.DownloadURL)

		archivePath, err := downloader.DownloadAndVerify(
			javaVersion.DownloadURL,
			javaVersion.FileName,
			javaVersion.Checksum,
		)
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		// Install
		fmt.Println("\nInstalling...")
		if err := installer.Install(archivePath, javaVersion.Version); err != nil {
			return fmt.Errorf("installation failed: %w", err)
		}

		fmt.Printf("\n✓ Java %s installed successfully!\n", javaVersion.Version)
		fmt.Printf("Run 'jvt use %d' to activate this version.\n", javaVersion.MajorVersion)

		return nil
	},
}
