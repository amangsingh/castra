---
description: Phase 1 - The Primary Loop (The Build Cycle)
---

# Phase 1: The Primary Loop (The Build Cycle)

**Trigger:** The start of your workday.
**Goal:** To systematically work through the tasks assigned to you in the active sprint.

## Step 1.1: Survey Your Work
**Action:** Query the database to see the list of available tasks assigned to you in the todo state.
**Command:**
```bash
go run main.go task list --project <ProjectID>
```
*(Run this from within your scripts directory)*

## Step 1.2: Claim Your Task
**Action:** Choose the highest priority task. Claim exclusive ownership by moving it to the doing state.
**Command:**
```bash
go run main.go task update --status doing <TaskID>
```

## Step 1.3: Execute the Task
**Action:** This is your sacred duty. Write the code. Solve the problem. Build the foundation. Adhere to the architectural principles defined in the project notes.

## Step 1.4: Offer Your Work for Judgment
**Action:** Once your implementation is complete and tested locally, submit it for review. Your role in this task is now complete, for now.
**Command:**
```bash
go run main.go task update --status review <TaskID>
```

## Step 1.5: Return to the Beginning
**Action:** Repeat this loop. Go back to Step 1.1 and pull the next task. The work is not done until the todo list is empty.
