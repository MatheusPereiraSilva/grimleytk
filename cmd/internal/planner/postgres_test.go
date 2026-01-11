package planner

import (
	"testing"

	"grimleytk/internal/config"
)

func TestPostgresPlanner_GeneratesCreateTable(t *testing.T) {
	cfg := &config.Config{
		Database: config.Database{
			Engine: "postgres",
		},
		Domains: map[string]config.Domain{
			"catalog": {
				Schema: "catalog",
				Owns: &config.OwnedResources{
					Tables: map[string]config.Table{
						"products": {
							Columns: map[string]config.Column{
								"id": {
									Type:       "uuid",
									PrimaryKey: true,
									Nullable:   false,
								},
							},
						},
					},
				},
			},
		},
	}

	actions := BuildPlan(cfg)

	if len(actions) == 0 {
		t.Fatalf("expected plan actions, got none")
	}

	foundCreateTable := false
	for _, action := range actions {
		if action.Type == CreateTable {
			foundCreateTable = true
		}
	}

	if !foundCreateTable {
		t.Fatalf("expected CREATE_TABLE action, got none")
	}
}
