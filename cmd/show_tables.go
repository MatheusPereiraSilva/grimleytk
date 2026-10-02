package cmd

import (
	"fmt"

	"github.com/MatheusPereiraSilva/grimleytk/internal/config"

	"github.com/spf13/cobra"
)

var showTablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "Show owned tables per domain",
	RunE: func(cmd *cobra.Command, args []string) error {

		cfg, err := config.Load("grimley.yaml")
		if err != nil {
			return err
		}

		fmt.Println("Tables:")
		for domainName, domain := range cfg.Domains {
			if domain.Owns == nil || len(domain.Owns.Tables) == 0 {
				continue
			}

			for tableName, table := range domain.Owns.Tables {
				fmt.Printf("- %s.%s\n", domainName, tableName)
				if table.Description != "" {
					fmt.Printf("  Description: %s\n", table.Description)
				}
				fmt.Printf("  Columns: %d\n\n", len(table.Columns))
			}
		}
		return nil
	},
}

func init() {
	showCmd.AddCommand(showTablesCmd)
}
