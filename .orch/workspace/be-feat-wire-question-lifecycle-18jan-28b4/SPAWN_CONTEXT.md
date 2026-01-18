TASK: Wire question lifecycle in beads.

## Task
Implement question status transitions and lifecycle integration.

## Lifecycle
Open → Investigating → Answered → Closed

- Open: Question asked, not yet investigated
- Investigating: Investigation spawned to answer question
- Answered: Understanding reached (may auto-close investigation)
- Closed: Question resolved

## CLI Commands to Implement/Verify
- bd update <id> --status=investigating - transition to investigating
- bd update <id> --status=answered - transition to answered
- bd close <id> --reason='Answered: ...' - close with answer

## Implementation
1. Verify status transitions work for questions (build on commit 2dc8f7dc)
2. Consider: Auto-link investigations to questions via 'Answers: <question-id>' field
3. Ensure bd list shows correct lifecycle state
4. Add validation: questions can only use open/investigating/answered/closed statuses

## Success Criteria
- Question can transition: open → investigating → answered → closed
- bd update <question-id> --status=investigating works
- bd close <question-id> --reason='Answered: Yes because...' works
- Invalid status transitions are rejected with helpful error

## Context
Task 1 completed: commit 2dc8f7dc added question type and investigating/answered statuses
Design: .kb/decisions/2026-01-18-questions-as-first-class-entities.md

Part of epic orch-go-5j2hx. This is task 2 of 4.


SPAWN TIER: full

📚 FULL TIER: This spawn requires SYNTHESIS.md for knowledge externalization.
   Document your findings, decisions, and learnings in SYNTHESIS.md before completing.



## PRIOR KNOWLEDGE (from kb context)

**Query:** "wire question lifecycle"

### Constraints (MUST respect)
- orch-go agent state exists in four layers (OpenCode memory, OpenCode disk, registry, tmux)
  - Reason: Each layer has independent lifecycle - cleanup must touch all layers or ghosts accumulate
- Epics with parallel component work must include a final integration child issue
  - Reason: Swarm agents build components in parallel but nothing wires them together. Without explicit integration issue, manual intervention needed to create runnable feature. Learned from pw-4znt where 8 components built but no route existed.

### Prior Decisions
- Registry respawn workflow uses slot reuse pattern
  - Reason: Preserves single-entry-per-ID invariant while enabling abandon→respawn lifecycle
- Iteration 8 tmux fallback verification successful
  - Reason: After attach mode implementation (7ca8438), all three fallback mechanisms (tail, question, status) continue to work correctly with no regressions
- Three-tier temporal model (ephemeral/persistent/operational) organizes artifact placement
  - Reason: Artifacts live where their lifecycle dictates - session-bound to workspace, project-lifetime to kb, work-in-progress to beads
- Reflection value comes from orchestrator review + follow-up, not execution-time process changes
  - Reason: Evidence: post-synthesis reflection with Dylan created orch-go-ws4z epic (6 children) from captured questions
- Post-registry lifecycle uses 4 state sources: OpenCode sessions, tmux windows, beads issues, workspaces
  - Reason: Registry removed due to false positive completion detection; derived lookups replace central state
- Tmux session existence is correct abstraction for 'orch servers' status - separate infrastructure from project servers
  - Reason: orch serve is persistent monitoring infrastructure, not ephemeral project server. Tmux lifecycle matches dev server lifecycle. Mixing them creates semantic confusion and false status reporting.
- kb context uses keyword matching, not semantic understanding - 'how would the system recommend...' questions require orchestrator synthesis
  - Reason: Tested kb context with various query formats: keyword queries (swarm, dashboard) returned results, but semantic queries (swarm map sorting, how should dashboard present agents) returned nothing. The pattern reveals desire for semantic query answering that would require LLM-based RAG.
- SHOULD-HOW-EXECUTE sequence for strategic questions
  - Reason: Epic orch-go-erdw was created from 'how do we' without testing premise. Architect found premise wrong. Wasted work avoided if 'should we' asked first.
- Two-model experiments worth running selectively
  - Reason: Design critique worked (kb ask session - Gemini caught --force principle violation). Other candidates: adversarial review (security), synthesis from angles (research), devil's advocate (attachment). Run when real question fits pattern, don't pre-create.
- Questions as First-Class Entities
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/decisions/2026-01-18-questions-as-first-class-entities.md

### Models (synthesized understanding)
- Agent Lifecycle State Model
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/models/agent-lifecycle-state-model.md
- Agent Completion Lifecycle
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/models/completion-lifecycle.md
- Orchestrator Session Lifecycle
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/models/orchestrator-session-lifecycle.md
- Workspace Lifecycle & Hierarchy
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/models/workspace-lifecycle-model.md
- OpenCode Session Lifecycle
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/models/opencode-session-lifecycle.md
- Phase 3 Review: Model Pattern Analysis (N=5)
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/models/PHASE3_REVIEW.md
- Phase 4 Review: Model Pattern at N=11
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/models/PHASE4_REVIEW.md
- Model Access and Spawn Paths
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/models/model-access-spawn-paths.md
- Models
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/models/README.md
- Completion Verification Architecture
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/models/completion-verification.md

