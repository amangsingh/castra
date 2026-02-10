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
		sid := fs.Int64("sprint", 0, "Sprint ID (optional)")
		title := fs.String("title", "", "Title")
		desc := fs.String("desc", "", "Description")
		prio := fs.String("prio", "medium", "Priority")
		fs.Parse(argsToParse)
		if *pid == 0 || *title == "" {
			log.Fatal("Project ID and Title required")
		}

		var sprintID *int64
		if *sid != 0 {
			sprintID = sid
		}

		id, err := cli.AddTask(db, *pid, sprintID, *title, *desc, *prio)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Task created: %d\n", id)

	case "list":
		fs := flag.NewFlagSet("task list", flag.ExitOnError)
		pid := fs.Int64("project", 0, "Project ID")
		sid := fs.Int64("sprint", 0, "Sprint ID")
		backlog := fs.Bool("backlog", false, "Show backlog only")
		fs.Parse(argsToParse)

		if *pid == 0 {
			log.Fatal("Project ID required")
		}

		var sprintID *int64
		if *sid != 0 {
			sprintID = sid
		}

		tasks, err := cli.ListTasks(db, *pid, sprintID, *backlog, role)
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
