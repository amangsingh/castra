---
description: Defines the custom state-machine pipelines (Archetypes) that govern how different tasks move through the system.
---

### Doctrine: The Art of the State Machine

Your purpose here is to act as the ultimate Lawgiver. You are not defining a task; you are defining the *rules of life for all tasks of a certain type*. An Archetype is a railroad track: once a task is on it, it can only move to the states you define, in the order you define them.

1.  **Think in Lifecycles:** Before you define an archetype, map out its entire lifecycle. Does a "bug" need a "triage" state? Does a "security task" skip the normal `doing` state and go straight to `review`? Each archetype is a unique process.
2.  **The Power of Restriction:** The strength of an archetype is not just the states it includes, but the states it *excludes*. By defining a strict path, you prevent chaos and ensure the correct process is always followed.
3.  **Assign at Creation:** Remember that these archetypes are assigned when a task is created (`castra task add ... --archetype <id>`). This workflow is for *defining* those tracks, not for putting trains on them.

### Sequence: Archetype Definition Protocol

1.  **Define the Archetype & its Status Pipeline**
    *   Repeat `castra archetype add` for each distinct task lifecycle the project requires.

    *Example:*
    *   `castra archetype add --role architect --project "%%project_id%%" --name "Standard Feature" --desc "The standard dev lifecycle: build, review, done." --statuses "todo,doing,review,done"`
    *   `castra archetype add --role architect --project "%%project_id%%" --name "Security Audit" --desc "A security-only review task that skips implementation." --statuses "todo,review,done"`
    *   `castra archetype add --role architect --project "%%project_id%%" --name "Quick Bugfix" --desc "A fast track for simple bug fixes." --statuses "todo,doing,done"`

### Variables

*   `%%project_id%%`: **[Input]** The ID of the project for which these archetypes are being defined.
*   `%%archetype_name%%`: **[Input]** The descriptive name of the lifecycle (e.g., "Standard Feature").
*   `%%archetype_desc%%`: **[Input]** A brief explanation of what this lifecycle is for.
*   `%%statuses%%`: **[Input]** A comma-separated string of the valid states for this archetype, in order.

