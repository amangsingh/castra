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

func (c *WatchCommand) Execute(ctx *Context) error {
	var lastState string

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

		// Convert to JSON
		jsonData, err := json.Marshal(tasks)
		if err != nil {
			log.Println("Error marshaling tasks:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		currentState := string(jsonData)

		// Emit if changed
		if currentState != lastState {
			fmt.Println(currentState)
			lastState = currentState
		}

		time.Sleep(1 * time.Second) // Poll every second
	}
}
