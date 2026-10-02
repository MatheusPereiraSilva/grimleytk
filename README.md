# GrimleyTK

GrimleyTK is an initial MVP for declarative data architecture and governance.
It supports **PostgreSQL only**. It is neither an ORM nor a migration framework.

The CLI edits `grimley.yaml`, validates declarations, generates SQL without
connecting to a database, and can execute that plan in one transaction.
Ownership is declared per domain; cross-domain reads name explicit source columns.

## Build and install

Requires Go 1.23 or newer. From the directory containing `go.mod`:

```sh
go mod tidy
go build -o bin/grimleytk .
go install .
```

`go install .` installs into `GOBIN` (or `$(go env GOPATH)/bin` by default).
Add that directory to PATH to use the commands below. Alternatively, use the
absolute path to `bin/grimleytk`. In this checkout, the module is in `grimleytk/`;
after cloning, enter the directory containing `go.mod` before building.
These instructions install the local source and do not require a published release.

## Workflow without a database

Run in a new working directory:

```sh
grimleytk init
grimleytk create domain catalog --schema catalog --owner catalog-service
grimleytk create table catalog.products --description "Main product table"
grimleytk create column catalog.products.id --type uuid --primary-key --nullable=false
grimleytk create column catalog.products.price --type numeric --nullable=false
grimleytk create domain wishlist --schema wishlist --owner wishlist-service
grimleytk create view wishlist.products_view --from catalog.products --columns id,price
grimleytk validate
grimleytk plan
grimleytk show
grimleytk show domains
grimleytk show tables
grimleytk show reads
```

`init` creates a valid starter domain named `example` and refuses to overwrite
an existing file. Creation commands edit YAML only; `validate`, `plan` and
`apply` check the completed declaration. Tables need at least one column.
`examples/grimley.yaml` is a complete two-domain example. A domain that only
reads emits a nonblocking warning.

## Apply

All commands read `grimley.yaml` from the current working directory. There is
no `--file`, `--config`, or positional file argument. To select a declaration,
place it in its own directory as `grimley.yaml`, enter that directory, and run
`grimleytk apply [--auto-approve]` using the installed executable or an absolute
path to the built `bin/grimleytk`.

Create the PostgreSQL database separately, configure its connection in YAML,
and set the environment variable named by `database.credentials.password_env`.
Passwords are read from the environment; the username is stored in YAML.

```sh
export DB_PASSWORD='your-local-password'
grimleytk apply                 # displays SQL and asks for yes/no
grimleytk apply --auto-approve  # skips confirmation
```

Apply executes the plan in one transaction and rolls back that transaction on
statement failure. This is not automatic rollback migration support. `ssl: true`
uses lib/pq's `require` mode; configure your database/network appropriately.

## Safety and current limitations

- No DROP, DELETE, data updates, column type changes or constraint removal are generated.
- Names must match `[a-z_][a-z0-9_]{0,62}` and are quoted in SQL. SQL types use
  a fixed allowlist: uuid, text, boolean/bool, smallint, integer/int, bigint,
  real, double precision, numeric/decimal, date, time, timestamp (with or
  without time zone), timestamptz, json/jsonb, bytea, varchar/character varying,
  and one-dimensional arrays of these. Custom types, defaults, SQL expressions
  and parameterized types are unsupported.
- Plans have deterministic ordering. Schemas, tables and columns use
  `IF NOT EXISTS`. Primary keys (including composite keys, ordered by column
  name) are defined only when a table is first created. Existing tables can
  receive missing columns, but existing column definitions and constraints
  are not reconciled. Adding NOT NULL or UNIQUE columns may fail on populated
  tables; the transaction then rolls back. Plan output is not a live schema diff.
- Views are created after source tables, only if no relation with that name
  exists. Existing views are not replaced or reconciled. Concurrent applies
  are not coordinated; use one apply process at a time.
- Ownership and read-only access are declaration checks, not database grants.
  PostgreSQL views may be updatable; manage actual database privileges separately.
  Sensitive-column blocking is a name-based heuristic, not data classification.
- Unknown YAML fields and multiple documents are rejected. Reserved schema fields
  for indexes, policies, documentation, domain databases, access grants, sync,
  materialized views and consistency modes are rejected as unsupported.
- No live schema diffing, drift detection, automatic rollback migrations, CDC,
  replication, multiple databases, RLS generation or evolution validation.

## Code layout and checks

```text
main.go
cmd/                  CLI orchestration
internal/config/      YAML structs and strict loader
internal/validator/   structure, references, ownership and security
internal/planner/     validated PostgreSQL SQL generation
internal/executor/    transactional execution
examples/grimley.yaml
```

```sh
gofmt -w main.go cmd internal
go mod tidy
go test ./...
go vet ./...
go build ./...
```

Run the repeatable CLI smoke test with an absolute binary path:

```sh
go build -o /tmp/grimleytk-cli .
sh scripts/smoke.sh /tmp/grimleytk-cli
```

Unit and CLI help tests require no database and build and invoke the public
`grimleytk` executable. The optional integration test also invokes that executable
and checks the resulting table and view. Use a disposable local PostgreSQL
database named `grimleytk_test`, with user `postgres` and SSL disabled:

```sh
GRIMLEYTK_TEST_POSTGRES_PORT=55432 \
GRIMLEYTK_TEST_POSTGRES_PASSWORD=local-test-password \
go test -v -run TestCLIApplyPostgres .
```

Without `GRIMLEYTK_TEST_POSTGRES_PORT`, this integration test is skipped.
MIT license; see `LICENSE`.
