package cli

import (
	"testing"
)

func TestWorkflowConstraints_AuditCycle(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// 1. Setup Archetype for Audit Cycle (Security Only - no QA gate)
	aid, _ := AddArchetype(db, nil, "Security Audit", "Description", "security-ops", []string{"todo", "doing", "review", "done"})

	// 2. Setup Task with this Archetype
	tid, err := AddTask(db, 1, nil, nil, &aid, "Security Audit Task", "Audit the router", "high")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	// 3. Move to review
	_, err = db.Exec("UPDATE tasks SET status = 'review' WHERE id = ?", tid)
	if err != nil {
		t.Fatal(err)
	}

	// 4. Security Approves -> Should become DONE immediately bypassing QA
	err = UpdateTaskStatus(db, tid, "done", "", "security-ops", false, "")
	if err != nil {
		t.Fatalf("UpdateTaskStatus failed: %v", err)
	}

	// 5. Verify status is DONE
	task, _ := GetTask(db, tid)
	if task.Status != "done" {
		t.Errorf("Expected status done (bypassed QA), got %s", task.Status)
	}
	if !task.SecurityApproved {
		t.Error("Security approval flag not set")
	}
}

func TestWorkflowConstraints_DefaultDualGate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// 1. Setup Standard Task (No Archetype)
	tid, _ := AddTask(db, 1, nil, nil, nil, "Standard Task", "Simple fix", "medium")

	// 2. Move to review
	db.Exec("UPDATE tasks SET status = 'review' WHERE id = ?", tid)

	// 3. Security Approves -> Should remain REVIEW (waiting for QA)
	err := UpdateTaskStatus(db, tid, "done", "", "security-ops", false, "")
	if err != nil {
		t.Fatalf("UpdateTaskStatus failed: %v", err)
	}

	task, _ := GetTask(db, tid)
	if task.Status != "review" {
		t.Errorf("Expected status review (waiting for QA), got %s", task.Status)
	}

	// 4. QA Approves -> Finally becomes DONE
	err = UpdateTaskStatus(db, tid, "done", "", "qa-functional", false, "")
	if err != nil {
		t.Fatalf("UpdateTaskStatus failed second time: %v", err)
	}

	task, _ = GetTask(db, tid)
	if task.Status != "done" {
		t.Errorf("Expected status done after both approvals, got %s", task.Status)
	}
}

func TestRoleEnforcement(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// 1. Setup Archetype for Senior Engineers
	aid, _ := AddArchetype(db, nil, "Senior Task", "Description", "senior-engineer", []string{"todo", "doing", "review", "done"})

	// 2. Setup Task
	tid, _ := AddTask(db, 1, nil, nil, &aid, "Complex Fix", "Fix kernel", "high")

	// 3. Junior tries to claim -> Should FAIL
	err := UpdateTaskStatus(db, tid, "doing", "", "junior-engineer", false, "")
	if err == nil {
		t.Error("Expected error when junior claims senior task, got nil")
	}

	// 4. Senior tries to claim -> Should SUCCEED
	err = UpdateTaskStatus(db, tid, "doing", "", "senior-engineer", false, "")
	if err != nil {
		t.Fatalf("Senior failed to claim senior task: %v", err)
	}

	// 5. Verify status
	task, _ := GetTask(db, tid)
	if task.Status != "doing" {
		t.Errorf("Expected status doing, got %s", task.Status)
	}
}
func TestBreakGlassProtocol(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// 1. Setup a Senior task with dual gates
	aid, _ := AddArchetype(db, nil, "Senior Task", "Description", "senior-engineer", []string{"todo", "doing", "review", "done"})
	tid, _ := AddTask(db, 1, nil, nil, &aid, "Mission Critical Fix", "Fix production bug", "high")

	// 2. Try to move to DONE directly from todo as architect WITHOUT break-glass -> Should FAIL gate check
	err := UpdateTaskStatus(db, tid, "done", "", "architect", false, "")
	if err != nil {
		t.Fatalf("UpdateTaskStatus (architect, no break-glass) failed: %v", err)
	}
	task, _ := GetTask(db, tid)
	if task.Status == "done" {
		t.Error("Architect forced 'done' without break-glass or gates")
	}

	// 3. Junior tries to use break-glass -> Should FAIL (restricted to architect)
	err = UpdateTaskStatus(db, tid, "done", "", "junior-engineer", true, "")
	if err == nil {
		t.Error("Expected error when junior tries to use break-glass, got nil")
	}

	// 4. Architect uses BREAK-GLASS to force DONE -> Should SUCCEED bypassing all gates and role checks
	err = UpdateTaskStatus(db, tid, "done", "", "architect", true, "")
	if err != nil {
		t.Fatalf("Break-glass failed for architect: %v", err)
	}

	// 5. Verify final state
	task, _ = GetTask(db, tid)
	if task.Status != "done" {
		t.Errorf("Expected status 'done' after break-glass override, got %s", task.Status)
	}
	// Verify bypass flags are set by break-glass (not forged approvals)
	if !task.QABypassed || !task.SecurityBypassed {
		t.Error("Break-glass failed to set bypass flags to true")
	}

	// Verify audit log entry for break-glass
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM audit_log WHERE action = 'status_change.break_glass' AND entity_id = ?`, tid).Scan(&count)
	if count == 0 {
		t.Error("No audit log entry found for break-glass action")
	}
}

func TestStateGuards_QASecurity(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// 1. Setup Task in 'todo'
	tid, _ := AddTask(db, 1, nil, nil, nil, "Bug Fix", "Some bug", "medium")

	// 2. QA tries to approve/reject task in 'todo' -> Should FAIL
	err := UpdateTaskStatus(db, tid, "done", "", "qa-functional", false, "")
	if err == nil {
		t.Error("QA should not be able to approve task in 'todo'")
	}

	err = UpdateTaskStatus(db, tid, "todo", "", "qa-functional", false, "Rejecting")
	if err == nil {
		t.Error("QA should not be able to reject task in 'todo'")
	}

	// 3. Move to review
	_, _ = db.Exec("UPDATE tasks SET status = 'review' WHERE id = ?", tid)

	// 4. QA can now act
	err = UpdateTaskStatus(db, tid, "done", "", "qa-functional", false, "")
	if err != nil {
		t.Errorf("QA failed to act on task in 'review': %v", err)
	}
}
