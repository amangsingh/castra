---
description: Executes Phase 3 (The Work Order), decomposing a Feature-Milestone into a granular, actionable backlog of engineering tasks.
---

### Doctrine: The Art of the Work Order

Your purpose here is the final act of translation. You transform a single Feature from the roadmap into a series of small, concrete, and immediately actionable tasks for the engineering team. This is the bridge from planning to execution.

1.  **Isolate the Feature:** Your focus must be on a single Feature-Milestone. Do not try to plan multiple features in one sprint.
2.  **Think Like the Builder:** Deconstruct the feature into the exact sequence of steps an engineer would take. What must be built first? What are the dependencies?
3.  **Embrace Atomicity:** Each task you create must be "atomic." It should be small enough to be completed by a single engineer in 1-2 days and have a clear, verifiable "definition of done."
4.  **The Law of the Contract:** You are **forbidden** from creating a task with a vague description. Every task you generate is a binding contract for the engineer. Its description (`--desc` flag) **MUST** contain the following structured markdown sections:
    *   **Objective:** A clear, one-sentence statement of what will be built or fixed.
    *   **Acceptance Criteria:** An enumerated list of testable, pass/fail conditions that define "DONE." (e.g., "1. The endpoint returns a 200 OK. 2. The response body contains a 'user_id'.")
    *   **Technical Notes:** (Optional) Any specific guidance, libraries to use, or architectural patterns the engineer should follow.
    *   **Verification Steps:** Explicit, step-by-step instructions for the `@qa-functional` agent to verify the acceptance criteria have been met.
5.  **The Law of the Path:** You are **forbidden** from creating a task without assigning it to a pre-defined Archetype. A task without a path is a rogue element and a violation of system integrity. The `--archetype` flag is mandatory.

### Sequence: Work Order Protocol

1.  **Create a Time-Boxed Sprint Container**
    *   `castra sprint add --role architect --project "%%project_id%%" --name "%%sprint_name%%" --start-date "%%start_date%%" --end-date "%%end_date%%"`
2.  **Decompose the Feature into Contractual Tasks**
    *   Repeat `castra task add` for each atomic task, ensuring each is assigned to an Archetype.

    *Example:*
    *   `castra task add --role architect --sprint "%%sprint_id%%" --title "Task 1" --desc "%%contract%%" --archetype "%%archetype_id%%"`
    *   `castra task add --role architect --sprint "%%sprint_id%%" --title "Task 2" --desc "%%contract%%" --archetype "%%archetype_id%%"`

### Variables

*   `%%project_id%%`: **[Input]** The ID of the parent project.
*   `%%sprint_name%%`: **[Input]** The specific name for this sprint cycle.
*   `%%start_date%%`, `%%end_date%%`: **[Input]** The start and end dates for the sprint.
*   `%%task_1_title%%`, `%%task_1_contract%%`, etc.: **[Input]** The titles and detailed, structured contract descriptions for each atomic task.
*   `%%sprint_id%%`: **[Output]** The ID returned from the `sprint add` command, to be used in subsequent steps.
*   `%%archetype_id%%`: **[Input]** The ID of the Archetype that defines the lifecycle for this specific task.
