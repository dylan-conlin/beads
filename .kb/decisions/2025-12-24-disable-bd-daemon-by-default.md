# Decision: Disable bd daemon by default

**Date:** 2025-12-24
**Status:** Accepted

## Context

The bd daemon provides background sync (file watching, git auto-commit, faster queries via Unix socket). However, we experienced severe CPU overload from daemon process accumulation - 9+ daemon processes spawned and ran concurrently, causing 163% CPU usage.

Root cause is a race condition in `daemon_autostart.go` where concurrent `bd` calls can each spawn a daemon before others complete.

## Decision

Disable bd daemon globally via `BEADS_NO_DAEMON=1` in shell config.

## Rationale

- **Single-user workflow** - Daemon designed for multi-user/multi-process scenarios
- **Direct mode is functional** - Slightly slower but stable
- **Manual sync acceptable** - `bd sync` when needed vs auto-sync
- **Stability > speed** - CPU meltdowns are unacceptable

## Implementation

```bash
# In ~/.zshrc
export BEADS_NO_DAEMON=1
```

## Consequences

- Each `bd` command reads/writes files directly (slower by ~50-100ms)
- No automatic git sync - run `bd sync` manually when needed
- No file watching - changes from external tools need manual import
- Avoids daemon accumulation bug entirely

## Reversal

If bd-qgrf is fixed and daemon is stable:
```bash
# Remove from ~/.zshrc
unset BEADS_NO_DAEMON
```
