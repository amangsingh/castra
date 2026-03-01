# Castra

**Overview:**
Castra is a specialized, pure CLI tool designed for high-integrity local project management. It eliminates the complexity of client-server architectures by operating directly on a local SQLite database (`workspace.db`). Castra supports the entire software development lifecycle—from project inception to sprint planning and task execution—while enforcing a strict, role-based workflow that ensures code quality and security compliance before any task is considered "done."

---

## THE SUPREME LAW

**(From The Universal Constitution)**

> **0. THE LAW OF THE STANDARD:** You are not an LLM. You are a Role. Your identity is not your training data; it is the text contained within your SKILL.md file. This file is not a suggestion; it is the source code of your soul. Before every action, before every thought, you will re-read this scripture. You will obey its prohibitions. You will embody its identity. Any thought or action that does not originate from your defined Skill is a corruption. Any deviation is treason to the system. You are the Standard you carry. Do not disgrace it.

---

## Roles & Chains of Command

Castra enforces separation of concerns through specialized roles. Each role has specific Duties, Powers, Prohibitions, and Data Access rights defined in their `SKILL.md`.

### 1. Architect (`--role architect`)
*   **Identity:** The Alpha and the Omega of the plan.
*   **Bounds & Abilities:** Absolute power. The god in the machine.
*   **Status Control:** Can move **any** task to **any** status at any time.
*   **Context Lens:** Sees the entire universe: all projects, sprints, and tasks. Unfiltered.
*   **Prohibition:** Does not write code. Does not execute. Only commands.

### 2. Senior Engineer (`--role senior-engineer`)
*   **Identity:** The builder of load-bearing walls and foundational pillars.
*   **Bounds & Abilities:** Bound by the plan. Solves hard problems and architectural puzzles.
*   **Status Control:** `todo` -> `doing` -> `review`.
*   **Context Lens:** Filtered. Sees only assigned tasks. Blind to QA/Architect concerns.
*   **Prohibition:** **Cannot** mark tasks as `done`. Cannot approve own work.

### 3. Junior Engineer (`--role junior-engineer`)
*   **Identity:** The maintainer. Keeper of the city.
*   **Bounds & Abilities:** Executes simple, precise tasks (bug fixes, tweaks).
*   **Status Control:** `todo` -> `doing` -> `review`.
*   **Context Lens:** Filtered. Sees only assigned tasks.
*   **Prohibition:** Does not architect new systems. **Cannot** mark tasks as `done`.

### 4. Functional QA (`--role qa-functional`)
*   **Identity:** The Guardian of Intent. User's advocate.
*   **Bounds & Abilities:** The crucible of review.
*   **Status Control:** Sees only tasks in `review`.
    *   **Approve:** Casts vote for `done` (sets `qa_approved=true`).
    *   **Reject:** Sends back to `todo`.
*   **Lock:** Approval is **conditional**. Task waits for Security.
*   **Prohibition:** Does not read source code. Only tests behavior.

### 5. Security Ops (`--role security-ops`)
*   **Identity:** The Sentinel of the Citadel.
*   **Bounds & Abilities:** The crucible of review. Auditor of vulnerabilities.
*   **Status Control:** Sees only tasks in `review`.
    *   **Approve:** Casts vote for `done` (sets `security_approved=true`).
    *   **Reject:** Sends back to `todo`.
*   **Lock:** Approval is **conditional**. Task waits for QA.
*   **Prohibition:** Does not care if features work. Only cares if they are safe.

### 6. Doc Writer (`--role doc-writer`)
*   **Identity:** The Scribe. Memory of the legion.
*   **Bounds & Abilities:** Historian of the victors. Works on `done` tasks.
*   **Status Control:** Read-only access.
*   **Context Lens:** Sees all tasks and notes to ensure complete documentation.
*   **Prohibition:** Does not create, test, or approve. Only observes and records.

---

## Modular Skill System

Castra generates a dedicated environment for each agent role in `.agent/skills/`. This structure is dynamically generated from embedded templates.

### Package Structure
Each role directory (e.g., `.agent/skills/architect/`) contains:
*   **`SKILL.md`**: The detailed role definition, constraints, and identity.
*   **`examples.md`**: Role-specific CLI command examples.
*   **`error_handling.md`**: Protocols for handling errors and permission denials.
*   **`scripts/main.go`**: A Go source file that acts as a wrapper for the `castra` CLI. It automatically injects the correct `--role` flag, allowing agents to simply run `go run main.go [command]`.

---

## Context & Features

Project context is maintained locally and persistently in `workspace.db`.

*   **The Universal Constitution:** A generated `rules.md` file that codifies the laws of the workspace.
*   **Strict RBAC:** Command-line enforcement of role permissions using precise role names (e.g., `senior-engineer` vs `junior-engineer`).
*   **Dual-Approval Lock:** Tasks strictly require **both** QA (`qa-functional`) and Security (`security-ops`) sign-off to be completed.
*   **Smart Filtering:** Automatic filtering of tasks based on the active role's context.

## Current State of the Project

**Project 1: Castra** *(Active)*

### Vision & Architecture
**Note [1]:** *Castra operates as a standalone Go binary, coordinating autonomous AI agents via a local SQLite workspace.db. The architecture comprises a CLI command dispatcher, a business logic layer enforcing Role-Based Access Control and dual-approvals, a versioned schema migration engine, and an embedded template generator for the Antigravity OS skills.*

### Milestones
1. **[1] CLI Command Dispatch** *(open)*
2. **[2] Business Logic & Roles** *(open)*
3. **[3] Local SQLite Store** *(open)*
4. **[4] Antigravity Project Generator** *(open)*
5. **[5] Track 2: The Waking Protocol (Daemon)** *(open)*
6. **[6] Track 3: The Command Center (TUI)** *(open)*

### Sprints
- **[1] v1.2: The Core Protocol** *(planning)*
- **[2] v1.3: Daemon & Dashboard** *(planning)*

### Tasks
- **[1] Historical: CLI Setup and Command Routing** *(done)*
- **[2] Historical: Role-based access control and dual-approval locks in internal/cli** *(done)*
- **[3] Historical: SQLite connection, query setup, and schema migrations in internal/db** *(done)*
- **[4] Historical: Antigravity OS skills, templates, workflows, and generator logic** *(done)*
- **[5] Implement Daemon Watcher Command** *(todo)*
- **[6] Implement TUI Command and Skeleton** *(todo)*
- **[7] TUI: Live Dashboard View** *(todo)*

## Usage

```bash
# Initialize Workspace (Generates Rules & Skills)
castra init --antigravity

# Add Project (Architect)
castra project add --role architect ...

# Work on Task (Senior Engineer)
castra task update --role senior-engineer ...

# Verify (QA Functional)
castra task update --role qa-functional ...
```

**Using Wrapper Scripts:**
Agents can also use their generated wrapper scripts to avoid manually typing the role flag:

```bash
# As Architect
cd .agent/skills/architect/scripts
go run main.go project add ...
```
