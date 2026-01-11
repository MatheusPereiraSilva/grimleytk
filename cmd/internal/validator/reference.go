package validator

import (
	"fmt"
	"strings"

	"grimleytk/internal/config"
)

// ValidateReferences validates cross-domain references
func ValidateReferences(cfg *config.Config) []Issue {
	var issues []Issue

	for domainName, domain := range cfg.Domains {

		for readName, read := range domain.Reads {

			// Validate "from" format: domain.table
			parts := strings.Split(read.From, ".")
			if len(parts) != 2 {
				issues = append(issues, Issue{
					Code:     "GRIMLEY-E101",
					Severity: Error,
					Message:  fmt.Sprintf("Invalid from reference '%s', expected format 'domain.table'", read.From),
					Path:     fmt.Sprintf("domains.%s.reads.%s.from", domainName, readName),
				})
				continue
			}

			sourceDomainName := parts[0]
			sourceTableName := parts[1]

			// Domain cannot read from itself via read model
			if sourceDomainName == domainName {
				issues = append(issues, Issue{
					Code:     "GRIMLEY-E102",
					Severity: Error,
					Message:  fmt.Sprintf("Domain '%s' defines a read model for its own table '%s'", domainName, sourceTableName),
					Path:     fmt.Sprintf("domains.%s.reads.%s", domainName, readName),
				})
				continue
			}

			sourceDomain, ok := cfg.Domains[sourceDomainName]
			if !ok {
				issues = append(issues, Issue{
					Code:     "GRIMLEY-E103",
					Severity: Error,
					Message:  fmt.Sprintf("Referenced domain '%s' does not exist", sourceDomainName),
					Path:     fmt.Sprintf("domains.%s.reads.%s.from", domainName, readName),
				})
				continue
			}

			if sourceDomain.Owns == nil || sourceDomain.Owns.Tables == nil {
				issues = append(issues, Issue{
					Code:     "GRIMLEY-E104",
					Severity: Error,
					Message:  fmt.Sprintf("Domain '%s' does not own any tables", sourceDomainName),
					Path:     fmt.Sprintf("domains.%s.reads.%s.from", domainName, readName),
				})
				continue
			}

			sourceTable, ok := sourceDomain.Owns.Tables[sourceTableName]
			if !ok {
				issues = append(issues, Issue{
					Code:     "GRIMLEY-E105",
					Severity: Error,
					Message:  fmt.Sprintf(
						"Table '%s' does not exist in domain '%s'",
						sourceTableName,
						sourceDomainName,
					),
					Path: fmt.Sprintf("domains.%s.reads.%s.from", domainName, readName),
				})
				continue
			}

			// Validate columns
			for _, col := range read.Columns {
				if _, ok := sourceTable.Columns[col]; !ok {
					issues = append(issues, Issue{
						Code:     "GRIMLEY-E106",
						Severity: Error,
						Message:  fmt.Sprintf(
							"Column '%s' does not exist in table '%s.%s'",
							col,
							sourceDomainName,
							sourceTableName,
						),
						Path: fmt.Sprintf("domains.%s.reads.%s.columns", domainName, readName),
					})
				}
			}
		}
	}

	return issues
}
