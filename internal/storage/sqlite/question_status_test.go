package sqlite

import (
	"testing"

	"github.com/steveyegge/beads/internal/types"
)

// TestQuestionAnsweredUnblocks verifies that when a question transitions from
// 'investigating' to 'answered' status, dependent work becomes unblocked.
// This is the core decidability substrate requirement for question workflow.
//
// Question lifecycle: open → investigating → answered → closed
// Blocking behavior:
//   - open: blocks dependent work
//   - investigating: blocks dependent work (question is being worked on)
//   - answered: UNBLOCKS dependent work (answer available, work can proceed)
//   - closed: unblocks dependent work (final state)
func TestQuestionAnsweredUnblocks(t *testing.T) {
	env := newTestEnv(t)

	// Create a question and a task that depends on it
	question := env.CreateIssueWith("Should we use X or Y?", types.StatusOpen, 1, types.TypeQuestion)
	task := env.CreateIssueWith("Implement the chosen approach", types.StatusOpen, 2, types.TypeTask)

	// Task depends on question
	env.AddDep(task, question)

	// Task should be blocked (question is open)
	env.AssertBlocked(task)

	// Move question to investigating status
	err := env.Store.UpdateIssue(env.Ctx, question.ID, map[string]interface{}{
		"status": types.StatusInvestigating,
	}, "test-user")
	if err != nil {
		t.Fatalf("UpdateIssue to investigating failed: %v", err)
	}

	// Task should still be blocked (question is being investigated)
	env.AssertBlocked(task)

	// Move question to answered status
	err = env.Store.UpdateIssue(env.Ctx, question.ID, map[string]interface{}{
		"status": types.StatusAnswered,
	}, "test-user")
	if err != nil {
		t.Fatalf("UpdateIssue to answered failed: %v", err)
	}

	// Task should now be READY (question is answered)
	env.AssertReady(task)
}

// TestQuestionInvestigatingBlocks verifies that 'investigating' status blocks
// dependent work, distinguishing it from 'answered'.
func TestQuestionInvestigatingBlocks(t *testing.T) {
	env := newTestEnv(t)

	// Create question in investigating status
	question := env.CreateIssueWith("How does X work?", types.StatusInvestigating, 1, types.TypeQuestion)
	task := env.CreateIssue("Use X for implementation")

	// Task depends on question
	env.AddDep(task, question)

	// Task should be blocked (question is investigating)
	env.AssertBlocked(task)
}

// TestQuestionTransitiveBlocking verifies that question blocking propagates
// through parent-child hierarchies.
func TestQuestionTransitiveBlocking(t *testing.T) {
	env := newTestEnv(t)

	// Create question
	question := env.CreateIssueWith("Architecture question", types.StatusInvestigating, 1, types.TypeQuestion)

	// Create epic that depends on question
	epic := env.CreateEpic("Implement feature")
	env.AddDep(epic, question)

	// Create task under epic
	task := env.CreateIssue("Subtask of feature")
	env.AddParentChild(task, epic)

	// Both epic and task should be blocked
	env.AssertBlocked(epic)
	env.AssertBlocked(task)

	// Answer the question
	err := env.Store.UpdateIssue(env.Ctx, question.ID, map[string]interface{}{
		"status": types.StatusAnswered,
	}, "test-user")
	if err != nil {
		t.Fatalf("UpdateIssue to answered failed: %v", err)
	}

	// Both should now be ready
	env.AssertReady(epic)
	env.AssertReady(task)
}

// TestMultipleQuestionsBlocking verifies that all blocking questions must be
// answered before dependent work is unblocked.
func TestMultipleQuestionsBlocking(t *testing.T) {
	env := newTestEnv(t)

	// Create two questions
	q1 := env.CreateIssueWith("Question 1", types.StatusInvestigating, 1, types.TypeQuestion)
	q2 := env.CreateIssueWith("Question 2", types.StatusInvestigating, 1, types.TypeQuestion)

	// Create task that depends on both
	task := env.CreateIssue("Task blocked by both questions")
	env.AddDep(task, q1)
	env.AddDep(task, q2)

	// Task should be blocked
	env.AssertBlocked(task)

	// Answer first question
	err := env.Store.UpdateIssue(env.Ctx, q1.ID, map[string]interface{}{
		"status": types.StatusAnswered,
	}, "test-user")
	if err != nil {
		t.Fatalf("UpdateIssue q1 to answered failed: %v", err)
	}

	// Task should still be blocked (q2 is still investigating)
	env.AssertBlocked(task)

	// Answer second question
	err = env.Store.UpdateIssue(env.Ctx, q2.ID, map[string]interface{}{
		"status": types.StatusAnswered,
	}, "test-user")
	if err != nil {
		t.Fatalf("UpdateIssue q2 to answered failed: %v", err)
	}

	// Task should now be ready
	env.AssertReady(task)
}

// TestClosedAlsoUnblocks verifies that 'closed' status also unblocks
// (not just 'answered').
func TestClosedAlsoUnblocks(t *testing.T) {
	env := newTestEnv(t)

	// Create question and task
	question := env.CreateIssueWith("Question", types.StatusInvestigating, 1, types.TypeQuestion)
	task := env.CreateIssue("Task")
	env.AddDep(task, question)

	// Task is blocked
	env.AssertBlocked(task)

	// Close the question
	env.Close(question, "Answered: Use option A")

	// Task should be ready
	env.AssertReady(task)
}

// TestDecisionSupersededUnblocks verifies that 'superseded' status
// (when added) would also unblock dependent work.
// Decisions have lifecycle: active → superseded
// A superseded decision no longer blocks work.
//
// TODO: This test is commented out until StatusSuperseded is added.
// func TestDecisionSupersededUnblocks(t *testing.T) {
//     // Similar pattern to answered status
// }
