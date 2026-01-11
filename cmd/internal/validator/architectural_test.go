package validator

import (
	"testing"

	"grimleytk/internal/config"
)

func TestArchitecturalValidator_TableWithoutOwner(t *testing.T) {
	cfg := &config.Config{
		Domains: map[string]config.Domain{
			"wishlist": {
				Reads: map[string]config.Read{
					"products_view": {
						From:    "catalog.products",
						Columns: []string{"id"},
					},
				},
			},
		},
	}

	issues := ValidateArchitecture(cfg)

	if len(issues) == 0 {
		t.Fatalf("expected architectural error for table without owner, got none")
	}
}

func TestArchitecturalValidator_DomainReadsOwnTable(t *testing.T) {
	cfg := &config.Config{
		Domains: map[string]config.Domain{
			"catalog": {
				Schema: "catalog",
				Owns: &config.OwnedResources{
					Tables: map[string]config.Table{
						"products": {
							Columns: map[string]config.Column{},
						},
					},
				},
				Reads: map[string]config.Read{
					"products_view": {
						From:    "catalog.products",
						Columns: []string{"id"},
					},
				},
			},
		},
	}

	issues := ValidateArchitecture(cfg)

	if len(issues) == 0 {
		t.Fatalf("expected architectural error when domain reads its own table")
	}
}

func TestArchitecturalValidator_ReadOnlyDomainWarning(t *testing.T) {
	cfg := &config.Config{
		Domains: map[string]config.Domain{
			"search": {
				Reads: map[string]config.Read{
					"products_view": {
						From:    "catalog.products",
						Columns: []string{"id"},
					},
				},
			},
			"catalog": {
				Schema: "catalog",
				Owns: &config.OwnedResources{
					Tables: map[string]config.Table{
						"products": {
							Columns: map[string]config.Column{},
						},
					},
				},
			},
		},
	}

	issues := ValidateArchitecture(cfg)

	foundWarning := false
	for _, issue := range issues {
		if issue.Severity == Warning {
			foundWarning = true
		}
	}

	if !foundWarning {
		t.Fatalf("expected warning for read-only domain, got none")
	}
}
