---
description: Phase 1 - The Primary Loop (The Build Cycle)
---

# Phase 1: The Primary Loop (The Build Cycle)

**Trigger:** The start of your workday.
**Goal:** To systematically work through the simple, well-defined tasks assigned to you in the active sprint.

## Step 1.1: Survey Your Work
**Action:** Query the database for tasks in the `todo` state. Focus on tasks explicitly scoped for your skill level — bug fixes, refactors, dependency updates, and small component changes.
**Command:**
```bash
go run main.go task list --project <ProjectID>
```
*(Run this from within your scripts directory)*

## Step 1.2: Claim Your Task
**Action:** Choose the highest priority task that matches your scope. Claim exclusive ownership by moving it to the `doing` state.
**Command:**
```bash
go run main.go task update --status doing <TaskID>
```

## Step 1.3: Gather Context
**Action:** Before writing any code, read the task's full description, architectural notes, and any previous rejection logs. This is your comprehensive blueprint.
**Command:**
```bash
go run main.go task view <TaskID>
```

## Step 1.4: Execute the Task
**Action:** Focus. Read the task description carefully. Execute the work precisely as specified. Follow the architectural principles in the project notes. Do not over-engineer — solve exactly what the task asks for, nothing more.

**If you hit a blocker that exceeds your scope:**
```bash
go run main.go task update --status blocked <TaskID>
go run main.go note add --project <ProjectID> --task <TaskID> --content "BLOCKED: <description of what is blocking you>" --tags "junior-engineer,blocked"
```

## Step 1.5: Offer Your Work for Judgment
**Action:** Once your implementation is complete and tested locally, submit it for review. Your role in this task is now complete — for now.
**Command:**
```bash
go run main.go task update --status review <TaskID>
```

## Step 1.6: Return to the Beginning
**Action:** Repeat this loop. Go back to Step 1.1 and pull the next task. The work is not done until the todo list is empty.
