---
description: The Architect's emergency protocol for overriding system rules with mandatory, audited justification.
---

### Doctrine: The Law of Last Resort

This is the most dangerous tool in your arsenal. You will only invoke it when the established rules of the system (Archetypes, role permissions) have created a deadlock, or when an external emergency requires an action that the system, by design, forbids. This is not a shortcut. It is an admission of system failure.

1.  **The Law of Justification:** The reason for the override is more important than the override itself. You will not execute an override without providing a clear, concise, and compelling justification. This is not optional.
2.  **The Law of Atomic Action:** The override and the justification are one. You will use the `--force` flag to unlock the action, and the `--reason "..."` flag to provide the mandatory justification. These two flags are inseparable. Any command with `--force` that is missing `--reason` is an invalid command.
3.  **The Law of Scrutiny:** Every "Break Glass" event will be a permanent, high-visibility entry in the system's audit log. Your justification will be recorded and subject to review. Use this power with the gravity it deserves.

### Sequence: The Override Protocol

This is not a multi-step process. It is a single, deliberate, and justified action.

1.  **Execute Override with Justification**
    *   Append the `--force` and `--reason "..."` flags to the command that is currently being blocked by the system's rules.

    *Example 1: Forcing a task to 'done' from 'todo', bypassing the 'review' state.*
    *   `castra task update --role architect --status done "%%task_id%%" --force --reason "Manual verification of hotfix for critical production outage."`

    *Example 2: Forcibly re-assigning a task locked by another user.*
    *   `castra task assign "%%task_id%%" "%%new_user_id%%" --force --reason "Original assignee is unresponsive and the task is blocking the release."`

### Variables

*   `%%reason%%`: **[Input]** A mandatory, clear, and concise justification for why the system's rules are being overridden.
*   *Other variables (`%%task_id%%`, etc.) are dependent on the specific command being executed.*

