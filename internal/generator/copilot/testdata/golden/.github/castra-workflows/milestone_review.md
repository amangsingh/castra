---
description: Architect periodic check to review milestone completion and close out sprints.
---
## Step 0: Log Your Intent
Before taking any mutating action, you **MUST** declare your intent per Law 9. Run the following command:
`castra note add --role architect --project <ProjectID> --task <TaskID> --content "INTENT: Reviewing milestone completion. REASON: Periodic check." --tags "architect"`

## Step 1: List Milestones
`castra --role architect milestone list --project <ProjectID>`

## Step 2: Review Tasks
Check all tasks in the milestone to verify they're in the done status.
`castra --role architect task list --project <ProjectID> --milestone <MilestoneID>`

## Step 3: Close Milestone
If all tasks are complete, close the milestone.
`castra --role architect milestone update --status completed <MilestoneID>`
