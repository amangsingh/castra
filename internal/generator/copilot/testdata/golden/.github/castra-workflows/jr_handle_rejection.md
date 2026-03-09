---
description: Phase 2 - The Crucible (Handling Rejection)
---

# Phase 2: The Crucible (Handling Rejection)

**Trigger:** A task you previously submitted to `review` reappears in your `todo` list.
**Goal:** To address the feedback from QA or Security with urgency and precision.

## Step 2.1: Prioritize the Rejected Task
**Action:** A rejected task is a defect in your work. It takes absolute priority over any new task. Drop what you are doing.

## Step 2.2: Understand the Failure
**Action:** Before writing any code, read the rejection notes attached to the task. Understand exactly what failed and why.
**Command:**
```bash
go run main.go note list --project <ProjectID> --task <TaskID>
```
Look for notes tagged with `qa,rejection` or `security,vulnerability`. These contain the failure details.

## Step 2.3: Fix the Root Cause
**Action:** Correct the issue. Do not patch the symptom — fix the underlying logic error or vulnerability. If the fix exceeds your scope (e.g., requires architectural changes), escalate by blocking the task:
```bash
go run main.go task update --status blocked <TaskID>
go run main.go note add --project <ProjectID> --task <TaskID> --content "ESCALATION: Fix requires architectural change beyond junior scope. <description>" --tags "junior-engineer,escalation"
```

## Step 2.4: Resubmit to the Crucible
**Action:** Once the fix is implemented and tested locally, move the task back to `review`. Both QA and Security will re-verify from scratch (all approval flags were reset on rejection).
**Command:**
```bash
go run main.go task update --status review <TaskID>
```
