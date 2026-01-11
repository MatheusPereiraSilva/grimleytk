package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"grimleytk/internal/config"
	"grimleytk/internal/executor"
	"grimleytk/internal/planner"
	"grimleytk/internal/validator"

	"github.com/spf13/cobra"
)

var autoApprove bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the planned database changes",
	Long: `Apply executes the SQL generated from the GrimleyTK plan.
This operation modifies the database and requires confirmation.`,
	Run: func(cmd *cobra.Command, args []string) {

		// 1. Load config
		cfg, err := config.Load("grimley.yaml")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 2. Validate before apply
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
			fmt.Println("No actions to apply.")
			return
		}

		// 4. Show plan
		fmt.Println("Execution Plan:\n")
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
