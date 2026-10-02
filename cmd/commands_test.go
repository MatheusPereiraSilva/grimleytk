package cmd

import "testing"

func TestCommandRegistration(t *testing.T) {
	for _, path := range [][]string{{"init"}, {"validate"}, {"plan"}, {"apply"}, {"create", "domain"}, {"create", "table"}, {"create", "column"}, {"create", "view"}, {"show"}, {"show", "domains"}, {"show", "tables"}, {"show", "reads"}} {
		c, remaining, e := rootCmd.Find(path)
		if e != nil || len(remaining) != 0 || c.Name() != path[len(path)-1] {
			t.Fatalf("missing command %v", path)
		}
	}
}
