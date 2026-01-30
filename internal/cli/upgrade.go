package cli

import (
	"fmt"
	"strconv"

	"github.com/rexqwer911/jvt/internal/config"
	"github.com/rexqwer911/jvt/internal/download"
	"github.com/rexqwer911/jvt/internal/install"
	"github.com/rexqwer911/jvt/internal/registry"
	"github.com/rexqwer911/jvt/internal/version"
	"github.com/spf13/cobra"
)

var (
	upgradeAll     bool
	upgradeDryRun  bool
	upgradeKeepOld bool
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [major-version]",
	Short: "Upgrade Java to the latest version",
	Long: `Upgrade installed Java versions to the latest available version.

Examples:
  jvt upgrade 17              # Upgrade Java 17 to latest
  jvt upgrade --all           # Upgrade all installed versions
  jvt upgrade --all --dry-run # Check for updates without installing`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		installer := install.NewInstaller(cfg.InstallDir)

		// Handle --all flag
		if upgradeAll {
			return upgradeAllVersions(cfg, installer)
		}

		// Require version argument if not using --all
		if len(args) == 0 {
			return fmt.Errorf("please specify a major version or use --all flag")
		}

		majorVersion, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid major version: %s", args[0])
		}

		return upgradeVersion(cfg, installer, majorVersion)
	},
}

func init() {
	upgradeCmd.Flags().BoolVar(&upgradeAll, "all", false, "Upgrade all installed Java versions")
	upgradeCmd.Flags().BoolVar(&upgradeDryRun, "dry-run", false, "Check for updates without installing")
	upgradeCmd.Flags().BoolVar(&upgradeKeepOld, "keep-old", false, "Keep old version after upgrade")
}

// upgradeAllVersions upgrades all installed major versions
func upgradeAllVersions(cfg *config.Config, installer *install.Installer) error {
	majorVersions, err := installer.GetInstalledMajorVersions()
	if err != nil {
		return fmt.Errorf("failed to get installed versions: %w", err)
	}

	if len(majorVersions) == 0 {
		fmt.Println("No Java versions installed.")
		return nil
	}

	if upgradeDryRun {
		fmt.Println("Checking for Java updates...")
		fmt.Println()
	}

	hasUpdates := false
	var updateCount int
	var upToDateCount int

	for _, major := range majorVersions {
		result, err := checkAndUpgradeVersion(cfg, installer, major)
		if err != nil {
			fmt.Printf("Error checking Java %d: %v\n", major, err)
			continue
		}

		if result == "updated" {
			updateCount++
			hasUpdates = true
		} else if result == "available" {
			hasUpdates = true
		} else if result == "up-to-date" {
			upToDateCount++
		}
	}

	if upgradeDryRun {
		fmt.Println()
		if hasUpdates {
			fmt.Println("Run 'jvt upgrade --all' to install updates.")
		} else {
			fmt.Println("All Java versions are up to date!")
		}
	} else {
		fmt.Println()
		if updateCount > 0 {
			fmt.Printf("✓ Successfully upgraded %d Java version(s)!\n", updateCount)
		}
		if upToDateCount > 0 {
			fmt.Printf("  %d version(s) already up to date.\n", upToDateCount)
		}
	}

	return nil
}

// upgradeVersion upgrades a specific major version
func upgradeVersion(cfg *config.Config, installer *install.Installer, majorVersion int) error {
	result, err := checkAndUpgradeVersion(cfg, installer, majorVersion)
	if err != nil {
		return err
	}

	if result == "not-installed" {
		return fmt.Errorf("Java %d is not installed. Use 'jvt install %d' first", majorVersion, majorVersion)
	}

	return nil
}

