package validator

import (
	"fmt"

	"grimleytk/internal/config"
)

// ValidateArchitecture enforces architectural rules
func ValidateArchitecture(cfg *config.Config) []Issue {
	var issues []Issue

	// Map of table ownership: domain.table -> domain
	tableOwners := make(map[string]string)

	// 1. Register owned tables
	for domainName, domain := range cfg.Domains {
		if domain.Owns == nil {
			continue
		}

		for tableName := range domain.Owns.Tables {
			key := fmt.Sprintf("%s.%s", domainName, tableName)

			if owner, exists := tableOwners[key]; exists {
				issues = append(issues, Issue{
					Code:     "GRIMLEY-E201",
					Severity: Error,
					Message:  fmt.Sprintf("Table '%s' is owned by multiple domains (%s, %s)", key, owner, domainName),
					Path:     fmt.Sprintf("domains.%s.owns.tables.%s", domainName, tableName),
				})
			} else {
				tableOwners[key] = domainName
			}
		}
	}

	// 2. Validate read-only access
	for domainName, domain := range cfg.Domains {
		for readName, read := range domain.Reads {

			source := read.From
			owner, exists := tableOwners[source]
			if !exists {
				issues = append(issues, Issue{
					Code:     "GRIMLEY-E202",
					Severity: Error,
					Message:  fmt.Sprintf("Read model '%s' references table '%s' with no owning domain", readName, source),
					Path:     fmt.Sprintf("domains.%s.reads.%s.from", domainName, readName),
				})
				continue
			}

			// A domain cannot own and read via read-model
			if owner == domainName {
				issues = append(issues, Issue{
					Code:     "GRIMLEY-E203",
					Severity: Error,
					Message:  fmt.Sprintf("Domain '%s' should access its own table '%s' directly, not via read model", domainName, source),
					Path:     fmt.Sprintf("domains.%s.reads.%s", domainName, readName),
				})
			}

			// Enforce read-only
			if read.Access.Mode != "" && read.Access.Mode != "read-only" {
				issues = append(issues, Issue{
					Code:     "GRIMLEY-E204",
					Severity: Error,
					Message:  fmt.Sprintf("Read model '%s' in domain '%s' must be read-only", readName, domainName),
					Path:     fmt.Sprintf("domains.%s.reads.%s.access.mode", domainName, readName),
				})
			}
		}
	}

	// 3. Warn about read-only domains
	for domainName, domain := range cfg.Domains {
		if domain.Owns == nil || len(domain.Owns.Tables) == 0 {
			if len(domain.Reads) > 0 {
				issues = append(issues, Issue{
					Code:     "GRIMLEY-W205",
					Severity: Warning,
					Message:  fmt.Sprintf("Domain '%s' has no owned tables (read-only domain)", domainName),
					Path:     fmt.Sprintf("domains.%s", domainName),
				})
			}
		}
	}

	return issues
}
