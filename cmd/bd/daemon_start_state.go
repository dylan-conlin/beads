package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DaemonStartState tracks daemon startup attempts for backoff.
// Stored in .beads/daemon-start-state.json (gitignored, local-only).
// This prevents rapid restart loops that cause SQLite WAL corruption.
// See: orch-go decision 2026-01-22-beads-daemon-rapid-restart-prevention.md
type DaemonStartState struct {
	LastAttempt  time.Time `json:"last_attempt,omitempty"`
	AttemptCount int       `json:"attempt_count"`
	BackoffUntil time.Time `json:"backoff_until,omitempty"`
	LastError    string    `json:"last_error,omitempty"`
}

const (
	daemonStartStateFile = "daemon-start-state.json"
	// minRestartInterval is the minimum time between daemon start attempts.
	// This is enforced even on first attempt to prevent rapid loops.
	minRestartInterval = 30 * time.Second
	// maxDaemonBackoff caps the exponential backoff.
	maxDaemonBackoff = 30 * time.Minute
	// clearDaemonStateThreshold clears state after successful operation period.
	clearDaemonStateThreshold = 24 * time.Hour
)

var (
	// daemonStartBackoffSchedule defines the exponential backoff durations.
	// First failure: 30s, then 1m, 2m, 5m, 10m, 30m (cap).
	daemonStartBackoffSchedule = []time.Duration{
		30 * time.Second,
		1 * time.Minute,
		2 * time.Minute,
		5 * time.Minute,
		10 * time.Minute,
		30 * time.Minute,
	}
	// daemonStartStateMu protects concurrent access to start state file.
	daemonStartStateMu sync.Mutex
)

// LoadDaemonStartState loads the daemon start state from .beads/daemon-start-state.json.
// Returns empty state if file doesn't exist or is stale.
func LoadDaemonStartState(beadsDir string) DaemonStartState {
	daemonStartStateMu.Lock()
	defer daemonStartStateMu.Unlock()

	statePath := filepath.Join(beadsDir, daemonStartStateFile)
	data, err := os.ReadFile(statePath) // #nosec G304 - path constructed from beadsDir
	if err != nil {
		return DaemonStartState{}
	}

	var state DaemonStartState
	if err := json.Unmarshal(data, &state); err != nil {
		return DaemonStartState{}
	}

	// Clear stale state (successful operation for >24h means daemon is healthy)
	if !state.LastAttempt.IsZero() && state.AttemptCount == 0 &&
		time.Since(state.LastAttempt) > clearDaemonStateThreshold {
		_ = os.Remove(statePath)
		return DaemonStartState{}
	}

	return state
}

// SaveDaemonStartState saves the daemon start state to .beads/daemon-start-state.json.
func SaveDaemonStartState(beadsDir string, state DaemonStartState) error {
	daemonStartStateMu.Lock()
	defer daemonStartStateMu.Unlock()

	statePath := filepath.Join(beadsDir, daemonStartStateFile)

	// If state is empty/reset, remove the file
	if state.AttemptCount == 0 && state.BackoffUntil.IsZero() {
		_ = os.Remove(statePath)
		return nil
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statePath, data, 0600)
}

// ClearDaemonStartState removes the daemon start state file.
// Called after successful daemon startup.
func ClearDaemonStartState(beadsDir string) error {
	daemonStartStateMu.Lock()
	defer daemonStartStateMu.Unlock()

	statePath := filepath.Join(beadsDir, daemonStartStateFile)
	err := os.Remove(statePath)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// CanStartDaemon checks if daemon can start based on backoff state.
// Returns (canStart bool, backoffRemaining time.Duration).
// If canStart is false, backoffRemaining indicates how long to wait.
func CanStartDaemon(beadsDir string) (bool, time.Duration) {
	state := LoadDaemonStartState(beadsDir)

	// Check if in backoff period
	if !state.BackoffUntil.IsZero() && time.Now().Before(state.BackoffUntil) {
		return false, time.Until(state.BackoffUntil)
	}

	// Check minimum restart interval (even for first attempt after failure)
	if !state.LastAttempt.IsZero() && time.Since(state.LastAttempt) < minRestartInterval {
		return false, minRestartInterval - time.Since(state.LastAttempt)
	}

	return true, 0
}

// RecordDaemonStartAttempt records a daemon start attempt.
// This should be called at the very beginning of daemon startup,
// BEFORE any database operations.
func RecordDaemonStartAttempt(beadsDir string) {
	state := LoadDaemonStartState(beadsDir)
	state.LastAttempt = time.Now()
	_ = SaveDaemonStartState(beadsDir, state)
}

// RecordDaemonStartFailure records a failed daemon start attempt.
// This should be called when daemon startup fails for any reason.
// Returns the duration until next retry is allowed.
func RecordDaemonStartFailure(beadsDir string, reason string) time.Duration {
	state := LoadDaemonStartState(beadsDir)

	state.LastAttempt = time.Now()
	state.AttemptCount++
	state.LastError = reason

	// Calculate backoff duration
	backoffIndex := state.AttemptCount - 1
	if backoffIndex >= len(daemonStartBackoffSchedule) {
		backoffIndex = len(daemonStartBackoffSchedule) - 1
	}
	backoff := daemonStartBackoffSchedule[backoffIndex]

	state.BackoffUntil = time.Now().Add(backoff)

	_ = SaveDaemonStartState(beadsDir, state)
	return backoff
}

// RecordDaemonStartSuccess clears the daemon start state after successful startup.
// This should be called after daemon is fully operational (socket listening, DB open).
func RecordDaemonStartSuccess(beadsDir string) {
	_ = ClearDaemonStartState(beadsDir)
}

// GetDaemonStartState returns the current daemon start state for diagnostics.
// Used by `bd doctor` and logging.
func GetDaemonStartState(beadsDir string) DaemonStartState {
	return LoadDaemonStartState(beadsDir)
}
