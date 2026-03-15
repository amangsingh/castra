package cli

import (
	"testing"
)

func TestArchetypeCRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// 1. Add
	id, err := AddArchetype(db, nil, "Backend", "Description", "senior-engineer", []string{"todo", "doing", "done"})
	if err != nil {
		t.Fatalf("AddArchetype failed: %v", err)
	}

	// 2. Get
	a, err := GetArchetype(db, id)
	if err != nil {
		t.Fatalf("GetArchetype failed: %v", err)
	}
	if a.Name != "Backend" {
		t.Errorf("Expected name Backend, got %s", a.Name)
	}

	// 3. List
	list, err := ListArchetypes(db, nil)
	if err != nil {
		t.Fatalf("ListArchetypes failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("Expected 1 archetype, got %d", len(list))
	}

	// 4. Soft Delete
	err = SoftDeleteArchetype(db, id)
	if err != nil {
		t.Fatalf("SoftDeleteArchetype failed: %v", err)
	}

	// 5. Verify deleted
	_, err = GetArchetype(db, id)
	if err == nil {
		t.Error("Expected error getting deleted archetype, got nil")
	}
}

func TestTaskWithArchetype(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	aid, _ := AddArchetype(db, nil, "Frontend", "Description", "junior-engineer", []string{"todo", "doing", "done"})

	// Add Task with Archetype
	tid, err := AddTask(db, 1, nil, nil, &aid, "UI Fix", "Fix button", "low")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	// Verify linkage
	task, err := GetTask(db, tid)
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}
	if task.ArchetypeID == nil || *task.ArchetypeID != aid {
		t.Errorf("Expected ArchetypeID %d, got %v", aid, task.ArchetypeID)
	}
}

func TestMilestoneWithArchetype(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	aid, _ := AddArchetype(db, nil, "Backend", "Description", "senior-engineer", []string{"todo", "doing", "done"})

	// Add Milestone with Archetype
	mid, err := AddMilestone(db, 1, nil, &aid, "Database Migration", "Description")
	if err != nil {
		t.Fatalf("AddMilestone failed: %v", err)
	}

	// Verify linkage
	milestones, err := ListMilestones(db, 1, "architect")
	if err != nil {
		t.Fatalf("ListMilestones failed: %v", err)
	}
	found := false
	for _, m := range milestones {
		if m.ID == mid {
			if m.ArchetypeID == nil || *m.ArchetypeID != aid {
				t.Errorf("Expected ArchetypeID %d, got %v", aid, m.ArchetypeID)
			}
			found = true
			break
		}
	}
	if !found {
		t.Error("Milestone not found in list")
	}
}
