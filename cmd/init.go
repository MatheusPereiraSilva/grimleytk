package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const defaultConfigFile = "grimley.yaml"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new GrimleyTK project",
	Run: func(cmd *cobra.Command, args []string) {

		// 1. Check if grimley.yaml already exists
		if _, err := os.Stat(defaultConfigFile); err == nil {
			fmt.Printf("✖ %s already exists. Initialization aborted.\n", defaultConfigFile)
			os.Exit(1)
		}

		// 2. Create the file
		file, err := os.Create(defaultConfigFile)
		if err != nil {
			fmt.Printf("Failed to create %s: %v\n", defaultConfigFile, err)
			os.Exit(1)
		}
		defer file.Close()

		// 3. Write template
		_, err = file.WriteString(defaultConfigTemplate)
		if err != nil {
			fmt.Printf("Failed to write %s: %v\n", defaultConfigFile, err)
			os.Exit(1)
		}

		fmt.Printf("✔ %s created successfully\n", defaultConfigFile)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

const defaultConfigTemplate = `# GrimleyTK Architecture Definition
# This file defines the data architecture of your system.
# It is the single source of truth for GrimleyTK.

version: "0.1"

project:
  # Project name (required)
  name: my-project

  # Environment: local | staging | prod
  environment: local

database:
  # Supported engines: postgres
  engine: postgres

  host: localhost
  port: 5432

  # Database name
  name: my_database

  ssl: false

  credentials:
    # Database user
    user: postgres

    # Environment variable containing the database password
    password_env: DB_PASSWORD

domains:
  example:
    description: Example domain
    schema: example
    owner: example-service

    owns:
      tables:
        example_table:
          description: Example table
          columns:
            id:
              type: uuid
              primary_key: true
`
