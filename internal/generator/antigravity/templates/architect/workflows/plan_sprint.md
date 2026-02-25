---
description: Phase 3 - Tactical Planning (Work Order Generation)
---

# Phase 3: Tactical Planning (Work Order Generation)

**Trigger:** A decision is made to begin active development on a specific Milestone from Phase 2.
**Goal:** To break down a high-level roadmap task into a series of small, concrete, and immediately actionable tasks for the engineering team, and schedule them into a time-box.

## Step 3.1: Create a Time-Boxed "Active Sprint"
**Action:** Create a new, short-term sprint for the upcoming development cycle (e.g., an "Iteration Batch" or "Session"). Sprints manage the *when*.
**Dates are Optional:** For AI agents, a sprint might take 10 minutes, so `start` and `end` dates/labels are purely optional strings for context (e.g., `--start "Session 1"` or omitted entirely).
**Command:**
```bash
castra sprint add --role architect --project <ProjectID> --name "Iteration 1: Backend Scaffolding"
```
**Capture:** Record the Active Sprint ID.

## Step 3.2: Decompose Roadmap Tasks into Granular Work Orders
**Action:** Select a high-level roadmap task from a Milestone (e.g., "API & Database Design"). Analyze it and break it down into small, well-defined engineering tasks that can be completed by one person.
**Task Granularity:** A good task should represent 1-2 days of work and have a clear, verifiable outcome.
**Linking:** Assign these granular tasks to *both* the active Sprint (the time-box) and the parent Milestone (the feature tracker).
**Examples:** "Create 'users' table migration," "Implement /auth/register endpoint," "Implement JWT generation service," "Add password hashing logic."
**Command (repeat for each granular task):**
```bash
castra task add --role architect --project <ProjectID> --milestone <MilestoneID> --sprint <ActiveSprintID> --title "Create 'users' table migration" --desc "The migration must include fields for id, email, password_hash, created_at, and updated_at."
```

## Step 3.3: Task Cleanup
**Action:** Once the high-level roadmap task has been fully decomposed into granular sprint tasks, delete the original high-level placeholder task to keep the backlog clean.
**Command:**
```bash
castra task delete --role architect <TaskID>
```

**Output:** A time-boxed "Active Sprint" is now filled with multiple granular tasks, all correctly tethered to their parent Milestone. The engineering team now has a ready and waiting backlog of work. The system is primed for execution.
