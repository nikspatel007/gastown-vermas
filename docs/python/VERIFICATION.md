# VerMAS Verification System

> Multi-agent verification through adversarial review

## Overview

VerMAS (Verification Multi-Agent System) adds structured verification to the work pipeline. When code is ready to merge, the QA Department triggers a verification pipeline that uses multiple LLM agents to evaluate whether the work meets requirements.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      VERIFICATION PIPELINE                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Work Complete                                                             │
│        │                                                                    │
│        ▼                                                                    │
│   ┌─────────────┐                                                           │
│   │  Designer   │  ← Elaborate requirements into spec                      │
│   └──────┬──────┘                                                           │
│          │                                                                  │
│          ▼                                                                  │
│   ┌─────────────┐                                                           │
│   │ Strategist  │  ← Create objective test criteria                        │
│   └──────┬──────┘                                                           │
│          │                                                                  │
│          ▼                                                                  │
│   ┌─────────────┐                                                           │
│   │  Verifier   │  ← Run shell tests (NO LLM)                              │
│   └──────┬──────┘                                                           │
│          │                                                                  │
│          ▼                                                                  │
│   ┌─────────────┐                                                           │
│   │   Auditor   │  ← LLM review for subjective criteria                    │
│   └──────┬──────┘                                                           │
│          │                                                                  │
│          ├─────────────────────────────────────┐                            │
│          ▼                                     ▼                            │
│   ┌─────────────┐                       ┌─────────────┐                     │
│   │  Advocate   │                       │   Critic    │                     │
│   │ (argue PASS)│                       │ (argue FAIL)│                     │
│   └──────┬──────┘                       └──────┬──────┘                     │
│          │                                     │                            │
│          └─────────────┬───────────────────────┘                            │
│                        ▼                                                    │
│                 ┌─────────────┐                                             │
│                 │    Judge    │  ← Final verdict                           │
│                 └──────┬──────┘                                             │
│                        │                                                    │
│                        ▼                                                    │
│                   PASS / FAIL                                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Design Philosophy

### Objective Over Subjective

The core principle: **prefer shell tests over LLM judgment**.

```
VERIFICATION HIERARCHY (most to least objective)

1. Shell commands (exit code 0/1)     ← PREFERRED
   - Deterministic
   - Repeatable
   - No LLM variability

2. Output matching (grep, diff)        ← GOOD
   - Verifiable
   - Evidence-based

3. LLM Auditor (structured review)     ← FALLBACK
   - For subjective criteria
   - Requires evidence

4. Adversarial review (Advocate/Critic/Judge)  ← TIEBREAKER
   - When Auditor is uncertain
   - Multiple perspectives
```

### Why Adversarial?

Single LLM can have blind spots. Adversarial structure:
- Forces consideration of both pass and fail cases
- Surfaces edge cases and assumptions
- Creates audit trail of reasoning
- Judge synthesizes multiple viewpoints

---

## Verification Roles

### Designer

**Purpose:** Transform raw requirements into detailed specifications.

**Input:** Work order description, user story, or feature request
**Output:** Structured specification with acceptance criteria

**Example prompt:**
```
You are the Designer in a verification system.

REQUIREMENT:
"Implement a function that checks if a number is prime"

Create a specification with:
1. Functional requirements (R1, R2, ...)
2. Edge cases to test
3. Acceptance criteria (AC-1, AC-2, ...)

Be thorough. The Strategist will create tests from your spec.
```

**Example output:**
```
SPECIFICATION: Prime Number Checker

FUNCTIONAL REQUIREMENTS:
R1: Accept an integer input
R2: Return True if prime, False if not prime
R3: Handle edge cases (0, 1, negative numbers)

ACCEPTANCE CRITERIA:
AC-1: is_prime(2) returns True
AC-2: is_prime(4) returns False
AC-3: is_prime(1) returns False
AC-4: is_prime(0) returns False
AC-5: is_prime(-5) returns False
AC-6: is_prime(97) returns True (large prime)
AC-7: Function handles integers up to 10^6
```

---

### Strategist

**Purpose:** Convert acceptance criteria into executable shell tests.

**Input:** Specification from Designer
**Output:** Shell commands that test each criterion

**Critical rule:** Tests must be **objective**
- Exit 0 = PASS
- Exit non-zero = FAIL
- No human interpretation needed

**Example prompt:**
```
You are the Strategist in a verification system.

SPECIFICATION:
[Designer's output]

Create shell commands to verify each acceptance criterion.

RULES:
1. Each test MUST be a shell command
2. Exit 0 = PASS, non-zero = FAIL
3. Tests must be deterministic
4. Capture output as evidence
```

