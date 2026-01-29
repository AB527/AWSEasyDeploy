package cmd

import (
	"easy-deploy-cli/internal/workflows"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your deployment project",
	Long:  "Interactive initialization of your deployment project with framework, repo, branch, and VPC setup.",
	Run: func(cmd *cobra.Command, args []string) {
		err := workflows.RunSetupWorkflow()
		if err != nil {
			fmt.Println("❌ setup failed:", err)
			os.Exit(1)
		}
		// ctx := context.Background()
		// ui.RunAppDetails(ctx)

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
