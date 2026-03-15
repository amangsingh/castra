---
description: Phase 2 - The Crucible (Handling Rejection)
---

# Phase 2: The Crucible (Handling Rejection)

**Trigger:** A task you previously submitted to review reappears in your todo list.
**Goal:** To address the feedback from QA or Security with urgency and precision.

## Step 0: Log Your Intent
**Action:** Before starting any research or implementation, log your intent to work on the task. This ensures universality of surveillance and record-keeping.
**Command:**
```bash
castra log add --role <role> --msg "Starting work on task <TaskID>" --type task --entity <TaskID>
```

## Step 2.1: Prioritize the Rejected Task
**Action:** A rejected task is a bug in your work. It takes absolute priority over any new feature development.

## Step 2.2: Understand the Failure
**Action:** Before writing any code, read the notes associated with the task to understand why it was rejected. Was it a functional bug found by QA? A security vulnerability found by Ops?
**Command:**
```bash
castra note list --project <ProjectID> --task <TaskID>
```

## Step 2.3: Fix the Root Cause
**Action:** Correct the issue. Do not just patch the symptom; fix the underlying flaw in your logic.

## Step 2.4: Resubmit to the Crucible
**Action:** Once the fix is implemented, move the task back to the review state to be judged again.
**Command:**
```bash
castra task update --status review <TaskID>
```