// checkAndUpgradeVersion checks and optionally upgrades a single major version
// Returns: "updated", "available", "up-to-date", "not-installed", "error"
func checkAndUpgradeVersion(cfg *config.Config, installer *install.Installer, majorVersion int) (string, error) {
	// Get installed versions for this major version
	installedVersions, err := installer.GetInstalledByMajor(majorVersion)
	if err != nil {
		return "error", fmt.Errorf("failed to get installed versions: %w", err)
	}

	if len(installedVersions) == 0 {
		return "not-installed", nil
	}

	// Find the newest installed version
	newestInstalled := installedVersions[0]
	for _, v := range installedVersions {
		cmp, err := version.CompareVersions(v, newestInstalled)
		if err != nil {
			continue
		}
		if cmp > 0 {
			newestInstalled = v
		}
	}

	// Fetch latest available version
	reg := registry.NewRegistry()
	if !upgradeDryRun {
		fmt.Printf("Checking for Java %d updates...\n", majorVersion)
	}

	latestAvailable, err := reg.FindLatestForMajor(majorVersion)
	if err != nil {
		return "error", fmt.Errorf("failed to fetch latest version: %w", err)
	}

	// Compare versions
	cmp, err := version.CompareVersions(newestInstalled, latestAvailable.Version)
	if err != nil {
		return "error", fmt.Errorf("failed to compare versions: %w", err)
	}

	if cmp >= 0 {
		// Already up to date
		if upgradeDryRun || upgradeAll {
			fmt.Printf("Java %d is up to date (%s)\n", majorVersion, newestInstalled)
		} else {
			fmt.Printf("Java %d is already up to date (%s)\n", majorVersion, newestInstalled)
		}
		return "up-to-date", nil
	}

	// Update available
	if upgradeDryRun {
		fmt.Printf("Java %d: %s → %s (update available)\n", majorVersion, newestInstalled, latestAvailable.Version)
		return "available", nil
	}

	// Perform upgrade
	fmt.Printf("\nUpgrading Java %d...\n", majorVersion)
	fmt.Printf("  Current version: %s\n", newestInstalled)
	fmt.Printf("  Latest version:  %s\n", latestAvailable.Version)
	fmt.Println()

	// Check if current version is active
	mgr := version.NewManager(cfg.InstallDir)
	isActive, err := mgr.IsVersionActive(newestInstalled)
	if err != nil {
		fmt.Printf("Warning: Failed to determine if current version %s is active: %v\n", newestInstalled, err)
		isActive = false
	}

	// Download
	downloader := download.NewDownloader(cfg.CacheDir)
	fmt.Printf("Downloading from: %s\n", latestAvailable.DownloadURL)

	archivePath, err := downloader.DownloadAndVerify(
		latestAvailable.DownloadURL,
		latestAvailable.FileName,
		latestAvailable.Checksum,
	)
	if err != nil {
		return "error", fmt.Errorf("download failed: %w", err)
	}

	// Install
	fmt.Println("\nInstalling...")
	if err := installer.Install(archivePath, latestAvailable.Version); err != nil {
		return "error", fmt.Errorf("installation failed: %w", err)
	}

	// If the old version was active, set environment to the new version
	if isActive {
		if err := mgr.SetUserEnvironment(latestAvailable.Version); err != nil {
			fmt.Printf("Warning: Failed to set environment: %v\n", err)
		} else {
			fmt.Println("✓ Java version updated")
		}
	}

	// Remove old version unless --keep-old
	if !upgradeKeepOld {
		fmt.Printf("Removing old version %s...\n", newestInstalled)
		if err := installer.Uninstall(newestInstalled); err != nil {
			fmt.Printf("Warning: Failed to remove old version: %v\n", err)
			fmt.Printf("You can manually remove it with: jvt uninstall %s\n", newestInstalled)
		} else {
			fmt.Println("✓ Old version removed")
		}
	}

	fmt.Printf("\n✓ Java %d upgraded successfully! (%s → %s)\n", majorVersion, newestInstalled, latestAvailable.Version)

	if isActive {
		fmt.Println("\nPlease restart your terminal or run:")
		fmt.Println("  source ~/.bashrc  (Linux/macOS)")
		fmt.Println("  or open a new terminal window")
	}

	return "updated", nil
}
