<!--
D.E.K.N. Summary - 30-second handoff for fresh Claude
Fill this at the END of your investigation, before marking Complete.
-->

## Summary (D.E.K.N.)

**Delta:** The reported bug "parent-child creates blocking dependency" is NOT present in current code; children of unblocked epics ARE ready.

**Evidence:** Unit tests pass (`TestEpicChildrenReady`); CLI testing confirms child issues are ready when parent epic is not blocked.

**Knowledge:** The constraint kb-8294e7 appears to reference stale behavior or misattributed friction. The `parent-child` dependency type only propagates blocking from ALREADY-blocked parents - it does not inherently block.

**Next:** Mark constraint kb-8294e7 as obsolete; close beads issue as "not reproducible".

**Promote to Decision:** recommend-no (bug investigation, not architectural)

---

# Investigation: Fix Epic Children Blocking Dependency

**Question:** Does `bd create --parent` create a blocking dependency that prevents children from being spawnable while the parent epic is open?

**Started:** 2026-01-08
**Updated:** 2026-01-08
**Owner:** Agent (spawned for orch-go-tuofe)
**Phase:** Complete
**Next Step:** None - bug not reproducible
**Status:** Complete

---

## Findings

### Finding 1: Code uses `parent-child` dependency type, not `blocks`

**Evidence:** Both direct mode (`cmd/bd/create.go:447`) and daemon mode (`internal/rpc/server_issues_epics.go:263`) create dependencies with `Type: types.DepParentChild`, not `types.DepBlocks`.

**Source:** 
- `cmd/bd/create.go:443-452` - Direct mode parent-child creation
- `internal/rpc/server_issues_epics.go:258-270` - Daemon mode parent-child creation

**Significance:** The dependency type is correct. Parent-child dependencies do NOT directly block - they only propagate blocking from parents that are already blocked by something else.

---

### Finding 2: Blocked cache SQL only blocks children of blocked parents

**Evidence:** The `blocked_cache.go` SQL has two parts:
1. `blocked_directly` - finds issues blocked by `blocks`, `conditional-blocks`, or `waits-for` dependencies
2. `blocked_transitively` - propagates blocking to children via `parent-child`

Parent-child only appears in the transitive propagation step, NOT in the direct blocking step.

**Source:** `internal/storage/sqlite/blocked_cache.go:138-234`

**Significance:** Parent-child dependencies are hierarchical, not blocking. A child is only blocked if its parent is blocked by something ELSE.

---

### Finding 3: Unit tests confirm children of unblocked epics are ready

**Evidence:** Created and ran `TestEpicChildrenReady` test:
```go
// Create unblocked epic (no dependencies)
// Create child task with parent-child dependency to epic
// Verify BOTH epic AND child appear in ready work
```
Test passes - both the epic and its child are ready.

**Source:** 
- `internal/storage/sqlite/epic_children_ready_test.go:12-63`
- Test output: `--- PASS: TestEpicChildrenReady (0.03s)`

**Significance:** The code works as designed. Children of unblocked epics ARE ready.

---

### Finding 4: CLI testing confirms expected behavior

**Evidence:** 
```bash
# Create epic and child
bd create "Test Epic" --type epic ...   # → test-gr5
bd create "Child 1" --parent test-gr5   # → test-gr5.1

# Check ready work
bd ready
# Output shows BOTH issues as ready:
# 1. [P2] [epic] test-gr5: Test Epic
# 2. [P2] [task] test-gr5.1: Child 1
```

**Source:** CLI testing in `/tmp/.beads/beads.db`

**Significance:** End-to-end confirmation that the bug is not present in the current codebase.

---

### Finding 5: Evidence issue orch-go-lv3yx.3 has no parent-child dependency

**Evidence:** Running `bd dep list orch-go-lv3yx.3` returns "has no dependencies", while sibling issues `.4`, `.5`, `.6`, `.7` all have parent-child dependencies.

**Source:** `bd dep list orch-go-lv3yx.{3,4,5,6,7}` in orch-go repo

**Significance:** The issue cited as evidence (`orch-go-lv3yx.3 required manual dep removal`) had its dependency removed manually. The current state cannot reproduce the original friction.

---

## Synthesis

**Key Insights:**

1. **Parent-child is NOT a blocking dependency type** - The `IsBlocking()` function returns true for parent-child for OTHER purposes (cache invalidation, ready work calculation), but the actual blocking logic only propagates blocking from already-blocked parents.

