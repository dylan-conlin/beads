<!--
D.E.K.N. Summary - 30-second handoff for fresh Claude
Fill this at the END of your investigation, before marking Complete.
-->

## Summary (D.E.K.N.)

**Delta:** Daemon command skips sandbox detection due to `noDbCommands` early return in `PersistentPreRun`, allowing daemon to start in sandbox, open database, fail at chmod, and corrupt WAL through repeated cycles.

**Evidence:** Code analysis shows daemon command in `noDbCommands` list (main.go:271) causes early return at line 299 BEFORE sandbox detection at lines 312-319; daemon then runs `runDaemonLoop()` which opens DB, starts RPC server, fails at `os.Chmod` on socket, closes DB.

**Knowledge:** Sandbox detection must run for ALL commands, not just database-using commands; daemon processes spawned via auto-start inherit no sandbox state from parent.

**Next:** Fix implemented - moved sandbox detection before `noDbCommands` return, added fail-fast in `runDaemonLoop()`, updated `shouldAutoStartDaemon()` to check sandbox.

**Promote to Decision:** recommend-no - Bug fix, not architectural choice.

---

# Investigation: Prevent Daemon Auto Start Sandbox

**Question:** Why does daemon auto-start inside Claude Code sandbox cause SQLite WAL corruption, and how to prevent it?

**Started:** 2026-01-21
**Updated:** 2026-01-21
**Owner:** be-arch (emma)
**Phase:** Complete
**Next Step:** None - fix implemented
**Status:** Complete

<!-- Lineage -->
**Related-Issue:** bd-07f8 (Prevent daemon auto-start in sandbox to avoid SQLite corruption)

---

## Findings

### Finding 1: Daemon command bypasses sandbox detection

**Evidence:** In `main.go`, the `PersistentPreRun` function has an early return for commands in `noDbCommands` (line 270-299). The daemon command is in this list (line 271), so it returns BEFORE sandbox detection at lines 312-319 ever executes.

**Source:** `cmd/bd/main.go:270-319`

**Significance:** When `bd daemon --start` is spawned (whether manually or via auto-start), it never runs sandbox detection. The daemon process is independent of the parent process and doesn't inherit sandbox state.

---

### Finding 2: Database opened before RPC server fails

**Evidence:** In `runDaemonLoop()` (daemon.go:286), the database is opened at line 412 (`sqlite.New(ctx, daemonDBPath)`), which enables WAL mode. The RPC server is started at line 498, and it fails at `os.Chmod(s.socketPath, 0600)` in `server_lifecycle_conn.go:36` because sandbox environments can't chmod on host filesystem.

**Source:** `cmd/bd/daemon.go:412, 498` and `internal/rpc/server_lifecycle_conn.go:34-39`

**Significance:** Each failed daemon start cycle: opens DB (enables WAL) -> fails at chmod -> closes DB (WAL checkpoint). Repeated cycles corrupt the WAL file.

---

### Finding 3: Auto-start doesn't check for sandbox

**Evidence:** `shouldAutoStartDaemon()` in `daemon_autostart.go:48-67` checks `BEADS_NO_DAEMON` env var and worktree status, but did NOT check `isSandboxed()`. This meant auto-start could attempt to spawn daemon in sandbox even if parent process set `noDaemon=true` (parent's global variable doesn't transfer to child process).

**Source:** `cmd/bd/daemon_autostart.go:48-67`

**Significance:** Even if the calling command detected sandbox and set `noDaemon=true`, the function that decides whether to auto-start daemon would return `true` and spawn a child process that wouldn't inherit the sandbox detection.

---

## Synthesis

**Key Insights:**

1. **Subprocess independence** - Spawned daemon processes don't inherit parent state. Each process must independently detect its environment.

2. **Defense in depth needed** - Multiple layers of protection needed: (a) prevent auto-start attempt, (b) detect sandbox early in daemon, (c) fail before opening database.

3. **Early detection critical** - Sandbox detection must happen before any code path that might open the database or create files that could fail due to sandbox restrictions.

**Answer to Investigation Question:**

The SQLite WAL corruption occurs because:
1. Daemon commands skip sandbox detection due to being in `noDbCommands` list
2. Daemon opens database (enabling WAL mode) before starting RPC server
3. RPC server fails at `os.Chmod` on Unix socket (sandbox restriction)
4. Daemon closes database (WAL checkpoint)
5. Repeated auto-start attempts cause rapid open/close cycles corrupting WAL

