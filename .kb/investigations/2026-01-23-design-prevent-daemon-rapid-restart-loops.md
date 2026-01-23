# Investigation: Prevent Daemon Rapid Restart Loops

**Question:** How should beads prevent rapid daemon restart loops that cause SQLite WAL corruption?

**Started:** 2026-01-23
**Updated:** 2026-01-23
**Owner:** Agent (be-arch-prevent-daemon-rapid-22jan-4e15)
**Phase:** Complete
**Status:** Complete

**Prior Work:**
- Investigation: orch-go/.kb/investigations/2026-01-21-inv-investigate-beads-sqlite-database-corruption.md
- Model: orch-go/.kb/models/beads-database-corruption.md
- Decision: orch-go/.kb/decisions/2026-01-22-beads-daemon-rapid-restart-prevention.md

---

## TLDR

Implementing three-layer protection against rapid daemon restarts: (1) persistent restart backoff stored in `.beads/daemon-start-state.json`, (2) pre-flight validation before database open, (3) graceful degradation to direct mode.

---

## Problem Statement

Beads daemon enters rapid restart loops when any failure occurs (sandbox chmod, legacy fingerprint validation, etc.). Each restart cycle opens/closes the database performing WAL checkpoint. High-frequency checkpoints create race conditions that corrupt the database.

**Evidence:**
- Jan 21: 57 daemon restarts in one day (sandbox chmod failures)
- Jan 22: 9+ restarts in 3 minutes (legacy fingerprint validation)
- All incidents show 0-byte WAL file = incomplete checkpoint operation

**Root Cause:**
1. No backoff between daemon restart attempts (launchd/process managers cause immediate restart)
2. Pre-flight validation happens AFTER database open, not before
3. Current backoff (`daemon_autostart.go`) is **in-memory only** - doesn't survive process restarts

---

## Design Approach

### Layer 1: Persistent Restart Backoff (Critical)

Create `daemon_start_state.go` modeled after `daemon_sync_state.go`:

```go
// DaemonStartState tracks daemon startup attempts for backoff.
// Stored in .beads/daemon-start-state.json (gitignored, local-only).
type DaemonStartState struct {
    LastAttempt    time.Time `json:"last_attempt,omitempty"`
    AttemptCount   int       `json:"attempt_count"`
    BackoffUntil   time.Time `json:"backoff_until,omitempty"`
    LastError      string    `json:"last_error,omitempty"`
}

// Backoff schedule: 30s, 1m, 2m, 5m, 10m, 30m (cap)
// Minimum restart interval: 30s even for first attempt
```

**Key behavior:**
- Check backoff state BEFORE any daemon operations
- Minimum 30s between any daemon start attempts (prevents rapid loops)
- Exponential backoff on repeated failures
- State persists across process restarts
- Clear state on successful daemon startup

### Layer 2: Pre-flight Validation (Important)

Move validation checks BEFORE database open in `runDaemonLoop()`:

```go
func runDaemonLoop(...) {
    // PHASE 1: Pre-flight validation (NO database operations)
    if err := validatePreflightChecks(beadsDir); err != nil {
        log.Error("pre-flight validation failed", "error", err)
        recordDaemonStartFailure(beadsDir, err.Error())
        return
    }

    // PHASE 2: Database operations (safe to open now)
    store, err := sqlite.New(ctx, daemonDBPath)
    ...
}

func validatePreflightChecks(beadsDir string) error {
    // Already done: Sandbox detection (exists at line 296-312)

    // Move: Fingerprint validation (currently at line 470-479)
    // Move: Database version check (currently at line 482-513)
    // Move: Any other validation that requires DB to fail

    return nil
}
```

### Layer 3: Graceful Degradation (Recommended)

When daemon cannot start (backoff or pre-flight failure), ensure clean degradation:

1. Write error to `.beads/daemon-error` file (already done for sandbox)
2. Ensure `shouldAutoStartDaemon()` respects persistent backoff
3. CLI commands continue working in direct mode

---

## Implementation Plan

### File Changes

1. **NEW: `cmd/bd/daemon_start_state.go`**
   - `DaemonStartState` struct
   - `LoadDaemonStartState()` / `SaveDaemonStartState()`
   - `RecordDaemonStartAttempt()` / `RecordDaemonStartSuccess()`
   - `CanStartDaemon()` - checks backoff state

2. **MODIFY: `cmd/bd/daemon.go`**
   - Add backoff check at start of `runDaemonLoop()`
   - Move pre-flight validations before database open
   - Call `RecordDaemonStartSuccess()` after successful start

3. **MODIFY: `cmd/bd/daemon_autostart.go`**
   - Update `shouldAutoStartDaemon()` to check persistent backoff
   - Update `canRetryDaemonStart()` to use persistent state

4. **NEW: `cmd/bd/daemon_start_state_test.go`**
   - Tests for backoff logic
   - Tests for state persistence

---

## Success Criteria

- [ ] Daemon start attempts are rate-limited to minimum 30s interval
- [ ] Backoff state survives process restarts
- [ ] Pre-flight validation fails fast (before DB open) for known failure modes
- [ ] CLI commands continue working in direct mode when daemon cannot start
- [ ] Existing tests pass

---

## References

- **Pattern:** `daemon_sync_state.go` - similar persistent state for sync backoff
- **Decision:** orch-go decision document defines the three-layer approach
- **Model:** orch-go model explains WAL corruption mechanism
