package cmd

import (
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show GrimleyTK architecture information",
	Long:  "Display different views of the declared data architecture.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default behavior: show domains
		return showDomains(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
