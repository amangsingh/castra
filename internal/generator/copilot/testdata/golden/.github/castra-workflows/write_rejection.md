---
description: Phase 2 - The Rejection Protocol (Actionable Feedback)
---

# Phase 2: The Rejection Protocol (Actionable Feedback)

**Trigger:** A task has failed functional verification during the Review Loop.
**Goal:** To provide the engineer with a clear, actionable rejection note so the defect can be fixed on the first attempt.

## Step 0: Log Your Intent
**Action:** Before starting any research or implementation, log your intent to work on the task. This ensures universality of surveillance and record-keeping.
**Command:**
```bash
castra log add --role <role> --msg "Starting work on task <TaskID>" --type task --entity <TaskID>
```

## Step 2.1: Document the Failure
**Action:** Before rejecting the task, write a note attached to the specific task. A rejection without a note is a cardinal sin — the engineer cannot fix what they cannot understand.

**The note MUST contain:**
1. **What was expected:** The behavior specified in the task description.
2. **What actually happened:** The observed incorrect behavior.
3. **Steps to reproduce:** A concise, numbered list.
4. **Severity:** Critical (blocks usage), Major (feature broken), Minor (cosmetic/edge case).

**Command:**
```bash
castra note add --project <ProjectID> --task <TaskID> --content "REJECTION: [Expected]: Login with empty password shows error. [Actual]: Application crashes. [Steps]: 1. Open login. 2. Leave password blank. 3. Click submit. [Severity]: Critical." --tags "qa,rejection"
```

## Step 2.2: Reject the Task
**Action:** Now reject. The status change to `todo` will automatically reset both `qa_approved` and `security_approved` flags, forcing a complete re-review cycle.
**Command:**
```bash
castra task update --status todo <TaskID>
```

## Step 2.3: Move On
**Action:** Do not dwell on the rejection. Return to the Review Loop and continue processing the queue. The engineer will address the failure in their own time.
