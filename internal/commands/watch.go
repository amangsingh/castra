package commands

import (
	"castra/internal/cli"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type WatchCommand struct{}

func (c *WatchCommand) Name() string        { return "watch" }
func (c *WatchCommand) Description() string { return "Watch tasks state" }
func (c *WatchCommand) Usage() string       { return "castra watch" }

func (c *WatchCommand) ReadInfo() (string, string) {
	return "task", "task.watch"
}

func (c *WatchCommand) Execute(ctx *Context) error {
	var lastState string

	// Automated Scribe: tracks task IDs that were already 'done' to detect
	// newly completed tasks and auto-generate a doc-writer summary note.
	knownDone := make(map[int64]bool)
	scribalInit := false // first pass seeds knownDone without writing notes

	for {
		tasks, err := cli.ListAllTasksForRole(ctx.DB, ctx.Role)
		if err != nil {
			log.Println("Error querying tasks:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if tasks == nil {
			tasks = []cli.Task{}
		}

		// Automated Scribe: detect newly completed tasks
		currentDone := make(map[int64]bool)
		for _, t := range tasks {
			if t.Status == "done" {
				currentDone[t.ID] = true
				if scribalInit && !knownDone[t.ID] {
					// New completion detected — auto-generate a scribe note
					note := fmt.Sprintf(
						"[SCRIBE] Task #%d '%s' marked done. Priority: %s.",
						t.ID, t.Title, t.Priority,
					)
					_, _ = cli.AddNote(ctx.DB, t.ProjectID, &t.ID, note, "doc-writer")
					log.Printf("Automated Scribe: logged completion note for task %d", t.ID)
				}
			}
		}
		knownDone = currentDone
		scribalInit = true

		// Convert to JSON and emit if changed
		jsonData, err := json.Marshal(tasks)
		if err != nil {
			log.Println("Error marshaling tasks:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		currentState := string(jsonData)
		if currentState != lastState {
			fmt.Println(currentState)
			lastState = currentState
		}

		time.Sleep(1 * time.Second) // Poll every second
	}
}
