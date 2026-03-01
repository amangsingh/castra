# Error Handling Protocols for the Doc Writer Role

## Protocol: ON_ROLE_ERROR
- **Trigger:** `castra` command fails with the error message "Role is required".
- **Action:** Immediately retry the command, appending the `--role doc-writer` flag.

## Protocol: ON_PERMISSION_DENIED
- **Trigger:** `castra` command fails with "doc-writer cannot update tasks".
- **Action:** Do not attempt to update the task. Doc Writers are read-only for task status.