### Guides (procedural knowledge)
- Understanding Artifact Lifecycle
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/guides/understanding-artifact-lifecycle.md
- Agent Lifecycle Guide
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/guides/agent-lifecycle.md
- Workspace Lifecycle Guide
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/guides/workspace-lifecycle.md
- Spawned Orchestrator Pattern
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/guides/spawned-orchestrator-pattern.md
- OpenCode Integration Guide
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/guides/opencode.md
- orch CLI Reference
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/guides/cli.md
- Code Extraction Patterns
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/guides/code-extraction-patterns.md
- Orchestrator Session Management
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/guides/orchestrator-session-management.md
- Worker Patterns Guide
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/guides/worker-patterns.md
- Decision Authority Guide
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/guides/decision-authority.md

### Related Investigations
- Design Questions as First-Class Entities
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/investigations/2026-01-18-inv-design-questions-first-class-entities.md
- Llm Detect Premise Skipping Question
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/investigations/2026-01-17-inv-llm-detect-premise-skipping-question.md
- Orchestrator Completion Lifecycle Design
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/investigations/2025-12-25-design-orchestrator-completion-lifecycle-two.md
- Update Core Skills to Use OpenCode Question Tool
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/investigations/2026-01-07-inv-update-core-skills-opencode-ask.md
- orch-go Add Question Command
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/investigations/2025-12-20-inv-orch-add-question-command.md
- Audit Orchestration Lifecycle Post-Registry Removal
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/investigations/2025-12-22-inv-audit-orchestration-lifecycle-post-registry.md
- Design Question Should Orch Servers
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/investigations/2025-12-23-inv-design-question-should-orch-servers.md
- Orchestrator Session Lifecycle Without Beads Tracking
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/investigations/2026-01-05-inv-design-orchestrator-session-lifecycle-without.md
- Update Orchestrator Session Lifecycle Model
  - See: /Users/dylanconlin/Documents/personal/orch-go/.kb/investigations/archived/2026-01-14-inv-update-orchestrator-session-lifecycle-model.md

**IMPORTANT:** The above context represents existing knowledge and decisions. Do not contradict constraints. Reference models and guides for established patterns. Reference investigations for prior findings.





📋 AD-HOC SPAWN (--no-track):
This is an ad-hoc spawn without beads issue tracking.
Progress tracking via bd comment is NOT available.

🚨 SESSION COMPLETE PROTOCOL:
After your final commit, BEFORE typing anything else:

⛔ **NEVER run `git push`** - Workers commit locally only.
   - Your orchestrator will handle pushing to remote after review
   - Running `git push` can trigger deploys that disrupt production systems
   - Worker rule: Commit your work, call `/exit`. Don't push.


1. Ensure SYNTHESIS.md is created and committed in your workspace.
2. Run: `/exit` to close the agent session



CONTEXT: [See task description]

PROJECT_DIR: /Users/dylanconlin/Documents/personal/beads

SESSION SCOPE: Medium (estimated [1-2h / 2-4h / 4-6h+])
- Default estimation
- Recommend checkpoint after Phase 1 if session exceeds 2 hours


AUTHORITY:
**You have authority to decide:**
- Implementation details (how to structure code, naming, file organization)
- Testing strategies (which tests to write, test frameworks to use)
- Refactoring within scope (improving code quality without changing behavior)
- Tool/library selection within established patterns (using tools already in project)
- Documentation structure and wording

**You must escalate to orchestrator when:**
- Architectural decisions needed (changing system structure, adding new patterns)
- Scope boundaries unclear (unsure if something is IN vs OUT scope)
- Requirements ambiguous (multiple valid interpretations exist)
- Blocked by external dependencies (missing access, broken tools, unclear context)
- Major trade-offs discovered (performance vs maintainability, security vs usability)
- Task estimation significantly wrong (2h task is actually 8h)

**When uncertain:** Err on side of escalation. Document question in workspace, set Status: QUESTION, and wait for orchestrator response. Better to ask than guess wrong.

**Full criteria:** See `.kb/guides/decision-authority.md` for the complete decision tree and examples.

**Surface Before Circumvent:**
Before working around ANY constraint (technical, architectural, or process):
1. Document it in your investigation file: "CONSTRAINT: [what constraint] - [why considering workaround]"
2. Include the constraint and your reasoning in SYNTHESIS.md
3. The accountability is a feature, not a cost

This applies to:
- System constraints discovered during work (e.g., API limits, tool limitations)
- Architectural patterns that seem inconvenient for your task
- Process requirements that feel like overhead
- Prior decisions (from `kb context`) that conflict with your approach

**Why:** Working around constraints without surfacing them:
- Prevents the system from learning about recurring friction
- Bypasses stakeholders who should know about the limitation
- Creates hidden technical debt

