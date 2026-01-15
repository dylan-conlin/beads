## Summary (D.E.K.N.)

**Delta:** Session close protocol now asks user before pushing, preventing accidental deploy triggers.

**Evidence:** `bd prime --full` shows step 6 as "ASK USER: Ready to push? (pushes can trigger deploys)" and `bd prime --mcp` shows "ASK USER before push".

**Knowledge:** Pushes can trigger deploys that disrupt production systems like collection runs. User should decide when to push.

**Promote to Decision:** recommend-no (tactical change per user request)

---

# Investigation: Update Bd Prime Session Close

**Question:** Where does bd prime output the session close protocol, and how should it be updated to require user confirmation before push?

**Started:** 2026-01-15
**Updated:** 2026-01-15
**Owner:** Worker agent
**Phase:** Complete
**Next Step:** None
**Status:** Complete

---

## Findings

### Finding 1: Session close protocol is defined in cmd/bd/prime.go

**Evidence:** The file contains two output functions:
- `outputCLIContext()` - full CLI mode (line 189)
- `outputMCPContext()` - minimal MCP mode (line 153)

Both have conditional logic for different modes (stealth, ephemeral, noPush, normal).

**Source:** `cmd/bd/prime.go:153-257`

**Significance:** Two places need updating for consistency: the normal mode CLI output and the MCP output.

---

### Finding 2: Normal mode had "git push (push to remote)" as step 6

**Evidence:** Original code at line 241-248:
```go
closeProtocol = `[ ] 1. git status              (check what changed)
...
[ ] 6. git push                (push to remote)`
closeNote = "**NEVER skip this.** Work is not done until pushed."
```

**Source:** `cmd/bd/prime.go:241-248`

**Significance:** This was the primary target for the change - step 6 and the note about "pushed".

---

### Finding 3: MCP mode also had "git push" in the protocol

**Evidence:** Original code at line 167:
```go
closeProtocol = "Before saying \"done\": git status → git add → bd sync → git commit → bd sync → git push"
```

**Source:** `cmd/bd/prime.go:167`

**Significance:** For consistency, MCP mode also needed updating to prompt user before push.

---

## Synthesis

**Key Insights:**

1. **Two output modes exist** - Both CLI and MCP modes needed updating for consistency.

2. **Config already supports no-push** - The code already had `noPush` mode via config, showing the pattern of user-controlled push behavior.

**Answer to Investigation Question:**

Session close protocol is defined in `cmd/bd/prime.go` in two functions. Updated both:
- CLI mode: Step 6 now says "ASK USER: Ready to push? (pushes can trigger deploys)"
- CLI mode: Note now says "Work is not done until committed. Push when ready to deploy."
- MCP mode: Protocol now ends with "ASK USER before push"

---

## Structured Uncertainty

**What's tested:**

- ✅ `bd prime --full` shows updated step 6 text (verified: ran command, saw output)
- ✅ `bd prime --mcp` shows updated protocol ending (verified: ran command, saw output)
- ✅ Code compiles successfully (verified: `go build ./cmd/bd`)

**What's untested:**

- ⚠️ Behavior in ephemeral branch mode (not changed, should still work)
- ⚠️ Behavior in stealth mode (not changed, should still work)
- ⚠️ Behavior in noPush mode (not changed, should still work)

**What would change this:**

- Finding would be wrong if MCP detection logic changes
- Finding would be wrong if new output modes are added

---

## Implementation Recommendations

### Recommended Approach: Direct code update

**Why this approach:**
- Minimal change to achieve goal
- Consistent across both output modes
- Follows existing pattern (noPush config mode)

**Implementation sequence:**
1. Update normal mode closeProtocol step 6 text
2. Update normal mode closeNote message
3. Update MCP mode closeProtocol ending

### Alternative Approaches Considered

**Option B: Add new config flag (no-auto-push)**
- **Pros:** Could be toggled per-user
- **Cons:** Task specifically asked to change default behavior, not add option
- **When to use instead:** If some users want old behavior

**Rationale for recommendation:** Task was explicit about changing the default message, not adding configuration.

---

## References

**Files Examined:**
- `cmd/bd/prime.go` - Main implementation file

**Commands Run:**
```bash
# Build and test full mode
go build ./cmd/bd && ./bd prime --full

# Test MCP mode
./bd prime --mcp
```

---

## Investigation History

**2026-01-15:** Investigation started
- Initial question: Where to update session close protocol to require user confirmation before push
- Context: Spawned from pw-h7yp to prevent accidental deploys from auto-push

**2026-01-15:** Investigation completed
- Status: Complete
- Key outcome: Updated both CLI and MCP mode outputs in cmd/bd/prime.go
