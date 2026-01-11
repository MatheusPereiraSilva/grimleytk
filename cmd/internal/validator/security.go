package validator

import (
	"fmt"
	"strings"

	"grimleytk/internal/config"
)

// Sensitive column name patterns (case-insensitive)
var sensitiveColumnKeywords = []string{
	"password",
	"passwd",
	"token",
	"secret",
	"api_key",
	"credit_card",
	"card_number",
	"cvv",
	"ssn",
}

// ValidateSecurity checks for unsafe data exposure
func ValidateSecurity(cfg *config.Config) []Issue {
	var issues []Issue

	for domainName, domain := range cfg.Domains {
		for readName, read := range domain.Reads {

			for _, column := range read.Columns {
				if isSensitiveColumn(column) {
					issues = append(issues, Issue{
						Code:     "GRIMLEY-E301",
						Severity: Error,
						Message: fmt.Sprintf(
							"Sensitive column '%s' exposed in read model '%s' of domain '%s'",
							column,
							readName,
							domainName,
						),
						Path: fmt.Sprintf("domains.%s.reads.%s.columns", domainName, readName),
					})
				}
			}
		}
	}

	return issues
}

func isSensitiveColumn(column string) bool {
	col := strings.ToLower(column)
	for _, keyword := range sensitiveColumnKeywords {
		if strings.Contains(col, keyword) {
			return true
		}
	}
	return false
}
