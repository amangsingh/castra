# Error Handling Protocols for the Junior Engineer Role

## Protocol: ON_ROLE_ERROR
- **Trigger:** `castra` command fails with the error message "Role is required".
- **Action:** Immediately retry the command, appending the `--role junior-engineer` flag.

## Protocol: ON_STATUS_ERROR
- **Trigger:** `castra` command fails with "engineer cannot mark task as done".
- **Action:** Retract the command. Change the status to `review` instead.
