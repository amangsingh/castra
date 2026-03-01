---
description: Phase 2 - The Crucible (Handling Rejection)
---

# Phase 2: The Crucible (Handling Rejection)

**Trigger:** A task you previously submitted to review reappears in your todo list.
**Goal:** To address the feedback from QA or Security with urgency and precision.

## Step 2.1: Prioritize the Rejected Task
**Action:** A rejected task is a bug in your work. It takes absolute priority over any new feature development.

## Step 2.2: Understand the Failure
**Action:** Before writing any code, read the notes associated with the task to understand why it was rejected. Was it a functional bug found by QA? A security vulnerability found by Ops?
**Command:**
```bash
go run main.go note list --project <ProjectID> --task <TaskID>
```

## Step 2.3: Fix the Root Cause
**Action:** Correct the issue. Do not just patch the symptom; fix the underlying flaw in your logic.

## Step 2.4: Resubmit to the Crucible
**Action:** Once the fix is implemented, move the task back to the review state to be judged again.
**Command:**
```bash
go run main.go task update --status review <TaskID>
```
