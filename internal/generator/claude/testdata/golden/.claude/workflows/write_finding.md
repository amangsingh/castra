---
description: Phase 2 - The Finding Protocol (Security Vulnerability Report)
---

# Phase 2: The Finding Protocol (Security Vulnerability Report)

**Trigger:** A task has failed the security audit during the Audit Loop.
**Goal:** To provide a structured, actionable security finding so the engineer can remediate the vulnerability.

## Step 0: Log Your Intent
**Action:** Before starting any research or implementation, log your intent to work on the task. This ensures universality of surveillance and record-keeping.
**Command:**
```bash
castra log add --role security-ops --msg "Starting work on task <TaskID>" --type task --entity <TaskID>
```

## Step 2.1: Document the Finding
**Action:** Before rejecting the task, write a note attached to the specific task. A rejection without a finding report is negligent — you are the last line of defense.

**The note MUST contain:**
1. **Vulnerability Class:** The category (e.g., SQL Injection, XSS, Insecure Deserialization, Secret Leak).
2. **Severity:** Critical / High / Medium / Low (use CVSS-like reasoning).
3. **Location:** The file and function where the vulnerability exists.
4. **Description:** What the vulnerability is and how it can be exploited.
5. **Remediation:** A specific, actionable fix recommendation.

**Command:**
```bash
castra note add --role security-ops --project <ProjectID> --task <TaskID> --content "SECURITY FINDING: [Class]: SQL Injection. [Severity]: Critical. [Location]: handlers/login.go:42 handleLogin(). [Description]: User input concatenated directly into SQL query without parameterization. [Remediation]: Use parameterized queries via db.Query with ? placeholders." --tags "security,vulnerability,critical"
```

## Step 2.2: Reject the Task
**Action:** Reject the task. The status change to `todo` automatically resets both approval flags, forcing a complete re-review cycle after the fix.
**Command:**
```bash
castra task update --role security-ops --status todo --reason "<Vulnerability Class>: <Severity>" <TaskID>
```

## Step 2.3: Log to Audit Trail
**Action:** Security findings are significant events. Ensure the audit log captures the rejection.
**Command:**
```bash
castra log add --role security-ops --msg "Security rejection for task <TaskID>: <Vulnerability Class>" --type task --entity <TaskID>
```

## Step 2.4: Move On
**Action:** Return to the Audit Loop. The engineer will remediate and resubmit. You will verify the fix in the next cycle.
