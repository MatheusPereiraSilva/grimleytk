package validator

import (
	"testing"

	"grimleytk/internal/config"
)

func TestStructuralValidator_NoDomains(t *testing.T) {
	cfg := &config.Config{
		Version: "0.1",
		Project: config.Project{
			Name:        "test",
			Environment: "local",
		},
		Database: config.Database{
			Engine: "postgres",
			Name:   "test",
		},
		Domains: map[string]config.Domain{},
	}

	issues := ValidateStructural(cfg)

	if len(issues) == 0 {
		t.Fatalf("expected validation errors, got none")
	}
}
