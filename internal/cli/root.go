package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jvt",
	Short: "Java Version Tool - Manage multiple Java installations",
	Long: `JVT (Java Version Tool) is a command-line utility that helps you
download, install, and switch between different Java versions easily.

Similar to nvm for Node.js, jvt simplifies Java version management.`,
	Version: "1.3.0",
}

// Execute runs the root command
func Execute() error {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(listRemoteCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(upgradeCmd)
}
