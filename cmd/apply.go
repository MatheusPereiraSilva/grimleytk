package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MatheusPereiraSilva/grimleytk/internal/config"
	"github.com/MatheusPereiraSilva/grimleytk/internal/executor"
	"github.com/MatheusPereiraSilva/grimleytk/internal/planner"
	"github.com/MatheusPereiraSilva/grimleytk/internal/validator"

	"github.com/spf13/cobra"
)

var autoApprove bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the planned database changes",
	Long: `Apply executes the SQL generated from the GrimleyTK plan.
This operation modifies the database and requires confirmation.`,
	Args: cobra.NoArgs,
	Example: `  # Run from the directory containing grimley.yaml:
  grimleytk apply
  grimleytk apply --auto-approve`,
	Run: func(cmd *cobra.Command, args []string) {

		// 1. Load config
		cfg, err := config.Load("grimley.yaml")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 2. Validate before apply
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
			fmt.Println("No actions to apply.")
			return
		}

		// 4. Show plan
		fmt.Println("Execution Plan:")
		for i, action := range actions {
			fmt.Printf("%d. %s\n", i+1, action.Description)
			fmt.Println(action.SQL)
			fmt.Println()
		}

		// 5. Confirmation
		if !autoApprove {
			if !askForConfirmation() {
				fmt.Println("Apply aborted.")
				return
			}
		}

		// 6. Create executor
		exec, err := executor.NewPostgresExecutor(cfg.Database)
		if err != nil {
			fmt.Printf("Failed to initialize executor: %v\n", err)
			os.Exit(1)
		}

		defer exec.Close()

		// 7. Execute plan with context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := exec.Execute(ctx, actions); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("✔ Apply completed successfully")
	},
}

func init() {
	applyCmd.Flags().BoolVar(
		&autoApprove,
		"auto-approve",
		false,
		"Apply changes without confirmation",
	)

	rootCmd.AddCommand(applyCmd)
}

func askForConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to apply these changes? (yes/no): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "yes"
}
