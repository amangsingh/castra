# Castra

**The Universal Protocol for Agentic Software Development**

Castra is a standalone Go binary that serves as the coordination protocol for autonomous digital work. It operates directly on a local SQLite database (`workspace.db`), eliminating the complexity of client-server architectures while enforcing a role-based, dual-approval workflow.

> **"The workspace.db is the only truth."** — *The Universal Constitution*

For a comprehensive technical deep-dive, see [CASTRA_OVERVIEW.md](CASTRA_OVERVIEW.md).

## Features

*   **Zero-Config:** Runs entirely locally. No servers, no clouds, no accounts. Single standalone binary.
*   **Role-Based Access Control (RBAC):** Enforces separation of concerns via the `--role` flag.
    *   **Architect (`architect`):** The planner (God mode).
    *   **Designer (`designer`):** The Visual Architect (UI mockups & user flows).
    *   **Senior Engineer (`senior-engineer`):** The builder (Complex problems).
    *   **Junior Engineer (`junior-engineer`):** The maintainer (Bug fixes, tweaks).
    *   **Functional QA (`qa-functional`):** The verifier (Review gatekeeper).
    *   **Security Ops (`security-ops`):** The auditor (Review gatekeeper).
    *   **Doc Writer (`doc-writer`):** The scribe (Documentation).
*   **Dual-Approval Locks:** Tasks cannot be marked `done` without explicit approval from **both** QA and Security. Rejections reset all approval flags.
*   **Task-Level Notes:** Notes scoped to specific tasks enable the rejection feedback loop — QA/Security attach structured rejection reasons directly to the task.
*   **Audit Trail:** Immutable audit log with auto-logging of status changes, approvals, and rejections.
*   **Versioned Schema Migrations:** Automatic database schema evolution with backward compatibility for pre-migration databases.
*   **Workflow System:** Step-by-step operational protocols for every role, generated into `.agent/workflows/` on init.
*   **The Universal Constitution:** The Three Gates of Conformance — a constraint architecture that eliminates LLM compliance drift.
*   **Terminal User Interface (TUI):** Includes a live dashboard view (`castra tui`) for seamless project state monitoring, alongside a headless daemon watcher (`castra watch`).
*   **Sprint Automation:** Sprints intelligently auto-start when work begins and auto-complete when all tasks are approved.
*   **Session Identity Enforcement:** Real-time systemic checks to prevent agent persona drift.
*   **Platform Extensibility:** Generator abstraction layer with out-of-the-box support for DeepMind Antigravity (`--antigravity`), GitHub Copilot (`--copilot`), and Gemini Code Assist (`--gemini`).

## Installation

### Pre-built Binaries
Download the latest binaries from [GitHub Releases](https://github.com/AmanSingh494/castra/releases).

**macOS / Linux:**
```bash
# Download and install
chmod +x castra-mac
xattr -d com.apple.quarantine castra-mac   # macOS only, if downloaded via browser
mv castra-mac /usr/local/bin/castra

# Verify
castra init --antigravity
```

**Windows:**
1. Rename `castra-windows.exe` to `castra.exe`.
2. Move it to a folder in your `%PATH%`.

### Build from Source
Requirements: [Go](https://go.dev/) 1.22+

```bash
git clone https://github.com/AmanSingh494/castra.git
cd castra
go build -o castra .
sudo mv castra /usr/local/bin/
```

### Troubleshooting: macOS Sandbox (Antigravity)
If you are running Castra within a sandboxed agent environment on macOS (e.g., using Antigravity) and encounter an `out of memory (14)` SQLite error, this is a **sandbox file system restriction**, not a true memory issue. 

To resolve this, ensure the agent or terminal environment running Castra has **sandbox mode disabled**. Castra is a local CLI developer tool and requires unrestricted file system access to create temporary SQLite journal files within the user's workspace.

## Usage

**1. Initialize a Workspace**
```bash
castra init --antigravity
```

**2. Create a Project (Architect)**
```bash
castra project add --role architect --name "Project Alpha" --desc "Next-gen AI platform"
```

**3. Define a Milestone (Architect)**
```bash
castra milestone add --role architect --project 1 --name "User Authentication"
```

**4. Schedule an Iteration / Sprint (Architect)**
```bash
castra sprint add --role architect --project 1 --name "Iteration 1"
```

**5. Add Tasks representing Work (Architect)**
```bash
# Add a task to both the milestone (the feature) and the sprint (the timeline)
castra task add --role architect --project 1 --milestone 1 --sprint 1 --title "Setup DB" --desc "SQLite schema for users"
```

**6. Engineer Works**
```bash
# Read full task blueprint (description, notes, logs)
castra task view --role senior-engineer 1

# Claim and execute
castra task update --role senior-engineer --status doing 1
castra task update --role senior-engineer --status review 1
```

**7. Review & Approve (QA & Security)**
```bash
# QA approves (waits for Security)
castra task update --role qa-functional --status done 1

# Security approves (task transitions to done)
castra task update --role security-ops --status done 1
```

**8. Audit Trail**
```bash
castra log list --role architect
castra log add --role architect --msg "Completed Sprint 1 planning"
```

## The Roles

| Role | Flag | Power | Prohibition |
|------|------|-------|-------------|
| **Architect** | `--role architect` | God mode. Manages everything. | Cannot write code. |
| **Designer** | `--role designer` | Visualizes intent into UI and application flow mockups. | Cannot write production code. |
| **Senior Engineer** | `--role senior-engineer` | Complex problem solving. | Cannot mark `done`. |
| **Junior Engineer** | `--role junior-engineer` | Bug fixes, tweaks. | Cannot mark `done`. Cannot architect. |
| **Functional QA** | `--role qa-functional` | First key to `done`. | Cannot read source code. |
| **Security Ops** | `--role security-ops` | Second key to `done`. Veto is absolute. | Only cares about security. |
| **Doc Writer** | `--role doc-writer` | Documents completed work. | Read-only. Cannot change status. |

## License

MIT
