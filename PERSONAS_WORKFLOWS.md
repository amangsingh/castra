# Castra v2.0 - Personas & Workflows Inventory

This document provides a categorical analysis of the personas and mandatory workflows defined within the Castra v2.0 templates.

## 1. Architect (The Lawgiver)
**Persona:** Translates the Sovereign's will into an actionable lattice of Milestones and Tasks using the `castra` CLI.

*   **Workflow [plan_project]:** Phase 1 Strategic Planning — Formal project inception and top-level goal setting.
*   **Workflow [plan_feature]:** Phase 2 Blueprinting — Decomposing features into executable Milestones and logical units.
*   **Workflow [plan_sprint]:** Phase 3 Tactical Planning — Generating work orders and assigning tasks to specific sprints.
*   **Workflow [plan_design]:** Phase 2 Interface Blueprinting — Planning UI/UX requirements before design execution.
*   **Workflow [milestone_review]:** Periodic audit to verify milestone completion and formal sprint closure.
*   **Workflow [archetype_setup]:** Defining custom status pipelines (HATEOAS) and assigning default roles to archetypes.
*   **Workflow [break_glass]:** High-authority protocol for forcing emergency gate overrides with mandatory audit logging.

## 2. Senior Engineer (The Core Builder)
**Persona:** Implements foundational, load-bearing code and complex features from architectural blueprints to the highest standard.

*   **Workflow [build_cycle]:** The Primary Build Loop — Multi-step process for planning, implementing, and self-testing core code.
*   **Workflow [handle_rejection]:** The Crucible — Formal protocol for processing code/security rejections and re-implementing fixes.
*   **Workflow [document_task]:** Phase 1 Chronicling — Creating technical documentation (feature docs) upon task completion.

## 3. Junior Engineer (The Maintainer)
**Persona:** Executes routine maintenance, bug fixes, and minor refactors without touching protected system-critical files.

*   **Workflow [jr_build_cycle]:** Simplified Build Loop — Targeted at maintenance and speed with strict scope boundaries.
*   **Workflow [jr_handle_rejection]:** Junior Rejection Protocol — Focused on iterative fixes for routine/UI-level issues.
*   **Workflow [escalation]:** Formal protocol for identifying and handing off tasks that require modification of protected load-bearing files.
*   **Workflow [pre_commit_checklist]:** Self-audit loop for verification of linting, tests, and documentation before submitting for review.

## 4. Security Ops (The Sentinel)
**Persona:** Provides absolute security oversight via code auditing and vulnerability detection; holds the final veto.

*   **Workflow [audit_cycle]:** The Audit Loop — Systematic verification of authentication, sanitization, and code-level vulnerability assessment.
*   **Workflow [write_finding]:** The Finding Protocol — Generating formal, structured vulnerability reports (SQLi, Auth, etc.) for engineers.

## 5. QA Functional (The Guardian of Intent)
**Persona:** Verifies observable system behavior and technical criteria against the Architect's original requirements.

*   **Workflow [review_cycle]:** The Review Loop — Explicit verification of acceptance criteria and automated test suite health.
*   **Workflow [write_rejection]:** The Rejection Protocol — Crafting actionable behavioral feedback when functional requirements are not met.

## 6. Designer (The Shaper)
**Persona:** Visualizes the interface and UX in Pencil-format deliverables, focusing on aesthetic and structural integrity.

*   **Workflow [execute_design]:** The Design Cycle — Iterative process for creating frames, interaction maps, and design rationale artifacts.

## 7. Doc Writer (The Chronicler)
**Persona:** Synthesizes project history and technical state into durable documentation like READMEs and Release Notes.

*   **Workflow [synthesize_project]:** Aggregate Documentation — Compiling cross-task history into comprehensive project-level manuals.
*   **Workflow [document_task]:** Feature Chronicle — Documenting the purpose, usage, and configuration changes of a specific feature.

## 8. System / Agnostic
**Persona:** Operational protocols that apply to the environment rather than a specific role.

*   **Workflow [ingest_project]:** Phase 0 Reality Sync — Unified protocol for analyzing an existing codebase to synchronize the engine's internal state.
