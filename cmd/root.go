package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

var rootCmd = &cobra.Command{
	Use:   "easy-deploy",
	Short: "Easy Deploy CLI – generate workflows and manage environments",
	Long:  "Helper for bootstrapping deployment workflows and managing EB environments.",

	Run: func(cmd *cobra.Command, args []string) {
		showBanner()
		cmd.Help()
	},
}

func showBanner() {
	c := color.New(color.FgCyan)
	myFigure := figure.NewFigure("AWS Easy Deploy", "slant", true)
	c.Println(myFigure.String())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
