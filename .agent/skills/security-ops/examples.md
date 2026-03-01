# `castra` CLI Examples for the Security Ops Role

## Listing Tasks for Audit
`castra task list --role security-ops`

## Approving a Task (Security Pass)
`castra task update --role security-ops --status done 123`
*(Note: This sets `security_approved=true`. Task status changes to `done` only if QA also approves.)*

## Rejecting a Task (Security Fail)
`castra task update --role security-ops --status todo 123`

## Logging a Vulnerability Report
`castra note add --role security-ops --project 1 --content "CRITICAL: SQL injection vulnerability found in login handler." --tags "security,vulnerability"`
