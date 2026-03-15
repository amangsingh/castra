package commands

import (
	"castra/internal/cli"
	"fmt"
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
	bypassedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("160")).
			PaddingLeft(1).
			PaddingRight(1).
			Bold(true)
)

type model struct {
	role       string
	projects   []cli.Project
	sprints    map[int64][]cli.Sprint
	tasks           map[int64][]cli.Task // mapped by project id
	archetypes      map[int64]cli.Archetype
	milestones      map[int64]cli.Milestone
	flatTasks       []cli.Task // flattened list for navigation
	cursor          int        // current selection index
	showDetails     bool       // toggle for detail view
	selectedTaskLog []cli.AuditEntry
	err             error
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel(role string) model {
	m := model{
		role:       role,
		sprints:    make(map[int64][]cli.Sprint),
		tasks:      make(map[int64][]cli.Task),
		archetypes: make(map[int64]cli.Archetype),
		milestones: make(map[int64]cli.Milestone),
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

	// Load archetypes (Task 140/143 alignment)
	archs, _ := cli.ListArchetypes(db, nil) // Show global archetypes in TUI for mapping

	for _, a := range archs {
		m.archetypes[a.ID] = a
	}

	// Load all milestones across projects
	for _, p := range projects {
		ms, _ := cli.ListMilestones(db, p.ID, m.role)
		for _, m_obj := range ms {
			m.milestones[m_obj.ID] = m_obj
		}
	}

	newFlatTasks := []cli.Task{}

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
					newFlatTasks = append(newFlatTasks, t)
				}
			}
		}
		m.tasks[p.ID] = projectTasks
	}
	m.flatTasks = newFlatTasks

	// Bounds check for cursor after refresh
	if m.cursor >= len(m.flatTasks) {
		m.cursor = len(m.flatTasks) - 1
	}
	if m.cursor < 0 && len(m.flatTasks) > 0 {
		m.cursor = 0
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
		case "up", "k":
			if !m.showDetails && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if !m.showDetails && m.cursor < len(m.flatTasks)-1 {
				m.cursor++
			}
		case "enter":
			if !m.showDetails && len(m.flatTasks) > 0 {
				m.showDetails = true
				m.loadAuditLogs()
			} else {
				m.showDetails = false
			}
		case "n":
			if !m.showDetails && len(m.flatTasks) > 0 {
				m.performTransition(false)
			}
		case "r":
			if !m.showDetails && len(m.flatTasks) > 0 {
				m.performTransition(true)
			}
		}
	case tickMsg:
		if !m.showDetails {
			m.loadData()
		}
		return m, tickCmd()
	}
	return m, nil
}

func (m *model) loadAuditLogs() {
	if len(m.flatTasks) == 0 || m.cursor >= len(m.flatTasks) {
		return
	}
	selected := m.flatTasks[m.cursor]
	db := GetDB()
	defer db.Close()

	entries, _ := cli.ListAuditEntries(db, "task", &selected.ID)
	if len(entries) > 5 {
		m.selectedTaskLog = entries[:5]
	} else {
		m.selectedTaskLog = entries
	}
}

