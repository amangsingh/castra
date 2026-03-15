# Castra

**The Universal Protocol for Agentic Software Development**

Castra v2.0 is a specialized, pure CLI tool and multi-vendor coordination protocol designed for high-integrity local project management. It eliminates the complexity of client-server architectures by operating directly on a local SQLite database (`workspace.db`). Castra supports the entire software development lifecycle (SDLC Tier-3)—from project inception to sprint planning and task execution—while enforcing a strict, role-based workflow and dual-approval gates.

> **"The workspace.db is the only truth."** — *The Universal Constitution*

## Documentation Index
- [TECHNICAL_SPEC.md](TECHNICAL_SPEC.md): Architecture, Schema, and Logic.
- [PROJECT_STATE.md](PROJECT_STATE.md): Current Release State and Audit results.
- [PERSONAS_WORKFLOWS.md](PERSONAS_WORKFLOWS.md): Detailed Persona definitions and Workflow protocols.
- [CHANGELOG.md](CHANGELOG.md): Project History.


## Features
*   **Zero-Config:** Runs entirely locally via a standalone binary.
*   **Role-Based Access Control (RBAC):** Enforces separation of concerns via the `--role` flag.
*   **Dual-Approval Locks:** Tasks require explicit approval from **both** QA and Security.
*   **HATEOAS Affordance Engine:** Dynamic command availability based on task state.
*   **Hierarchical Milestones**: Support for complex feature roadmaps and nesting.
*   **Agnostic multi-vendor initialization**: Support for Antigravity, Claude, Copilot, and Gemini.
*   **Terminal UI (TUI):** Real-time project monitoring and management.

## Personas & Roles
| Role | Identity | Authority |
| :--- | :--- | :--- |
| **Architect** | The Lawgiver | Plans, schedules, and commands (God mode). |
| **Designer** | The Shaper | Visualizes intent into UI and user flows. |
| **Senior Engineer** | The Core Builder | Implements complex blueprints into load-bearing code. |
| **Junior Engineer** | The Maintainer | Executes routine tasks and maintenance. |
| **Functional QA** | The Guardian of Intent | Verifies behavior against requirements. |
| **Security Ops** | The Sentinel | Audits code for security vulnerabilities. |
| **Doc Writer** | The Chronicler | Records the evolution of the system. |

## Workflows
Castra standardizes the following operational protocols:
- **Planning**: `plan_project` → `plan_feature` → `plan_sprint`.
- **Execution**: `build_cycle` (Survey → Claim → Execute → Submit).
- **Verification**: `review_cycle` (QA) & `audit_cycle` (Security).
- **Synthesis**: `document_task` & `synthesize_project`.

## Getting Started

### 1. Initialize Workspace
```bash
castra init
```

### 2. Create a Project (Architect)
```bash
castra project add --role architect --name "Project Alpha" --desc "Next-gen AI platform"
```

### 3. Work on Task (Senior Engineer)
```bash
castra task update --role senior-engineer --status doing <id>
castra task update --role senior-engineer --status review <id>
```

### 4. Approve (QA & Security)
```bash
castra task update --role qa-functional --status done <id>
castra task update --role security-ops --status done <id>
```

## License
MIT