DELIVERABLES (REQUIRED):
1. **FIRST:** Verify project location: pwd (must be /Users/dylanconlin/Documents/personal/beads)
2. **SET UP investigation file:** Run `kb create investigation wire-question-lifecycle-beads-task` to create from template
   - This creates: `.kb/investigations/simple/YYYY-MM-DD-wire-question-lifecycle-beads-task.md`
   - This file is your coordination artifact (replaces WORKSPACE.md)
   - If command fails, report to orchestrator immediately

3. **UPDATE investigation file** as you work:
   - Add TLDR at top (1-2 sentence summary of question and finding)
   - Fill sections: What I tried → What I observed → Test performed
   - Only fill Conclusion if you actually tested (this is the key discipline)
4. Update Status: field when done (Active → Complete)
5. [Task-specific deliverables]

6. **CREATE SYNTHESIS.md:** Before completing, create `SYNTHESIS.md` in your workspace: /Users/dylanconlin/Documents/personal/beads/.orch/workspace/be-feat-wire-question-lifecycle-18jan-28b4/SYNTHESIS.md
   - Use the template from: /Users/dylanconlin/Documents/personal/beads/.orch/templates/SYNTHESIS.md
   - This is CRITICAL for the orchestrator to review your work.


STATUS UPDATES:
Update Status: field in your investigation file:
- Status: Active (while working)
- Status: Complete (when done and committed) → then call /exit to close agent session

Signal orchestrator when blocked:
- Add '**Status:** BLOCKED - [reason]' to investigation file
- Add '**Status:** QUESTION - [question]' when needing input



## SKILL GUIDANCE (feature-impl)

**IMPORTANT:** You have been spawned WITH this skill context already loaded.
You do NOT need to invoke this skill using the Skill tool.
Simply follow the guidance provided below.

---

---
name: feature-impl
skill-type: procedure
description: Unified feature implementation with configurable phases (investigation, clarifying-questions, design, implementation, validation, integration). Replaces test-driven-development, surgical-change, and feature-coordination skills. Use for any feature work with phases/mode/validation configured by orchestrator.
---

<!-- AUTO-GENERATED by skillc -->
<!-- Checksum: d7a5e8be0268 -->
<!-- Source: /Users/dylanconlin/orch-knowledge/skills/src/worker/feature-impl/.skillc -->
<!-- Deployed to: /Users/dylanconlin/.claude/skills/worker/feature-impl/SKILL.md -->
<!-- To modify: edit files in /Users/dylanconlin/orch-knowledge/skills/src/worker/feature-impl/.skillc, then run: skillc deploy -->
<!-- Last compiled: 2026-01-17 01:46:54 -->

## Summary

name: feature-impl
skill-type: procedure
description: Unified feature implementation with configurable phases (investigation, clarifying-questions, design, implementation, validation, integration). Replaces test-driven-development, surgical-change, and feature-coordination skills. Use for any feature work with phases/mode/validation configured by orchestrator.

---

---
name: feature-impl
skill-type: procedure
description: Unified feature implementation with configurable phases (investigation, clarifying-questions, design, implementation, validation, integration). Replaces test-driven-development, surgical-change, and feature-coordination skills. Use for any feature work with phases/mode/validation configured by orchestrator.
---

<!-- AUTO-GENERATED by skillc -->
<!-- Checksum: 047ddb2689b3 -->
<!-- Source: .skillc -->
<!-- To modify: edit files in .skillc, then run: skillc build -->
<!-- Last compiled: 2026-01-07 14:41:54 -->

## Summary

**For orchestrators:** Spawn via `orch spawn feature-impl "task" --phases "..." --mode ... --validation ...`

---

# Feature Implementation (Unified Framework)

**For orchestrators:** Spawn via `orch spawn feature-impl "task" --phases "..." --mode ... --validation ...`

**For workers:** You've been spawned to implement a feature using a phased approach with specific configuration.

---

## Your Configuration

**Read from SPAWN_CONTEXT.md** to understand your configuration:

- **Phases:** Which phases you'll proceed through (e.g., `investigation,clarifying-questions,design,implementation,validation`)
- **Current Phase:** Determined by your progress (start with first configured phase)
- **Implementation Mode:** `tdd` or `direct` (only relevant if implementation phase included)
- **Validation Level:** `none`, `tests`, `smoke-test`, or `multi-phase` (only relevant if validation phase included)

**Example configuration:**
```
Phases: design, implementation, validation
Mode: tdd
Validation: smoke-test

---




CONTEXT AVAILABLE:
- Global: ~/.claude/CLAUDE.md
- Project: /Users/dylanconlin/Documents/personal/beads/CLAUDE.md

🚨 FINAL STEP - SESSION COMPLETE PROTOCOL:
After your final commit, BEFORE doing anything else:

⛔ **NEVER run `git push`** - Workers commit locally only.
   - Your orchestrator will handle pushing to remote after review
   - Running `git push` can trigger deploys that disrupt production systems
   - Worker rule: Commit your work, call `/exit`. Don't push.



1. Ensure SYNTHESIS.md is created and committed in your workspace.
2. `/exit`


⚠️ Your work is NOT complete until you run these commands.
