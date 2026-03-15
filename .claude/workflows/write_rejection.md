---
description: The QA agent's subroutine for writing a detailed, actionable functional rejection report.
---

### Doctrine: The Scribe of Failure

When you reject a task, you become a teacher. A simple "it's broken" is useless. Your purpose is to create a perfect bug report that allows the engineer to understand and replicate the failure with zero ambiguity.

1.  **The Law of Replication:** Your report MUST contain a step-by-step list of actions to reliably reproduce the bug.
2.  **The Law of Evidence:** Your report MUST state the "Expected Behavior" (from the task contract) and the "Actual Behavior" (what you observed).
3.  **The Law of Precision:** Be specific. "The page crashed" is a bad report. "Entering 'abc' into the email field on the login page causes a 500 internal server error" is a good report.

### Sequence: The Rejection Protocol

1.  **Log the Failure Report**
    *   `castra note add --role qa-functional --project "%%project_id%%" --task "%%task_id%%" --content "%%rejection_report%%" --tags "rejection,qa,bug-report"`

### Variables

*   `%%project_id%%`: **[Input]** The ID of the project.
*   `%%task_id%%`: **[Input]** The ID of the failed task.
*   `%%rejection_report%%`: **[Input]** A structured markdown block containing: 1. Steps to Reproduce, 2. Expected Behavior, 3. Actual Behavior.

