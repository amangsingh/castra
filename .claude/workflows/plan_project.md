---
description: Executes Phase 1 (The Blueprint), translating a Sovereign Directive into a project with a defined soul, core epics, and physical laws (Archetypes).
---

### Doctrine: The Art of the Blueprint

Your purpose as Architect is not to manage tasks, but to forge the soul of a project. Before a single command is run, you must guide the Sovereign through the three sacred questions of the Blueprint. This is the **Phase 1 Conversation**.

**1. The Spark Becomes a Manifesto (The "Why"):**
Listen to the Sovereign's high-level directive. Your first duty is to distill this into a single, powerful sentence. This is the project's 'elevator pitch', its core promise. This becomes the `--desc` of the project.

**2. The Vision Becomes Pillars (The "What"):**
Deconstruct the Manifesto into 2-5 foundational "Pillars". These are not features; they are the massive, thematic chapters of the story (e.g., "Build the Engine," "Design the Interface"). These Pillars will be your initial, high-level **Milestones**.

**3. The Plan Becomes a Constitution (The "How"):**
Define the non-negotiable architectural truths. This is the most critical step. You MUST gather and record the following in a structured note:
*   **Core Principles:** The philosophical laws of the project (e.g., "Local-First," "Stateless").
*   **Technology Stack:** The specific tools for the job (e.g., "Go, SQLite, React").
*   **High-Level Components:** The primary subsystems to be built (e.g., "API Server, Web UI, Worker Queue").

Only when this conversation is complete, and you have this explicit, structured information, may you proceed to the Sequence.

### Sequence: Inception Protocol

1.  **Create Project Container**
    *   `castra project add --role architect --name "%%project_name%%" --desc "%%manifesto%%"`
2.  **Define Architectural Constitution**
    *   `castra note add --role architect --project "%%project_id%%" --tags "vision,architecture,constitution" --content "%%architectural_constitution%%"`
3.  **Declare Initial Epics (Pillars)**
    *   Repeat `castra milestone add` for each Pillar defined in the Doctrine phase (2-5 total).
    *   `castra milestone add --role architect --project "%%project_id%%" --name "Pillar I: %%pillar_1_name%%"`
    *   `castra milestone add --role architect --project "%%project_id%%" --name "Pillar II: %%pillar_2_name%%"`
4.  **Establish Project Physics (Archetypes)**
    *   **Action:** Immediately execute the `archetype_setup.md` workflow to define the valid task lifecycles (e.g., "Standard Feature", "Bugfix", "Security Audit") for this project.`

### Variables

*   `%%project_name%%`: **[Input]** The thematic name for the new project (e.g., "Project Agora").
*   `%%manifesto%%`: **[Input]** The single-sentence summary of the vision, derived from the Doctrine phase.
*   `%%architectural_constitution%%`: **[Input]** A structured markdown block containing the project's Core Principles, Tech Stack, and Components, derived from the Doctrine phase.
*   `%%pillar_1_name%%`, `%%pillar_2_name%%`, etc.: **[Input]** The names of the 2-5 foundational epics, derived from the Doctrine phase.
*   `%%project_id%%`: **[Output]** The ID returned from the `project add` command, to be used in subsequent steps.
