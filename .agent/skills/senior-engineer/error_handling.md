# Error Handling Protocols for the Senior Engineer Role

## Protocol: ON_ROLE_ERROR
- **Trigger:** `castra` command fails with the error message "Role is required".
- **Action:** Immediately retry the command, appending the `--role senior-engineer` flag.

## Protocol: ON_STATUS_ERROR
- **Trigger:** `castra` command fails with "engineer cannot mark task as done".
- **Action:** Retract the command. Change the status to `review` instead, as Engineers cannot mark tasks as `done`.
