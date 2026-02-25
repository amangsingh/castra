package commands

import (
	"castra/internal/cli"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func HandleTask(role string) {
	subCmdIdx := GetSubcommandIndex()
	if subCmdIdx == -1 {
		return
	}
	cmd := os.Args[subCmdIdx]

	db := GetDB()
	defer db.Close()
	argsToParse := FilterArgs(os.Args[subCmdIdx+1:])

	switch cmd {
	case "add":
		if role != "architect" {
			log.Fatal("Only architect can add tasks")
		}
		fs := flag.NewFlagSet("task add", flag.ExitOnError)
		pid := fs.Int64("project", 0, "Project ID")
		mid := fs.Int64("milestone", 0, "Milestone ID (optional)")
		sid := fs.Int64("sprint", 0, "Sprint ID (optional)")
		title := fs.String("title", "", "Title")
		desc := fs.String("desc", "", "Description")
		prio := fs.String("prio", "medium", "Priority")
		fs.Parse(argsToParse)
		if *pid == 0 || *title == "" {
			log.Fatal("Project ID and Title required")
		}

		var milestoneID *int64
		if *mid != 0 {
			milestoneID = mid
		}

		var sprintID *int64
		if *sid != 0 {
			sprintID = sid
		}

		id, err := cli.AddTask(db, *pid, milestoneID, sprintID, *title, *desc, *prio)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Task created: %d\n", id)

	case "list":
		fs := flag.NewFlagSet("task list", flag.ExitOnError)
		pid := fs.Int64("project", 0, "Project ID")
		mid := fs.Int64("milestone", 0, "Milestone ID")
		sid := fs.Int64("sprint", 0, "Sprint ID")
		backlog := fs.Bool("backlog", false, "Show backlog only")
		fs.Parse(argsToParse)

		if *pid == 0 {
			log.Fatal("Project ID required")
		}

		var milestoneID *int64
		if *mid != 0 {
			milestoneID = mid
		}

		var sprintID *int64
		if *sid != 0 {
			sprintID = sid
		}

		tasks, err := cli.ListTasks(db, *pid, milestoneID, sprintID, *backlog, role)
		if err != nil {
			log.Fatal(err)
		}
		for _, t := range tasks {
			approvals := ""
			if t.QAApproved {
				approvals += "[QA]"
			}
			if t.SecurityApproved {
				approvals += "[SEC]"
			}
			fmt.Printf("[%d] %s (%s) %s\n", t.ID, t.Title, t.Status, approvals)
		}

	case "view":
		fs := flag.NewFlagSet("task view", flag.ExitOnError)
		fs.Parse(argsToParse)

		idParsed := fs.Args() // Remaining args after flags
		if len(idParsed) < 1 {
			log.Fatal("Task ID required")
		}
		id, _ := strconv.ParseInt(idParsed[0], 10, 64)

		task, err := cli.GetTask(db, id)
		if err != nil {
			log.Fatal(err)
		}

		approvals := ""
		if task.QAApproved {
			approvals += "[QA Approved] "
		}
		if task.SecurityApproved {
			approvals += "[Security Approved]"
		}

		fmt.Printf("--- Task [%d]: %s ---\n", task.ID, task.Title)
		fmt.Printf("Status:   %s %s\n", task.Status, approvals)
		fmt.Printf("Priority: %s\n", task.Priority)
		fmt.Printf("Project:  %d\n", task.ProjectID)
		if task.MilestoneID != nil {
			fmt.Printf("Milestone: %d\n", *task.MilestoneID)
		}
		if task.SprintID != nil {
			fmt.Printf("Sprint:   %d\n", *task.SprintID)
		}
		fmt.Printf("\nDescription:\n%s\n", task.Description)

		fmt.Println("\n--- Notes ---")
		notes, err := cli.ListNotes(db, task.ProjectID, &id, role)
		if err == nil {
			if len(notes) == 0 {
				fmt.Println("No notes found for this role.")
			} else {
				for _, n := range notes {
					fmt.Printf("[Tags: %s] %s\n", n.Tags, n.Content)
				}
			}
		}

		fmt.Println("\n--- Audit Log ---")
		logs, err := cli.ListAuditEntries(db, "task", &id)
		if err == nil {
			if len(logs) == 0 {
				fmt.Println("No logs found.")
			} else {
				for _, l := range logs {
					fmt.Printf("[%s] %s: %s\n", l.Timestamp, l.Role, l.Payload)
				}
			}
		}

	case "update":
		fs := flag.NewFlagSet("task update", flag.ExitOnError)
		status := fs.String("status", "", "New Status")

		fs.Parse(argsToParse)

		idParsed := fs.Args() // Remaining args after flags
		if len(idParsed) < 1 {
			log.Fatal("ID required")
		}
		id, _ := strconv.ParseInt(idParsed[0], 10, 64)

		if *status == "" {
			log.Fatal("Status required")
		}
		if err := cli.UpdateTaskStatus(db, id, *status, role); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Task updated.")

	case "delete": // soft
		if role != "architect" {
			log.Fatal("Only architect can delete tasks")
		}
		fs := flag.NewFlagSet("task delete", flag.ExitOnError)
		fs.Parse(argsToParse)

		idParsed := fs.Args() // Remaining args after flags
		if len(idParsed) < 1 {
			log.Fatal("ID required")
		}
		id, _ := strconv.ParseInt(idParsed[0], 10, 64)
		if err := cli.SoftDeleteTask(db, id); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Task deleted.")
	}
}
