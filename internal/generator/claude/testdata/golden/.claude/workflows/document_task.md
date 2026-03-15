---
description: Phase 1 - Task Documentation (Feature Chronicle)
---

# Phase 1: Task Documentation (Feature Chronicle)

**Trigger:** A task transitions to the `done` state (both QA and Security have approved).
**Goal:** To produce clear, human-readable documentation for the completed feature.

## Step 0: Log Your Intent
**Action:** Before starting any research or implementation, log your intent to work on the task. This ensures universality of surveillance and record-keeping.
**Command:**
```bash
castra log add --role <role> --msg "Starting work on task <TaskID>" --type task --entity <TaskID>
```

## Step 1.1: Survey Completed Work
**Action:** Query the database for tasks in the `done` state that have not yet been documented. Cross-reference with your own notes tagged `docs-link` to find undocumented tasks.
**Command:**
```bash
castra task list --project <ProjectID>
```

## Step 1.2: Gather the Full Context
**Action:** For the target task, read its complete context including title, description, attached notes, and audit logs. This contains the architect's intent, the engineer's implementation notes, and any QA/Security feedback. This is your source material.
**Command:**
```bash
castra task view <TaskID>
```

## Step 1.3: Produce the Documentation
**Action:** Using the chat interface, produce the documentation artifact as raw markdown. The document should contain:
- **What was built:** A clear, jargon-free description of the feature.
- **Why it was built:** The business or architectural reasoning (from architect notes).
- **How it works:** A user-facing explanation — not implementation details, but observable behavior.
- **Any caveats or known limitations.**

Write the documentation to the appropriate file in the project (e.g., `docs/`, `README.md` section, or a standalone doc).

## Step 1.4: Log the Documentation Link
**Action:** Record that this task has been documented by adding a note with the file path or URL.
**Command:**
```bash
castra note add --project <ProjectID> --task <TaskID> --content "Documented at: docs/<feature-name>.md" --tags "docs-link"
```
