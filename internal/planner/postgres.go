package planner

import (
	"fmt"
	"github.com/MatheusPereiraSilva/grimleytk/internal/config"
	"github.com/lib/pq"
	"sort"
	"strings"
)

func keys[T any](m map[string]T) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
func qualified(schema, name string) string {
	return pq.QuoteIdentifier(schema) + "." + pq.QuoteIdentifier(name)
}
func columnSQL(name string, c config.Column) string {
	s := pq.QuoteIdentifier(name) + " " + c.Type
	if !c.Nullable {
		s += " NOT NULL"
	}
	if c.Unique {
		s += " UNIQUE"
	}
	return s
}
func buildPostgresPlan(cfg *config.Config) []Action {
	var actions []Action
	for _, dn := range keys(cfg.Domains) {
		d := cfg.Domains[dn]
		actions = append(actions, Action{Type: CreateSchema, Description: "Create schema " + d.Schema, SQL: "CREATE SCHEMA IF NOT EXISTS " + pq.QuoteIdentifier(d.Schema) + ";"})
		if d.Owns == nil {
			continue
		}
		for _, tn := range keys(d.Owns.Tables) {
			t := d.Owns.Tables[tn]
			full := qualified(d.Schema, tn)
			var cols, pk []string
			for _, cn := range keys(t.Columns) {
				c := t.Columns[cn]
				cols = append(cols, columnSQL(cn, c))
				if c.PrimaryKey {
					pk = append(pk, pq.QuoteIdentifier(cn))
				}
			}
			if len(pk) > 0 {
				cols = append(cols, "PRIMARY KEY ("+strings.Join(pk, ", ")+")")
			}
			// Primary keys are created only with new tables. Existing constraints are never altered.
			actions = append(actions, Action{Type: CreateTable, Description: "Create table " + full, SQL: fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", full, strings.Join(cols, ", "))})
			for _, cn := range keys(t.Columns) {
				actions = append(actions, Action{Type: AddColumn, Description: "Add column " + full + "." + cn, SQL: fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s;", full, columnSQL(cn, t.Columns[cn]))})
			}
		}
	}
	// Views follow all tables, independent of domain name or map iteration order.
	for _, dn := range keys(cfg.Domains) {
		d := cfg.Domains[dn]
		for _, rn := range keys(d.Reads) {
			r := d.Reads[rn]
			parts := strings.Split(r.From, ".")
			source := qualified(cfg.Domains[parts[0]].Schema, parts[1])
			target := qualified(d.Schema, rn)
			cols := make([]string, len(r.Columns))
			for i, c := range r.Columns {
				cols[i] = pq.QuoteIdentifier(c)
			}
			// PostgreSQL lacks CREATE VIEW IF NOT EXISTS. Do not replace an existing relation.
			sql := fmt.Sprintf("DO $grimley$ BEGIN IF pg_catalog.to_regclass(%s) IS NULL THEN EXECUTE %s; END IF; END $grimley$;", pq.QuoteLiteral(target), pq.QuoteLiteral("CREATE VIEW "+target+" AS SELECT "+strings.Join(cols, ", ")+" FROM "+source))
			actions = append(actions, Action{Type: CreateView, Description: "Create view " + target, SQL: sql})
		}
	}
	return actions
}
