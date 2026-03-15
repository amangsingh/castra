---
description: Create archetypes with explicit status pipelines and default roles.
---
## Step 0: Log Your Intent
Before taking any mutating action, you **MUST** declare your intent per Law 9. Run the following command:
`castra note add --role architect --project <ProjectID> --task <TaskID> --content "INTENT: Creating project archetypes. REASON: [why]" --tags "architect"`

## Step 1: Add Archetypes
Define archetype chains representing specific workflows.
`castra archetype add --role architect --project <ProjectID> --name "Feature" --desc "Standard engineering feature task" --statuses "todo,doing,review,done"`
`castra archetype add --role architect --project <ProjectID> --name "Security Audit" --desc "Security-only review task" --statuses "todo,review,done"`

## Step 2: Assign Archetypes to Tasks
When creating tasks, ensure you assign them to the correct archetype:
`castra task add --role architect --project <ProjectID> --milestone <MilestoneID> --sprint <SprintID> --title "..." --desc "..." --prio high --archetype <ArchetypeID>`