2. **The constraint may have been misattributed** - The friction event that led to kb-8294e7 may have been caused by:
   - A `blocks` dependency that was incorrectly added alongside the parent-child
   - User error in how the epic was set up
   - The issue being blocked by something else in the dependency graph

3. **The codebase has comprehensive tests** - Existing tests like `TestParentChildTransitiveBlocking`, `TestBlockerClosedUnblocksChildren`, and `TestRelatedDoesNotPropagate` all confirm the correct behavior.

**Answer to Investigation Question:**

No, `bd create --parent` does NOT create a blocking dependency. It creates a `parent-child` dependency which only propagates blocking from parents that are blocked by OTHER dependencies. Children of unblocked epics are ready to work on.

---

## Structured Uncertainty

**What's tested:**

- ✅ Unit test confirms children of unblocked epics are ready (`TestEpicChildrenReady` passes)
- ✅ CLI testing confirms both epic and children appear in `bd ready` output
- ✅ SQL in `blocked_cache.go` only uses parent-child for transitive propagation
- ✅ Both direct and daemon code paths use `DepParentChild` type

**What's untested:**

- ⚠️ Original friction scenario that led to kb-8294e7 cannot be reproduced
- ⚠️ Whether there was a temporary bug that has since been fixed
- ⚠️ Whether specific database state could trigger blocking

**What would change this:**

- Finding would be wrong if there's a code path that adds `blocks` dependency alongside parent-child
- Finding would be wrong if there's a database migration that changed behavior

---

## Implementation Recommendations

### Recommended Approach ⭐

**Mark constraint obsolete and close issue** - The bug is not reproducible in current code.

**Why this approach:**
- All tests pass confirming correct behavior
- CLI testing confirms expected behavior
- Code review shows correct dependency type usage

**Trade-offs accepted:**
- We don't know the original root cause of the friction
- The constraint may have been valid for an older version

**Implementation sequence:**
1. Mark kb-8294e7 as obsolete: `kb quick obsolete kb-8294e7 --reason "Bug not reproducible in current code"`
2. Close beads issue with "not reproducible" status
3. Keep new test `epic_children_ready_test.go` as regression protection

### Alternative Approaches Considered

**Option B: Add defensive code to prevent accidental blocking**
- **Pros:** Extra safety net
- **Cons:** Code already works correctly; adding code for non-existent bug
- **When to use instead:** If friction recurs

**Rationale for recommendation:** The code works correctly. Adding defensive code would add complexity without addressing a real problem.

---

## Implementation Details

**What to implement first:**
- Remove the test file `epic_children_ready_test.go` or keep it as documentation
- Mark constraint as obsolete

**Things to watch out for:**
- ⚠️ If friction recurs, investigate the specific scenario more carefully

**Success criteria:**
- ✅ Constraint kb-8294e7 marked obsolete
- ✅ Beads issue closed with clear reason

---

## References

**Files Examined:**
- `cmd/bd/create.go:443-452` - Direct mode parent-child dependency creation
- `internal/rpc/server_issues_epics.go:258-270` - Daemon mode parent-child dependency creation
- `internal/storage/sqlite/blocked_cache.go:138-234` - Blocking SQL logic
- `internal/types/types.go:589` - AffectsReadyWork() function

**Commands Run:**
```bash
# Test that children are ready
cd /tmp && bd init --prefix test
bd create "Test Epic" --type epic --no-understanding --reason "Testing"
bd create "Child 1" --parent test-gr5
bd ready  # Shows both as ready

# Check dependency types in orch-go
bd dep list orch-go-lv3yx.{3,4,5,6,7}  # .3 has no deps, others have parent-child
```

**Related Artifacts:**
- **Constraint:** kb-8294e7 - "bd create --parent creates blocking dependency on parent"
- **Issue:** orch-go-tuofe - Original issue tracking this bug

---

## Investigation History

**2026-01-08 16:00:** Investigation started
- Initial question: Does --parent create blocking dependency?
- Context: Constraint kb-8294e7 reported this as a bug

**2026-01-08 16:30:** Found code uses correct dependency type
- Both direct and daemon modes use `DepParentChild`

**2026-01-08 16:45:** Confirmed blocking logic is correct
- Parent-child only propagates blocking, doesn't cause it

**2026-01-08 17:00:** Investigation completed
- Status: Complete - Bug not reproducible
- Key outcome: Current code correctly handles epic children; constraint appears stale
