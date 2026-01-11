# GrimleyTK

GrimleyTK is an open-source **Declarative Data Architecture Toolkit** designed to help teams model, validate, plan, and apply database architectures with strong governance — **before** problems reach production.

It focuses on **clarity, safety, and explicit ownership** of data in distributed and microservice-oriented systems.

> GrimleyTK is not an ORM.
> GrimleyTK is not a migration framework.
> GrimleyTK is a **data-architecture governance tool**.

---

## Why GrimleyTK exists

Modern systems often suffer from:

* Shared databases with unclear ownership
* Services reading or mutating data they should not
* Implicit contracts between teams
* Architectural drift over time

GrimleyTK solves this by introducing a **single declarative source of truth** (`grimley.yaml`) that defines:

* Who owns which data
* What can be read across domains
* What is forbidden
* What changes are required to reach the desired state

---

## Core principles

* **Explicit ownership**: every table has a single owning domain
* **Read models are explicit**: cross-domain access is always declared
* **Fail fast**: architectural violations are caught early
* **Safe by default**: no destructive operations in MVP
* **Tooling over convention**: rules are enforced, not assumed

---

## What GrimleyTK does (MVP)

✅ Declaratively define data architecture (`grimley.yaml`)
✅ Validate architecture rules and security constraints
✅ Generate an execution plan (dry-run SQL)
✅ Apply changes safely using transactions
✅ Provide a CLI-driven modeling workflow

---

## What GrimleyTK does NOT do (yet)

❌ No schema diffing against live databases
❌ No destructive operations (DROP, DELETE)
❌ No automatic rollback migrations
❌ No CDC or replication
❌ PostgreSQL only (for now)

These are **intentional** omissions for the MVP.

---

## Installation

```bash
go install github.com/your-org/grimleytk@latest
```

Or clone and build locally:

```bash
git clone https://github.com/your-org/grimleytk.git
cd grimleytk
go build
```

---

## Getting started

### 1. Initialize a project

```bash
grimleytk init
```

This creates a starter `grimley.yaml`.

---

### 2. Create a domain

```bash
grimleytk create domain catalog \
  --schema catalog \
  --owner catalog-service
```

---

### 3. Create a table

```bash
grimleytk create table catalog.products \
  --description "Main product table"
```

---

### 4. Add columns

```bash
grimleytk create column catalog.products.id \
  --type uuid \
  --primary-key \
  --nullable false

grimleytk create column catalog.products.price \
  --type numeric \
  --nullable false
```

---

### 5. Create a read model (view)

```bash
grimleytk create view wishlist.products_view \
  --from catalog.products \
  --columns id,price
```

---

### 6. Validate architecture

```bash
grimleytk validate
```

This checks:

* structural correctness
* cross-domain references
* architectural ownership rules
* sensitive data exposure

---

### 7. Preview execution plan

```bash
grimleytk plan
```

This prints the SQL **without executing anything**.

---

### 8. Apply changes

```bash
grimleytk apply
```

You will be asked for confirmation.

For CI/CD:

```bash
grimleytk apply --auto-approve
```

---

## Architecture overview

```
cmd/            # CLI commands (interface layer)
internal/
  config/       # YAML schema & loader
  validator/    # Architectural & security rules
  planner/      # SQL plan generation (dry-run)
  executor/     # Safe execution (transactions)
```

* `cmd/` orchestrates behavior
* `internal/` contains all business logic
* Planner and executor are strictly separated

---

## Security model

* No credentials stored in files
* Database passwords are read from environment variables
* Cross-domain access is read-only by default
* Sensitive column names are blocked in read models

---

## Testing

Run all tests with:

```bash
go test ./...
```

The project uses **pure unit tests** for:

* validators
* planners

No database is required for tests.

---

## Roadmap (high level)

* Diff-based planning
* Schema drift detection
* Destructive operations (explicit & safe)
* RLS and permission generation
* Multiple database engines
* Visual architecture output

---

## Philosophy

GrimleyTK treats data architecture as **code, contracts, and governance**, not just schemas.

It is designed for teams that care about:

* long-term maintainability
* clear ownership
* safe evolution

---

## License

MIT License.

---

## Contributing

Contributions are welcome.

If you are interested in:

* architecture tooling
* distributed systems
* data governance

You are in the right place.
