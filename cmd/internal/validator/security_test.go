func TestSecurityValidator_SensitiveColumn(t *testing.T) {
	cfg := &config.Config{
		Domains: map[string]config.Domain{
			"auth": {
				Reads: map[string]config.Read{
					"user_view": {
						From:    "users.users",
						Columns: []string{"id", "password"},
					},
				},
			},
		},
	}

	issues := ValidateSecurity(cfg)

	if len(issues) == 0 {
		t.Fatalf("expected security validation error, got none")
	}
}
