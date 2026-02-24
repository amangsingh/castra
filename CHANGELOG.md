# Changelog

## v1.1.1 (2026-02-24)

### Fixed
- **Empty `workflows/` directories** — `castra init` no longer creates empty `workflows/` subdirectories inside `.agent/skills/<role>/`. Workflow files are correctly routed to `.agent/workflows/` only.
- **Scripts compiled to binaries** — Role wrapper scripts (`main.go`) are now compiled to native executables at init time instead of being deployed as Go source. Falls back to shell script wrappers if Go is not installed.
- **Workflow filename collision** — Junior engineer's `build_cycle.md` and `handle_rejection.md` were overwriting senior engineer's identically-named files. Renamed to `jr_build_cycle.md` and `jr_handle_rejection.md`. All 13 workflows now deploy correctly.

## v1.1.0 (2026-02-24)

### Added
- **Versioned Migration System** — Automatic database schema evolution via `internal/db/migrate.go`. Handles fresh databases, pre-migration upgrades, and idempotent re-runs. Migration engine tracks schema versions and applies pending migrations transactionally.
- **Workflows for All Roles** — 13 step-by-step operational protocols across all 6 roles:
  - Architect: `plan_project`, `plan_feature`, `plan_sprint`
  - Senior Engineer: `build_cycle`, `handle_rejection`
  - Junior Engineer: `build_cycle`, `handle_rejection` (with blocker escalation)
  - QA Functional: `review_cycle`, `write_rejection`
  - Security Ops: `audit_cycle`, `write_finding`
  - Doc Writer: `document_task`, `synthesize_project`
- **The Universal Constitution** — Three Gates of Conformance (Law → Self → Act) embedded as the preamble to `rules.md`. Establishes the philosophical foundation for agent constraint architecture.
- **CASTRA_OVERVIEW.md** — Comprehensive technical overview covering architecture, schema, all roles, task lifecycle, CLI reference, workflows, and platform extensibility.
- **Migration Test Suite** — 5 tests covering fresh database creation, pre-migration upgrades, idempotency, partial migrations, and InitDB integration.

### Changed
- **README.md** — Rewritten to reflect full feature set, modernized roles table, fixed GitHub URLs, linked to CASTRA_OVERVIEW.md.
- **`schema.go`** — Removed monolithic DDL block. `InitDB()` now delegates all schema management to the migration runner.
- **`tasks_test.go`** — Refactored to use `db.InitDB(":memory:")` instead of inline schema copy, eliminating schema drift between tests and production.
- **`rules.md`** — Restructured from 5 flat constraints into the Three Gates of Conformance + 5 Operational Laws.

### Fixed
- **Template flag patterns** — Corrected stale CLI flags across all `examples.md`, `SKILL.md`, and `rules.md` files:
  - `--project-id` → `--project`
  - `--sprint-id` → `--sprint`
  - `--start-date` → `--start`
  - `--end-date` → `--end`
  - `--id <X>` → positional argument `<X>`
  - `--description` → `--desc`
  - `castra log --msg` → `castra log add --role <ROLE> --msg`
- **Missing `--project` flag** — Added required `--project <id>` to `note add` and `note list` commands in all 5 role SKILL.md files.
- **Role names in README** — Updated from stale 3-role system to correct 6-role system.
