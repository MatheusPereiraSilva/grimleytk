package cmd

import (
	"fmt"

	"grimleytk/internal/config"

	"github.com/spf13/cobra"
)

var showReadsCmd = &cobra.Command{
	Use:   "reads",
	Short: "Show read models (views)",
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.Load("grimley.yaml")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Read Models:")
		for domainName, domain := range cfg.Domains {
			for readName, read := range domain.Reads {
				fmt.Printf("- %s.%s\n", domainName, readName)
				fmt.Printf("  From   : %s\n", read.From)
				fmt.Printf("  Columns: %v\n\n", read.Columns)
			}
		}
	},
}

func init() {
	showCmd.AddCommand(showReadsCmd)
}
