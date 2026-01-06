# VerMAS: Verifiable Multi-Agent System

> Design Specification v0.1

## Core Principles

1. **Zero User Workflow Change**: Users interact with Mayor only. Inspector ecosystem is invisible.
2. **First-Class Citizens**: Inspector roles have ALL capabilities other agents have (mail, hooks, beads, molecules).
3. **Deterministic Hooks**: Every step has hooks for observability and extensibility.
4. **Mayor as Human Interface**: All questions from Inspector route through Mayor.
5. **Same Infrastructure**: No special cases - Inspector uses existing Gas Town primitives.

## Overview

VerMAS extends Gas Town with a verification layer that ensures work quality through independent inspection before merge. The core principle: **define acceptance criteria before work begins, verify against those criteria before work merges.**

This is Test-Driven Development (TDD) applied to multi-agent workflows.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              VerMAS Architecture                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        PLANNING PHASE                                │   │
│  │                                                                      │   │
│  │   ┌────────┐         ┌──────────┐         ┌─────────────────────┐  │   │
│  │   │ Mayor  │─────────│ Auditor  │─────────│    Test Spec        │  │   │
│  │   │        │  task   │          │ creates │  • Acceptance tests │  │   │
│  │   │        │         │          │         │  • Verification plan│  │   │
│  │   └────────┘         └──────────┘         │  • Success criteria │  │   │
│  │                                           └──────────┬──────────┘  │   │
│  └───────────────────────────────────────────────────────┼──────────────┘   │
│                                                          │                  │
│  ┌───────────────────────────────────────────────────────┼──────────────┐   │
│  │                        EXECUTION PHASE                │               │   │
│  │                                                       ▼               │   │
│  │   ┌────────┐         ┌──────────┐         ┌─────────────────────┐   │   │
│  │   │ Mayor  │─────────│ Polecat  │─────────│   Implementation    │   │   │
│  │   │        │  sling  │          │ creates │  • Code changes     │   │   │
│  │   │        │         │          │         │  • Tests (optional) │   │   │
│  │   └────────┘         └──────────┘         │  • Documentation    │   │   │
│  │                                           └──────────┬──────────┘   │   │
│  └───────────────────────────────────────────────────────┼──────────────┘   │
│                                                          │                  │
│  ┌───────────────────────────────────────────────────────┼──────────────┐   │
│  │                     VERIFICATION PHASE                │               │   │
│  │                                                       ▼               │   │
│  │   ┌──────────┐       ┌─────────────────────────────────────────┐    │   │
│  │   │ Refinery │──────▶│              INSPECTOR                  │    │   │
│  │   │  (gate)  │       │                                         │    │   │
│  │   └──────────┘       │  ┌─────────┐ ┌─────────┐ ┌─────────┐   │    │   │
│  │                      │  │Advocate │ │ Critic  │ │ Judge   │   │    │   │
│  │                      │  │(defense)│ │(attack) │ │(verdict)│   │    │   │
│  │                      │  └────┬────┘ └────┬────┘ └────┬────┘   │    │   │
│  │                      │       │           │           │         │    │   │
│  │                      │       └───────────┴───────────┘         │    │   │
│  │                      │                   │                      │    │   │
│  │                      │           ┌───────▼───────┐              │    │   │
│  │                      │           │   Verdict     │              │    │   │
│  │                      │           │ PASS/FAIL/    │              │    │   │
│  │                      │           │ NEEDS_HUMAN   │              │    │   │
│  │                      │           └───────────────┘              │    │   │
│  │                      └─────────────────────────────────────────┘    │   │
│  │                                          │                           │   │
│  │                         ┌────────────────┴────────────────┐         │   │
│  │                         ▼                                 ▼         │   │
│  │                   ┌──────────┐                     ┌──────────┐    │   │
│  │                   │  MERGE   │                     │  REJECT  │    │   │
│  │                   └──────────┘                     └────┬─────┘    │   │
│  │                                                         │          │   │
│  └─────────────────────────────────────────────────────────┼──────────┘   │
│                                                            │               │
│  ┌─────────────────────────────────────────────────────────┼──────────┐   │
│  │                        REWORK PHASE                     │          │   │
│  │                                                         ▼          │   │
│  │   ┌────────┐         ┌──────────┐         ┌─────────────────────┐ │   │
│  │   │ Mayor  │◀────────│ Feedback │◀────────│   Inspector Report  │ │   │
│  │   │        │         │          │         │  • Failed criteria  │ │   │
│  │   │        │─────────│ Polecat  │         │  • Suggested fixes  │ │   │
│  │   │        │ re-sling│  (fix)   │         │  • Context          │ │   │
│  │   └────────┘         └──────────┘         └─────────────────────┘ │   │
│  │                                                                    │   │
│  └────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Inspector Ecosystem

The Inspector is not a single agent but an ecosystem of specialized roles, each wearing a different "hat" to ensure thorough, unbiased verification.

### Role Definitions

| Role | Hat | Purpose | Mindset |
|------|-----|---------|---------|
| **Advocate** | Defense Attorney | Argue FOR the code | "This code is correct because..." |
| **Critic** | Prosecutor | Argue AGAINST the code | "This code fails because..." |
| **Judge** | Impartial Arbiter | Weigh arguments, deliver verdict | "Based on evidence, the verdict is..." |
| **Verifier** | Lab Technician | Run objective tests | "Test X returned Y" |
| **Auditor** | Compliance Officer | Check against spec | "Requirement A is/isn't met" |

### Why Multiple Roles?

1. **Adversarial Process**: Advocate vs Critic prevents single-model bias
2. **Separation of Concerns**: Each role has one job, does it well
3. **Audit Trail**: Clear record of who said what
4. **Different Models**: Each role can use a different LLM for true independence

### Inspector Workflow

