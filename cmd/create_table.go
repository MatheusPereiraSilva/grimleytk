package cmd

import (
	"fmt"
	"os"
	"strings"

	"grimleytk/internal/config"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var tableDescription string

var createTableCmd = &cobra.Command{
	Use:   "create table <domain>.<table_name>",
	Short: "Create a table owned by a domain in the Grimley architecture file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Parse domain.table
		parts := strings.Split(args[0], ".")
		if len(parts) != 2 {
			fmt.Println("Invalid table name. Expected format: <domain>.<table_name>")
			os.Exit(1)
		}

		domainName := parts[0]
		tableName := parts[1]

		// Load config
		cfg, err := config.Load("grimley.yaml")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		domain, ok := cfg.Domains[domainName]
		if !ok {
			fmt.Printf("Domain '%s' does not exist\n", domainName)
			os.Exit(1)
		}

		if domain.Owns == nil {
			domain.Owns = &config.OwnedResources{
				Tables: map[string]config.Table{},
			}
		}

		if _, exists := domain.Owns.Tables[tableName]; exists {
			fmt.Printf("Table '%s.%s' already exists\n", domainName, tableName)
			os.Exit(1)
		}

		// Create table
		domain.Owns.Tables[tableName] = config.Table{
			Description: tableDescription,
			Columns:     map[string]config.Column{},
		}

		cfg.Domains[domainName] = domain

		// Write back to file
		data, err := yaml.Marshal(cfg)
		if err != nil {
			fmt.Printf("Failed to serialize config: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile("grimley.yaml", data, 0644); err != nil {
			fmt.Printf("Failed to write grimley.yaml: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✔ Table '%s.%s' created successfully\n", domainName, tableName)
	},
}

func init() {
	createTableCmd.Flags().StringVar(&tableDescription, "description", "", "Table description")

	rootCmd.AddCommand(createTableCmd)
}
