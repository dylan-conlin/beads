package sqlite

import (
	"context"
	"testing"

	"github.com/steveyegge/beads/internal/types"
)

// TestEpicChildrenReady tests that children of an UNBLOCKED epic are ready to work on.
// This is the core epic workflow: create epic, create children, children should be workable.
func TestEpicChildrenReady(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create an unblocked epic (no dependencies)
	epic := &types.Issue{
		Title:     "Test Epic",
		Status:    types.StatusOpen,
		Priority:  1,
		IssueType: types.TypeEpic,
	}
	if err := store.CreateIssue(ctx, epic, "test-user"); err != nil {
		t.Fatalf("Failed to create epic: %v", err)
	}

	// Create a child task
	child := &types.Issue{
		Title:     "Child Task",
		Status:    types.StatusOpen,
		Priority:  1,
		IssueType: types.TypeTask,
	}
	if err := store.CreateIssue(ctx, child, "test-user"); err != nil {
		t.Fatalf("Failed to create child: %v", err)
	}

	// Add parent-child dependency (child -> epic)
	dep := &types.Dependency{
		IssueID:     child.ID,
		DependsOnID: epic.ID,
		Type:        types.DepParentChild,
	}
	if err := store.AddDependency(ctx, dep, "test-user"); err != nil {
		t.Fatalf("Failed to add parent-child dependency: %v", err)
	}

	// Check ready work
	ready, err := store.GetReadyWork(ctx, types.WorkFilter{Status: types.StatusOpen})
	if err != nil {
		t.Fatalf("GetReadyWork failed: %v", err)
	}

	// Build set of ready IDs
	readyIDs := make(map[string]bool)
	for _, issue := range ready {
		readyIDs[issue.ID] = true
	}

	// Epic should be ready (it has no blockers)
	if !readyIDs[epic.ID] {
		t.Errorf("Expected epic %s to be ready (no blockers)", epic.ID)
	}

	// Child should ALSO be ready (parent is not blocked)
	if !readyIDs[child.ID] {
		t.Errorf("Expected child %s to be ready (parent is not blocked)", child.ID)
	}
}

// TestMultipleEpicChildrenReady tests that multiple children of an unblocked epic are ready.
func TestMultipleEpicChildrenReady(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create an unblocked epic
	epic := &types.Issue{
		Title:     "Test Epic",
		Status:    types.StatusOpen,
		Priority:  1,
		IssueType: types.TypeEpic,
	}
	if err := store.CreateIssue(ctx, epic, "test-user"); err != nil {
		t.Fatalf("Failed to create epic: %v", err)
	}

	// Create multiple children
	var children []*types.Issue
	for i := 0; i < 5; i++ {
		child := &types.Issue{
			Title:     "Child " + string(rune('A'+i)),
			Status:    types.StatusOpen,
			Priority:  1,
			IssueType: types.TypeTask,
		}
		if err := store.CreateIssue(ctx, child, "test-user"); err != nil {
			t.Fatalf("Failed to create child: %v", err)
		}
		children = append(children, child)

		// Add parent-child dependency
		dep := &types.Dependency{
			IssueID:     child.ID,
			DependsOnID: epic.ID,
			Type:        types.DepParentChild,
		}
		if err := store.AddDependency(ctx, dep, "test-user"); err != nil {
			t.Fatalf("Failed to add parent-child dependency: %v", err)
		}
	}

	// Check ready work
	ready, err := store.GetReadyWork(ctx, types.WorkFilter{Status: types.StatusOpen})
	if err != nil {
		t.Fatalf("GetReadyWork failed: %v", err)
	}

	// Build set of ready IDs
	readyIDs := make(map[string]bool)
	for _, issue := range ready {
		readyIDs[issue.ID] = true
	}

	// All children should be ready
	for _, child := range children {
		if !readyIDs[child.ID] {
			t.Errorf("Expected child %s to be ready (parent is not blocked)", child.ID)
		}
	}
}
