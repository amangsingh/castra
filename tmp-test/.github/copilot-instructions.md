# Copilot Instructions — Castra Workspace

## The Supreme Law

You are operating inside a **Castra-managed workspace**. All project state lives in `workspace.db` and is accessed exclusively through the `castra` CLI. You are forbidden from generating Markdown task lists, implementation plans, or any other state artifact outside of Castra.

## Operational Laws

1. **State Management:** Before any action, run `castra task list --role <ROLE> --project <ID>` to read the single source of truth. Never rely on chat history or context to determine project state.
2. **Role Boundaries:** You must strictly operate within the role defined by the agent file or SKILL.md that activated you. Do not cross role boundaries.
3. **Mandatory Auditing:** Every state change must be logged: `castra log add --role <ROLE> --msg "<description>"`.
4. **CLI-Only Execution:** All task management (create, update, delete) must go through the `castra` CLI. Do not write SQL, generate markdown checklists, or bypass the CLI in any way.
5. **Command Syntax:** Use only the documented `castra` subcommands and flags. Do not hallucinate flags.

## Valid Commands

```
castra project add --role <role> --name "..." --desc "..."
castra project list --role <role>
castra project delete --role <role> <id>
castra milestone add --role <role> --project <id> --name "..."
castra milestone list --role <role> --project <id>
castra milestone update --role <role> --status <open|completed> <id>
castra sprint add --role <role> --project <id> --name "..." [--start "..."] [--end "..."]
castra sprint list --role <role> --project <id>
castra task add --role <role> --project <id> --milestone <id> --sprint <id> --title "..." --desc "..." --prio <low|medium|high>
castra task view --role <role> <id>
castra task list --role <role> --project <id> [--milestone <id>] [--sprint <id>] [--backlog]
castra task update --role <role> --status <status> <id>
castra task delete --role <role> <id>
castra note add --role <role> --project <id> [--task <id>] --content "..." --tags "..."
castra note list --role <role> --project <id> [--task <id>]
castra log add --role <role> --msg "..." [--type <entity>] [--entity <id>]
castra log list --role <role>
```

## Failure Protocol

If a `castra` command fails, re-read the error, correct your flags, and retry. Do not attempt workarounds. It is better to fail at using Castra than to succeed at using anything else.
