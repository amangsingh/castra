# Changelog

## v2.0.0 (2026-03-15)

### Added
- **SDLC Tier-3 Maturity**: Achieving full compliance with high-integrity "Airtight State" requirements.
  - **Hierarchical Milestones**: Support for nested milestones and complex feature roadmaps.
  - **Scoped Archetypes**: Project-specific status pipelines enabling granular workflow enforcement.
  - **Universal Audit Log**: Transactional recording of every state change with immutable role signatures.
  - **Gate 0 Alignment**: Mandatory pre-verification of system status before state modifications.
- **HATEOAS Interaction Protocol**: Introduced a dynamic Affordance Engine that surfaces available CLI commands based on current task state and role permissions.
- **Multi-Vendor Coordination Matrix**: Unified project initialization and coordination across **Antigravity**, **Claude Code**, **GitHub Copilot**, and **Google Gemini Code Assist**.
- **Terminal UI (TUI) Modernization**:
  - **Live Hierarchy Visualization**: Real-time rendering of hierarchical milestones in the dashboard.
  - **Interactive Affordance Bar**: Visual feedback of HATEOAS-driven available actions.
  - **Audit Drill-down**: Real-time inspection of transactional history for any task.
  - **Break-Glass Highlighting**: Visual indicators for emergency overrides and post-incident reviews.
- **Persona Linter (Gate 3)**: Systemic enforcement of `SKILL.md` constraints to eliminate agent persona drift.
- **Multi-Vendor project initialization**: Unified `castra init` command with agnostic workspace generation for all supported platforms.

### Changed
- **Database Schema (v11)**: Migrated core storage to support hierarchical structures and scoped status pipelines.
- **Core Library Refactor**: Centralized business logic (RBAC, dual-approvals) to use the new Hierarchical SDLC engine.
- **Simplified CLI Experience**: Positioned as a single-entry protocol (`castra init`) with dynamic vendor-agnostic scaffolding.

### Fixed
- **NULL Scan Regressions**: Fixed a critical crash in the core library when scanning older tasks with NULL descriptions (legacy v1.x data).
- **Hierarchical Integrity**: Resolved an issue where child milestones could lose their parent pointers during project restoration.
- **Audit Consistency**: Ensured that rejections correctly reset both functional and security approval flags in a single atomic transaction.


## v1.4.0 (2026-03-10)

### Added
- **Terminal UI (TUI) & Daemon Watcher**: Integrated a live dashboard view (`castra tui`) and headless background monitoring (`castra watch`).
- **New AI Generators**: Added support for Copilot (`--copilot`) and Gemini Code Assist (`--gemini`) project scaffolding.
- **Designer Persona**: Introduced the `designer` role for UI/UX mocks and integrated it across CLI roles and platform generators.
- **Sprint Automation**: Sprints automatically start when the first task is picked up and auto-complete when all tasks are done.
- **Session Identity Enforcement**: Systemic checks enforce role boundaries to prevent agent persona drift.
- **Test Coverage Setup**: Added Golden File framework testing for generators and comprehensive unit testing for command router and CLI logic.
- **Auto-Audit via Router (`MutatingCommand` interface)**: All state-mutating commands (`task.add`, `task.update`, `project.add`, etc.) now automatically write a best-effort audit entry on success without any caller-side instrumentation. Implemented via the new optional `MutatingCommand` interface in `internal/commands/router.go`.
- **Integration Test Suite for Auto-Audit**: Added `internal/commands/auto_audit_test.go` with 6 tests covering: mutating commands produce entries, read-only commands do not, failed commands do not produce spurious entries, and the role is correctly captured in the audit record.
- **Hardened Universal Constitution (Law 2 & Law 9)**: Rewrote Law 2 (Role Boundaries) from a 3-bullet generic list into per-role prohibitions covering all 7 system roles. Added Law 9 (The Law of Dual Approval) explicitly mandating the QA-first → Security-second approval sequence with prohibition on out-of-order approvals. Updated `qa-functional`, `security-ops`, and `junior-engineer` SKILL.md files to surface these constraints inline.
- **Law 3 Updated (Mandatory Auditing)**: Distinguished between automatic audit entries (written by the router) and manual context log entries (written by the agent via `castra log add`), eliminating ambiguity for agents receiving dual-audit instructions.

### Changed
- **CLI Architecture Refactor**: Separated business logic from the CLI CRUD layer by introducing a Command Router Interface.
- **Central Shared Templates Package**: Consolidated templates and workflows to dry up generation code, adding YAML frontmatter and removing legacy directories.

### Fixed
- **Database Stability**: Resolved SQLite 'out of memory (14)' errors within Antigravity Sandbox isolation.

## v1.2.2 (2026-02-25)

### Added
- **Command:** `castra task view <id>` provides a comprehensive context snapshot of a task, aggregating its properties, description, attached role-filtered notes, and audit logs.
- **Agent Directives:** Added Law 6 (Command Structure) to `rules.md` to prevent agent syntax hallucination.
- **Workflow Mandate:** Explicitly commanded all agents in their `SKILL.md` to strictly follow defined workflows, reducing "guessing."
- **Context Gathering:** Updated all worker workflows (`build_cycle`, `review_cycle`, `audit_cycle`, `document_task`) to mandate running `castra task view` before taking action.

## v1.2.1 (2026-02-25)

### Added
- **Codebase Ingestion Workflow**: Introduced `ingest_project.md` to guide Architects working with existing codebases to retroactively map history to `workspace.db` using God Mode.
- **AI-Speed Sprints**: Reframed Sprint workflows, documentation, and command examples to treat Sprints as minutes-long "Iteration Batches" rather than mandatory calendar weeks.

## v1.2.0 (2026-02-25)

### Added
- **Native Milestones Component** — Separated feature planning ("What") from time-boxed scheduling ("When").
  - Added new `castra milestone add|list|update` commands (Architect only).
  - Added `milestone_id` to `tasks`. Tasks can now concurrently belong to both a Milestone and a Sprint.
  - Added `--milestone <id>` flag to `castra task add` and `list` commands.
- **Architect Workflows Rewritten**
  - `plan_feature.md`: Decomposes roadmap into Milestones with high-level Milestone tasks (no more pseudo "Feature Sprints").
  - `plan_sprint.md`: Distills high-level Milestone tasks into Sprint-assigned granular work.
- **Documentation**: Updated `CASTRA_OVERVIEW.md` and `README.md` to reflect the new three-tier structure (Project -> Milestone -> Sprint -> Task).

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
