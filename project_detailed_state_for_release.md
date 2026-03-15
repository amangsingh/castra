# Castra v2.0: Comprehensive Release Master State

This document serves as the absolute source of truth for all project details, features, and technical specifications for Castra v2.0. Use this for marketing, promotion, and technical documentation synthesis.

---

## 1. Repository Structure (Clean State)

| Path | Purpose | Key Files |
| :--- | :--- | :--- |
| **`internal/config/`** | Configuration Management | `config.go` (Parses `castra.yaml`) |
| **`internal/db/`** | Persistence Layer | `migrate.go` (Schema v11 engine), `db.go` |
| **`internal/commands/`** | CLI Logic & Routing | `router.go` (RBAC), `registry.go` (Command Map) |
| **`internal/generator/`** | Scaffolding Engine | `generator.go`, `templates/` (Skill/Workflow source) |
| **`internal/persona/`** | Behavioral Enforcement | `linter.go` (Gate 3 persona validation) |
| **`internal/cli/`** | UI Components | `tui.go` (Bubbletea dashboard), `table.go` |
| **`root/`** | Project Metadata | `main.go`, `go.mod`, `README.md`, `TECHNICAL_SPEC.md` |

---

## 2. Definitive Feature List (v2.0)

### **A. SDLC Tier-3 Maturity (The "Airtight State")**
*   **Hierarchical Milestones**: First-class support for parent/child milestones, allowing complex project roadmaps.
*   **Scoped Archetypes**: Project-specific status pipelines that govern how tasks move from `todo` to `done`.
*   **Universal Audit Log**: Every database mutation is transactionally logged with the active role and timestamp.
*   **Gate 0 Alignment**: Mandatory pre-verification protocol ensuring the system state is valid before any change.

### **B. Interaction & Intelligence**
*   **HATEOAS Affordance Engine**: The CLI dynamically tells the agent which commands are valid based on the task's current status and the agent's role.
*   **Persona Linter**: System-level check that prevents agents from acting outside their `SKILL.md` constraints.
*   **Terminal UI (TUI) Dashboard**: A real-time, high-fidelity dashboard for monitoring project health, milestones, and task flow.
*   **Break-Glass Protocol**: A governed emergency override system for Architects that auto-generates Post-Incident Review tasks.

### **C. Compatibility**
*   **Multi-Vendor Coordination**: Unified initialization for Antigravity, Claude Code, GitHub Copilot, and Gemini Code Assist.
*   **Local-First Architecture**: Pure Go + SQLite. No external dependencies, zero-config, maximum privacy.

---

## 3. Categorized Command Reference

### **Workspace Management**
*   `castra init`: Initialize the workspace and generate vendor-specific rules/skills.
*   `castra tui`: Launch the real-time visual dashboard.
*   `castra watch`: Headless daemon for monitoring workspace changes.

### **Project & Roadmap**
*   `castra project add|list|view|update|delete`
*   `castra milestone add|list|view|update|delete`: Supports `--project` and optional `--parent` for hierarchy.
*   `castra sprint add|list`: Time-boxed iteration management.

### **Task & Execution**
*   `castra task add`: Arguments: `--project`, `--milestone`, `--sprint`, `--title`, `--desc`, `--prio`.
*   `castra task update`: The core engine. Flags: `--status`, `--role`.
*   `castra task view`: Aggregated view of description, notes, logs, and metadata.
*   `castra task list`: Filterable by project, sprint, milestone, or role.

### **Knowledge & Audit**
*   `castra note add|list`: Attached to projects or specific tasks.
*   `castra log add|list`: Manual context logs (separate from auto-audit logs).
*   `castra archetype add|list|delete`: Customize status pipelines.

---

## 4. The Role Matrix (Detailed)

| Role | Responsibility | Authority Level |
| :--- | :--- | :--- |
| **Architect** | Roadmap & Strategy | God Mode. Bypasses HATEOAS restrictions. |
| **Designer** | UI/UX & Flow | Definitional. Creates mocks and flows. |
| **Senior Engineer** | Core Implementation | Load-bearing code. Solves complex puzzles. |
| **Junior Engineer** | Maintenance | Tweaks, fixes, and routine operations. |
| **Functional QA** | Behavioral Verification | Gate 1. Approves functional correctness. |
| **Security Ops** | Vulnerability Audit | Gate 2. Final security veto. |
| **Doc Writer** | System Memory | Synthesis of history and project state. |

---

## 5. Technical Ground Truth (Internal Protocols)

### **The Dual-Approval Sequence**
To reach `status: done`, a task must achieve:
1. `qa_approved = 1`
2. `security_approved = 1`
**Note**: If an engineer updates a task in review, or a reviewer rejects it, both flags are instantly reset to zero to prevent stale approvals.

### **Database Migration History**
*   **v1-v5**: Initial MVP schema.
*   **v6-v8**: Added Audit logging and project-scoped archetypes.
*   **v11 (current)**: Introduced Hierarchical Milestones and refactored task-approval bypass flags for the Break-Glass protocol.

### **HATEOAS Status Flow**
Standard Flow: `todo` → `doing` → `review` → `done`.
*   `review` status is binary-locked: it is the only state where QA/Security can vote.
*   `done` status is legally unattainable without dual signatures (or Break-Glass).

---

## 6. Marketing Hook & Promotion Angles

1.  **"Stop Guessing, Start Governing"**: Market Castra as the OS for AI Agent squads.
2.  **"The Audit-Ready SDLC"**: Focus on the universal audit log for SOC 2 and regulated industries.
3.  **"Zero-Trust Agentic Work"**: Highlight how Castra prevents agent "hallucination" by grounding them in a persistent SQLite reality.
4.  **"Multi-Vendor Orchestration"**: Position Castra as the bridge between Claude, Gemini, and Copilot.
