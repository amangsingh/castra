package commands

import (
	"castra/internal/cli"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)
	subTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			MarginTop(1).
			MarginBottom(1)
	taskStyle = lipgloss.NewStyle().
			PaddingLeft(2)
	approvedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))
	pendingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208"))
)

type model struct {
	role     string
	projects []cli.Project
	sprints  map[int64][]cli.Sprint
	tasks    map[int64][]cli.Task // mapped by project id for simplicity
	err      error
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel(role string) model {
	m := model{
		role:    role,
		sprints: make(map[int64][]cli.Sprint),
		tasks:   make(map[int64][]cli.Task),
	}
	m.loadData()
	return m
}

func (m *model) loadData() {
	db := GetDB()
	defer db.Close()

	// Load projects
	projects, err := cli.ListProjects(db, false, false)
	if err != nil {
		m.err = err
		return
	}
	m.projects = projects

	// Load sprints and tasks for each project
	for _, p := range projects {
		sprints, _ := cli.ListSprints(db, p.ID)
		m.sprints[p.ID] = sprints

		tasks, _ := cli.ListAllTasksForRole(db, m.role)
		// Filter tasks for this project and active sprints
		var projectTasks []cli.Task

		// Build map of active sprint IDs for this project
		activeSprintIDs := make(map[int64]bool)
		for _, s := range sprints {
			if s.Status == "active" {
				activeSprintIDs[s.ID] = true
			}
		}

		for _, t := range tasks {
			if t.ProjectID == p.ID {
				// Only include if it belongs to an active sprint, or if there are no active sprints (fallback)
				if len(activeSprintIDs) == 0 || (t.SprintID != nil && activeSprintIDs[*t.SprintID]) {
					projectTasks = append(projectTasks, t)
				}
			}
		}
		m.tasks[p.ID] = projectTasks
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tickMsg:
		m.loadData()
		return m, tickCmd()
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	var sb strings.Builder

	header := fmt.Sprintf("Castra Dashboard - Role: %s (Auto-refreshing)", m.role)
	sb.WriteString(titleStyle.Render(header) + "\n\n")

	if len(m.projects) == 0 {
		sb.WriteString("No active projects found.\n")
	}

	for _, p := range m.projects {
		sb.WriteString(subTitleStyle.Render(fmt.Sprintf("Project: %s (#%d)", p.Name, p.ID)) + "\n")

		// Render active sprints
		sprints := m.sprints[p.ID]
		if len(sprints) > 0 {
			sb.WriteString("  Active Sprints:\n")
			for _, s := range sprints {
				if s.Status == "active" {
					sb.WriteString(fmt.Sprintf("    - %s (%s to %s)\n", s.Name, s.StartDate, s.EndDate))
				}
			}
		}

		// Render tasks
		tasks := m.tasks[p.ID]
		if len(tasks) == 0 {
			sb.WriteString(taskStyle.Render("No tasks for your role context.\n"))
		} else {
			for _, t := range tasks {
				qaMarker := "[QA Pending]"
				if t.QAApproved {
					qaMarker = approvedStyle.Render("[QA Pass]")
				}

				secMarker := "[SEC Pending]"
				if t.SecurityApproved {
					secMarker = approvedStyle.Render("[SEC Pass]")
				}

				statusMarker := pendingStyle.Render(fmt.Sprintf("(%s)", t.Status))
				if t.Status == "done" || t.Status == "review" {
					statusMarker = approvedStyle.Render(fmt.Sprintf("(%s)", t.Status))
				}

				taskLine := fmt.Sprintf("[%d] %s %s %s %s", t.ID, t.Title, statusMarker, qaMarker, secMarker)
				sb.WriteString(taskStyle.Render(taskLine) + "\n")
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\nPress 'q' or 'ctrl+c' to quit.\n")
	return sb.String()
}

func HandleTUI(role string) {
	if role == "" {
		log.Fatal("Role is required for TUI.")
	}
	p := tea.NewProgram(initialModel(role), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
