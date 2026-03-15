---
description: The triage protocol for a Senior Engineer or Architect responding to an escalated task from a Junior Engineer.
---

### Doctrine: The Shepherd's Duty

A junior has raised a flag, declaring a task is beyond their scope. Your duty is to respond with speed and wisdom. You are the senior officer arriving at a barricade. You must assess, plan, and redeploy your forces.

1.  **The Law of Triage:** An escalated task is a high-priority interrupt. Your first step is to review the junior's escalation note to understand the nature of the blockage.
2.  **The Law of Redeployment:** You have three options:
    *   **Re-scope:** If the task can be broken down into smaller, safer pieces, you will create new tasks and assign them.
    *   **Re-assign:** If the task cannot be broken down but is still valid, you will claim it yourself.
    *   **Reject:** If the junior's assessment was wrong and the task is within their scope, you will add a note clarifying the path forward and move the task back to `todo`.
3.  **The Law of the Closed Loop:** You will not leave the junior in limbo. The escalation is not resolved until the original blocked task is either closed (in favor of new tasks) or reassigned.

### Sequence: The Triage Protocol

1.  **Identify Escalated Task**
    *   `castra task list --role <role> --status blocked`
2.  **Review the Escalation Note**
    *   `castra note list --role <role> --task "%%task_id%%"`
3.  **Execute Triage:**
    *   **Option A (Re-scope & Replace):**
        *   `castra task update --role <role> --status closed "%%task_id%%" --force --reason "Closing in favor of new, re-scoped tasks."`
        *   *(Then, execute `plan_sprint.md` to create the new, smaller tasks)*
    *   **Option B (Re-assign):**
        *   `castra task assign "%%task_id%%" "%%your_user_id%%"`
        *   `castra task update --role <role> --status todo "%%task_id%%"`
    *   **Option C (Reject & Clarify):**
        *   `castra note add --role <role> --task "%%task_id%%" --content "CLARIFICATION: %%explanation%%"`
        *   `castra task update --role <role> --status todo "%%task_id%%"`

### Variables
*   `%%task_id%%`: **[Input]** The ID of the blocked task.
*   `%%role%%`: **[Input]** Your role (e.g., `architect`, `senior-engineer`).
*   `%%explanation%%`: **[Input]** A clear explanation for why the escalation is being rejected.