func (m *model) performTransition(isReject bool) {
	if m.cursor >= len(m.flatTasks) {
		return
	}
	selected := m.flatTasks[m.cursor]
	db := GetDB()
	defer db.Close()

	var newStatus string
	if isReject {
		newStatus = "todo"
	} else {
		// Use cli logic to find next status
		statuses := cli.GetTaskStatuses(db, selected.ID)
		newStatus = cli.NextStatus(statuses, selected.Status)
	}

	if newStatus == "" {
		return
	}

	err := cli.UpdateTaskStatus(db, selected.ID, newStatus, "", m.role, false, "TUI quick action")
	if err != nil {
		m.err = err
	} else {
		m.loadData() // Refresh list
	}
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

				// Archetype & Pipeline (Task 131)
				archetypeLabel := ""
				pipelineBar := ""
				if t.ArchetypeID != nil {
					if arch, ok := m.archetypes[*t.ArchetypeID]; ok {
						archetypeLabel = fmt.Sprintf("[%s] ", arch.Name)
						
						// Render pipeline progress bar
						for _, s := range arch.Statuses {
							if s == t.Status {
								pipelineBar += "●"
							} else {
								pipelineBar += "○"
							}
						}
						pipelineBar = fmt.Sprintf(" %s", pipelineBar)
					}
				}

				// Break-glass Highlight (Task 134)
				bypassBadge := ""
				if t.QABypassed || t.SecurityBypassed {
					bypassBadge = " " + bypassedStyle.Render("BYPASSED")
				}

				// Find global index to check if selected
				globalIndex := -1
				for i, ft := range m.flatTasks {
					if ft.ID == t.ID {
						globalIndex = i
						break
					}
				}

				cursorIndicator := "  "
				taskInfo := fmt.Sprintf("[%d] %s%s %s%s%s %s%s", t.ID, archetypeLabel, t.Title, statusMarker, pipelineBar, bypassBadge, qaMarker, secMarker)
				lineContent := taskInfo

				if globalIndex == m.cursor {
					cursorIndicator = "> "
					lineContent = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("235")).Render(lineContent)
				}

				sb.WriteString(taskStyle.Render(cursorIndicator + lineContent) + "\n")
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\nPress 'q' or 'ctrl+c' to quit.\n")

	if m.showDetails {
		m.renderDetails(&sb)
	} else {
		m.renderAffordanceBar(&sb)
	}

	return sb.String()
}

func (m model) renderDetails(sb *strings.Builder) {
	if m.cursor >= len(m.flatTasks) {
		return
	}
	selected := m.flatTasks[m.cursor]

	sb.WriteString("\n" + lipgloss.NewStyle().Bold(true).Underline(true).Render("Task Details (Esc or Enter to close)") + "\n")
	sb.WriteString(fmt.Sprintf("ID: %d | Title: %s\n", selected.ID, selected.Title))
	if selected.MilestoneID != nil {
		if ms, ok := m.milestones[*selected.MilestoneID]; ok {
			hier := ms.Name
			if ms.ParentID != nil {
				if parent, ok := m.milestones[*ms.ParentID]; ok {
					hier = parent.Name + " / " + hier
				}
			}
			sb.WriteString(fmt.Sprintf("Milestone: %s\n", hier))
			if ms.Description != "" {
				sb.WriteString(fmt.Sprintf("Milestone Desc: %s\n", ms.Description))
			}
		}
	}
	sb.WriteString(fmt.Sprintf("Description: %s\n", selected.Description))
	sb.WriteString("\n" + subTitleStyle.Render("Latest Audit Entries:") + "\n")

	if len(m.selectedTaskLog) == 0 {
		sb.WriteString("  No audit entries found.\n")
	} else {
		for _, e := range m.selectedTaskLog {
			sb.WriteString(fmt.Sprintf("  [%s] %s: %s\n", e.Timestamp[11:19], e.Role, e.Payload))
		}
	}
}

func (m model) renderAffordanceBar(sb *strings.Builder) {
	if len(m.flatTasks) == 0 || m.cursor >= len(m.flatTasks) {
		return
	}
	selected := m.flatTasks[m.cursor]

	sb.WriteString("\n" + lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		PaddingTop(1).
		Render("Actions: [Enter] View Details | [↑/↓] Navigate") + " ")

	// Contextual actions based on role and status
	// (Simplified for now, just showing potential)
	if selected.Status != "done" {
		sb.WriteString(" | [n] Next Status")
	}
	if selected.Status == "review" {
		sb.WriteString(" | [r] Reject")
	}
}

type TUICommand struct{}

func (c *TUICommand) Name() string        { return "tui" }
func (c *TUICommand) Description() string { return "Launch Castra TUI dashboard" }
func (c *TUICommand) Usage() string       { return "castra tui" }

func (c *TUICommand) ReadInfo() (string, string) {
	return "project", "tui.launch"
}

func (c *TUICommand) Execute(ctx *Context) error {
	if ctx.Role == "" {
		return fmt.Errorf("role is required for TUI")
	}

	p := tea.NewProgram(initialModel(ctx.Role), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("alas, there's been an error: %v", err)
	}
	return nil
}
