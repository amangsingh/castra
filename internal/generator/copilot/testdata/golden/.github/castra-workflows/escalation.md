---
description: Engineer workflow for escalating tasks that require protected file modifications.
---
## Step 0: Log Your Intent
Before taking any mutating action, you **MUST** declare your intent per Law 9. Run the following command:
`castra note add --role <role> --project <ProjectID> --task <TaskID> --content "INTENT: Escalating task for protected file modifications. REASON: Required to touch router.go/migrate.go/logic.go" --tags "<role>"`

## Step 1: Block the Task
Update the task status to `blocked` to stop progress and notify the team that architect or senior engineer intervention is needed.
`castra --role <role> task update --status blocked <TaskID>`

## Step 2: Handoff
Notify the Architect or Senior Engineer to review the note and reassign the task or break it into safe components.
