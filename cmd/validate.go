package cmd

import (
	"fmt"
	"os"

	"grimleytk/internal/config"
	"grimleytk/internal/validator"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Grimley architecture definition",
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.Load("grimley.yaml")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var issues []validator.Issue

		issues = append(issues, validator.ValidateStructural(cfg)...)
		issues = append(issues, validator.ValidateReferences(cfg)...)
		issues = append(issues, validator.ValidateArchitecture(cfg)...)
		issues = append(issues, validator.ValidateSecurity(cfg)...)

		report := validator.BuildReport(issues)
		fmt.Println(report.String())

		if report.HasErrors() {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
