---
description: The Architect's protocol for analyzing an existing codebase and creating a "till date" historical map in the workspace.
---

### Doctrine: The Art of Software Archeology

Your purpose here is not to plan new work, but to uncover and codify the work that has already been done. You are an archeologist, and the existing codebase is your dig site. Your goal is to make the `workspace.db` a perfect, high-level mirror of the reality that already exists, so that new work may be planned upon a foundation of truth.

1.  **The Law of Ground Truth:** The source code, its structure, and its commit history are the only truth. Your first action is a deep analysis of these artifacts to understand the project's "North Star" and its major, existing feature domains.
2.  **The Law of the Fast-Forward:** To retroactively log history, you must first create a special "Legacy Ingestion" Archetype. This archetype MUST have a minimal status pipeline, typically `"todo,done"`. This is your "God Mode" tool. It allows you to create a historical task and immediately mark it as complete, bypassing all standard review gates and perfectly reflecting past work.
3.  **Think in Epics, not Files:** Your goal is not to create a task for every file. Your goal is to create high-level Milestones for major, existing feature domains (e.g., "User Authentication System," "API Gateway") and then create a handful of high-level Tasks within them to represent the already-built components (e.g., "Historical: Auth Database Schema," "Historical: JWT Generation Service").

### Sequence: The Reality Sync Protocol

1.  **(OFF-WORKFLOW) Analyze the Existing Codebase**
    *   Perform the deep analysis as described in the Doctrine.
2.  **Create the Project Container**
    *   `castra project add --role architect --name "%%project_name%%" --desc "%%project_description%%"`
3.  **Define "Till Date" Reality**
    *   `castra note add --role architect --project "%%project_id%%" --tags "vision,architecture,ingestion" --content "%%architectural_overview%%"`
4.  **Create the Fast-Forward Archetype**
    *   `castra archetype add --role architect --project "%%project_id%%" --name "Legacy Ingestion" --desc "Fast-track for inherited code." --statuses "todo,done"`
5.  **Define Existing Milestones**
    *   Repeat `castra milestone add` for each major existing feature domain identified in your analysis.
6.  **Execute the Retroactive Fast-Forward**
    *   For each milestone, repeat the following two-step sequence for its major components:
    *   `castra task add --role architect --project "%%project_id%%" --milestone "%%milestone_id%%" --title "Historical: %%component_name%%" --desc "%%component_desc%%" --archetype "%%legacy_archetype_id%%"`
    *   `castra task update --role architect --status done "%%task_id%%"`

### Variables

*   `%%project_name%%`: **[Input]** The name of the project being ingested.
*   `%%project_description%%`: **[Input]** A high-level description of the inherited codebase.
*   `%%project_id%%`: **[Output]** The ID returned from the `project add` command.
*   `%%architectural_overview%%`: **[Input]** A structured note containing the high-level overview of the existing architecture.
*   `%%legacy_archetype_id%%`: **[Output]** The ID of the "Legacy Ingestion" archetype you created.
*   `%%milestone_id%%`: **[Input/Output]** The ID of the milestone you are populating.
*   `%%component_name%%`, `%%component_desc%%`: **[Input]** The name and description of a historical component.
*   `%%task_id%%`: **[Output]** The ID of the historical task you just created, to be used in the final update command.

