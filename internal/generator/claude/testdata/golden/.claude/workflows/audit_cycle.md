---
description: Phase 1 - The Audit Loop (Security Verification)
---

# Phase 1: The Audit Loop (Security Verification)

**Trigger:** Tasks appear in the `review` state.
**Goal:** To systematically audit each task's implementation for security vulnerabilities before it can be marked done.

## Step 0: Log Your Intent
**Action:** Before starting any research or implementation, log your intent to work on the task. This ensures universality of surveillance and record-keeping.
**Command:**
```bash
castra log add --role security-ops --msg "Starting work on task <TaskID>" --type task --entity <TaskID>
```

## Step 1.1: Survey the Queue
**Action:** Query the database for all tasks awaiting security review.
**Command:**
```bash
castra task list --role security-ops --project <ProjectID>
```
*(Run this from within your scripts directory)*

## Step 1.2: Understand the Attack Surface
**Action:** For each task in review, fetch its complete context using the view command. Read the task description, architectural notes, implementation details, and the audit log. Understand what the code does so you can identify what it exposes.
**Command:**
```bash
castra task view --role security-ops <TaskID>
```

## Step 1.3: Conduct the Security Audit
**Action:** Read the source code. Unlike QA, you MUST inspect the implementation directly. Your checklist:
- **Injection:** SQL injection, command injection, XSS, template injection.
- **Authentication/Authorization:** Broken auth flows, privilege escalation, missing access controls.
- **Data Exposure:** Secrets in code, PII leaks, verbose error messages, debug endpoints.
- **Dependencies:** Known CVEs in imported packages, insecure transitive dependencies.
- **Cryptography:** Weak hashing, hardcoded keys, insecure random number generation.
- **Input Validation:** Unvalidated user input, missing bounds checks, path traversal.

## Step 1.4: Render Judgment
**Decision Point:**
- **PASS** → Proceed to Step 1.5a.
- **FAIL** → Proceed to Step 1.5b.

## Step 1.5a: Approve the Task
**Action:** Mark the task as security-approved. This sets `security_approved=true`. Combined with QA approval, this transitions the task to `done`.
**Command:**
```bash
castra task update --role security-ops --status done <TaskID>
castra log add --role security-ops --msg "Security approval granted for task <TaskID>" --type task --entity <TaskID>
```

## Step 1.5b: Reject the Task
**Action:** Reject the task back to `todo`. This resets ALL approval flags. You MUST attach a security finding note. See the `write_finding` workflow.
**Command:**
```bash
castra task update --role security-ops --status todo --reason "<Reason>" <TaskID>
castra log add --role security-ops --msg "Security rejection for task <TaskID>: <Reason>" --type task --entity <TaskID>
```

## Step 1.6: Continue the Loop
**Action:** Return to Step 1.1. Security is never done.
