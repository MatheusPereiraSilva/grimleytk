func TestReferenceValidator_InvalidTable(t *testing.T) {
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
			"catalog": {
				Owns: &config.OwnedResources{
					Tables: map[string]config.Table{},
				},
			},
		},
	}

	issues := ValidateReferences(cfg)

	if len(issues) == 0 {
		t.Fatalf("expected reference validation error, got none")
	}
}
