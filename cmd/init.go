package cmd

import (
	"easy-deploy/internal/ui"
	"easy-deploy/internal/workflows"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Gather project details and emit CI workflow",
	Long: `Interactively collects framework/branch information, updates the
.easy-deploy configuration file and writes a CI workflow (GitHub or GitLab)
into the repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		ui.RunAppDetails(ctx)

		if err := workflows.RunSetupWorkflow(); err != nil {
			fmt.Println("❌ setup failed:", err)
			os.Exit(1)
		}
		fmt.Println("✅ workflow created successfully")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
