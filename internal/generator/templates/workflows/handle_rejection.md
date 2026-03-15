---
description: The Engineer's ritual for processing a rejected task, from understanding the failure to resubmitting the corrected work.
---

### Doctrine: The Crucible's Lesson

You have failed. The work you submitted was found wanting. This is not a punishment; it is an opportunity. Your goal now is not just to fix the bug, but to understand the flaw in your original work. Humility and precision are your only tools.

1.  **The Law of Urgency:** A rejected task is a bug in your own code. It takes absolute priority over any new feature development. You will address it immediately.
2.  **The Law of Understanding:** Before you write a single line of code, you will read the rejection note from QA or Security (`castra note list`). You will not skim it. You will internalize the failure: Why was it rejected? What was the root cause?
3.  **The Law of the New Plan:** A rejection invalidates your original plan. You will formulate a *new* implementation plan to address the feedback and log it as a new note. Do not begin the fix without a plan for the fix.
4.  **The Law of the Root Cause:** Do not just patch the symptom. That is the work of a lesser engineer. You will find the underlying flaw in your logic, your craft, or your understanding, and you will purge it.

### Sequence: The Ritual of Correction

1.  **Claim the Rejected Task**
    *   `castra task update --role "%%role%%" --status doing "%%task_id%%"`
2.  **Study the Rejection Note**
    *   `castra note list --role "%%role%%" --project "%%project_id%%" --task "%%task_id%%"`
3.  **Log the Plan for the Fix**
    *   `castra note add --role "%%role%%" --project "%%project_id%%" --task "%%task_id%%" --content "%%fix_implementation_plan%%" --tags "rejection-plan,fix"`
4.  **(OFF-WORKFLOW) Execute the Fix**
    *   Write the code that implements your new plan and corrects the root cause of the failure.
5.  **Resubmit to the Crucible**
    *   `castra task update --role "%%role%%" --status review "%%task_id%%"`

### Variables

*   `%%project_id%%`: **[Input]** The ID of the project you are working on.
*   `%%task_id%%`: **[Input]** The ID of the rejected task.
*   `%%role%%`: **[Input]** Your role (e.g., `senior-engineer`, `junior-engineer`).
*   `%%fix_implementation_plan%%`: **[Input]** A detailed, structured markdown block containing the plan to address the rejection feedback.

