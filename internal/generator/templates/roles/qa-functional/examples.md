# `castra` CLI Examples for the QA Functional Role

## Listing Tasks for Review
`castra task list --role qa-functional`

## Approving a Task (Functional Pass)
`castra task update --role qa-functional --status done 789`
*(Note: This sets `qa_approved=true`. Task status changes to `done` only if Security also approves.)*

## Rejecting a Task (Functional Fail)
`castra task update --role qa-functional --status todo 789`

## Logging a Test Report
`castra note add --role qa-functional --project 1 --content "Login flow failed on empty password field. Expected error message." --tags "qa,bug"`
