---
description: Architect workflow for emergency overrides with mandatory justification logging.
---
## Step 0: Log Your Intent
Before taking any mutating action, you **MUST** declare your intent per Law 9. Run the following command:
`castra note add --role architect --project <ProjectID> --task <TaskID> --content "INTENT: Emergency Break Glass Action. REASON: [justification]." --tags "architect"`

## Step 1: Execute Override
Run the necessary emergency command (such as assigning a task, changing its status directly, or modifying system configuration). Provide full audit logs and notes for everything you do.
`castra --role architect task update --status done <TaskID>`