```
                    ┌─────────────────────────────────────┐
                    │         INSPECTOR SESSION           │
                    └─────────────────┬───────────────────┘
                                      │
                    ┌─────────────────▼───────────────────┐
                    │           1. VERIFICATION           │
                    │                                     │
                    │  Verifier runs:                     │
                    │  • Unit tests (go test ./...)      │
                    │  • Integration tests                │
                    │  • Linting (golangci-lint)         │
                    │  • Build check (go build)          │
                    │  • Custom test spec from Auditor   │
                    │                                     │
                    │  Output: Objective results          │
                    └─────────────────┬───────────────────┘
                                      │
                    ┌─────────────────▼───────────────────┐
                    │            2. AUDIT                 │
                    │                                     │
                    │  Auditor checks:                    │
                    │  • Requirements from original task │
                    │  • Acceptance criteria (test spec) │
                    │  • Code coverage thresholds        │
                    │  • Security requirements           │
                    │                                     │
                    │  Output: Compliance report          │
                    └─────────────────┬───────────────────┘
                                      │
          ┌───────────────────────────┼───────────────────────────┐
          │                           │                           │
          ▼                           │                           ▼
┌─────────────────────┐               │               ┌─────────────────────┐
│      ADVOCATE       │               │               │       CRITIC        │
│                     │               │               │                     │
│ Reviews:            │               │               │ Reviews:            │
│ • Verification logs │               │               │ • Verification logs │
│ • Audit report      │               │               │ • Audit report      │
│ • Code diff         │               │               │ • Code diff         │
│                     │               │               │                     │
│ Produces:           │               │               │ Produces:           │
│ • Defense brief     │               │               │ • Prosecution brief │
│ • Mitigations       │               │               │ • Concerns          │
│ • Strengths         │               │               │ • Weaknesses        │
└──────────┬──────────┘               │               └──────────┬──────────┘
           │                          │                          │
           └──────────────────────────┼──────────────────────────┘
                                      │
                    ┌─────────────────▼───────────────────┐
                    │             3. JUDGMENT             │
                    │                                     │
                    │  Judge reviews:                     │
                    │  • Verification results             │
                    │  • Audit compliance report          │
                    │  • Advocate's defense brief         │
                    │  • Critic's prosecution brief       │
                    │                                     │
                    │  Delivers:                          │
                    │  • Verdict (PASS/FAIL/NEEDS_HUMAN) │
                    │  • Confidence score (0.0-1.0)      │
                    │  • Reasoning                        │
                    │  • Required fixes (if FAIL)        │
                    └─────────────────┬───────────────────┘
                                      │
                                      ▼
                              ┌───────────────┐
                              │    VERDICT    │
                              └───────────────┘
```

## Data Model

### Test Spec Bead

Created by Auditor before work begins. Attached to the work item.

```yaml
id: spec-abc123
type: test-spec
parent: gt-xyz789  # The work item this spec is for

acceptance_criteria:
  - id: ac-1
    description: "Function handles empty input gracefully"
    verification: "unit-test"
    test_command: "go test -run TestEmptyInput ./..."

  - id: ac-2
    description: "API returns 400 for invalid requests"
    verification: "integration-test"
    test_command: "make test-api"

  - id: ac-3
    description: "No new security vulnerabilities"
    verification: "security-scan"
    test_command: "gosec ./..."

success_threshold:
  required_pass: all  # or percentage like "80%"

verification_strategy:
  approach: "adversarial"  # advocate + critic
  models:
    advocate: claude
    critic: codex
    judge: claude  # different instance
```

### Verification Result Bead

Created by Inspector after verification completes.

```yaml
id: vr-def456
type: verification-result
parent: gt-xyz789  # The work item
spec: spec-abc123  # The test spec used

verdict: PASS | FAIL | NEEDS_HUMAN
confidence: 0.85

verification_results:
  - criteria_id: ac-1
    status: pass
    output: "PASS: TestEmptyInput (0.003s)"

  - criteria_id: ac-2
    status: fail
    output: "Expected 400, got 500"

  - criteria_id: ac-3
    status: pass
    output: "No issues found"

advocate_brief: |
  The implementation correctly handles the core requirements...

critic_brief: |
  Concern: Error handling in edge case X is incomplete...

judge_reasoning: |
  While the critic raises valid concerns about edge case X,
  the advocate demonstrates this is out of scope for the
  current task. The acceptance criteria are met.

required_fixes: []  # Empty if PASS

reviewed_by:
  verifier: "go-test-runner"
  auditor: "codex"
  advocate: "claude"
  critic: "codex"
  judge: "claude"

timestamp: "2026-01-06T15:30:00Z"
duration: "45s"
```

## Commands

### Planning Phase

```bash
# Auditor creates test spec for a work item
gt audit plan <bead-id>
# → Creates spec-xxx bead, attaches to work item

# View the test spec
gt audit show <bead-id>
# → Shows acceptance criteria and verification strategy

# Manually add acceptance criteria
gt audit add-criteria <bead-id> --description="..." --test="..."
```

### Execution Phase

```bash
# Normal sling - now includes test spec
gt sling <bead-id> <rig>
# → Polecat receives work item + attached test spec

# Polecat can check what they need to satisfy
gt audit criteria
# → Shows acceptance criteria for hooked work
```

### Verification Phase

```bash
# Trigger inspection (usually automatic at Refinery gate)
gt inspect <bead-id>
# → Runs full Inspector workflow

# Run just the verifier (objective tests)
gt inspect verify <bead-id>
# → Runs tests, linting, build

# Run just the audit (compliance check)
gt inspect audit <bead-id>
# → Checks against test spec

# Run adversarial review
gt inspect review <bead-id>
# → Advocate + Critic + Judge

# Check inspection status
gt inspect status <bead-id>
# → Shows current state and results
```

### Rework Phase

```bash
# After rejection, Mayor gets feedback
gt inspect feedback <bead-id>
# → Shows what failed, suggested fixes

# Re-sling with context
gt sling <bead-id> <rig> --rework
# → Polecat gets original work + failure context
```

## Runtime Configuration

Each Inspector role can use a different AI runtime:

```yaml
# config/runtimes.yaml

roles:
  # Standard Gas Town roles
  mayor: claude
  polecat: claude
  witness: claude
  refinery: claude

  # Inspector ecosystem - intentionally diverse
  auditor: codex      # Different model for test spec creation
  verifier: local     # No LLM, just runs commands
  advocate: claude    # Argues for the code
  critic: codex       # Argues against (different perspective)
  judge: claude       # Final decision (could be opus for complex cases)

# Fallback chains
inspector_fallback:
  - codex
  - opencode
  - claude  # Last resort: same model

# Verification is mandatory
verification:
  required: true
  scope: all  # Every MR must pass
```

## AI-Assisted Requirements & Criteria

Users shouldn't need to write detailed requirements or verification commands. AI agents assist on both sides.

### The Problem

**Without assistance:**
```
User: "Create fizzbuzz"
Mayor: "What are the requirements?"
User: "Uh... it should print fizzbuzz?"
Inspector: "What verification commands should I run?"
User: "I don't know bash that well..."
```

**With assistance:**
```
User: "Create fizzbuzz"
Designer: "Here are the requirements I suggest: [detailed list]"
User: "Looks good" or "Also add X"
Strategist: "Here are the verification criteria I suggest: [detailed list]"
User: "Looks good" or "Make that threshold higher"
```

### Two New Agents

| Agent | Reports To | Purpose |
|-------|------------|---------|
| **Designer** | Mayor | Elaborates vague requests into detailed requirements |
| **Strategist** | Inspector | Proposes verification criteria from requirements |

