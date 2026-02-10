# Castra

**The High-Integrity, Role-Based Project Management CLI**

Castra is a specialized, pure CLI tool designed for strict, local project management. It operates directly on a local SQLite database (`workspace.db`), eliminating the complexity of client-server architectures while enforcing a rigid, role-based workflow.

> **"The workspace.db is the only truth."** — *The Universal Constitution*

## Features

*   **Zero-Config:** Runs entirely locally. No servers, no clouds, no accounts.
*   **Role-Based Access Control (RBAC):** Enforces separation of concerns via the `--role` flag.
    *   **Architect:** The planner (God mode).
    *   **Engineer:** The builder (Execution only).
    *   **QA:** The verifier (Review gatekeeper).
    *   **Security:** The auditor (Review gatekeeper).
*   **Dual-Approval Locks:** Tasks cannot be marked `done` without explicit approval from **both** QA and Security.
*   **Contextual Notes:** Tag-based notes for role-specific communication (e.g., `#security` checklists).
*   **Universal Constitution:** Generates a set of immutable rules and role definitions (`.agent/rules/`) to guide AI agents or human users.

## Installation

### Pre-built Binaries
You can download the latest pre-built binaries for macOS, Linux, and Windows from our [GitHub Releases](https://github.com/yourusername/castra/releases) page.

**Installation Steps:**

1.  **Download** the binary for your platform (e.g., `castra-mac`, `castra-linux`, or `castra-windows.exe`).
2.  **Rename** the file to `castra` (or `castra.exe` on Windows).
3.  **Move** it to a directory in your system's `PATH`.

    *   **macOS / Linux:**
        ```bash
        mv castra-mac castra
        chmod +x castra
        sudo mv castra /usr/local/bin/
        ```
    *   **Windows:**
        Rename `castra-windows.exe` to `castra.exe` and move it to a folder in your `%PATH%` (e.g., `C:\Program Files\Castra\`).

4.  **Verify:** Run `castra` in your terminal. You should see the help message.

### Build from Source
Requirements: [Go](https://go.dev/) 1.22+

```bash
git clone https://github.com/yourusername/castra.git
cd castra
go build -o castra main.go
```

### Path Setup (Critical)
Castra relies on the binary being accessible in your system's `PATH`.

**Mac/Linux:**
```bash
sudo mv castra /usr/local/bin/
```

**Windows:**
Move `castra.exe` to a folder in your `%PATH%` (e.g., `C:\Windows\System32` or a custom tools folder).

## Usage

**1. Initialize a Workspace**
```bash
# Creates workspace.db and .agent/ rules (for Antigravity platform)
./castra init --antigravity
```

**2. Create a Project (Architect)**
```bash
./castra project add --role architect --name "Project Alpha" --desc "Next-gen AI"
```

**3. Define a Sprint (Architect)**
```bash
./castra sprint add --role architect --project 1 --name "Sprint 1" --start "2024-01-01" --end "2024-01-14"
```

**4. Add a Task (Architect)**
```bash
./castra task add --role architect --project 1 --sprint 1 --title "Setup DB" --desc "Use SQLite"
```

**5. Work on Task (Engineer)**
```bash
# Engineer moves task to 'doing'
./castra task update --role engineer --status doing 1

# Engineer moves task to 'review' (Requesting approval)
./castra task update --role engineer --status review 1
```

**6. Review & Approve (QA & Security)**
```bash
# QA approves (Task stays in 'review' until Security also approves)
./castra task update --role qa --status done 1

# Security approves (Task transitions to 'done')
./castra task update --role security --status done 1
```

## The Roles

*   **Architect (`--role architect`):** The planner. Manages projects, sprints, and tasks.
*   **Senior Engineer (`--role engineer`):** The builder. Solves complex problems. Cannot mark `done`.
*   **Junior Engineer (`--role engineer`):** The maintainer. Fixes bugs and tweaks. Cannot mark `done`.
*   **Functional QA (`--role qa`):** The verifier. Approves functionality.
*   **Security Ops (`--role security`):** The auditor. Approves compliance.
*   **Doc Writer (`--role doc-writer`):** The scribe. Read-only access to `done` tasks for documentation.

## License

MIT
