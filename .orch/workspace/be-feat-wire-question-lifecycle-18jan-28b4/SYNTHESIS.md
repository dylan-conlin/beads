# Session Synthesis

**Agent:** be-feat-wire-question-lifecycle-18jan-28b4
**Issue:** (ad-hoc spawn, no tracking)
**Duration:** 2026-01-18
**Outcome:** success

---

## TLDR

Implemented question lifecycle status validation in beads CLI. Questions can now only use open/investigating/answered/closed statuses, and reject assignee/estimate updates.

---

## Delta (What Changed)

### Files Created
- None

### Files Modified
- `internal/validation/bead.go` - Added ValidQuestionStatuses, IsValidQuestionStatus(), ValidateQuestionStatus()
- `internal/validation/bead_test.go` - Added TestIsValidQuestionStatus, TestValidateQuestionStatus
- `cmd/bd/show_unit_helpers.go` - Added validateQuestionUpdate() helper function
- `cmd/bd/update.go` - Added question validation calls in both daemon and direct modes

### Commits
- (To be committed - all changes ready)

---

## Evidence (What Was Observed)

- create.go already validates question-specific flags at lines 207-228, establishing CLI-level validation pattern
- types.go has StatusInvestigating and StatusAnswered constants added in Task 1 (commit 2dc8f7dc)
- update.go has two code paths (daemon RPC, direct) that both needed validation

### Tests Run
```bash
# Validation unit tests
go test ./internal/validation/... -v
# PASS: all tests including new TestIsValidQuestionStatus and TestValidateQuestionStatus

# Manual CLI testing
bd create --type question "Test question"
bd update test-7ia --status=investigating  # SUCCESS
bd update test-7ia --status=in_progress    # ERROR: invalid status
bd update test-7ia --assignee=someone      # ERROR: not applicable
bd close test-7ia --reason="Answered: ..." # SUCCESS
```

---

## Knowledge (What Was Learned)

### New Artifacts
- `.kb/investigations/2026-01-18-inv-wire-question-lifecycle-beads-task.md` - Investigation record

### Decisions Made
- CLI-level validation: Validation happens at CLI layer (update.go) rather than store layer, consistent with existing patterns in create.go

### Constraints Discovered
- RPC server-side validation not added - CLI validation only. If API clients bypass CLI, they could set invalid statuses.

### Externalized via `kn`
- N/A - tactical implementation following existing patterns

---

## Next (What Should Happen)

**Recommendation:** close

### If Close
- [x] All deliverables complete
- [x] Tests passing
- [x] Investigation file has `**Phase:** Complete`
- [ ] Ready for commit

---

## Unexplored Questions

**Questions that emerged during this session that weren't directly in scope:**
- Should status transitions be enforced as ordered? (e.g., must go open→investigating→answered→closed)
- Should server-side (RPC/store layer) validation be added for API clients?

**Areas worth exploring further:**
- Auto-linking investigations to questions via "Answers: <question-id>" field (mentioned in task spec)

**What remains unclear:**
- Whether the parallel investigation entity should auto-close when question is answered

---

## Session Metadata

**Skill:** feature-impl
**Model:** claude-opus-4-5
**Workspace:** `.orch/workspace/be-feat-wire-question-lifecycle-18jan-28b4/`
**Investigation:** `.kb/investigations/2026-01-18-inv-wire-question-lifecycle-beads-task.md`
**Beads:** (ad-hoc spawn)
