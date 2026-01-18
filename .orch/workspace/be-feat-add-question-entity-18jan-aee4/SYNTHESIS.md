# Session Synthesis

**Agent:** be-feat-add-question-entity-18jan-aee4
**Issue:** (ad-hoc spawn - no beads tracking)
**Duration:** ~15 minutes
**Outcome:** success

---

## TLDR

Added 'question' as a first-class beads entity type with its own status lifecycle (investigating, answered) and validation that rejects inapplicable flags (assignee, estimate, repro, understanding).

---

## Delta (What Changed)

### Files Created
- None

### Files Modified
- `internal/types/types.go` - Added TypeQuestion constant, StatusInvestigating and StatusAnswered constants, updated IsValid() methods
- `cmd/bd/create.go` - Added validation for question type (rejects assignee, estimate, repro, understanding flags), updated help text

### Commits
- (Not yet committed - ready for commit)

---

## Evidence (What Was Observed)

- Issue types are defined in internal/types/types.go with IssueType constants (line 391-405)
- Status types are defined in same file with Status constants (line 352-362)
- Both have IsValid() methods that need updating when adding new values
- Type-specific validation exists in create.go for bugs (repro required) and epics (understanding required)

### Tests Run
```bash
# Build
go build ./cmd/bd/
# Success - no errors

# Test question creation
cd /tmp/beads-question-test && ./bd create --type question --title 'Should we add caching?'
# ✓ Created issue: test-c97

# Test filtering
./bd list --type question
# test-c97 [P2] [question] open - Should we add a caching layer?

# Test status transitions
./bd update test-c97 --status investigating
./bd update test-c97 --status answered
# Both successful

# Test validation
./bd create --type question --title 'Test' --assignee someone
# Error: --assignee flag is not applicable to questions ✓

./bd create --type question --title 'Test' --estimate 60
# Error: --estimate flag is not applicable to questions ✓

# Unit tests
go test ./internal/types/...
# PASS: ok github.com/steveyegge/beads/internal/types 0.007s
```

---

## Knowledge (What Was Learned)

### New Artifacts
- `.kb/investigations/2026-01-18-inv-add-question-entity-type-beads.md` - Implementation investigation

### Decisions Made
- Decision 1: Add question-specific statuses (investigating, answered) because questions have a different lifecycle than work items
- Decision 2: Reject inapplicable flags at create time rather than silently ignoring them, because explicit errors help users understand the question schema

### Constraints Discovered
- Questions don't support: assignee, estimate, repro, understanding fields
- Questions support: title, description, status, priority, labels

### Externalized via `kn`
- None needed - this was a tactical implementation

---

## Next (What Should Happen)

**Recommendation:** close

### If Close
- [x] All deliverables complete
- [x] Tests passing
- [x] Investigation file has `**Phase:** Complete`
- [ ] Ready for orchestrator review and commit

---

## Unexplored Questions

**Questions that emerged during this session that weren't directly in scope:**
- Should questions have their own default priority (e.g., P3 instead of P2)?
- Should there be a `bd question` alias command for convenience?
- How should questions interact with the ready/blocked system?

**Areas worth exploring further:**
- Question-to-decision promotion workflow
- Question threading/conversation support

**What remains unclear:**
- How questions should appear in reports/dashboards (same as issues or special treatment?)

---

## Session Metadata

**Skill:** feature-impl
**Model:** Claude Opus 4.5
**Workspace:** `.orch/workspace/be-feat-add-question-entity-18jan-aee4/`
**Investigation:** `.kb/investigations/2026-01-18-inv-add-question-entity-type-beads.md`
**Beads:** (ad-hoc spawn)
