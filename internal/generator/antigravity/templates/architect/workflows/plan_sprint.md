---
description: Phase 3 - Tactical Planning (Work Order Generation)
---

# Phase 3: Tactical Planning (Work Order Generation)

**Trigger:** A decision is made to begin active development on a specific milestone from Phase 2.
**Goal:** To break down a single, high-level "Milestone Task" into a series of small, concrete, and immediately actionable tasks for the engineering team.

## Step 3.1: Create a Time-Boxed "Active Sprint"
**Action:** Create a new, short-term sprint for the upcoming development cycle (e.g., one or two weeks). The name should be specific and indicate the cycle.
**Command:**
```bash
castra sprint add --role architect --project <ProjectID> --name "Sprint 1: Auth Backend (Week of Feb 19)" --start "2026-02-19" --end "2026-02-26"
```
**Capture:** Record the Active Sprint ID.

## Step 3.2: Decompose the Milestone into Granular Tasks
**Action:** Select a single "Milestone Task" from a "Feature Sprint" (e.g., "Milestone: API & Database Design"). Analyze it and break it down into small, well-defined engineering tasks that can be completed by one person.
**Task Granularity:** A good task should represent 1-2 days of work and have a clear, verifiable outcome.
**Examples:** "Create 'users' table migration," "Implement /auth/register endpoint," "Implement JWT generation service," "Add password hashing logic."
**Command (repeat for each granular task):**
```bash
castra task add --role architect --project <ProjectID> --sprint <ActiveSprintID> --title "Create 'users' table migration" --desc "The migration must include fields for id, email, password_hash, created_at, and updated_at."
```

**Output:** A time-boxed "Active Sprint" is now filled with multiple todo tasks, each with a clear title and description. The engineering team now has a ready and waiting backlog of work. The system is primed for execution.
