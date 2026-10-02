package planner

import (
	"fmt"
	"github.com/MatheusPereiraSilva/grimleytk/internal/config"
	"github.com/MatheusPereiraSilva/grimleytk/internal/validator"
)

// BuildPlan validates declarations before generating any SQL, including direct callers.
func BuildPlan(cfg *config.Config) ([]Action, error) {
	report := validator.BuildReport(validator.Validate(cfg))
	if report.HasErrors() {
		return nil, fmt.Errorf("invalid configuration: %s", report.String())
	}
	return buildPostgresPlan(cfg), nil
}
