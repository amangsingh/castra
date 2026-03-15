---
description: Phase 1 - The Review Loop (Functional Verification)
---

# Phase 1: The Review Loop (Functional Verification)

**Trigger:** Tasks appear in the `review` state.
**Goal:** To systematically verify each task's implementation against its defined requirements.

## Step 0: Log Your Intent
**Action:** Before starting any research or implementation, log your intent to work on the task. This ensures universality of surveillance and record-keeping.
**Command:**
```bash
castra log add --role qa-functional --msg "Starting work on task <TaskID>" --type task --entity <TaskID>
```

## Step 1.1: Survey the Queue
**Action:** Query the database for all tasks awaiting review.
**Command:**
```bash
castra task list --role qa-functional --project <ProjectID>
```
*(Run this from within your scripts directory)*

## Step 1.2: Read the Specification
**Action:** For each task in review, fetch its complete context. This includes the task description (your test plan), architectural notes, engineer implementation details, and the audit log.
**Command:**
```bash
castra task view --role qa-functional <TaskID>
```

## Step 1.3: Execute Functional Tests
**Action:** Test the observable behavior of the implementation against the specification. You do not read the source code. You verify:
- Does the feature do what the task description says it should?
- Do the expected inputs produce the expected outputs?
- Do edge cases (empty inputs, boundary values) behave correctly?
- Are error states handled gracefully?

## Step 1.4: Render Judgment
**Decision Point:**
- **PASS** → Proceed to Step 1.5a.
- **FAIL** → Proceed to Step 1.5b.

## Step 1.5a: Approve the Task
**Action:** Mark the task as functionally approved. This sets `qa_approved=true`. The task remains in `review` until Security Ops also approves.
**Command:**
```bash
castra task update --role qa-functional --status done <TaskID>
castra log add --role qa-functional --msg "Functional approval granted for task <TaskID>" --type task --entity <TaskID>
```

## Step 1.5b: Reject the Task
**Action:** Reject the task back to `todo`. This resets ALL approval flags. You MUST attach a rejection note explaining the failure. See the `write_rejection` workflow.
**Command:**
```bash
castra task update --role qa-functional --status todo --reason "<Failure Summary>" <TaskID>
castra log add --role qa-functional --msg "Functional rejection for task <TaskID>: <Failure Summary>" --type task --entity <TaskID>
```

## Step 1.6: Continue the Loop
**Action:** Return to Step 1.1. The queue is never empty until the sprint is complete.