**Example output:**
```
VERIFICATION TESTS:

AC-1: is_prime(2) returns True
Command: python -c "from solution import is_prime; assert is_prime(2) == True"

AC-2: is_prime(4) returns False
Command: python -c "from solution import is_prime; assert is_prime(4) == False"

AC-3: is_prime(1) returns False
Command: python -c "from solution import is_prime; assert is_prime(1) == False"

AC-4: is_prime(0) returns False
Command: python -c "from solution import is_prime; assert is_prime(0) == False"

AC-5: is_prime(-5) returns False
Command: python -c "from solution import is_prime; assert is_prime(-5) == False"

AC-6: is_prime(97) returns True
Command: python -c "from solution import is_prime; assert is_prime(97) == True"

AC-7: Handles large numbers
Command: timeout 5 python -c "from solution import is_prime; is_prime(999983)"
```

---

### Verifier

**Purpose:** Execute shell tests and capture evidence. **NO LLM.**

**Input:** Test commands from Strategist
**Output:** Pass/fail results with stdout/stderr

**This is NOT an LLM agent.** It's a shell executor:

```python
def verify(tests: list[Test]) -> list[TestResult]:
    results = []
    for test in tests:
        result = subprocess.run(
            test.command,
            shell=True,
            capture_output=True,
            timeout=60
        )
        results.append(TestResult(
            criterion=test.criterion,
            command=test.command,
            exit_code=result.returncode,
            stdout=result.stdout.decode(),
            stderr=result.stderr.decode(),
            passed=result.returncode == 0
        ))
    return results
```

**Example output:**
```
VERIFICATION RESULTS:

AC-1: PASS (exit 0)
AC-2: PASS (exit 0)
AC-3: PASS (exit 0)
AC-4: FAIL (exit 1)
  stderr: AssertionError
AC-5: PASS (exit 0)
AC-6: PASS (exit 0)
AC-7: PASS (exit 0)

Summary: 6/7 passed
```

---

### Auditor

**Purpose:** LLM review for criteria that can't be shell-tested.

**When used:**
- Shell test impossible (e.g., "code is readable")
- Shell test failed but evidence is ambiguous
- Subjective quality assessment needed

**Input:** Evidence from Verifier, criterion to evaluate
**Output:** PASS/FAIL with reasoning

**Example prompt:**
```
You are the Auditor in a verification system.

CRITERION: AC-8 "Code follows PEP 8 style guidelines"

EVIDENCE:
File: solution.py
```python
def is_prime(n):
    if n < 2:
        return False
    for i in range(2, int(n**0.5) + 1):
        if n % i == 0:
            return False
    return True
```

Evaluate whether the code meets PEP 8 guidelines.

Provide:
1. ASSESSMENT: PASS or FAIL
2. REASONING: Why you reached this conclusion
3. EVIDENCE QUOTES: Specific code supporting assessment
```

---

### Advocate

**Purpose:** Argue the strongest case for PASS.

**Input:** Evidence and criterion
**Output:** Argument for why verification should pass

**Example prompt:**
```
You are the Advocate. Argue that the criterion IS met.

CRITERION: AC-4 "is_prime(0) returns False"

EVIDENCE:
Test failed with AssertionError.
Actual: is_prime(0) raised ValueError("Input must be positive")

Make the strongest case for PASS:
- Quote specific evidence
- Address potential objections
- Explain why passing is justified

You are the defense attorney. Argue your case.
```

**Example output:**
```
ARGUMENT FOR PASS:

The criterion states is_prime(0) should return False. While the
implementation raises ValueError instead, this behavior is SUPERIOR
to the requirement:

1. Raising ValueError for 0 is mathematically correct - 0 is not
   in the domain of primality testing.

2. The function still correctly identifies that 0 is not prime -
   it just communicates this more explicitly via exception.

3. Returning False for 0 could mask bugs where callers accidentally
   pass invalid inputs.

CONCLUSION: The implementation exceeds the requirement by providing
better error handling. This should be considered a PASS.
```

---

### Critic

**Purpose:** Argue the strongest case for FAIL.

**Input:** Same evidence as Advocate
**Output:** Argument for why verification should fail

**Example prompt:**
```
You are the Critic. Argue that the criterion is NOT met.

[Same evidence as Advocate]

Make the strongest case for FAIL:
- Identify gaps in evidence
- Challenge assumptions
- Point out what's missing

You are the prosecutor. Find the flaws.
```

**Example output:**
```
ARGUMENT FOR FAIL:

The criterion is clear: "is_prime(0) returns False"

1. The implementation does NOT return False for 0. It raises
   ValueError. This is objectively different behavior.

2. The specification said nothing about raising exceptions.
   If exception handling was acceptable, it would have stated
   "returns False or raises an exception."

3. Code that expects False will break. Example:
   ```
   if not is_prime(n):  # Crashes on ValueError
       handle_non_prime()
   ```

4. This is a contract violation. The function signature implies
   it returns bool, not raises exceptions.

CONCLUSION: The implementation violates the explicit requirement.
This must be marked FAIL.
```

---

### Judge

**Purpose:** Render final verdict after hearing both sides.

