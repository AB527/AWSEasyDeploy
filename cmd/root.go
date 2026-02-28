package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "easy-deploy",
	Short: "Easy Deploy CLI – generate workflows and manage environments",
	Long:  "Helper for bootstrapping deployment workflows and managing EB environments.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
