---
description: Engineer self-review before submitting code to review status.
---
## Step 0: Log Your Intent
Before taking any mutating action, you **MUST** declare your intent per Law 9. Run the following command:
`castra note add --role senior-engineer --project <ProjectID> --task <TaskID> --content "INTENT: Conducting pre-commit checks. REASON: Moving to review." --tags "senior-engineer"`

## Step 1: Run Automated Tests
Verify that the code compiles and tests pass locally.
Run the project's test suite or compiler directly (e.g. `go test ./...`).

## Step 2: Validate Acceptance Criteria
Ensure that every numbered criterion from the architect's description is satisfied.

## Step 3: Submit for Review
If everything passes, push the task to QA for functional testing.
`castra --role <role> task update --status review <TaskID>`
