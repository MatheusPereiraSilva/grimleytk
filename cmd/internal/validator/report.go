package validator

import (
	"fmt"
	"strings"
)

// Report represents a formatted validation report
type Report struct {
	Errors   []Issue
	Warnings []Issue
}

// BuildReport separates issues by severity
func BuildReport(issues []Issue) Report {
	report := Report{}

	for _, issue := range issues {
		switch issue.Severity {
		case Error:
			report.Errors = append(report.Errors, issue)
		case Warning:
			report.Warnings = append(report.Warnings, issue)
		}
	}

	return report
}

// HasErrors returns true if there are blocking errors
func (r Report) HasErrors() bool {
	return len(r.Errors) > 0
}

// String formats the report for terminal output
func (r Report) String() string {
	var sb strings.Builder

	if len(r.Errors) == 0 && len(r.Warnings) == 0 {
		sb.WriteString("✔ Validation successful\n")
		sb.WriteString("0 errors, 0 warnings\n")
		return sb.String()
	}

	if len(r.Errors) > 0 {
		sb.WriteString("✖ Validation failed\n\n")
		sb.WriteString("Errors:\n")
		for _, err := range r.Errors {
			sb.WriteString(formatIssue(err))
		}
	}

	if len(r.Warnings) > 0 {
		if len(r.Errors) > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString("Warnings:\n")
		for _, warn := range r.Warnings {
			sb.WriteString(formatIssue(warn))
		}
	}

	sb.WriteString(fmt.Sprintf(
		"\nSummary: %d errors, %d warnings\n",
		len(r.Errors),
		len(r.Warnings),
	))

	return sb.String()
}

func formatIssue(issue Issue) string {
	return fmt.Sprintf(
		"- [%s] %s\n  Path: %s\n",
		issue.Code,
		issue.Message,
		issue.Path,
	)
}
