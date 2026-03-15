---
description: The Engineer's mandatory subroutine for creating feature-level documentation upon task completion.
---

### Doctrine: The Scribe's Duty

You are not just a builder; you are a historian. Code without documentation is a ruin. Before you submit your work for judgment, you will chronicle what you have built. The memory of your work is as important as the work itself.

1.  **The Law of Immediacy:** Documentation is written when the knowledge is fresh. This is the final step you take *before* moving a task to `review`.
2.  **The Law of Utility:** Your documentation must be useful. It must explain the *purpose* of the feature, the *usage* (e.g., API endpoints, new commands), and any necessary *configuration*.

### Sequence: The Chronicler's Protocol

1.  **Log the Feature Document**
    *   `castra note add --role senior-engineer --project "%%project_id%%" --task "%%task_id%%" --content "%%feature_documentation%%" --tags "documentation,feature-doc"`

### Variables

*   `%%project_id%%`: **[Input]** The ID of the project.
*   `%%task_id%%`: **[Input]** The ID of the task being documented.
*   `%%feature_documentation%%`: **[Input]** A structured markdown block containing: 1. Purpose, 2. Usage/Examples, 3. Configuration.

