package planner

import (
	"fmt"
	"strings"

	"grimleytk/internal/config"
)

func buildPostgresPlan(cfg *config.Config) []Action {
	var actions []Action

	for domainName, domain := range cfg.Domains {

		// 1. CREATE SCHEMA
		actions = append(actions, Action{
			Type:        CreateSchema,
			Description: fmt.Sprintf("Create schema %s", domain.Schema),
			SQL:         fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", domain.Schema),
		})

		// 2. CREATE TABLE + COLUMNS
		if domain.Owns == nil {
			continue
		}

		for tableName, table := range domain.Owns.Tables {

			fullTable := fmt.Sprintf("%s.%s", domain.Schema, tableName)

			actions = append(actions, Action{
				Type:        CreateTable,
				Description: fmt.Sprintf("Create table %s", fullTable),
				SQL:         buildCreateTableSQL(domain.Schema, tableName, table),
			})

			// Columns (ADD COLUMN)
			for columnName, column := range table.Columns {
				actions = append(actions, Action{
					Type:        AddColumn,
					Description: fmt.Sprintf("Add column %s.%s", fullTable, columnName),
					SQL:         buildAddColumnSQL(domain.Schema, tableName, columnName, column),
				})
			}
		}
	}

	return actions
}

func buildCreateTableSQL(schema, table string, tableDef config.Table) string {
	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s.%s ();",
		schema,
		table,
	)
}

func buildAddColumnSQL(schema, table, column string, colDef config.Column) string {
	var parts []string

	parts = append(parts, colDef.Type)

	if !colDef.Nullable {
		parts = append(parts, "NOT NULL")
	}

	if colDef.Unique {
		parts = append(parts, "UNIQUE")
	}

	sql := fmt.Sprintf(
		"ALTER TABLE %s.%s ADD COLUMN IF NOT EXISTS %s %s;",
		schema,
		table,
		column,
		strings.Join(parts, " "),
	)

	if colDef.PrimaryKey {
		sql += fmt.Sprintf(
			"\nALTER TABLE %s.%s ADD PRIMARY KEY (%s);",
			schema,
			table,
			column,
		)
	}

	return sql
}
