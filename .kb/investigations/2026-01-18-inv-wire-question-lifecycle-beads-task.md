<!--
D.E.K.N. Summary - 30-second handoff for fresh Claude
Fill this at the END of your investigation, before marking Complete.
-->

## Summary (D.E.K.N.)

**Delta:** Question lifecycle status validation successfully wired into beads CLI.

**Evidence:** Manual testing confirms: (1) questions can only use open/investigating/answered/closed statuses, (2) invalid statuses rejected with helpful error, (3) assignee/estimate updates rejected for questions, (4) bd close --reason works correctly.

**Knowledge:** CLI-level validation pattern (in update.go) is consistent with create.go's question validation. Validation function added to validation package for reuse.

**Next:** Close investigation - all success criteria met.

**Promote to Decision:** recommend-no - Tactical implementation following existing patterns, not architectural.

<!--
Example D.E.K.N.:
Delta: Test-running guidance is missing from spawn prompts and CLAUDE.md.
Evidence: Searched 5 agent sessions - none ran tests; guidance exists in separate docs but isn't loaded.
Knowledge: Agents follow documentation literally; guidance must be in loaded context to be followed.
Next: Add test-running instruction to SPAWN_CONTEXT.md template.
Promote to Decision: recommend-no (tactical fix, not architectural)

Guidelines:
- Keep each line to ONE sentence
- Delta answers "What did we find?"
- Evidence answers "How do we know?"
- Knowledge answers "What does this mean?"
- Next answers "What should happen now?"
- Promote to Decision: flag for orchestrator/human - recommend-yes if this establishes a pattern, constraint, or architectural choice worth preserving
- Enable 30-second understanding for fresh Claude
-->

---

# Investigation: Wire Question Lifecycle Beads Task

**Question:** How do we validate and enforce question-specific status transitions in beads?

**Started:** 2026-01-18
**Updated:** 2026-01-18
**Owner:** Worker Agent
**Phase:** Complete
**Next Step:** None
**Status:** Complete

<!-- Lineage (fill only when applicable) -->
**Patches-Decision:** [Path to decision document this investigation patches/extends, if applicable - enables review triggers]
**Extracted-From:** [Project/path of original artifact, if this was extracted from another project]
**Supersedes:** [Path to artifact this replaces, if applicable]
**Superseded-By:** [Path to artifact that replaced this, if applicable]

---

## Findings

### Finding 1: CLI-level validation is the established pattern

**Evidence:** create.go validates question-specific flags at lines 207-228, rejecting assignee, estimate, repro, and understanding flags for questions. This pattern is consistent throughout the CLI.

**Source:** cmd/bd/create.go:207-228

**Significance:** Following this pattern ensures consistency - validation happens at the CLI before reaching the store layer.

---

### Finding 2: Question statuses must be restricted to lifecycle states

**Evidence:** Questions follow lifecycle: open → investigating → answered → closed. General statuses like in_progress, blocked, deferred are not appropriate for questions.

**Source:** internal/types/types.go:361-363 (StatusInvestigating, StatusAnswered constants)

**Significance:** Restricting statuses prevents questions from being treated like tasks while maintaining the question-specific workflow.

---

### Finding 3: Both daemon and direct modes need validation

**Evidence:** update.go has two code paths: daemon mode (RPC) and direct mode. Both need validation before applying updates.

**Source:** cmd/bd/update.go:142-222 (daemon), cmd/bd/update.go:224-337 (direct)

**Significance:** Validation added to both paths ensures consistent behavior regardless of how bd is invoked.

---

## Synthesis

**Key Insights:**

1. **[Insight title]** - [Explanation of the insight, connecting multiple findings]

2. **[Insight title]** - [Explanation of the insight, connecting multiple findings]

3. **[Insight title]** - [Explanation of the insight, connecting multiple findings]

**Answer to Investigation Question:**

[Clear, direct answer to the question posed at the top of this investigation. Reference specific findings that support this answer. Acknowledge any limitations or gaps.]

---

## Structured Uncertainty

**What's tested:**

- ✅ Status validation rejects in_progress for questions (verified: bd update test-7ia --status=in_progress returns error)
- ✅ Valid statuses work (verified: open → investigating → answered → closed transitions succeed)
- ✅ Assignee update rejected for questions (verified: bd update test-ode --assignee=someone returns error)
- ✅ Estimate update rejected for questions (verified: bd update test-ode --estimate=60 returns error)
- ✅ bd close --reason works (verified: close_reason displayed in show output)
- ✅ Unit tests pass for validation functions (verified: go test ./internal/validation/...)

**What's untested:**

- ⚠️ RPC server-side validation (only CLI-level validation added, server accepts any status)
- ⚠️ Multi-issue update with mixed types (e.g., bd update task-id question-id --status=investigating)

**What would change this:**

- If server-side validation is needed for API clients bypassing CLI
- If status transitions need to be ordered (e.g., can't go from open directly to closed)

---

## Implementation Recommendations

**Purpose:** Bridge from investigation findings to actionable implementation using directive guidance pattern (strong recommendations + visible reasoning).

### Recommended Approach ⭐

**[Approach Name]** - [One sentence stating the recommended implementation]

**Why this approach:**
- [Key benefit 1 based on findings]
- [Key benefit 2 based on findings]
- [How this directly addresses investigation findings]

**Trade-offs accepted:**
- [What we're giving up or deferring]
- [Why that's acceptable given findings]

**Implementation sequence:**
1. [First step - why it's foundational]
2. [Second step - why it comes next]
3. [Third step - builds on previous]

### Alternative Approaches Considered

**Option B: [Alternative approach]**
- **Pros:** [Benefits]
- **Cons:** [Why not recommended - reference findings]
- **When to use instead:** [Conditions where this might be better]

**Option C: [Alternative approach]**
- **Pros:** [Benefits]
- **Cons:** [Why not recommended - reference findings]
- **When to use instead:** [Conditions where this might be better]

**Rationale for recommendation:** [Brief synthesis of why Option A beats alternatives given investigation findings]

---

### Implementation Details

**What to implement first:**
- [Highest priority change based on findings]
- [Quick wins or foundational work]
- [Dependencies that need to be addressed early]

**Things to watch out for:**
- ⚠️ [Edge cases or gotchas discovered during investigation]
- ⚠️ [Areas of uncertainty that need validation during implementation]
- ⚠️ [Performance, security, or compatibility concerns to address]

**Areas needing further investigation:**
- [Questions that arose but weren't in scope]
- [Uncertainty areas that might affect implementation]
- [Optional deep-dives that could improve the solution]

**Success criteria:**
- ✅ [How to know the implementation solved the investigated problem]
- ✅ [What to test or validate]
- ✅ [Metrics or observability to add]

---

## References

**Files Examined:**
- [File path] - [What you looked at and why]
- [File path] - [What you looked at and why]

**Commands Run:**
```bash
# [Command description]
[command]

# [Command description]
[command]
```

**External Documentation:**
- [Link or reference] - [What it is and relevance]

**Related Artifacts:**
- **Decision:** [Path to related decision document] - [How it relates]
- **Investigation:** [Path to related investigation] - [How it relates]
- **Workspace:** [Path to related workspace] - [How it relates]

---

## Investigation History

**[YYYY-MM-DD HH:MM]:** Investigation started
- Initial question: [Original question as posed]
- Context: [Why this investigation was initiated]

**[YYYY-MM-DD HH:MM]:** [Milestone or significant finding]
- [Description of what happened or was discovered]

**[YYYY-MM-DD HH:MM]:** Investigation completed
- Status: [Complete/Paused with reason]
- Key outcome: [One sentence summary of result]