**Input:** Evidence, Advocate's argument, Critic's argument
**Output:** PASS or FAIL with reasoning

**Example prompt:**
```
You are the Judge. Render a verdict.

CRITERION: AC-4 "is_prime(0) returns False"

EVIDENCE:
[Test results]

ADVOCATE'S ARGUMENT:
[Advocate's case]

CRITIC'S ARGUMENT:
[Critic's case]

Consider both arguments carefully.

VERDICT: PASS or FAIL
REASONING: Why you reached this conclusion

Your decision is final.
```

**Example output:**
```
VERDICT: FAIL

REASONING:

Both arguments have merit, but the Critic's position is correct.

The specification explicitly requires "returns False" for input 0.
Raising an exception is categorically different behavior, regardless
of whether it might be "better" from a design perspective.

If the Advocate believes ValueError is superior, the correct process is:
1. Update the specification to allow exceptions
2. Re-run verification against new spec
3. Document the design decision

As written, the implementation does not meet the criterion.

The Advocate's argument about "exceeding requirements" is sophistry.
Exceeding requirements means doing MORE than required, not doing
something DIFFERENT than required.

FINAL VERDICT: FAIL
```

---

## Running Verification

### As Part of QA

Verification runs automatically when QA receives READY_FOR_QA:

```
READY_FOR_QA received
        │
        ▼
    Tests pass?
        │
        ├── No → REWORK_REQUEST
        │
        └── Yes
              │
              ▼
    VerMAS enabled?
        │
        ├── No → Merge
        │
        └── Yes
              │
              ▼
    Run verification pipeline
        │
        ├── PASS → Merge
        │
        └── FAIL → REWORK_REQUEST
```

### Standalone Verification

```bash
# Verify a specific work order
vermas verify wo-abc123

# Verify with verbose output
vermas verify wo-abc123 --verbose

# Run only objective tests (skip Advocate/Critic/Judge)
vermas verify wo-abc123 --objective-only

# Generate verification report
vermas verify wo-abc123 --output=report.json
```

---

## Configuration

### Enabling VerMAS

```bash
# Per-factory
touch <factory>/.work/.vermas-enabled

# Configure strictness
echo "objective_only = true" > <factory>/.work/vermas.toml
```

### Verification TOML

```toml
# .work/vermas.toml

[verification]
enabled = true
objective_only = false          # Skip Advocate/Critic/Judge
require_all_pass = true         # All criteria must pass
timeout_seconds = 300           # Max verification time

[thresholds]
min_criteria = 3                # Minimum acceptance criteria
max_fail_rate = 0.1             # Max 10% criteria can fail

[agents]
designer_model = "claude-sonnet"
strategist_model = "claude-sonnet"
auditor_model = "claude-sonnet"
advocate_model = "claude-sonnet"
critic_model = "claude-sonnet"
judge_model = "claude-opus"     # Most capable for final decision
```

---

## Evidence Storage

Verification produces evidence stored in `.work/evidence/`:

```
.work/evidence/
├── wo-abc123/
│   ├── spec.json           # Designer output
│   ├── tests.json          # Strategist output
│   ├── results.json        # Verifier output
│   ├── audits/
│   │   └── AC-8.json       # Auditor reviews
│   ├── arguments/
│   │   ├── advocate.json   # Advocate's argument
│   │   └── critic.json     # Critic's argument
│   └── verdict.json        # Judge's decision
```

### Evidence Schema

```json
{
  "work_order_id": "wo-abc123",
  "timestamp": "2026-01-06T12:00:00.000Z",
  "verdict": "PASS",
  "criteria": [
    {
      "id": "AC-1",
      "description": "is_prime(2) returns True",
      "test_type": "shell",
      "command": "python -c ...",
      "passed": true,
      "evidence": {
        "exit_code": 0,
        "stdout": "",
        "stderr": ""
      }
    }
  ],
  "summary": {
    "total": 7,
    "passed": 7,
    "failed": 0,
    "audited": 1
  }
}
```

---

## Events Emitted

See [EVENTS.md](./EVENTS.md) for full event schema.

| Event | When | Data |
|-------|------|------|
| `verify.started` | Pipeline begins | work_order_id, process_id |
| `verify.spec_created` | Designer done | spec summary |
| `verify.tests_generated` | Strategist done | test count |
| `verify.test_executed` | Each test runs | criterion, passed |
| `verify.audited` | Auditor reviews | criterion, assessment |
| `verify.verdict` | Judge decides | PASS/FAIL, reasoning |

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [AGENTS.md](./AGENTS.md) - Verification role definitions
- [WORKFLOWS.md](./WORKFLOWS.md) - Verification process
- [EVENTS.md](./EVENTS.md) - Verification events
- [SCHEMAS.md](./SCHEMAS.md) - Evidence data specifications
- [EVALUATION.md](./EVALUATION.md) - Verification accuracy metrics
