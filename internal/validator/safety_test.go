package validator

import (
	"github.com/MatheusPereiraSilva/grimleytk/internal/config"
	"strings"
	"testing"
)

func TestPhysicalOwnership(t *testing.T) {
	tables := &config.OwnedResources{Tables: map[string]config.Table{"products": {}}}
	cfg := &config.Config{Domains: map[string]config.Domain{"a": {Schema: "shared", Owns: tables}, "b": {Schema: "shared", Owns: tables}}}
	found := false
	for _, i := range ValidateArchitecture(cfg) {
		if i.Code == "GRIMLEY-E201" {
			found = true
		}
	}
	if !found {
		t.Fatal("duplicate physical ownership accepted")
	}
	d := cfg.Domains["b"]
	d.Schema = "other"
	cfg.Domains["b"] = d
	if BuildReport(ValidateArchitecture(cfg)).HasErrors() {
		t.Fatal("same table name in different schemas should be allowed")
	}
}
func TestIdentifiers(t *testing.T) {
	for _, s := range []string{"id", "select", "_name", strings.Repeat("a", 63)} {
		if !ValidIdentifier(s) {
			t.Fatal(s)
		}
	}
	for _, s := range []string{"", "A", "a.b", "a-b", "a\"", strings.Repeat("a", 64), "1a", "é"} {
		if ValidIdentifier(s) {
			t.Fatal(s)
		}
	}
}
func TestExampleAndInvalidStructures(t *testing.T) {
	load := func() *config.Config {
		c, e := config.Load("../../examples/grimley.yaml")
		if e != nil {
			t.Fatal(e)
		}
		return c
	}
	if r := BuildReport(Validate(load())); r.HasErrors() {
		t.Fatal(r.String())
	}
	cases := map[string]func(*config.Config){
		"engine":      func(c *config.Config) { c.Database.Engine = "mysql" },
		"empty table": func(c *config.Config) { c.Domains["catalog"].Owns.Tables["empty"] = config.Table{} },
		"reference": func(c *config.Config) {
			d := c.Domains["wishlist"]
			r := d.Reads["products_view"]
			r.From = "missing.products"
			d.Reads["products_view"] = r
		},
		"missing column": func(c *config.Config) {
			d := c.Domains["wishlist"]
			r := d.Reads["products_view"]
			r.Columns = []string{"missing"}
			d.Reads["products_view"] = r
		},
		"sensitive": func(c *config.Config) {
			d := c.Domains["wishlist"]
			r := d.Reads["products_view"]
			r.Columns = []string{"api_token"}
			d.Reads["products_view"] = r
			c.Domains["catalog"].Owns.Tables["products"].Columns["api_token"] = config.Column{Type: "text"}
		},
		"nullable key": func(c *config.Config) {
			c.Domains["catalog"].Owns.Tables["products"].Columns["id"] = config.Column{Type: "uuid", PrimaryKey: true, Nullable: true}
		},
	}
	for n, change := range cases {
		t.Run(n, func(t *testing.T) {
			c := load()
			change(c)
			if !BuildReport(Validate(c)).HasErrors() {
				t.Fatal("invalid config accepted")
			}
		})
	}
}