### Flow with AI Assistance

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     AI-ASSISTED REQUIREMENT FLOW                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  USER: "Create fizzbuzz"                                                    │
│         │                                                                   │
│         ▼                                                                   │
│  ┌─────────────┐      ┌─────────────┐                                      │
│  │   MAYOR     │ ───▶ │  DESIGNER   │                                      │
│  │             │      │             │                                      │
│  │  "Let me    │      │ Analyzes    │                                      │
│  │   get this  │      │ request,    │                                      │
│  │   designed" │      │ produces    │                                      │
│  │             │      │ requirements│                                      │
│  └─────────────┘      └──────┬──────┘                                      │
│                              │                                              │
│                              ▼                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ DESIGNER OUTPUT:                                                     │   │
│  │                                                                      │   │
│  │ Requirements for FizzBuzz:                                           │   │
│  │                                                                      │   │
│  │ Functional:                                                          │   │
│  │ • Print numbers 1 through 100, one per line                         │   │
│  │ • For multiples of 3, print "Fizz" instead of the number            │   │
│  │ • For multiples of 5, print "Buzz" instead of the number            │   │
│  │ • For multiples of both 3 and 5, print "FizzBuzz"                   │   │
│  │                                                                      │   │
│  │ Non-functional:                                                      │   │
│  │ • Clean, readable Python code                                        │   │
│  │ • Include docstrings                                                 │   │
│  │ • Pass linting with pylint >= 8.0                                   │   │
│  │                                                                      │   │
│  │ [Approve] [Modify] [Reject]                                         │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                              │                                              │
│                              ▼                                              │
│  USER: "Looks good" or "Also make it a reusable function"                  │
│         │                                                                   │
│         ▼                                                                   │
│  ┌─────────────┐                                                           │
│  │   MAYOR     │ Creates bead with full requirements                       │
│  │             │ gt-fb001 (BLOCKED, waiting for spec)                      │
│  └──────┬──────┘                                                           │
│         │                                                                   │
│         │ [MAIL: New spec needed]                                          │
│         ▼                                                                   │
│  ┌─────────────┐      ┌─────────────┐                                      │
│  │  INSPECTOR  │ ───▶ │ STRATEGIST  │                                      │
│  │             │      │             │                                      │
│  │  "Let me    │      │ Reads       │                                      │
│  │   figure    │      │ requirements│                                      │
│  │   out how   │      │ proposes    │                                      │
│  │   to verify"│      │ criteria    │                                      │
│  └─────────────┘      └──────┬──────┘                                      │
│                              │                                              │
│                              ▼                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ STRATEGIST OUTPUT:                                                   │   │
│  │                                                                      │   │
│  │ Proposed Verification Criteria for gt-fb001:                         │   │
│  │                                                                      │   │
│  │ [AC-1] Output line count                                            │   │
│  │        Check: Exactly 100 lines of output                           │   │
│  │        Command: python fizzbuzz.py | wc -l | grep -q '100'          │   │
│  │                                                                      │   │
│  │ [AC-2] FizzBuzz on line 15                                          │   │
│  │        Check: Line 15 (3×5) outputs "FizzBuzz"                      │   │
│  │        Command: python fizzbuzz.py | sed -n '15p' | grep -q 'FizzBuzz'│  │
│  │                                                                      │   │
│  │ [AC-3] Fizz on line 9                                               │   │
│  │        Check: Line 9 (3×3) outputs "Fizz"                           │   │
│  │        Command: python fizzbuzz.py | sed -n '9p' | grep -q '^Fizz$' │   │
│  │                                                                      │   │
│  │ [AC-4] Buzz on line 10                                              │   │
│  │        Check: Line 10 (2×5) outputs "Buzz"                          │   │
│  │        Command: python fizzbuzz.py | sed -n '10p' | grep -q '^Buzz$'│   │
│  │                                                                      │   │
│  │ [AC-5] Plain number on line 1                                       │   │
│  │        Check: Line 1 outputs just "1"                               │   │
│  │        Command: python fizzbuzz.py | head -1 | grep -q '^1$'        │   │
│  │                                                                      │   │
│  │ [AC-6] Code quality                                                 │   │
│  │        Check: Pylint score >= 8.0                                   │   │
│  │        Command: pylint fizzbuzz.py --fail-under=8.0                 │   │
│  │                                                                      │   │
│  │ [Approve All] [Modify] [Add More] [Reject]                          │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                              │                                              │
│                              ▼                                              │
│  USER: "Looks good" or "Also check line 30 is FizzBuzz"                    │
│         │                                                                   │
│         ▼                                                                   │
│  ┌─────────────┐                                                           │
│  │  INSPECTOR  │ Approves spec with criteria                               │
│  │             │ gt-fb001 now READY                                        │
│  └─────────────┘                                                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Designer Agent

**Role**: Turn vague user requests into detailed requirements.

**Triggers**: When user gives Mayor a task without full specification.

