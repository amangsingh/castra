---
description: Phase 2 - Project Synthesis (Aggregate Documentation)
---

# Phase 2: Project Synthesis (Aggregate Documentation)

**Trigger:** A direct command from the user/Architect to produce project-level documentation.
**Goal:** To synthesize the full project state into a comprehensive document — README, release notes, or project overview.

## Step 0: Log Your Intent
**Action:** Before starting any research or implementation, log your intent to work on the task. This ensures universality of surveillance and record-keeping.
**Command:**
```bash
castra log add --role <role> --msg "Starting work on task <TaskID>" --type task --entity <TaskID>
```

## Step 2.1: Read the Project Vision
**Action:** Start with the project description and the architect's foundational vision notes. These define the "why" and "what" at the highest level.
**Command:**
```bash
castra project list
castra note list --project <ProjectID>
```

## Step 2.2: Survey All Sprints and Tasks
**Action:** Understand the full scope of work that has been completed, is in progress, and is planned.
**Command:**
```bash
castra sprint list --project <ProjectID>
castra task list --project <ProjectID>
```

## Step 2.3: Produce the Synthesis
**Action:** Using the chat interface, produce the requested document. Common outputs:

**README.md:**
- Project name and one-line description
- Installation / setup instructions
- Usage examples
- Architecture overview (from architect notes)
- Contributing guidelines

**Release Notes (per sprint):**
- Sprint name and dates
- List of completed features (from `done` tasks)
- Known issues (from `blocked` or rejected tasks)
- Breaking changes (from architect notes tagged `breaking`)

**PROJECT_OVERVIEW.md:**
- Vision and goals
- Role descriptions
- Workflow summary
- Current state of the project

## Step 2.4: Log the Artifact
**Action:** Record that the synthesis has been produced.
**Command:**
```bash
castra note add --project <ProjectID> --content "Synthesized: <document name> at <file path>" --tags "docs-link,synthesis"
```
