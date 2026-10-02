package cmd

import (
	"fmt"
	"os"

	"github.com/MatheusPereiraSilva/grimleytk/internal/config"
	"github.com/MatheusPereiraSilva/grimleytk/internal/planner"
	"github.com/MatheusPereiraSilva/grimleytk/internal/validator"

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
		issues := validator.Validate(cfg)

		report := validator.BuildReport(issues)
		if report.HasErrors() {
			fmt.Println(report.String())
			os.Exit(1)
		}

		// 3. Build plan
		actions, err := planner.BuildPlan(cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(actions) == 0 {
			fmt.Println("No actions to perform.")
			return
		}

		// 4. Print plan
		fmt.Println("Execution Plan:")

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
