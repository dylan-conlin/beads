---
linked_issues:
  - bd-ef1a
---
## Summary (D.E.K.N.)

**Delta:** `bd list --json` was returning `null` for `source_repo` because the running daemon binary was built before the fix commit.

**Evidence:** Restarting daemon with current binary fixed the issue - `source_repo: "."` now appears correctly.

**Knowledge:** Daemon processes persist across code changes; version mismatches between running binary and source cause subtle bugs.

**Next:** No code changes needed; issue resolved by daemon restart. Consider documenting daemon restart requirement after builds.

**Confidence:** Very High (98%) - Direct observation of fix via daemon restart.

---

# Investigation: bd list --json source_repo returns null

**Question:** Why does `bd list --json` return `source_repo: null` when the database has valid values?

**Started:** 2025-12-21
**Updated:** 2025-12-21
**Owner:** Agent
**Phase:** Complete
**Next Step:** None
**Status:** Complete
**Confidence:** Very High (98%)

---

## Findings

### Finding 1: Database has correct values

**Evidence:** SQLite query confirms `source_repo` column has values:
```bash
sqlite3 .beads/beads.db "SELECT id, source_repo FROM issues LIMIT 5"
# Returns: bd-03r|., bd-0b2|., etc.
```

**Source:** Direct database query

**Significance:** Ruled out storage layer as the problem.

---

### Finding 2: Storage layer scanIssues works correctly

**Evidence:** Unit test confirms SourceRepo is properly scanned and populated:
- Created test issue with `SourceRepo: "test-repo"`
- `SearchIssues` returned issue with `SourceRepo: "test-repo"`
- JSON marshaling included `"source_repo": "test-repo"`

**Source:** `internal/storage/sqlite/source_repo_test.go` (temporary test file)

**Significance:** Confirmed storage layer is correct; issue must be in daemon/RPC layer.

---

### Finding 3: Daemon binary version mismatch

**Evidence:**
- Running binary: `bd version 0.33.2 (dev: fix-repo-empty-config@d1cace404ade)`
- Socket created: `Dec 21 20:24` (8:24 PM)
- Fix committed: `Dec 21 22:46` (10:46 PM) in commit 973c5206
- After daemon restart, `source_repo: "."` appears correctly

**Source:** `bd version`, `ls -la .beads/bd.sock`, `git log`

**Significance:** Root cause identified - daemon was running code from before the fix was applied.

---

## Synthesis

**Key Insights:**

1. **Daemon persistence across code changes** - The bd daemon runs continuously and doesn't reload when code changes are made and rebuilt. This creates version skew between the source code and the running binary.

2. **The fix approach in 973c5206** - The fix added `json:"-"` to `Issue.SourceRepo` (to prevent JSONL export) and added a separate `IssueWithCounts.SourceRepo` field that's explicitly populated. This design is intentional to separate storage format from API format.

3. **Uncommitted changes suggest simpler approach** - There are uncommitted changes that revert to using embedded `Issue.SourceRepo` with `json:"source_repo,omitempty"`. This is a cleaner approach but requires consideration of JSONL export implications.

**Answer to Investigation Question:**

The `source_repo` field was returning null because the running daemon binary was built before commit 973c5206 which added the `IssueWithCounts.SourceRepo` field population. The daemon's socket was created at 8:24 PM, but the fix wasn't committed until 10:46 PM. Restarting the daemon with the current binary resolved the issue immediately.

---

## Confidence Assessment

**Current Confidence:** Very High (98%)

**Why this level?**

Direct observation of the fix: before daemon restart, `source_repo` was null; after restart, it shows `.` correctly.

**What's certain:**

- ✅ Database has correct `source_repo` values (`.` for local repo)
- ✅ Storage layer correctly populates `Issue.SourceRepo` during scan
- ✅ Daemon restart resolved the issue
- ✅ Binary version mismatch was the root cause

**What's uncertain:**

- ⚠️ Whether the uncommitted changes (simpler approach) are the intended final fix

---

## Test performed

**Test:** Restarted daemon and verified JSON output
```bash
bd daemon --stop && sleep 1 && bd daemon --start && sleep 1 && bd list --limit 1 --json | jq '.[0] | {id, source_repo}'
```

**Result:** 
```json
{
  "id": "bd-p5za",
  "source_repo": "."
}
```

The `source_repo` field now correctly shows `.` after daemon restart.

---

## Conclusion

The issue was caused by a version mismatch between the running daemon binary and the source code. The daemon was started before commit 973c5206 was made, which added the `source_repo` field to JSON API output. Restarting the daemon with the current binary (`go install` + `bd daemon --stop && bd daemon --start`) resolved the issue.

No code changes are required - the fix is already in place in the current source code.

---

## References

**Files Examined:**
- `internal/types/types.go` - Checked Issue and IssueWithCounts structs
- `internal/storage/sqlite/dependencies.go:701-790` - Examined scanIssues function
- `internal/storage/sqlite/queries.go:1461-1657` - Examined SearchIssues function
- `internal/rpc/server_issues_epics.go:760-800` - Examined handleList function
- `cmd/bd/list.go:515-543` - Examined JSON output code

**Commands Run:**
```bash
# Check database values
sqlite3 .beads/beads.db "SELECT id, source_repo FROM issues LIMIT 5"

# Check running binary version
/Users/dylanconlin/go/bin/bd version

# Check socket creation time
ls -la .beads/bd.sock

# Check git commit timestamps
git log --format="%H %ci" d1cace404ade -1

# Restart daemon and verify fix
bd daemon --stop && sleep 1 && bd daemon --start && sleep 1 && bd list --limit 1 --json | jq '.[0] | {id, source_repo}'
```

---

## Self-Review

- [x] Real test performed (not code review)
- [x] Conclusion from evidence (not speculation)
- [x] Question answered
- [x] File complete

**Self-Review Status:** PASSED
