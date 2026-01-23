package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDaemonStartState_EmptyState(t *testing.T) {
	tempDir := t.TempDir()

	// Load state from empty directory
	state := LoadDaemonStartState(tempDir)

	if state.AttemptCount != 0 {
		t.Errorf("Expected AttemptCount=0 for empty state, got %d", state.AttemptCount)
	}
	if !state.LastAttempt.IsZero() {
		t.Errorf("Expected zero LastAttempt for empty state, got %v", state.LastAttempt)
	}
	if !state.BackoffUntil.IsZero() {
		t.Errorf("Expected zero BackoffUntil for empty state, got %v", state.BackoffUntil)
	}
}

func TestDaemonStartState_SaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()

	// Create and save state
	state := DaemonStartState{
		LastAttempt:  time.Now(),
		AttemptCount: 3,
		BackoffUntil: time.Now().Add(5 * time.Minute),
		LastError:    "test error",
	}

	err := SaveDaemonStartState(tempDir, state)
	if err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}

	// Load and verify
	loaded := LoadDaemonStartState(tempDir)

	if loaded.AttemptCount != state.AttemptCount {
		t.Errorf("AttemptCount mismatch: expected %d, got %d", state.AttemptCount, loaded.AttemptCount)
	}
	if loaded.LastError != state.LastError {
		t.Errorf("LastError mismatch: expected %q, got %q", state.LastError, loaded.LastError)
	}
}

func TestDaemonStartState_ClearState(t *testing.T) {
	tempDir := t.TempDir()

	// Create state file
	state := DaemonStartState{
		AttemptCount: 5,
		LastError:    "some error",
	}
	_ = SaveDaemonStartState(tempDir, state)

	// Verify file exists
	statePath := filepath.Join(tempDir, daemonStartStateFile)
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Fatal("State file should exist before clear")
	}

	// Clear state
	err := ClearDaemonStartState(tempDir)
	if err != nil {
		t.Fatalf("Failed to clear state: %v", err)
	}

	// Verify file is gone
	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		t.Error("State file should be removed after clear")
	}
}

func TestCanStartDaemon_NoState(t *testing.T) {
	tempDir := t.TempDir()

	canStart, backoff := CanStartDaemon(tempDir)

	if !canStart {
		t.Error("Expected canStart=true when no state exists")
	}
	if backoff != 0 {
		t.Errorf("Expected backoff=0 when no state exists, got %v", backoff)
	}
}

func TestCanStartDaemon_InBackoff(t *testing.T) {
	tempDir := t.TempDir()

	// Set up state with active backoff
	state := DaemonStartState{
		LastAttempt:  time.Now(),
		AttemptCount: 1,
		BackoffUntil: time.Now().Add(5 * time.Minute),
	}
	_ = SaveDaemonStartState(tempDir, state)

	canStart, backoff := CanStartDaemon(tempDir)

	if canStart {
		t.Error("Expected canStart=false when in backoff period")
	}
	if backoff <= 0 {
		t.Errorf("Expected positive backoff remaining, got %v", backoff)
	}
	if backoff > 5*time.Minute {
		t.Errorf("Expected backoff <= 5 minutes, got %v", backoff)
	}
}

func TestCanStartDaemon_BackoffExpired(t *testing.T) {
	tempDir := t.TempDir()

	// Set up state with expired backoff
	state := DaemonStartState{
		LastAttempt:  time.Now().Add(-10 * time.Minute),
		AttemptCount: 1,
		BackoffUntil: time.Now().Add(-5 * time.Minute), // Expired 5 minutes ago
	}
	_ = SaveDaemonStartState(tempDir, state)

	canStart, _ := CanStartDaemon(tempDir)

	if !canStart {
		t.Error("Expected canStart=true when backoff has expired")
	}
}

