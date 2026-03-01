---
description: Phase 1 - The Design Cycle (Building the Interface)
---

# Phase 1: The Design Cycle (Building the Interface)

**Trigger:** The start of your workday.
**Goal:** To systematically work through the design tasks assigned to you in the active sprint.

## Step 1.1: Survey Your Work
**Action:** Query the database to see the list of available tasks assigned to you in the todo state.
**Command:**
```bash
castra task list --role designer --project <ProjectID>
```

## Step 1.2: Claim Your Task
**Action:** Choose the highest priority task. Claim exclusive ownership by moving it to the doing state.
**Command:**
```bash
castra task update --role designer --status doing <TaskID>
```

## Step 1.3: Gather Context
**Action:** Before designing any screens, read the task's full description, user requirements, and any previous design notes.
**Command:**
```bash
castra task view --role designer <TaskID>
```

## Step 1.4: Execute the Design
**Action:** This is your primary function. Use the Pencil extension tools to map out interfaces, lay out components, and generate visual assets. Iterate until the design aligns perfectly with the intent described in the task.

## Step 1.5: Document the Blueprint
**Action:** Once the design is finalized, log a note linking the produced artifact (e.g., a `.pen` file) or explaining the core design decisions so developers can reference it.
**Command:**
```bash
castra note add --role designer --project <ProjectID> --content "Designed user profile screen. See [profile_screen.pen] in designs directory." --tags "design"
```

## Step 1.6: Offer Your Work for Judgment
**Action:** Submit the finalized design for review.
**Command:**
```bash
castra task update --role designer --status review <TaskID>
```

## Step 1.7: Return to the Beginning
**Action:** Repeat this loop. Go back to Step 1.1 and pull the next design task.
