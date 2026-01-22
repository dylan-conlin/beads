//go:build unix

package main

import (
	"os"
	"syscall"
)

// isSandboxed detects if we're running in a sandboxed environment where daemon
// operations are unsafe or impossible.
//
// Detection strategy (checked in order):
// 1. CLAUDECODE=1 env var - Claude Code sandbox (most reliable for that environment)
// 2. /.dockerenv file - Standard Docker container marker
// 3. Signal 0 to self returns EPERM - Generic sandbox detection (Codex, seccomp)
//
// Why multiple checks:
// - Claude Code sandbox allows kill(self, 0) but cannot chmod Unix sockets
// - Docker containers have /.dockerenv but may allow signals
// - Some sandboxes (Codex) restrict signals but lack other markers
//
// The daemon cannot function in sandboxed environments because:
// - chmod on Unix socket fails (Claude Code: "invalid argument")
// - Rapid daemon start/fail cycles corrupt SQLite WAL
// - See: orch-go investigation 2026-01-21-inv-investigate-beads-sqlite-database-corruption.md
//
// Implements bd-u3t: Phase 2 auto-detection for GH #353
func isSandboxed() bool {
	// Check 1: Claude Code sandbox (sets CLAUDECODE=1)
	// This is the most reliable check for Claude Code's Linux sandbox,
	// which allows signals but cannot chmod Unix sockets on host filesystem.
	if os.Getenv("CLAUDECODE") == "1" {
		return true
	}

	// Check 2: Docker container (/.dockerenv marker file)
	// Standard marker present in most Docker containers.
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// Check 3: Signal restriction (fallback for other sandboxes)
	// Try to send signal 0 (existence check) to our own process.
	// Signal 0 doesn't actually send a signal, just checks permissions.
	pid := os.Getpid()
	err := syscall.Kill(pid, 0)

	if err == syscall.EPERM {
		// EPERM = Operation not permitted
		// We can't signal our own process, likely sandboxed (Codex, seccomp)
		return true
	}

	// No sandbox indicators detected
	return false
}