func TestCanStartDaemon_MinRestartInterval(t *testing.T) {
	tempDir := t.TempDir()

	// Set up state with recent attempt but no explicit backoff
	state := DaemonStartState{
		LastAttempt:  time.Now().Add(-5 * time.Second), // 5 seconds ago
		AttemptCount: 0,                                // No failures, just a recent attempt
	}
	_ = SaveDaemonStartState(tempDir, state)

	canStart, backoff := CanStartDaemon(tempDir)

	if canStart {
		t.Error("Expected canStart=false due to minimum restart interval")
	}
	if backoff <= 0 {
		t.Errorf("Expected positive backoff remaining, got %v", backoff)
	}
}

func TestRecordDaemonStartAttempt(t *testing.T) {
	tempDir := t.TempDir()

	before := time.Now()
	RecordDaemonStartAttempt(tempDir)
	after := time.Now()

	state := LoadDaemonStartState(tempDir)

	if state.LastAttempt.Before(before) || state.LastAttempt.After(after) {
		t.Errorf("LastAttempt not set correctly: %v (expected between %v and %v)",
			state.LastAttempt, before, after)
	}
}

func TestRecordDaemonStartFailure_ExponentialBackoff(t *testing.T) {
	tempDir := t.TempDir()

	// Record first failure
	backoff1 := RecordDaemonStartFailure(tempDir, "first error")
	if backoff1 != 30*time.Second {
		t.Errorf("First failure should have 30s backoff, got %v", backoff1)
	}

	state := LoadDaemonStartState(tempDir)
	if state.AttemptCount != 1 {
		t.Errorf("AttemptCount should be 1 after first failure, got %d", state.AttemptCount)
	}
	if state.LastError != "first error" {
		t.Errorf("LastError should be 'first error', got %q", state.LastError)
	}

	// Record second failure
	backoff2 := RecordDaemonStartFailure(tempDir, "second error")
	if backoff2 != 1*time.Minute {
		t.Errorf("Second failure should have 1m backoff, got %v", backoff2)
	}

	// Record third failure
	backoff3 := RecordDaemonStartFailure(tempDir, "third error")
	if backoff3 != 2*time.Minute {
		t.Errorf("Third failure should have 2m backoff, got %v", backoff3)
	}

	// Verify backoff schedule continues correctly
	expectedBackoffs := []time.Duration{
		5 * time.Minute,  // 4th failure
		10 * time.Minute, // 5th failure
		30 * time.Minute, // 6th failure
		30 * time.Minute, // 7th failure (capped)
	}

	for i, expected := range expectedBackoffs {
		actual := RecordDaemonStartFailure(tempDir, "error")
		if actual != expected {
			t.Errorf("Failure %d should have %v backoff, got %v", i+4, expected, actual)
		}
	}
}

func TestRecordDaemonStartSuccess_ClearsState(t *testing.T) {
	tempDir := t.TempDir()

	// Set up some failure state
	_ = RecordDaemonStartFailure(tempDir, "test error")
	_ = RecordDaemonStartFailure(tempDir, "another error")

	state := LoadDaemonStartState(tempDir)
	if state.AttemptCount != 2 {
		t.Fatalf("Expected 2 attempts before success, got %d", state.AttemptCount)
	}

	// Record success
	RecordDaemonStartSuccess(tempDir)

	// Verify state is cleared
	state = LoadDaemonStartState(tempDir)
	if state.AttemptCount != 0 {
		t.Errorf("AttemptCount should be 0 after success, got %d", state.AttemptCount)
	}
}

func TestGetDaemonStartState(t *testing.T) {
	tempDir := t.TempDir()

	// Record a failure
	_ = RecordDaemonStartFailure(tempDir, "diagnostic test")

	// Get state for diagnostics
	state := GetDaemonStartState(tempDir)

	if state.AttemptCount != 1 {
		t.Errorf("Expected AttemptCount=1, got %d", state.AttemptCount)
	}
	if state.LastError != "diagnostic test" {
		t.Errorf("Expected LastError='diagnostic test', got %q", state.LastError)
	}
}