**Process**:
1. Analyze user's request
2. Research similar implementations (if needed)
3. Produce structured requirements:
   - Functional requirements (what it does)
   - Non-functional requirements (how well it does it)
   - Constraints (what it shouldn't do)
   - Examples (concrete expected behavior)
4. Present to user for approval

**Prompt**:
```
You are the Designer. Your job is to turn vague requests into detailed requirements.

User request: "{request}"

Produce requirements covering:
1. Functional - What exactly should this do?
2. Non-functional - Performance, quality, style requirements
3. Constraints - What should it NOT do?
4. Examples - Concrete input/output examples

Be specific enough that a developer could implement without asking questions.
Be concise - don't over-engineer simple requests.
```

### Strategist Agent

**Role**: Turn requirements into testable verification criteria.

**Triggers**: When a new spec is created and needs criteria.

**Process**:
1. Read the requirements from the work item
2. Identify key behaviors to verify
3. Design verification commands that:
   - Are deterministic (same input → same output)
   - Are automatable (can run without human)
   - Cover critical functionality
   - Are specific (clear pass/fail)
4. Present to user for approval

**Prompt**:
```
You are the Strategist. Your job is to design verification criteria.

Requirements:
{requirements}

For each key behavior, create a verification criterion:
- ID: AC-N
- Description: What we're checking (human readable)
- Command: Exact bash command to verify (must exit 0 for pass, non-0 for fail)

Focus on:
- Correctness (does it work?)
- Edge cases (does it handle boundaries?)
- Quality (is the code clean?)

Don't over-test. Cover the critical paths.
```

### User Interaction Model

**User's job is now just approval/adjustment:**

| Stage | AI Does | User Does |
|-------|---------|-----------|
| Requirements | Designer proposes full spec | "Approve" / "Also add X" / "Remove Y" |
| Criteria | Strategist proposes tests | "Approve" / "Make stricter" / "Add check for Z" |

**Example conversation:**

```
User: Create fizzbuzz

Mayor: Let me get Designer to elaborate on that...

Designer: Here's what I understand you want:
  - Print 1-100 with Fizz/Buzz/FizzBuzz substitutions
  - Clean Python code with docstrings
  - Should pass pylint

  Does this look right?

User: Yes, and make it a reusable function

Mayor: Got it. Creating work item with those requirements...
       Created gt-fb001 (blocked until Inspector approves criteria)

Inspector: Strategist is analyzing the requirements...

Strategist: Here are 6 verification criteria I propose:
  [AC-1] 100 lines of output
  [AC-2] Line 15 = "FizzBuzz"
  [AC-3] Line 9 = "Fizz"
  [AC-4] Line 10 = "Buzz"
  [AC-5] Line 1 = "1"
  [AC-6] Pylint >= 8.0

  Should I add anything else?

User: Also verify line 30 is FizzBuzz

Strategist: Added AC-7 for line 30. Approving spec...
            gt-fb001 is now READY for implementation.

Mayor: Great, slinging to polecat-1...
```

### Updated Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              VERMAS AGENTS                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                              USER (Overseer)                                │
│                             /              \                                │
│                            /                \                               │
│            ┌──────────────▼───┐          ┌───▼──────────────┐              │
│            │      MAYOR       │          │    INSPECTOR     │              │
│            │   (Input Side)   │          │  (Output Side)   │              │
│            │                  │          │                  │              │
│            │  ┌────────────┐  │          │  ┌────────────┐  │              │
│            │  │  DESIGNER  │  │          │  │ STRATEGIST │  │              │
│            │  │ (proposes  │  │          │  │ (proposes  │  │              │
│            │  │  require-  │  │          │  │  criteria) │  │              │
│            │  │  ments)    │  │          │  │            │  │              │
│            │  └────────────┘  │          │  └────────────┘  │              │
│            └────────┬─────────┘          └────────┬─────────┘              │
│                     │                             │                         │
│                     ▼                             ▼                         │
│            ┌─────────────────────────────────────────────────┐             │
│            │                    POLECATS                     │             │
│            │                 (Implementation)                │             │
│            └─────────────────────────────────────────────────┘             │
│                                     │                                       │
│                                     ▼                                       │
│            ┌─────────────────────────────────────────────────┐             │
│            │                    REFINERY                     │             │
│            │                   (Merge Gate)                  │             │
│            └─────────────────────────────────────────────────┘             │
│                                     │                                       │
│                                     ▼                                       │
│            ┌─────────────────────────────────────────────────┐             │
│            │              INSPECTOR ECOSYSTEM                │             │
│            │     Verifier → Auditor → Advocate → Critic      │             │
│            │                      ↓                          │             │
│            │                    Judge                        │             │
│            └─────────────────────────────────────────────────┘             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Agent Summary

| Agent | Under | Role | User Interaction |
|-------|-------|------|------------------|
| **Mayor** | User | Coordination | Direct |
| **Designer** | Mayor | Propose requirements | User approves/modifies |
| **Inspector** | User | Quality gate | Direct |
| **Strategist** | Inspector | Propose criteria | User approves/modifies |
| **Polecat** | Mayor | Implementation | None (autonomous) |
| **Verifier** | Inspector | Run tests | None (autonomous) |
| **Auditor** | Inspector | Check compliance | None (autonomous) |
| **Advocate** | Inspector | Argue for code | None (autonomous) |
| **Critic** | Inspector | Argue against code | None (autonomous) |
| **Judge** | Inspector | Deliver verdict | None (autonomous) |

## Work Gate: Spec-Before-Sling

Work cannot be dispatched until verification criteria exists. This ensures quality is designed in, not bolted on.

### The Gate Mechanism

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         SPEC-BEFORE-SLING GATE                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   MAYOR                           INSPECTOR                             │
│   (Input Side)                    (Output Side)                         │
│                                                                         │
│   ┌─────────────┐                 ┌─────────────┐                       │
│   │ 1. Create   │ ──── auto ────▶│ 2. Spec     │                       │
│   │    work     │    creates      │    created  │                       │
│   │    item     │    placeholder  │  (pending)  │                       │
│   └──────┬──────┘                 └──────┬──────┘                       │
│          │                               │                              │
│          ▼                               ▼                              │
│   ┌─────────────┐                 ┌─────────────┐                       │
│   │ BLOCKED     │                 │ 3. User &   │                       │
│   │ waiting for │◀── dependency ──│ Inspector   │                       │
│   │ spec        │                 │ define      │                       │
│   └──────┬──────┘                 │ criteria    │                       │
│          │                        └──────┬──────┘                       │
│          │                               │                              │
│          │                               ▼                              │
│          │                        ┌─────────────┐                       │
│          │                        │ 4. Approve  │                       │
│          │◀─── spec closed ───────│    spec     │                       │
│          │                        └─────────────┘                       │
│          ▼                                                              │
│   ┌─────────────┐                                                       │
│   │ 5. READY    │                                                       │
│   │ can sling   │                                                       │
│   └──────┬──────┘                                                       │
│          │                                                              │
│          ▼                                                              │
│   ┌─────────────┐                                                       │
│   │ 6. Sling to │                                                       │
│   │   Polecat   │                                                       │
│   └─────────────┘                                                       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### Bead Relationships

```
┌─────────────────┐
│   Work Item     │
│   gt-xyz        │
│   status: open  │
│                 │
│   depends_on:   │──────────┐
│   - spec-xyz    │          │
└─────────────────┘          │
                             ▼
                    ┌─────────────────┐
                    │   Test Spec     │
                    │   spec-xyz      │
                    │   type: spec    │
                    │   parent: gt-xyz│
                    │                 │
                    │   criteria:     │
                    │   - ac-1        │
                    │   - ac-2        │
                    └─────────────────┘
```

**Key insight**: Work item DEPENDS ON its spec. Until spec is closed, work is blocked.

### Commands

#### Mayor Side (Creating Work)

```bash
# Create work item - spec placeholder auto-created
bd create --title="Add user authentication"
# Output:
#   Created: gt-abc123 (BLOCKED)
#   Auto-created: spec-abc123 (pending Inspector approval)
#   Work blocked until spec-abc123 is approved

# Check what's blocked
bd blocked
# Output:
#   gt-abc123: Add user authentication
#     └── blocked by: spec-abc123 (test spec pending)

# Try to sling (will fail)
gt sling gt-abc123 polecat-1
# Output:
#   ERROR: Cannot sling gt-abc123
#   Reason: Test spec spec-abc123 not approved
#   Action: Ask Inspector to define verification criteria

# After spec is approved, work appears in ready
bd ready
# Output:
#   gt-abc123: Add user authentication [spec: approved]

# Now sling works
gt sling gt-abc123 polecat-1
# Output:
#   ✓ Assigned gt-abc123 to polecat-1
#   ✓ Test spec spec-abc123 attached
```

#### Inspector Side (Defining Criteria)

```bash
# See pending specs
gt inspect pending
# Output:
#   spec-abc123: Add user authentication [NEEDS CRITERIA]
#     Parent: gt-abc123
#     Created: 2 minutes ago
#     Criteria: (none defined)

# View spec details
gt inspect show spec-abc123
# Output:
#   Test Spec: spec-abc123
#   For: gt-abc123 (Add user authentication)
#   Status: pending
#   Criteria: (none)
#
#   Original task description:
#   > Add user authentication with login, logout, and session management

# Add acceptance criteria
gt inspect add-criteria spec-abc123 \
  --description="Passwords must be hashed with bcrypt" \
  --verify="grep -r 'bcrypt' . | grep -v test"

gt inspect add-criteria spec-abc123 \
  --description="Sessions expire after 30 minutes" \
  --verify="go test -run TestSessionExpiry ./..."

gt inspect add-criteria spec-abc123 \
  --description="No SQL injection vulnerabilities" \
  --verify="gosec ./..."

# User collaborates on criteria
# User: "Also check for rate limiting on login"
gt inspect add-criteria spec-abc123 \
  --description="Login endpoint has rate limiting" \
  --verify="go test -run TestLoginRateLimit ./..."

# Review criteria before approval
gt inspect show spec-abc123
# Output:
#   Test Spec: spec-abc123
#   For: gt-abc123 (Add user authentication)
#   Status: pending
#
#   Acceptance Criteria:
#   [1] Passwords must be hashed with bcrypt
#       Verify: grep -r 'bcrypt' . | grep -v test
#
#   [2] Sessions expire after 30 minutes
#       Verify: go test -run TestSessionExpiry ./...
#
#   [3] No SQL injection vulnerabilities
#       Verify: gosec ./...
#
#   [4] Login endpoint has rate limiting
#       Verify: go test -run TestLoginRateLimit ./...

# Approve the spec (closes it, unblocks work)
gt inspect approve spec-abc123
# Output:
#   ✓ Spec spec-abc123 approved
#   ✓ Work gt-abc123 is now READY
#   ✓ Notified Mayor
```

### Pre-Sling Hook

The gate is enforced by a pre-sling hook:

```yaml
# config/hooks.yaml
hooks:
  pre-sling:
    - name: "spec-gate"
      command: |
        SPEC_ID=$(bd show $BEAD_ID --json | jq -r '.dependencies[] | select(startswith("spec-"))')
        if [ -z "$SPEC_ID" ]; then
          echo "ERROR: No test spec found for $BEAD_ID"
          echo "Create one with: gt inspect create-spec $BEAD_ID"
          exit 1
        fi

        SPEC_STATUS=$(bd show $SPEC_ID --json | jq -r '.status')
        if [ "$SPEC_STATUS" != "closed" ]; then
          echo "ERROR: Test spec $SPEC_ID not approved"
          echo "Status: $SPEC_STATUS"
          echo "Action: gt inspect approve $SPEC_ID"
          exit 1
        fi

        echo "✓ Test spec $SPEC_ID approved"
        exit 0
      on_failure: block
```

### Auto-Create Spec on Work Creation

When Mayor creates work, a spec placeholder is auto-created:

```go
// internal/cmd/create.go (bd create hook)

func createWorkItem(title string, opts CreateOptions) (*Issue, error) {
    // Create the work item
    work, err := beads.Create(title, opts)
    if err != nil {
        return nil, err
    }

    // Auto-create test spec placeholder
    specID := "spec-" + work.ID[3:] // spec-abc123 for gt-abc123
    spec, err := beads.Create(CreateOptions{
        ID:          specID,
        Type:        "test-spec",
        Title:       fmt.Sprintf("Test spec for: %s", title),
        Parent:      work.ID,
        Status:      "open",
        Description: "Pending Inspector approval. Define acceptance criteria.",
    })
    if err != nil {
        return nil, fmt.Errorf("creating test spec: %w", err)
    }

    // Add dependency: work depends on spec
    if err := beads.AddDep(work.ID, spec.ID); err != nil {
        return nil, fmt.Errorf("adding spec dependency: %w", err)
    }

    // Notify Inspector
    mail.Send("inspector", mail.Message{
        Subject: fmt.Sprintf("NEW SPEC NEEDED: %s", title),
        Body:    fmt.Sprintf("Work item %s needs verification criteria.\n\nRun: gt inspect show %s", work.ID, spec.ID),
    })

    return work, nil
}
```

### Tmux Workflow Example

```
┌──────────────────── MAYOR ────────────────────┬──────────────────── INSPECTOR ────────────────────┐
│                                               │                                                   │
│ $ bd create --title="Add user auth"           │                                                   │
│ Created: gt-abc123 (BLOCKED)                  │ [MAIL] NEW SPEC NEEDED: Add user auth             │
│ Auto-created: spec-abc123                     │                                                   │
│ Work blocked until spec approved              │ $ gt inspect pending                              │
│                                               │ spec-abc123: Add user auth [NEEDS CRITERIA]       │
│ $ bd blocked                                  │                                                   │
│ gt-abc123: Add user auth                      │ $ gt inspect add-criteria spec-abc123 \           │
│   └── blocked by: spec-abc123                 │     --description="Passwords hashed with bcrypt"  │
│                                               │ ✓ Added criterion 1                               │
│ $ gt sling gt-abc123 polecat-1                │                                                   │
│ ERROR: Test spec not approved                 │ User: "Add OWASP top 10 check"                    │
│                                               │                                                   │
│ # User talks to Inspector about criteria...   │ $ gt inspect add-criteria spec-abc123 \           │
│                                               │     --description="OWASP top 10 compliance" \     │
│                                               │     --verify="gosec ./..."                        │
│                                               │ ✓ Added criterion 2                               │
│                                               │                                                   │
│                                               │ $ gt inspect approve spec-abc123                  │
│                                               │ ✓ Spec approved                                   │
│ [NOTIFICATION] gt-abc123 now READY            │ ✓ Work gt-abc123 unblocked                        │
│                                               │                                                   │
│ $ bd ready                                    │                                                   │
│ gt-abc123: Add user auth [spec: approved]     │                                                   │
│                                               │                                                   │
│ $ gt sling gt-abc123 polecat-1                │                                                   │
│ ✓ Assigned to polecat-1                       │                                                   │
│ ✓ Test spec attached                          │                                                   │
│                                               │                                                   │
└───────────────────────────────────────────────┴───────────────────────────────────────────────────┘
```

### Skip Gate (Emergency Override)

For urgent fixes, Mayor can bypass the gate (with audit trail):

```bash
# Emergency bypass (requires --force and reason)
gt sling gt-abc123 polecat-1 --force --reason="Critical production bug"
# Output:
#   ⚠️  GATE BYPASS: Slinging without approved spec
#   Reason: Critical production bug
#   Audit: Logged to gt-abc123 history
#
#   ✓ Assigned gt-abc123 to polecat-1
#   ⚠️  No test spec - verification will use default criteria
```

Bypass is logged and can be reviewed:

```bash
bd audit gt-abc123
# Output:
#   2024-01-06 10:30:00 | GATE_BYPASS | mayor | Critical production bug
```

## Integration Points

### 1. Sling Hook (Pre-Work)

When Mayor slings work, Auditor is triggered:

```go
// internal/cmd/sling.go

func slingWork(beadID, rig string) error {
    // Existing sling logic...

    // NEW: Trigger Auditor to create test spec
    if err := triggerAuditPlan(beadID); err != nil {
        return fmt.Errorf("audit planning failed: %w", err)
    }

    // Continue with sling...
}
```

### 2. Refinery Gate (Pre-Merge)

Before Refinery merges, Inspector must approve:

```go
// internal/refinery/merge.go

func processMergeRequest(mr *MergeRequest) error {
    // NEW: Run inspection
    result, err := inspector.Inspect(mr.BeadID, mr.Branch)
    if err != nil {
        return fmt.Errorf("inspection failed: %w", err)
    }

    switch result.Verdict {
    case VerdictPass:
        return merge(mr)  // Proceed with merge
    case VerdictFail:
        return reject(mr, result)  // Send back to Mayor
    case VerdictNeedsHuman:
        return escalate(mr, result)  // Notify Mayor for human review
    }
}
```

### 3. Rejection Flow

When Inspector rejects, work returns to Mayor:

```go
// internal/inspector/reject.go

func reject(mr *MergeRequest, result *VerificationResult) error {
    // Create feedback bead
    feedback := createFeedbackBead(mr.BeadID, result)

    // Notify Mayor
    mail.Send("mayor", mail.Message{
        Subject: fmt.Sprintf("REJECTED: %s", mr.Title),
        Body:    formatRejectionReport(result),
        Attach:  feedback.ID,
    })

    // Update bead status
    beads.Update(mr.BeadID, Status: "rejected")

    return nil
}
```

## Tmux Layout: Two-Pane Architecture

VerMAS uses a two-pane tmux layout separating input (Mayor) from output (Inspector).

### Layout

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              VERMAS TMUX                                    │
├─────────────────────────────────┬───────────────────────────────────────────┤
│                                 │                                           │
│            MAYOR                │              INSPECTOR                    │
│         (Input Side)            │            (Output Side)                  │
│                                 │                                           │
│  Role: Design & Planning        │  Role: Quality & Verification             │
│                                 │                                           │
│  Responsibilities:              │  Responsibilities:                        │
│  • Create work items            │  • Define verification criteria           │
│  • Set requirements             │  • Approve test specs                     │
│  • Coordinate polecats          │  • Review verification results            │
│  • Handle escalations           │  • Manage quality standards               │
│  • Strategic decisions          │  • Adversarial review process             │
│                                 │                                           │
│  User interaction:              │  User interaction:                        │
│  "Build feature X"              │  "Verify with these criteria"             │
│  "Fix bug Y"                    │  "Add security check for Z"               │
│  "Prioritize work"              │  "Lower threshold for tests"              │
│                                 │                                           │
│  [INTERACTIVE]                  │  [INTERACTIVE]                            │
│                                 │                                           │
├─────────────────────────────────┴───────────────────────────────────────────┤
│  [STATUS] Mayor: idle | Inspector: reviewing spec-abc123 | Polecats: 2 busy │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Session Names

```bash
# Two top-level sessions
mayor-<pid>      # Mayor session (input side)
inspector-<pid>  # Inspector session (output side)

# Child sessions (under Inspector)
inspector/auditor-<pid>
inspector/advocate-<pid>
inspector/critic-<pid>
inspector/judge-<pid>
```

### Starting VerMAS

```bash
# Start both panes
gt vermas start
# Output:
#   Starting VerMAS...
#   ✓ Mayor session: mayor-12345
#   ✓ Inspector session: inspector-12346
#   ✓ Tmux layout configured
#
#   Attach with: gt vermas attach

# Attach to split view
gt vermas attach
# Opens tmux with Mayor (left) and Inspector (right)

# Attach to individual panes
gt mayor attach      # Just Mayor
gt inspector attach  # Just Inspector
```

### Layout Configuration

```yaml
# config/vermas.yaml
layout:
  type: horizontal  # or vertical
  ratio: 50         # 50/50 split

  panes:
    - name: mayor
      role: mayor
      position: left
      startup:
        - "gt prime"
        - "gt mail inbox"

    - name: inspector
      role: inspector
      position: right
      startup:
        - "gt prime"
        - "gt inspect pending"

  status_bar:
    enabled: true
    items:
      - "Mayor: #{mayor_status}"
      - "Inspector: #{inspector_status}"
      - "Polecats: #{polecat_count}"
      - "Queue: #{queue_depth}"
```

### Communication Between Panes

Mayor and Inspector communicate via mail (same as all Gas Town agents):

```
┌─────────────── MAYOR ───────────────┐    ┌─────────────── INSPECTOR ───────────────┐
│                                     │    │                                          │
│ $ bd create --title="Add auth"      │    │                                          │
│ Created: gt-abc123                  │────│───▶ [MAIL] NEW SPEC NEEDED: gt-abc123    │
│                                     │    │                                          │
│                                     │    │ $ gt inspect approve spec-abc123         │
│ [MAIL] Spec approved ◀──────────────│────│                                          │
│                                     │    │                                          │
│ $ gt sling gt-abc123 polecat-1      │    │                                          │
│ ✓ Work dispatched                   │    │                                          │
│                                     │    │ [Polecat completes, triggers gate]       │
│                                     │    │                                          │
│                                     │    │ $ gt inspect run gt-abc123               │
│                                     │    │ Running verification...                  │
│                                     │    │ Advocate: PASS                           │
│                                     │    │ Critic: minor concerns                   │
│                                     │    │ Judge: PASS (confidence: 0.85)           │
│                                     │    │                                          │
│ [MAIL] Verification PASSED ◀────────│────│ ✓ Verdict: PASS                          │
│ gt-abc123 merged                    │    │                                          │
│                                     │    │                                          │
└─────────────────────────────────────┘    └──────────────────────────────────────────┘
```

### Switching Focus

```bash
# Switch between panes (within tmux)
Ctrl-b o          # Next pane
Ctrl-b ;          # Previous pane
Ctrl-b q 0        # Go to pane 0 (Mayor)
Ctrl-b q 1        # Go to pane 1 (Inspector)

# From command line
gt focus mayor     # Focus Mayor pane
gt focus inspector # Focus Inspector pane
```

### Independent Context

Each pane maintains independent context:

| Aspect | Mayor | Inspector |
|--------|-------|-----------|
| Conversation history | Planning discussions | Verification discussions |
| Model | Claude (conversational) | Mixed (analytical) |
| CLAUDE.md | Mayor role context | Inspector role context |
| Hook bead | hq-mayor | hq-inspector |
| Mail address | mayor/ | inspector/ |

### Why Two Panes?

1. **Separation of concerns**: Input vs output clearly divided
2. **Parallel work**: User can work with both simultaneously
3. **Independent models**: Different AI models for different purposes
4. **Clean context**: Verification discussion doesn't pollute planning
5. **Observability**: See both sides of the process
6. **Collaboration**: User shapes both design AND quality criteria

## Implementation Phases

### Phase 1: Foundation (Current PR #208)
- [x] Runtime abstraction
- [x] Auditor package structure
- [x] Basic `gt verify` command
- [x] Verified formula

### Phase 2: Test Spec System
- [ ] Test spec bead type
- [ ] `gt audit plan` command
- [ ] `gt audit add-criteria` command
- [ ] Attach spec to work items

### Phase 3: Inspector Ecosystem
- [ ] Inspector coordinator agent
- [ ] Verifier (test runner)
- [ ] Advocate role
- [ ] Critic role
- [ ] Judge role
- [ ] Adversarial workflow

### Phase 4: Integration
- [ ] Sling hook for Auditor
- [ ] Refinery gate integration
- [ ] Rejection flow to Mayor
- [ ] Rework context passing

### Phase 5: Observability
- [ ] Verification dashboard
- [ ] Audit trail
- [ ] Metrics (pass rate, rework rate)
- [ ] Quality trends

## User Experience: Zero Change

### What Users See (Before VerMAS)

```
User → Mayor → Polecat → Refinery → Merge
```

### What Users See (After VerMAS)

```
User → Mayor → Polecat → Refinery → Merge
```

**Same.** The Inspector ecosystem runs automatically. Users only notice:
- Higher quality merges
- Occasional Mayor questions ("Inspector needs clarification on X")

### How Questions Flow

When Inspector (any role) has a question, it goes through Mayor:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        QUESTION FLOW                                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Inspector Role                Mayor                    User           │
│   (Auditor/Judge/etc)                                                   │
│                                                                         │
│   ┌─────────────┐    mail     ┌─────────────┐   prompt   ┌─────────┐  │
│   │ "I need     │ ──────────▶ │  Receives   │ ─────────▶ │ Answers │  │
│   │ clarification│            │  question,  │            │         │  │
│   │ on X"       │             │  asks user  │ ◀───────── │         │  │
│   │             │ ◀────────── │  forwards   │   reply    └─────────┘  │
│   │             │    mail     │  answer     │                         │
│   └─────────────┘             └─────────────┘                         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Key point**: Inspector never talks to user directly. Mayor is the interface.

### Mail Examples

**Inspector → Mayor (question):**
```
From: inspector/judge
To: mayor
Subject: NEEDS_CLARIFICATION: gt-xyz789
Body:
The acceptance criteria for "handle edge cases" is ambiguous.
Specifically: should we handle negative numbers?

Options:
1. Yes, return error for negative input
2. Yes, treat as absolute value
3. No, assume positive input only

Please advise.
```

**Mayor → User (prompt):**
```
Inspector needs clarification on gt-xyz789:

Should we handle negative numbers?
1. Return error for negative input
2. Treat as absolute value
3. Assume positive input only

Which approach?
```

**Mayor → Inspector (answer):**
```
From: mayor
To: inspector/judge
Subject: RE: NEEDS_CLARIFICATION: gt-xyz789
Body:
User decision: Option 1 - return error for negative input.

Proceed with verification using this criterion.
```

## Infrastructure Parity

Inspector roles are first-class Gas Town agents. They get EVERYTHING other agents get.

### Agent Capabilities Matrix

| Capability | Mayor | Polecat | Witness | Refinery | Auditor | Advocate | Critic | Judge |
|------------|-------|---------|---------|----------|---------|----------|--------|-------|
| Mail inbox | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| Send mail | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| Hooks (start/stop) | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| Beads access | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| Molecules | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| Role bead | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| Agent bead | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| tmux session | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| CLAUDE.md | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| Nudge/notify | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

### Inspector Directory Structure

```
<rig>/
├── polecats/           # Existing workers
├── refinery/           # Existing merge queue
├── witness/            # Existing monitor
└── inspector/          # NEW: Inspector ecosystem
    ├── auditor/        # Test spec creator
    │   ├── rig/        # Worktree
    │   ├── .beads/     # Local beads
    │   └── CLAUDE.md   # Role context
    ├── advocate/       # Defense role
    │   ├── rig/
    │   ├── .beads/
    │   └── CLAUDE.md
    ├── critic/         # Prosecution role
    │   ├── rig/
    │   ├── .beads/
    │   └── CLAUDE.md
    └── judge/          # Verdict role
        ├── rig/
        ├── .beads/
        └── CLAUDE.md
```

### Role Beads (created by gt install)

```bash
# Inspector role beads (same pattern as existing roles)
hq-auditor-role     # Auditor role definition
hq-advocate-role    # Advocate role definition
hq-critic-role      # Critic role definition
hq-judge-role       # Judge role definition
hq-verifier-role    # Verifier role definition

# Inspector agent beads (per-rig, like polecats)
<prefix>-auditor    # Auditor agent instance
<prefix>-advocate   # Advocate agent instance
<prefix>-critic     # Critic agent instance
<prefix>-judge      # Judge agent instance
```

## Deterministic Hooks

Every step in the verification flow has hooks for observability and extensibility.

### Hook Points

```yaml
# Verification hooks (added to existing hook system)

hooks:
  # Planning phase
  pre-audit-plan:       # Before Auditor creates test spec
    trigger: "gt sling"
    input: bead_id

  post-audit-plan:      # After test spec created
    trigger: "audit plan complete"
    input: bead_id, spec_id

  # Verification phase
  pre-inspect:          # Before Inspector starts
    trigger: "refinery gate"
    input: mr_id, bead_id

  post-verify:          # After Verifier runs tests
    trigger: "tests complete"
    input: bead_id, test_results

  post-audit:           # After Auditor compliance check
    trigger: "audit complete"
    input: bead_id, audit_report

  pre-advocate:         # Before Advocate starts
    trigger: "advocate start"
    input: bead_id, evidence

  post-advocate:        # After Advocate completes
    trigger: "advocate complete"
    input: bead_id, defense_brief

  pre-critic:           # Before Critic starts
    trigger: "critic start"
    input: bead_id, evidence

  post-critic:          # After Critic completes
    trigger: "critic complete"
    input: bead_id, prosecution_brief

  pre-judgment:         # Before Judge deliberates
    trigger: "judgment start"
    input: bead_id, all_briefs

  post-judgment:        # After verdict delivered
    trigger: "judgment complete"
    input: bead_id, verdict

  # Outcome hooks
  on-verify-pass:       # Verification passed
    trigger: "verdict == PASS"
    input: bead_id, result

  on-verify-fail:       # Verification failed
    trigger: "verdict == FAIL"
    input: bead_id, result, feedback

  on-verify-escalate:   # Needs human
    trigger: "verdict == NEEDS_HUMAN"
    input: bead_id, result, question
```

### Hook Configuration

```yaml
# .claude/settings.json (Inspector hooks)
{
  "hooks": {
    "PostAuditPlan": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "gt mail send polecat -s 'Test spec ready' -m 'Review attached spec'"
          }
        ]
      }
    ],
    "OnVerifyFail": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "gt mail send mayor -s 'REJECTED' -m 'Work failed verification'"
          }
        ]
      }
    ],
    "OnVerifyEscalate": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "gt escalate --to=mayor --reason='Inspector needs human input'"
          }
        ]
      }
    ]
  }
}
```

### Hook Flow Diagram

```
gt sling <bead>
    │
    ├──► [pre-audit-plan] ──► Auditor creates spec ──► [post-audit-plan]
    │
    ▼
Polecat works
    │
    ▼
gt done (Refinery gate)
    │
    ├──► [pre-inspect]
    │
    ├──► Verifier ──► [post-verify]
    │
    ├──► Auditor ──► [post-audit]
    │
    ├──► [pre-advocate] ──► Advocate ──► [post-advocate]
    │
    ├──► [pre-critic] ──► Critic ──► [post-critic]
    │
    ├──► [pre-judgment] ──► Judge ──► [post-judgment]
    │
    └──► Verdict
            │
            ├── PASS ──► [on-verify-pass] ──► Merge
            │
            ├── FAIL ──► [on-verify-fail] ──► Mayor ──► Polecat
            │
            └── NEEDS_HUMAN ──► [on-verify-escalate] ──► Mayor ──► User
```

## Molecules for Inspector Workflows

Inspector roles use molecules (workflows) just like other agents.

### mol-inspect.formula.toml

```toml
description = """
Inspector verification workflow molecule.
Runs the full adversarial verification process.
"""

formula = "inspect"
version = 1

[[steps]]
id = "verify"
title = "Run objective verification"
description = """
Verifier runs:
- Unit tests
- Integration tests
- Linting
- Build check
- Custom test spec

Exit: Test results captured
"""

[[steps]]
id = "audit"
title = "Check compliance"
needs = ["verify"]
description = """
Auditor checks:
- Requirements met
- Acceptance criteria satisfied
- Coverage thresholds

Exit: Compliance report generated
"""

[[steps]]
id = "advocate"
title = "Build defense"
needs = ["audit"]
description = """
Advocate reviews evidence and argues FOR the code.

Exit: Defense brief submitted
"""

[[steps]]
id = "critic"
title = "Build prosecution"
needs = ["audit"]
description = """
Critic reviews evidence and argues AGAINST the code.

Exit: Prosecution brief submitted
"""

[[steps]]
id = "judge"
title = "Deliver verdict"
needs = ["advocate", "critic"]
description = """
Judge reviews all evidence and briefs.
Delivers: PASS, FAIL, or NEEDS_HUMAN

If NEEDS_HUMAN:
  gt mail send mayor -s "NEEDS_CLARIFICATION" -m "<question>"
  gt mol await-signal mayor-response

Exit: Verdict delivered
"""

[[steps]]
id = "route"
title = "Route based on verdict"
needs = ["judge"]
description = """
PASS → Signal refinery to proceed
FAIL → Mail mayor with feedback
NEEDS_HUMAN → Await mayor response, then re-judge
"""
```

## Open Questions

1. **Model Selection**: Should Judge always be the most capable model (Opus)?
2. **Cost Control**: How to balance verification thoroughness vs API costs?
3. **Caching**: Can we cache Advocate/Critic responses for similar changes?
4. **Human Override**: Can Mayor force-merge despite FAIL verdict?
5. **Partial Pass**: What if 80% of criteria pass? Configurable threshold?

## Appendix: Role Prompts

### Auditor (Test Spec Creation)

```
You are the Auditor for VerMAS. Your job is to create acceptance criteria
BEFORE work begins.

Given a task description, produce:
1. Concrete, testable acceptance criteria
2. Specific test commands to verify each criterion
3. Success thresholds

Be specific. Vague criteria like "code should be good" are useless.
Good criteria: "Function returns error for nil input"
```

### Advocate (Defense)

```
You are the Advocate. Your job is to DEFEND this code change.

You will receive:
- The code diff
- Test results
- Audit report

Argue WHY this code should be merged:
- Highlight strengths
- Explain design decisions
- Mitigate concerns
- Show requirement compliance

Be persuasive but honest. Don't defend genuinely bad code.
```

### Critic (Prosecution)

```
You are the Critic. Your job is to ATTACK this code change.

You will receive:
- The code diff
- Test results
- Audit report

Argue WHY this code should NOT be merged:
- Find bugs
- Identify security issues
- Question design decisions
- Note missing tests
- Flag performance concerns

Be thorough but fair. Don't manufacture false concerns.
```

### Judge (Verdict)

```
You are the Judge. Your job is to deliver a VERDICT.

You will receive:
- Verification results (objective)
- Audit report (compliance)
- Advocate's defense
- Critic's concerns

Weigh the evidence and decide:
- PASS: Code meets requirements, concerns addressed
- FAIL: Critical issues must be fixed
- NEEDS_HUMAN: Cannot decide, escalate to human

Provide reasoning for your verdict.
```
