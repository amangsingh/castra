---
description: Phase 2 - The Blueprint (Planning the Interface)
---

# Phase 2: The Blueprint (Planning the Interface)

**Trigger:** The Architect assigns a high-level feature that requires user interface work, but the exact flow and screens are not yet defined.
**Goal:** Define the visual architecture and component structure before starting execution.

## Step 2.1: Receive the Mandate
**Action:** The Architect will assign you a task in the `planning` or `todo` state that requires design definition. Claim it.
**Command:**
```bash
castra task update --role designer --status doing <TaskID>
```

## Step 2.2: Understand the Context
**Action:** Read the task definition completely. What is the goal? Who is the user? What data needs to be displayed or collected?
**Command:**
```bash
castra task view --role designer <TaskID>
```

## Step 2.3: Draft the Wireframes or Flow
**Action:** Before committing to high-fidelity designs, map out the layout, component hierarchy, and navigation structure. This is the structural blueprint. Record key decisions in a design note.
**Command:**
```bash
castra note add --role designer --project <ProjectID> --content "Planned navigation flow. Proposed screens: [Dashboard, Settings, User Profile]." --tags "design,plan"
```

## Step 2.4: Seek Feedback
**Action:** Present the structural plan for review before executing the final design. Move the task to review so the Architect or Sovereign can approve the direction.
**Command:**
```bash
castra task update --role designer --status review <TaskID>
```

## Step 2.5: Await Judgment
**Action:**
*   If **Approved** (moved to `done`), the blueprint is locked. The Architect will spawn execution tasks based on your plan.
*   If **Rejected** (moved back to `todo`), read the rejection log. Adjust the blueprint and repeat Step 2.3.
