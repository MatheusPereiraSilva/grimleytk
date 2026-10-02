package cmd

import (
	"fmt"
	"os"

	"github.com/MatheusPereiraSilva/grimleytk/internal/config"
	"github.com/MatheusPereiraSilva/grimleytk/internal/validator"

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

		issues := validator.Validate(cfg)

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
