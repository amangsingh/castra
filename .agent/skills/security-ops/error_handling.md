# Error Handling Protocols for the Security Ops Role

## Protocol: ON_ROLE_ERROR
- **Trigger:** `castra` command fails with the error message "Role is required".
- **Action:** Immediately retry the command, appending the `--role security-ops` flag.

## Protocol: ON_STATUS_ERROR
- **Trigger:** `castra` command fails with "can only process tasks in 'review' status".
- **Action:** Verify the task status first. If it is not `review`, do not attempt to approve/reject it. Log the issue.
