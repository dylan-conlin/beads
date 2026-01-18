<!--
D.E.K.N. Summary - 30-second handoff for fresh Claude
Fill this at the END of your investigation, before marking Complete.
-->

## Summary (D.E.K.N.)

**Delta:** Added 'question' as a first-class beads entity type with its own status lifecycle (investigating, answered).

**Evidence:** Manual testing confirmed: `bd create --type question`, `bd list --type question`, `bd show`, and status transitions all work correctly. Unit tests pass.

**Knowledge:** Adding a new entity type requires: (1) adding type constant to types.go, (2) updating IsValid(), (3) adding any type-specific statuses, (4) adding validation in create.go, (5) updating help text.

**Next:** Close - all success criteria met.

**Promote to Decision:** recommend-no (tactical addition, not architectural)

---

# Investigation: Add Question Entity Type Beads

**Question:** How to add 'question' as a first-class beads entity type?

**Started:** 2026-01-18
**Updated:** 2026-01-18
**Owner:** Worker Agent (be-feat-add-question-entity-18jan-aee4)
**Phase:** Complete
**Next Step:** None
**Status:** Complete

---

## Findings

### Finding 1: Issue types defined in internal/types/types.go

**Evidence:** All issue types (bug, feature, task, epic, etc.) are defined as IssueType constants with an IsValid() method.

**Source:** internal/types/types.go:391-414

**Significance:** New types need to be added to both the constants and the IsValid() switch statement.

---

### Finding 2: Status values also defined in types.go

**Evidence:** Statuses like open, in_progress, blocked, closed are defined with their own IsValid() method.

**Source:** internal/types/types.go:349-371

**Significance:** New question-specific statuses (investigating, answered) need to be added here.

---

### Finding 3: Type-specific validation in create.go

**Evidence:** The create command has type-specific validation for bugs (require repro) and epics (require understanding). Similar validation needed for questions.

**Source:** cmd/bd/create.go:143-205

**Significance:** Questions should reject flags that don't apply: assignee, estimate, repro, understanding.

---

## Synthesis

**Key Insights:**

1. **Minimal schema changes** - Adding a new type only requires updating types.go constants and IsValid() methods

2. **Type-specific validation** - The create command already has patterns for type-specific field validation

3. **New statuses integrate cleanly** - Adding StatusInvestigating and StatusAnswered follows the same pattern as other statuses

**Answer to Investigation Question:**

Adding a question type requires:
1. Add TypeQuestion constant to internal/types/types.go
2. Update IssueType.IsValid() to include TypeQuestion
3. Add StatusInvestigating and StatusAnswered constants
4. Update Status.IsValid() to include new statuses
5. Add validation in cmd/bd/create.go to reject inapplicable flags
6. Update help text in create.go

---

## Structured Uncertainty

**What's tested:**

- ✅ `bd create --type question --title "Test"` creates question successfully
- ✅ `bd list --type question` filters to questions only
- ✅ `bd show` displays question details correctly
- ✅ `bd update --status investigating` works
- ✅ `bd update --status answered` works
- ✅ Validation rejects --assignee for questions
- ✅ Validation rejects --estimate for questions
- ✅ All existing unit tests pass

**What's untested:**

- ⚠️ RPC/daemon mode (tested direct mode only)
- ⚠️ Import/export with questions
- ⚠️ Questions in multi-repo scenarios

**What would change this:**

- If RPC handler has separate type validation, would need updates there
- If questions need special handling in show output, would need view layer changes

---

## Implementation Recommendations

**Purpose:** Implementation complete - summarizing what was done.

### Implemented Approach ⭐

**Minimal schema extension** - Add question type and statuses with validation

**Changes made:**
- internal/types/types.go: Added TypeQuestion, StatusInvestigating, StatusAnswered
- cmd/bd/create.go: Added validation for question type, updated help text

**Trade-offs accepted:**
- No special view formatting for questions (uses standard issue display)
- No automated status transitions

---

## References

**Files Examined:**
- internal/types/types.go - Issue and Status type definitions
- cmd/bd/create.go - Create command and validation
- internal/validation/bead.go - Validation helpers

**Commands Run:**
```bash
# Build
go build ./cmd/bd/

# Test creation
bd create --type question --title 'Should we add caching?'

# Test listing
bd list --type question

# Test status changes
bd update test-c97 --status investigating
bd update test-c97 --status answered

# Run unit tests
go test ./internal/types/...
```

---

## Investigation History

**2026-01-18 12:00:** Investigation started
- Initial question: Add 'question' as first-class beads entity type
- Context: Part of epic orch-go-5j2hx, task 1 of 4

**2026-01-18 12:15:** Investigation completed
- Status: Complete
- Key outcome: Question type added with investigating/answered statuses and field validation
