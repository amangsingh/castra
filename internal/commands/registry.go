package commands

func NewDefaultRegistry() *Registry {
	r := NewRegistry()

	r.Register(&InitCommand{})

	logCmd := NewSubCommand("log", "View and add audit log entries")
	logCmd.Register(&LogAddCommand{})
	logCmd.Register(&LogListCommand{})
	r.Register(logCmd)

	projectCmd := NewSubCommand("project", "Manage projects")
	projectCmd.Register(&ProjectAddCommand{})
	projectCmd.Register(&ProjectListCommand{})
	projectCmd.Register(&ProjectViewCommand{})
	projectCmd.Register(&ProjectUpdateCommand{})
	projectCmd.Register(&ProjectDeleteCommand{})
	r.Register(projectCmd)

	milestoneCmd := NewSubCommand("milestone", "Manage milestones")
	milestoneCmd.Register(&MilestoneAddCommand{})
	milestoneCmd.Register(&MilestoneListCommand{})
	milestoneCmd.Register(&MilestoneViewCommand{})
	milestoneCmd.Register(&MilestoneUpdateCommand{})
	milestoneCmd.Register(&MilestoneDeleteCommand{})
	r.Register(milestoneCmd)

	sprintCmd := NewSubCommand("sprint", "Manage sprints")
	sprintCmd.Register(&SprintAddCommand{})
	sprintCmd.Register(&SprintListCommand{})
	r.Register(sprintCmd)

	taskCmd := NewSubCommand("task", "Manage tasks")
	taskCmd.Register(&TaskAddCommand{})
	taskCmd.Register(&TaskListCommand{})
	taskCmd.Register(&TaskViewCommand{})
	taskCmd.Register(&TaskUpdateCommand{})
	taskCmd.Register(&TaskDeleteCommand{})
	r.Register(taskCmd)

	noteCmd := NewSubCommand("note", "Manage project notes")
	noteCmd.Register(&NoteAddCommand{})
	noteCmd.Register(&NoteListCommand{})
	r.Register(noteCmd)

	archetypeCmd := NewSubCommand("archetype", "Manage task archetypes")
	archetypeCmd.Register(&ArchetypeAddCommand{})
	archetypeCmd.Register(&ArchetypeListCommand{})
	archetypeCmd.Register(&ArchetypeDeleteCommand{})
	r.Register(archetypeCmd)

	r.Register(&TUICommand{})
	r.Register(&WatchCommand{})

	return r
}
