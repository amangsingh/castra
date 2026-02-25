---
description: Phase 0 - Reality Sync (Existing Codebase Ingestion)
---

# Phase 0: Reality Sync (Existing Codebase Ingestion)

**Trigger:** Castra is initialized in a directory that already contains a significant, pre-existing codebase. 
**Goal:** To establish the "till date" structural map of reality into the `workspace.db`, allowing you to plan *new* work on top of what already exists.

## Step 0.1: Analyze & Synthesize
**Action:** Review the existing codebase. Understand its architecture, major components, and current state.
**Output:** You should have a mental map of the system's "North Star" and its major feature domains.

## Step 0.2: Create the Project Container
**Action:** Execute the `castra project add` command to establish the project.
**Command:**
```bash
castra project add --role architect --name "Project Name" --desc "High-level description of the inherited codebase."
```
**Capture:** Record the Project ID.

## Step 0.3: Define "Till Date" Reality
**Action:** Add a foundational note containing the high-level architectural overview of what you've just analyzed.
**Command:**
```bash
castra note add --role architect --project <ProjectID> --tags "vision,architecture,ingestion" --content "..."
```

## Step 0.4: Define Existing Milestones
**Action:** For each major feature domain that *already exists* in the codebase, create a Milestone.
**Command (repeat for each domain):**
```bash
castra milestone add --role architect --project <ProjectID> --name "Existing Feature: User Auth"
```
**Capture:** Record the Milestone IDs.

## Step 0.5: Retroactive Tasking (The 'God Mode' Fast-Forward)
**Action:** For each Milestone created in Step 0.4, create high-level tasks representing the structural components that have *already been built*. Then, use your God Mode authority (`Any -> Any` status control) to immediately mark them as `done`.
**Command Sequence (repeat for each major built component):**
```bash
# 1. Create the historical task
castra task add --role architect --project <ProjectID> --milestone <MilestoneID> --title "Historical: API & Database Auth Schema" --desc "Completed prior to Castra ingestion."

# 2. Immediately mark it done (Note the Task ID from the previous command)
castra task update --role architect --status done <TaskID>
```

**Output:** The `workspace.db` now perfectly mirrors the reality of the existing codebase. The historical work is tracked as `done`, and you are now historically grounded and positioned to execute Phase 2 (Blueprinting) for *new* features.
