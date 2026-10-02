package main

import (
	"context"
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/MatheusPereiraSilva/grimleytk/internal/config"
	"gopkg.in/yaml.v3"
)

func buildCLI(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "grimleytk")
	command := exec.Command("go", "build", "-o", binary, ".")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, output)
	}
	return binary
}

func TestCLIApplyAvailable(t *testing.T) {
	binary := buildCLI(t)
	for _, args := range [][]string{{"--help"}, {"apply", "--help"}} {
		command := exec.Command(binary, args...)
		command.Dir = t.TempDir() // Help must work without grimley.yaml.
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("%v: %v\n%s", args, err, output)
		}
		if len(args) == 1 {
			found := false
			for _, line := range strings.Split(string(output), "\n") {
				fields := strings.Fields(line)
				if len(fields) > 1 && fields[0] == "apply" {
					found = true
				}
			}
			if !found {
				t.Fatalf("apply missing from root help:\n%s", output)
			}
		} else if !strings.Contains(string(output), "grimleytk apply [flags]") || !strings.Contains(string(output), "--auto-approve") {
			t.Fatalf("unexpected apply help:\n%s", output)
		}
	}
	command := exec.Command(binary, "apply", "other.yaml")
	output, err := command.CombinedOutput()
	if err == nil || !strings.Contains(string(output), "unknown command") {
		t.Fatalf("apply must reject positional config paths: %v\n%s", err, output)
	}
}

// Use only a disposable PostgreSQL database: this test creates catalog and wishlist.
func TestCLIApplyPostgres(t *testing.T) {
	portText := os.Getenv("GRIMLEYTK_TEST_POSTGRES_PORT")
	if portText == "" {
		t.Skip("set GRIMLEYTK_TEST_POSTGRES_PORT for a disposable local PostgreSQL database")
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatal(err)
	}
	binary := buildCLI(t)
	cfg, err := config.Load("examples/grimley.yaml")
	if err != nil {
		t.Fatal(err)
	}
	cfg.Database.Host = "127.0.0.1"
	cfg.Database.Port = port
	cfg.Database.Name = "grimleytk_test"
	cfg.Database.Credentials.PasswordEnv = "GRIMLEYTK_TEST_POSTGRES_PASSWORD"
	data, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "grimley.yaml"), data, 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, binary, "apply", "--auto-approve")
	command.Dir = dir
	output, err := command.CombinedOutput()
	if err != nil || !strings.Contains(string(output), "Apply completed successfully") {
		t.Fatalf("public CLI apply: %v\n%s", err, output)
	}
	// Inspect the database independently of the executor to verify public CLI effects.
	db, err := sql.Open("postgres", "host=127.0.0.1 port="+portText+" user=postgres dbname=grimleytk_test sslmode=disable password="+os.Getenv("GRIMLEYTK_TEST_POSTGRES_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, relation := range []string{"catalog.products", "wishlist.products_view"} {
		var exists bool
		if err := db.QueryRowContext(ctx, "SELECT to_regclass($1) IS NOT NULL", relation).Scan(&exists); err != nil || !exists {
			t.Fatalf("CLI did not create %s: exists=%v err=%v", relation, exists, err)
		}
	}
	t.Log("public CLI apply created catalog.products and wishlist.products_view")
}
