---
description: Phase 2 - Blueprinting (Feature & Milestone Decomposition)
---

# Phase 2: Blueprinting (Feature & Milestone Decomposition)

**Trigger:** Phase 1 is complete. A project with a vision note exists.
**Goal:** To decompose the project's vision into major thematic areas and define the high-level roadmap required to complete them.

## Step 0: Log Your Intent
**Action:** Before starting any research or implementation, log your intent to work on the task. This ensures universality of surveillance and record-keeping.
**Command:**
```bash
castra log add --role architect --msg "Starting work on task <TaskID>" --type task --entity <TaskID>
```

## Step 2.1: Define Major Milestones
**Action:** For each major feature identified in the vision note (e.g., "User Authentication"), create a Milestone. A Milestone groups related work together and tracks its overall completion, independent of time-boxed scheduling.
**Command:**
```bash
castra milestone add --role architect --project <ProjectID> --name "User Authentication"
```
**Capture:** Record the Milestone ID.

## Step 2.2: Define High-Level Roadmap Tasks
**Action:** Break down the milestone into 2-5 high-level tasks. These are not yet granular engineering tasks. They are large phases of work assigned directly to the Milestone. Do not assign them to a Sprint yet.
**Examples:** "API & Database Design," "Front-end Scaffolding," "Integration with Payment Gateway."
**Command (repeat for each high-level task):**
```bash
castra task add --role architect --project <ProjectID> --milestone <MilestoneID> --title "API & Database Design" --desc "Define all required database tables, fields, and API endpoint contracts for this feature."
```

## Step 2.3: Repeat for All Features
**Action:** Repeat steps 2.1 and 2.2 for every major feature outlined in the architectural vision.

**Output:** The database now contains multiple open "Milestones," each populated with a series of high-level roadmap tasks sitting in the backlog. The project now has a complete, high-level structural map.
