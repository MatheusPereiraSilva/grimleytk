package validator

import (
	"fmt"

	"grimleytk/internal/config"
)

// IssueSeverity defines the level of a validation issue
type IssueSeverity string

const (
	Error   IssueSeverity = "ERROR"
	Warning IssueSeverity = "WARNING"
)

// Issue represents a validation problem
type Issue struct {
	Code     string
	Severity IssueSeverity
	Message  string
	Path     string
}

// ValidateStructural performs basic structural validation
func ValidateStructural(cfg *config.Config) []Issue {
	var issues []Issue

	// Version
	if cfg.Version == "" {
		issues = append(issues, Issue{
			Code:     "GRIMLEY-E001",
			Severity: Error,
			Message:  "Missing required field: version",
			Path:     "version",
		})
	}

	// Project
	if cfg.Project.Name == "" {
		issues = append(issues, Issue{
			Code:     "GRIMLEY-E002",
			Severity: Error,
			Message:  "Missing required field: project.name",
			Path:     "project.name",
		})
	}

	if cfg.Project.Environment == "" {
		issues = append(issues, Issue{
			Code:     "GRIMLEY-E003",
			Severity: Error,
			Message:  "Missing required field: project.environment",
			Path:     "project.environment",
		})
	}

	// Database
	if cfg.Database.Engine == "" {
		issues = append(issues, Issue{
			Code:     "GRIMLEY-E004",
			Severity: Error,
			Message:  "Missing required field: database.engine",
			Path:     "database.engine",
		})
	}

	if cfg.Database.Name == "" {
		issues = append(issues, Issue{
			Code:     "GRIMLEY-E005",
			Severity: Error,
			Message:  "Missing required field: database.name",
			Path:     "database.name",
		})
	}

	// Domains
	if len(cfg.Domains) == 0 {
		issues = append(issues, Issue{
			Code:     "GRIMLEY-E006",
			Severity: Error,
			Message:  "No domains defined",
			Path:     "domains",
		})
	}

	// Domain uniqueness and schema presence
	schemas := make(map[string]string)

	for domainName, domain := range cfg.Domains {
		if domain.Schema == "" {
			issues = append(issues, Issue{
				Code:     "GRIMLEY-E007",
				Severity: Error,
				Message:  fmt.Sprintf("Domain '%s' has no schema defined", domainName),
				Path:     fmt.Sprintf("domains.%s.schema", domainName),
			})
		}

		if existing, ok := schemas[domain.Schema]; ok {
			issues = append(issues, Issue{
				Code:     "GRIMLEY-E008",
				Severity: Error,
				Message:  fmt.Sprintf("Schema '%s' is used by multiple domains (%s, %s)", domain.Schema, existing, domainName),
				Path:     fmt.Sprintf("domains.%s.schema", domainName),
			})
		} else if domain.Schema != "" {
			schemas[domain.Schema] = domainName
		}
	}

	return issues
}
