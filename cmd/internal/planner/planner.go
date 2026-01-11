package planner

import "grimleytk/internal/config"

// BuildPlan builds a list of actions from the configuration
func BuildPlan(cfg *config.Config) []Action {
	var actions []Action

	if cfg.Database.Engine == "postgres" {
		actions = append(actions, buildPostgresPlan(cfg)...)
	}

	return actions
}
