package cmd

import (
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show GrimleyTK architecture information",
	Long:  "Display different views of the declared data architecture.",
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior: show domains
		showDomains(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
