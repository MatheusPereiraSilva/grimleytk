package validator

import (
	"fmt"
	"github.com/MatheusPereiraSilva/grimleytk/internal/config"
	"regexp"
	"strings"
)

var identifier = regexp.MustCompile(`^[a-z_][a-z0-9_]{0,62}$`)

// ValidIdentifier restricts names to ASCII identifiers within PostgreSQL's byte limit.
// SQL generation also quotes them, allowing reserved words without ambiguity.
func ValidIdentifier(s string) bool { return identifier.MatchString(s) }

var sqlType = regexp.MustCompile(`^(uuid|text|boolean|bool|smallint|integer|int|bigint|real|double precision|numeric|decimal|date|timestamp|timestamp with time zone|timestamp without time zone|timestamptz|time|json|jsonb|bytea|varchar|character varying)(\[\])?$`)

func ValidType(s string) bool { return sqlType.MatchString(s) }

func validateSQL(cfg *config.Config) []Issue {
	var issues []Issue
	fail := func(path, msg string) {
		issues = append(issues, Issue{Code: "GRIMLEY-E009", Severity: Error, Path: path, Message: msg})
	}
	name := func(path, value string) {
		if !ValidIdentifier(value) {
			fail(path, fmt.Sprintf("Invalid identifier %q: use lowercase ASCII letters, digits and underscores, maximum 63 bytes", value))
		}
	}
	if cfg.Database.Engine != "postgres" {
		fail("database.engine", "Only postgres is supported")
	}
	if cfg.Version != "0.1" {
		fail("version", "Only configuration version 0.1 is supported")
	}
	if len(cfg.Policies) > 0 || cfg.Docs != nil {
		fail("policies/documentation", "Policies and documentation generation are not implemented")
	}
	for dn, d := range cfg.Domains {
		path := "domains." + dn
		name(path, dn)
		name(path+".schema", d.Schema)
		if strings.TrimSpace(d.Owner) == "" {
			fail(path+".owner", "Domain owner is required")
		}
		if d.Database != nil || d.Access != nil || d.Sync != nil {
			fail(path, "Domain databases, access grants and sync are not implemented")
		}
		if d.Owns != nil {
			for tn, t := range d.Owns.Tables {
				tp := path + ".owns.tables." + tn
				name(tp, tn)
				if len(t.Columns) == 0 {
					fail(tp, "At least one column is required")
				}
				if len(t.Indexes) > 0 {
					fail(tp, "Indexes are not implemented")
				}
				for cn, c := range t.Columns {
					name(tp+".columns."+cn, cn)
					if !ValidType(c.Type) {
						fail(tp+".columns."+cn, "Unsupported SQL type")
					}
					if c.PrimaryKey && c.Nullable {
						fail(tp+".columns."+cn, "Primary key columns cannot be nullable")
					}
				}
			}
		}
		for rn, r := range d.Reads {
			rp := path + ".reads." + rn
			name(rp, rn)
			if d.Owns != nil {
				if _, ok := d.Owns.Tables[rn]; ok {
					fail(rp, "View conflicts with an owned table")
				}
			}
			if r.Materialized || r.Consistency.Type != "" {
				fail(rp, "Materialized views and consistency modes are not implemented")
			}
			if len(r.Columns) == 0 {
				fail(rp, "Read models require columns")
			}
			seen := map[string]bool{}
			for _, c := range r.Columns {
				name(rp+".columns", c)
				if seen[c] {
					fail(rp, "Duplicate read column")
				}
				seen[c] = true
			}
		}
	}
	return issues
}
