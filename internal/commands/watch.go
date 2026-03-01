package commands

import (
	"castra/internal/cli"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func HandleWatch(role string) {
	db := GetDB()
	defer db.Close()

	var lastState string

	for {
		tasks, err := cli.ListAllTasksForRole(db, role)
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
