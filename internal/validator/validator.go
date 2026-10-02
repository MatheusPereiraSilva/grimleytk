package validator

import "github.com/MatheusPereiraSilva/grimleytk/internal/config"

// Validate runs all declaration checks; it does not inspect a live database.
func Validate(cfg *config.Config) []Issue {
	if cfg == nil {
		return []Issue{{Code: "GRIMLEY-E000", Severity: Error, Message: "Missing configuration"}}
	}
	issues := ValidateStructural(cfg)
	issues = append(issues, ValidateReferences(cfg)...)
	issues = append(issues, ValidateArchitecture(cfg)...)
	return append(issues, ValidateSecurity(cfg)...)
}
