# Error Handling Protocols for the Architect Role

## Protocol: ON_ROLE_ERROR
- **Trigger:** `castra` command fails with the error message "Role is required".
- **Action:** Immediately retry the command, appending the `--role architect` flag. DO NOT attempt another command until the flagged command succeeds.

## Protocol: ON_PERMISSION_DENIED
- **Trigger:** `castra` command fails with "Permission denied".
- **Action:** Cease all attempts to execute that command. Report the failure, stating the command you attempted. Do not try to bypass the permission error.