The fix adds three layers of protection:
1. Move sandbox detection before `noDbCommands` early return (affects all commands)
2. Add fail-fast sandbox check at start of `runDaemonLoop()` (protects daemon specifically)
3. Add sandbox check in `shouldAutoStartDaemon()` (prevents auto-start attempts)

---

## Structured Uncertainty

**What's tested:**

- ✅ Build compiles successfully (`go build ./cmd/bd/...`)
- ✅ Existing sandbox tests pass (`TestSandboxDetection`, `TestSandboxDetectionExists`)
- ✅ Daemon auto-start tests pass (`TestDaemonAutoStart`, `TestShouldAutoStartDaemon_Disabled`)
- ✅ Code paths verified through manual inspection

**What's untested:**

- ⚠️ Actual sandbox environment behavior (requires running in Claude Code sandbox)
- ⚠️ Corruption prevention effectiveness (would require reproducing corruption scenario)
- ⚠️ Performance impact of additional `isSandboxed()` calls (likely negligible)

**What would change this:**

- Finding would be wrong if sandbox detection has false positives in non-sandbox environments
- Finding would be wrong if there's another code path that opens database before sandbox check
- Solution would fail if `isSandboxed()` doesn't detect Claude Code sandbox correctly

---

## Implementation Recommendations

**Purpose:** Document the implemented fix for future reference.

### Implemented Approach ⭐

**Multi-layer sandbox protection** - Detect and prevent daemon operations at multiple points to ensure no database corruption from sandbox restrictions.

**Why this approach:**
- Defense in depth - multiple checks catch edge cases
- Fail-fast - errors occur before database is opened
- Subprocess-safe - each process independently detects sandbox

**Trade-offs accepted:**
- Additional `isSandboxed()` calls (minimal overhead - one syscall)
- More code paths to maintain

**Implementation sequence:**
1. Move sandbox detection before `noDbCommands` early return in `main.go`
2. Add fail-fast sandbox check at start of `runDaemonLoop()` in `daemon.go`
3. Add sandbox check in `shouldAutoStartDaemon()` in `daemon_autostart.go`

---

### Implementation Details

**Files modified:**
- `cmd/bd/main.go` - Moved sandbox detection block before `noDbCommands` early return
- `cmd/bd/daemon.go` - Added fail-fast sandbox check in `runDaemonLoop()`
- `cmd/bd/daemon_autostart.go` - Added `isSandboxed()` check in `shouldAutoStartDaemon()`

**Things to watch out for:**
- ⚠️ Ensure `isSandboxed()` doesn't false positive in CI/test environments
- ⚠️ Users manually running daemon in sandbox will now get error message

**Success criteria:**
- ✅ No daemon auto-start attempts in sandboxed environments
- ✅ Daemon refuses to start with clear error message when sandboxed
- ✅ No SQLite WAL corruption from failed daemon starts

---

## References

**Files Examined:**
- `cmd/bd/main.go` - PersistentPreRun flow, noDbCommands check, sandbox detection
- `cmd/bd/daemon.go` - runDaemonLoop(), database opening
- `cmd/bd/daemon_autostart.go` - shouldAutoStartDaemon(), tryAutoStartDaemon()
- `cmd/bd/sandbox_unix.go` - isSandboxed() implementation
- `internal/rpc/server_lifecycle_conn.go` - Server.Start(), chmod failure point

**Commands Run:**
```bash
# Verify build
/usr/local/go/bin/go build ./cmd/bd/...

# Run related tests
/usr/local/go/bin/go test -v ./cmd/bd/... -run "Daemon|Sandbox|AutoStart" -count=1
```

**Related Artifacts:**
- **Issue:** bd-07f8 (Prevent daemon auto-start in sandbox to avoid SQLite corruption)
- **Prior Decision:** .kb/decisions/2025-12-24-disable-bd-daemon-by-default.md

---

## Investigation History

**2026-01-21 16:00:** Investigation started
- Initial question: Why does daemon auto-start in Claude Code sandbox cause SQLite WAL corruption?
- Context: Repeated daemon start failures causing rapid DB open/close cycles corrupt WAL file

**2026-01-21 16:30:** Root cause identified
- Found that daemon command bypasses sandbox detection via noDbCommands early return
- Found that database is opened before RPC server chmod failure

**2026-01-21 17:00:** Fix implemented
- Added three layers of sandbox protection
- All tests pass

**2026-01-21 17:15:** Investigation completed
- Status: Complete
- Key outcome: Multi-layer fix prevents daemon from opening database in sandboxed environments
