package cmd

import (
	"fmt"
	"os"

	"grimleytk/internal/config"
	"grimleytk/internal/planner"
	"grimleytk/internal/validator"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate a database execution plan (dry-run)",
	Long: `Generate a list of SQL statements required to apply
the declared architecture. No SQL is executed.`,
	Run: func(cmd *cobra.Command, args []string) {

		// 1. Load config
		cfg, err := config.Load("grimley.yaml")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 2. Validate architecture before planning
		var issues []validator.Issue
		issues = append(issues, validator.ValidateStructural(cfg)...)
		issues = append(issues, validator.ValidateReferences(cfg)...)
		issues = append(issues, validator.ValidateArchitecture(cfg)...)
		issues = append(issues, validator.ValidateSecurity(cfg)...)

		report := validator.BuildReport(issues)
		if report.HasErrors() {
			fmt.Println(report.String())
			os.Exit(1)
		}

		// 3. Build plan
		actions := planner.BuildPlan(cfg)

		if len(actions) == 0 {
			fmt.Println("No actions to perform.")
			return
		}

		// 4. Print plan
		fmt.Println("Execution Plan:\n")

		for i, action := range actions {
			fmt.Printf("%d. %s\n", i+1, action.Description)
			fmt.Println(action.SQL)
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
