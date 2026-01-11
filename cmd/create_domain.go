package cmd

import (
	"fmt"
	"os"

	"grimleytk/internal/config"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	domainSchema      string
	domainOwner       string
	domainDescription string
)

var createDomainCmd = &cobra.Command{
	Use:   "create domain <domain_name>",
	Short: "Create a new domain in the Grimley architecture file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		domainName := args[0]

		if domainSchema == "" {
			fmt.Println("Missing required flag --schema <schema_name>")
			os.Exit(1)
		}

		if domainOwner == "" {
			fmt.Println("Missing required flag --owner <owner_name>")
			os.Exit(1)
		}

		// Load config
		cfg, err := config.Load("grimley.yaml")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if _, exists := cfg.Domains[domainName]; exists {
			fmt.Printf("Domain '%s' already exists\n", domainName)
			os.Exit(1)
		}

		// Create domain
		cfg.Domains[domainName] = config.Domain{
			Description: domainDescription,
			Schema:      domainSchema,
			Owner:       domainOwner,
			Owns: &config.OwnedResources{
				Tables: map[string]config.Table{},
			},
			Reads: map[string]config.Read{},
		}

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

		fmt.Printf("✔ Domain '%s' created successfully\n", domainName)
	},
}

func init() {
	createDomainCmd.Flags().StringVar(&domainSchema, "schema", "", "Database schema for the domain")
	createDomainCmd.Flags().StringVar(&domainOwner, "owner", "", "Owning service or team")
	createDomainCmd.Flags().StringVar(&domainDescription, "description", "", "Domain description")

	rootCmd.AddCommand(createDomainCmd)
}
