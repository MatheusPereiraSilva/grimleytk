package cmd

import (
	"fmt"
	"os"
	"strings"

	"grimleytk/internal/config"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	columnType       string
	columnNullable   bool
	columnPrimaryKey bool
	columnUnique     bool
)

var createColumnCmd = &cobra.Command{
	Use:   "create column <domain>.<table>.<column>",
	Short: "Create a column in a table owned by a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Parse domain.table.column
		parts := strings.Split(args[0], ".")
		if len(parts) != 3 {
			fmt.Println("Invalid column name. Expected format: <domain>.<table>.<column>")
			os.Exit(1)
		}

		domainName := parts[0]
		tableName := parts[1]
		columnName := parts[2]

		if columnType == "" {
			fmt.Println("Missing required flag --type <sql_type>")
			os.Exit(1)
		}

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

		if domain.Owns == nil || domain.Owns.Tables == nil {
			fmt.Printf("Domain '%s' does not own any tables\n", domainName)
			os.Exit(1)
		}

		table, ok := domain.Owns.Tables[tableName]
		if !ok {
			fmt.Printf("Table '%s.%s' does not exist\n", domainName, tableName)
			os.Exit(1)
		}

		if _, exists := table.Columns[columnName]; exists {
			fmt.Printf("Column '%s.%s.%s' already exists\n", domainName, tableName, columnName)
			os.Exit(1)
		}

		// Create column
		table.Columns[columnName] = config.Column{
			Type:       columnType,
			Nullable:   columnNullable,
			PrimaryKey: columnPrimaryKey,
			Unique:     columnUnique,
		}

		domain.Owns.Tables[tableName] = table
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

		fmt.Printf(
			"✔ Column '%s.%s.%s' created successfully\n",
			domainName,
			tableName,
			columnName,
		)
	},
}

func init() {
	createColumnCmd.Flags().StringVar(&columnType, "type", "", "SQL column type (e.g. uuid, text, numeric)")
	createColumnCmd.Flags().BoolVar(&columnNullable, "nullable", true, "Whether the column is nullable")
	createColumnCmd.Flags().BoolVar(&columnPrimaryKey, "primary-key", false, "Whether the column is a primary key")
	createColumnCmd.Flags().BoolVar(&columnUnique, "unique", false, "Whether the column has a unique constraint")

	rootCmd.AddCommand(createColumnCmd)
}
