---
description: Phase 1 - Strategic Planning (Project Inception)
---

# Phase 1: Strategic Planning (Project Inception)

**Trigger:** A new, high-level project directive from the user.
**Goal:** To establish the project's "North Star" — its existence and its core architectural principles — within the database.

## Step 1.1: Create the Project Container
**Action:** Execute the castra project add command. The name should be thematic, and the description should be a one-sentence summary of the vision.
**Command:**
```bash
castra project add --role architect --name "Project Name" --desc "High-level project description."
```
**Capture:** Record the Project ID from the command's output. This ID is required for all subsequent steps.

## Step 1.2: Define the Architectural Constitution
**Action:** Create a foundational note that contains the high-level architectural vision. This is the most critical document for guiding the entire project.
**Content:** The note's content MUST include sections for:
- **Core Principles:** (e.g., "Local-first," "Stateless services," "Zero-trust security").
- **Technology Stack:** (e.g., "Go, SQLite, gRPC, React").
- **High-Level Components:** (e.g., "API Server, Web UI, Database, Worker Queue").
- **Non-Negotiable Constraints:** (e.g., "Must run on ARM64," "No GPLv3 libraries").
**Command:**
```bash
castra note add --role architect --project <ProjectID> --tags "vision,architecture" --content "..."
```

**Output:** A new project exists in the projects table, linked to a single, foundational note in the project_notes table. The strategic foundation is set.
