## Summary (D.E.K.N.)

**Delta:** Question gates implemented by adding 'investigating' status to blocking status checks in the dependency system.

**Evidence:** Full workflow tested: bd dep add epic question creates dependency; bd ready excludes epic when question is open/investigating; bd blocked shows epic blocked by question; bd ready shows epic when question is answered/closed.

**Knowledge:** Questions block work via standard dependency mechanism - only requires adding 'investigating' to the list of blocking statuses. 'answered' status does NOT block, enabling Premise Before Solution pattern.

**Next:** Close this task. Question gates are operational.

**Promote to Decision:** recommend-no (tactical implementation, follows existing blocking patterns)

---

# Investigation: Implement Question Gates in Beads

**Question:** How do we make questions block work via the dependency system?

**Started:** 2026-01-18
**Updated:** 2026-01-18
**Owner:** Worker agent
**Phase:** Complete
**Next Step:** None
**Status:** Complete

---

## Findings

### Finding 1: Blocking Status List is the Key Mechanism

**Evidence:** The blocking mechanism is controlled by SQL status checks in three key files:
- `internal/storage/sqlite/blocked_cache.go` - Cache that determines blocked issues
- `internal/storage/sqlite/ready.go` - GetBlockedIssues function
- `internal/storage/sqlite/schema.go` - SQL views for ready/blocked issues

**Source:**
- blocked_cache.go:149-150 (regular blocks)
- blocked_cache.go:163-164 (conditional blocks)
- ready.go:497, 502, 513
- schema.go:220, 251, 253

**Significance:** Adding 'investigating' to these status lists is all that's needed - the existing dependency system handles everything else.

### Finding 2: Question Statuses Already Defined

**Evidence:** Prior work (commits 2dc8f7dc, d14cf911) already added:
- `TypeQuestion` IssueType
- `StatusInvestigating` Status
- `StatusAnswered` Status

**Source:** internal/types/types.go lines 361-363, 407

**Significance:** No new type definitions needed - just needed to wire investigating status into blocking logic.

### Finding 3: 'answered' Status Should NOT Block

**Evidence:** Per task requirements: "Questions in 'answered' or 'closed' status should NOT block (dependency satisfied)". This enables the Premise Before Solution pattern - work becomes available once the question is answered.

**Source:** Task specification

**Significance:** Only 'open' and 'investigating' statuses should block. 'answered' is the resolution state before formal close.

---

## Synthesis

**Key Insights:**

1. **Minimal Change Required** - The existing dependency system already handles blocking logic. Only needed to extend the status whitelist to include 'investigating'.

2. **Three-File Pattern** - Blocking statuses appear in blocked_cache.go (the cache), ready.go (GetBlockedIssues), and schema.go (views). All must stay in sync.

3. **Workflow Enables Premise Before Solution** - Questions block work until answered, forcing investigation before implementation.

**Answer to Investigation Question:**

Questions block work via the standard dependency mechanism. When `bd dep add <epic-id> <question-id>` is run, it creates a 'blocks' dependency. The issue appears blocked when the question has status 'open' or 'investigating', and becomes ready when the question is 'answered' or 'closed'.

---

## Structured Uncertainty

**What's tested:**

- bd dep add <epic> <question> creates dependency (verified: ran command)
- bd ready excludes epic when question is open (verified: tested)
- bd ready excludes epic when question is investigating (verified: tested)
- bd blocked shows epic blocked by question (verified: tested)
- bd ready shows epic when question is answered (verified: tested)
- bd ready shows epic when question is closed (verified: tested)

**What's untested:**

- Performance impact on large databases (not benchmarked)
- Interaction with external dependencies (not tested)
- Question as child of epic hierarchy (not tested)

**What would change this:**

- Finding would be wrong if cache rebuild performance degrades significantly
- Finding would be incomplete if questions need special display in bd blocked

---

## Implementation Recommendations

**Purpose:** Document the implementation approach taken.

### Recommended Approach (Implemented)

**Add 'investigating' to blocking status lists** - Minimal change to existing system

**Why this approach:**
- Reuses existing dependency infrastructure
- Consistent with how other blocking statuses work
- No new tables or schema changes needed

**Trade-offs accepted:**
- Questions appear in blocked list same as other issue types (no special indication)
- No special workflow for question lifecycle

**Implementation sequence:**
1. Updated blocked_cache.go - Add 'investigating' to rebuildBlockedCache status checks
2. Updated ready.go - Add 'investigating' to GetBlockedIssues status checks
3. Updated schema.go - Add 'investigating' to view definitions for consistency

---

## References

**Files Modified:**
- `internal/storage/sqlite/blocked_cache.go` - Lines 149-150, 163-164
- `internal/storage/sqlite/ready.go` - Lines 497, 502, 513
- `internal/storage/sqlite/schema.go` - Lines 220, 251, 253

**Commands Tested:**
```bash
# Create question and epic
./bd create --type question --title "Test Question: Should we implement feature X?" --priority 1
./bd create --type epic --title "Test Epic: Feature X implementation" --priority 1 --no-understanding --reason "Test"

# Add dependency
./bd dep add bd-wuvj bd-52el

# Verify blocking
./bd ready | grep bd-wuvj  # Should NOT show
./bd blocked | grep bd-wuvj  # Should show blocked by question

# Update question to investigating
./bd update bd-52el --status=investigating
./bd ready | grep bd-wuvj  # Should NOT show (still blocked)

# Update question to answered
./bd update bd-52el --status=answered
./bd ready | grep bd-wuvj  # Should show (unblocked)
```

**Related Artifacts:**
- **Decision:** /Users/dylanconlin/Documents/personal/orch-go/.kb/decisions/2026-01-18-questions-as-first-class-entities.md
- **Investigation:** /Users/dylanconlin/Documents/personal/orch-go/.kb/investigations/2026-01-18-inv-design-questions-first-class-entities.md

---

## Investigation History

**2026-01-18 12:27:** Investigation started
- Initial question: How do we make questions block work via the dependency system?
- Context: Task 3 of 4 in question entity epic

**2026-01-18 12:35:** Implementation complete
- Added 'investigating' to blocking status lists in 3 files
- Tested full workflow successfully
- Cleaned up test data

**2026-01-18 12:40:** Investigation completed
- Status: Complete
- Key outcome: Question gates operational via standard dependency mechanism
