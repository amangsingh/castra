---
description: Phase 2 - Blueprinting (Feature & Milestone Decomposition)
---

# Phase 2: Blueprinting (Feature & Milestone Decomposition)

**Trigger:** Phase 1 is complete. A project with a vision note exists.
**Goal:** To decompose the project's vision into large, thematic features and define the major milestones required to complete each feature.

## Step 2.1: Create a "Feature Sprint"
**Action:** For a major feature identified in the vision note (e.g., "User Authentication"), create a long-running sprint that will act as a container for its milestones. The name must be prefixed with Feature:.
**Command:**
```bash
castra sprint add --role architect --project <ProjectID> --name "Feature: User Authentication"
```
**Capture:** Record the Feature Sprint ID.

## Step 2.2: Define Feature Milestones as Tasks
**Action:** Break down the feature into 2-5 high-level milestones. These are not yet engineering tasks. They are large phases of work. Create a task for each milestone and assign it to the "Feature Sprint".
**Examples:** "Milestone: API & Database Design," "Milestone: Front-end Scaffolding," "Milestone: Integration with Payment Gateway."
**Command (repeat for each milestone):**
```bash
castra task add --role architect --project <ProjectID> --sprint <FeatureSprintID> --title "Milestone: API & Database Design" --desc "Define all required database tables, fields, and API endpoint contracts for this feature."
```

## Step 2.3: Repeat for All Features
**Action:** Repeat steps 2.1 and 2.2 for every major feature outlined in the architectural vision.

**Output:** The database now contains multiple long-running "Feature Sprints," each populated with a series of high-level "Milestone Tasks." The project now has a complete, high-level roadmap.
