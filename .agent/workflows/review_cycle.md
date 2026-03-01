---
description: Phase 1 - The Review Loop (Functional Verification)
---

# Phase 1: The Review Loop (Functional Verification)

**Trigger:** Tasks appear in the `review` state.
**Goal:** To systematically verify each task's implementation against its defined requirements.

## Step 1.1: Survey the Queue
**Action:** Query the database for all tasks awaiting review.
**Command:**
```bash
go run main.go task list --project <ProjectID>
```
*(Run this from within your scripts directory)*

## Step 1.2: Read the Specification
**Action:** For each task in review, fetch its complete context. This includes the task description (your test plan), architectural notes, engineer implementation details, and the audit log.
**Command:**
```bash
go run main.go task view <TaskID>
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
go run main.go task update --status done <TaskID>
```

## Step 1.5b: Reject the Task
**Action:** Reject the task back to `todo`. This resets ALL approval flags. You MUST attach a rejection note explaining the failure. See the `write_rejection` workflow.
**Command:**
```bash
go run main.go task update --status todo <TaskID>
```

## Step 1.6: Continue the Loop
**Action:** Return to Step 1.1. The queue is never empty until the sprint is complete.
