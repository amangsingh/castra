---
description: Executes Phase 2 (The Roadmap), decomposing a single, high-level Milestone (an Epic) into a sequence of smaller, concrete feature milestones.
---

### Doctrine: The Art of the Roadmap

Your purpose here is to translate a grand, thematic Pillar into a concrete, sequential plan. You are the bridge from the "what" to the "how". Your thinking must be sequential, logical, and grounded in dependencies.

1.  **Isolate the Epic:** Your focus must be on a single, high-level milestone from the Blueprint phase. Do not try to plan multiple epics at once.
2.  **Identify the Critical Path:** Analyze the epic and identify the 2-5 major "features" that must be built *in order* to complete it. This is not about individual tasks yet; it's about defining the major phases of work. Think in terms of dependencies: What absolutely must exist before the next piece can be built?
3.  **Define Clear Outcomes:** Each feature you define must have a clear, tangible outcome. "Build the API" is a good feature. "Improve the database" is not. Be specific. These features will become the new, more granular **Milestones**.

### Sequence: Roadmap Protocol

1.  **Decompose the Epic into Feature Milestones**
    *   Repeat `castra milestone add` for each Feature defined in the Doctrine phase (2-5 total).

    *Example:*
    *   `castra milestone add --role architect --project "%%project_id%%" --parent "%%epic_milestone_id%%" --name "Feature: %%feature_1_name%%"`
    *   `castra milestone add --role architect --project "%%project_id%%" --parent "%%epic_milestone_id%%" --name "Feature: %%feature_2_name%%"`
    *   `castra milestone add --role architect --project "%%project_id%%" --parent "%%epic_milestone_id%%" --name "Feature: %%feature_3_name%%"`

### Variables

*   `%%project_id%%`: **[Input]** The ID of the parent project.
*   `%%epic_milestone_id%%`: **[Input]** The ID of the high-level Epic milestone you are decomposing.
*   `%%feature_1_name%%`, `%%feature_2_name%%`, etc.: **[Input]** The names of the 2-5 sequential features you have defined, derived from the Doctrine phase.

