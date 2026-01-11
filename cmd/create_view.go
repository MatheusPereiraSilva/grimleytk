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
	fromRef string
	columns string
)

var createViewCmd = &cobra.Command{
	Use:   "create view <domain>.<view_name>",
	Short: "Create a read model (view) in the Grimley architecture file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Parse domain.view
		parts := strings.Split(args[0], ".")
		if len(parts) != 2 {
			fmt.Println("Invalid view name. Expected format: <domain>.<view_name>")
			os.Exit(1)
		}

		domainName := parts[0]
		viewName := parts[1]

		if fromRef == "" {
			fmt.Println("Missing required flag --from <domain.table>")
			os.Exit(1)
		}

		if columns == "" {
			fmt.Println("Missing required flag --columns col1,col2,col3")
			os.Exit(1)
		}

		cols := strings.Split(columns, ",")

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

		if domain.Reads == nil {
			domain.Reads = make(map[string]config.Read)
		}

		// Create / overwrite view
		domain.Reads[viewName] = config.Read{
			From:    fromRef,
			Columns: cols,
			Access: config.AccessMode{
				Mode: "read-only",
			},
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

		fmt.Printf("✔ View '%s.%s' created successfully\n", domainName, viewName)
	},
}

func init() {
	createViewCmd.Flags().StringVar(&fromRef, "from", "", "Source table (domain.table)")
	createViewCmd.Flags().StringVar(&columns, "columns", "", "Comma-separated list of columns")

	rootCmd.AddCommand(createViewCmd)
}
