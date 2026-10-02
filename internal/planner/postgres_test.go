package planner

import (
	"testing"

	"github.com/MatheusPereiraSilva/grimleytk/internal/config"
)

func TestPostgresPlanner_GeneratesCreateTable(t *testing.T) {
	cfg := &config.Config{
		Version: "0.1", Project: config.Project{Name: "test", Environment: "local"},
		Database: config.Database{
			Engine: "postgres", Name: "test",
		},
		Domains: map[string]config.Domain{
			"catalog": {
				Schema: "catalog", Owner: "catalog-service",
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

	actions, err := BuildPlan(cfg)
	if err != nil {
		t.Fatal(err)
	}

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
