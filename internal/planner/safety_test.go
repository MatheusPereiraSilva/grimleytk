package planner

import (
	"github.com/MatheusPereiraSilva/grimleytk/internal/config"
	"reflect"
	"strings"
	"testing"
)

func example(t *testing.T) *config.Config {
	t.Helper()
	c, e := config.Load("../../examples/grimley.yaml")
	if e != nil {
		t.Fatal(e)
	}
	return c
}
func TestPlanDeterministicAndPrimaryKeySafe(t *testing.T) {
	c := example(t)
	first, e := BuildPlan(c)
	if e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 20; i++ {
		next, e := BuildPlan(c)
		if e != nil || !reflect.DeepEqual(first, next) {
			t.Fatal("unstable plan", e)
		}
	}
	pk, views := 0, 0
	for i, a := range first {
		if strings.Contains(a.SQL, "ADD PRIMARY KEY") {
			t.Fatal("non-idempotent primary key alteration")
		}
		if strings.Contains(a.SQL, "PRIMARY KEY") {
			pk++
			if a.Type != CreateTable || !strings.Contains(a.SQL, "CREATE TABLE IF NOT EXISTS") {
				t.Fatal(a.SQL)
			}
		}
		if a.Type == CreateView {
			views++
			if i != len(first)-1 || !strings.Contains(a.SQL, "pg_catalog.to_regclass") || !strings.Contains(a.SQL, `FROM "catalog"."products"`) {
				t.Fatal(a.SQL)
			}
		}
		for _, bad := range []string{"DROP ", "DELETE ", "TRUNCATE ", "CREATE OR REPLACE"} {
			if strings.Contains(a.SQL, bad) {
				t.Fatal(a.SQL)
			}
		}
	}
	if pk != 1 || views != 1 {
		t.Fatalf("pk=%d views=%d", pk, views)
	}
}
func TestPlanRejectsUnsafeInput(t *testing.T) {
	for _, value := range []string{"bad; DROP TABLE x", "MixedCase", strings.Repeat("a", 64), "bad.name", "bad\x00name"} {
		c := example(t)
		d := c.Domains["catalog"]
		d.Schema = value
		c.Domains["catalog"] = d
		if a, e := BuildPlan(c); e == nil || len(a) != 0 {
			t.Fatalf("accepted %q", value)
		}
	}
	c := example(t)
	table := c.Domains["catalog"].Owns.Tables["products"]
	table.Columns["id"] = config.Column{Type: "uuid); DROP TABLE products;--"}
	if _, e := BuildPlan(c); e == nil {
		t.Fatal("accepted injected type")
	}
	if _, e := BuildPlan(nil); e == nil {
		t.Fatal("accepted nil")
	}
}
func TestCompositePrimaryKeyAndQuotedKeywords(t *testing.T) {
	c := example(t)
	d := c.Domains["catalog"]
	d.Schema = "select"
	table := d.Owns.Tables["products"]
	table.Columns["price"] = config.Column{Type: "numeric", PrimaryKey: true}
	c.Domains["catalog"] = d
	a, e := BuildPlan(c)
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(a[1].SQL, `"select"."products"`) || !strings.Contains(a[1].SQL, `PRIMARY KEY ("id", "price")`) {
		t.Fatal(a[1].SQL)
	}
}
