# VerMAS Dry Run: FizzBuzz App

> A complete walkthrough of the VerMAS workflow using a simple FizzBuzz example.
> This version shows the AI-assisted flow where Designer and Strategist do the heavy lifting.

## Overview

This document walks through building a FizzBuzz app in Python using the VerMAS (Verifiable Multi-Agent System) workflow.

**Key principle**: User just approves/adjusts. AI agents propose the details.

## Setup

```
┌─────────────── MAYOR ───────────────┬─────────────── INSPECTOR ───────────────┐
│                                     │                                          │
│ User's terminal (Input Side)        │ User's terminal (Output Side)            │
│                                     │                                          │
│ Mayor + Designer                    │ Inspector + Strategist                   │
│ - Design & planning                 │ - Quality & verification                 │
│ - AI proposes requirements          │ - AI proposes criteria                   │
│ - User approves/modifies            │ - User approves/modifies                 │
│                                     │                                          │
└─────────────────────────────────────┴──────────────────────────────────────────┘
```

---

## Phase 1: User Request (Mayor Pane)

**User types a simple request:**

```
User: Create a fizzbuzz app in Python
```

That's it. User doesn't need to specify details.

---

## Phase 2: Designer Elaborates (Mayor Pane)

**Mayor triggers Designer agent:**

```
Mayor: Let me get Designer to elaborate on that...
```

**Designer analyzes and proposes requirements:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│ DESIGNER: Proposed Requirements for "FizzBuzz App"                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ FUNCTIONAL REQUIREMENTS:                                                │
│ ┌─────────────────────────────────────────────────────────────────────┐│
│ │ • Print numbers 1 through 100, one per line                         ││
│ │ • For multiples of 3, print "Fizz" instead of the number            ││
│ │ • For multiples of 5, print "Buzz" instead of the number            ││
│ │ • For multiples of both 3 and 5, print "FizzBuzz"                   ││
│ │ • Numbers that aren't multiples of 3 or 5 print as-is               ││
│ └─────────────────────────────────────────────────────────────────────┘│
│                                                                         │
│ NON-FUNCTIONAL REQUIREMENTS:                                            │
│ ┌─────────────────────────────────────────────────────────────────────┐│
│ │ • Clean, readable Python code                                        ││
│ │ • Include docstrings for functions                                   ││
│ │ • Pass pylint with score >= 8.0                                     ││
│ │ • Use a reusable function (not just inline logic)                   ││
│ └─────────────────────────────────────────────────────────────────────┘│
│                                                                         │
│ EXPECTED OUTPUT (sample):                                               │
│ ┌─────────────────────────────────────────────────────────────────────┐│
│ │ Line 1:  "1"                                                         ││
│ │ Line 3:  "Fizz"                                                      ││
│ │ Line 5:  "Buzz"                                                      ││
│ │ Line 15: "FizzBuzz"                                                  ││
│ │ Line 98: "98"                                                        ││
│ │ Line 100: "Buzz"                                                     ││
│ └─────────────────────────────────────────────────────────────────────┘│
│                                                                         │
│ Does this capture what you want?                                        │
│                                                                         │
│ [Approve] [Modify] [Start Over]                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**User reviews and approves (or modifies):**

```
User: Looks good. Also make sure it can be imported as a module.

Designer: Got it. Adding requirement:
  • Code should be importable (use if __name__ == "__main__" guard)

Updated requirements ready. Proceeding...
```

**Mayor creates the work item:**

```
Mayor: Creating work item with these requirements...

$ bd create --title="Create FizzBuzz app in Python" \
    --description="[Full requirements from Designer attached]" \
    --type=feature

Created: gt-fb001 (BLOCKED)
Auto-created: spec-fb001 (pending Inspector approval)

Work is blocked until Inspector defines verification criteria.
Notifying Inspector...
```

---

## Phase 3: Strategist Proposes Criteria (Inspector Pane)

**Inspector receives notification:**

```
[MAIL] NEW SPEC NEEDED: Create FizzBuzz app in Python
       Work item: gt-fb001
       Requirements attached
```

**Inspector triggers Strategist agent:**

```
Inspector: Strategist is analyzing the requirements...
```

**Strategist proposes verification criteria:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│ STRATEGIST: Proposed Verification Criteria for gt-fb001                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ Based on the requirements, here are the acceptance criteria I propose:  │
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────────┐│
│ │ [AC-1] OUTPUT LINE COUNT                                            ││
│ │        Verify: Exactly 100 lines of output                          ││
│ │        Why: Requirement says "1 through 100"                        ││
│ │        Command: python fizzbuzz.py | wc -l | grep -q '^100$'        ││
│ └─────────────────────────────────────────────────────────────────────┘│
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────────┐│
│ │ [AC-2] FIZZBUZZ ON LINE 15                                          ││
│ │        Verify: Line 15 (multiple of both 3 and 5) = "FizzBuzz"      ││
│ │        Why: Core FizzBuzz logic - 15 is LCM(3,5)                    ││
│ │        Command: python fizzbuzz.py | sed -n '15p' | grep -qx 'FizzBuzz'│
│ └─────────────────────────────────────────────────────────────────────┘│
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────────┐│
│ │ [AC-3] FIZZ ON LINE 9                                               ││
│ │        Verify: Line 9 (multiple of 3, not 5) = "Fizz"               ││
│ │        Why: Tests the "multiple of 3" branch                        ││
│ │        Command: python fizzbuzz.py | sed -n '9p' | grep -qx 'Fizz'  ││
│ └─────────────────────────────────────────────────────────────────────┘│
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────────┐│
│ │ [AC-4] BUZZ ON LINE 10                                              ││
│ │        Verify: Line 10 (multiple of 5, not 3) = "Buzz"              ││
│ │        Why: Tests the "multiple of 5" branch                        ││
│ │        Command: python fizzbuzz.py | sed -n '10p' | grep -qx 'Buzz' ││
│ └─────────────────────────────────────────────────────────────────────┘│
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────────┐│
│ │ [AC-5] PLAIN NUMBER ON LINE 1                                       ││
│ │        Verify: Line 1 (not multiple of 3 or 5) = "1"                ││
│ │        Why: Tests the "print number as-is" branch                   ││
│ │        Command: python fizzbuzz.py | head -1 | grep -qx '1'         ││
│ └─────────────────────────────────────────────────────────────────────┘│
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────────┐│
│ │ [AC-6] CODE QUALITY                                                 ││
│ │        Verify: Pylint score >= 8.0                                  ││
│ │        Why: Non-functional requirement                              ││
│ │        Command: pylint fizzbuzz.py --fail-under=8.0                 ││
│ └─────────────────────────────────────────────────────────────────────┘│
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────────┐│
│ │ [AC-7] IMPORTABLE MODULE                                            ││
│ │        Verify: Can import without executing                         ││
│ │        Why: User added this requirement                             ││
│ │        Command: python -c "import fizzbuzz" 2>&1 | grep -qv '.'     ││
│ └─────────────────────────────────────────────────────────────────────┘│
│                                                                         │
│ Coverage Analysis:                                                      │
│ • All 4 branches tested (FizzBuzz, Fizz, Buzz, plain number)          │
│ • Boundary tested (line 1, line 100)                                   │
│ • Code quality verified                                                 │
│ • Importability verified                                                │
│                                                                         │
│ Do these criteria look complete?                                        │
│                                                                         │
│ [Approve All] [Add More] [Modify] [Reject]                              │
└─────────────────────────────────────────────────────────────────────────┘
```

**User reviews and approves (or adds more):**

```
User: Good coverage. Also verify line 30 is FizzBuzz (another multiple of 15).

Strategist: Adding criterion:

  [AC-8] FIZZBUZZ ON LINE 30
         Verify: Line 30 (another multiple of 15) = "FizzBuzz"
         Why: User requested additional coverage
         Command: python fizzbuzz.py | sed -n '30p' | grep -qx 'FizzBuzz'

User: Approve all.

Inspector: Approving spec with 8 criteria...

$ gt inspect approve spec-fb001

✓ Spec spec-fb001 approved with 8 criteria
✓ Work gt-fb001 is now READY
✓ Notified Mayor
```

---

## Phase 4: Work Dispatch (Mayor Pane)

**Mayor receives notification:**

```
[MAIL] Spec approved: spec-fb001
       Work gt-fb001 is now ready (8 acceptance criteria)
```

**Mayor dispatches automatically (or user confirms):**

```
Mayor: Spec is approved. Shall I assign this to a polecat?

User: Yes, go ahead.

Mayor: Assigning to polecat-1...

$ gt sling gt-fb001 polecat-1

✓ Pre-sling check: spec-fb001 approved (8 criteria)
✓ Assigned gt-fb001 to polecat-1
✓ Requirements + test spec attached
✓ Polecat notified
```

---

## Phase 5: Implementation (Polecat - Autonomous)

**Polecat receives work with full context:**

```
[HOOK] Work assigned: gt-fb001

Title: Create FizzBuzz app in Python

Requirements (from Designer):
• Print numbers 1-100
• Multiples of 3 → "Fizz"
• Multiples of 5 → "Buzz"
• Multiples of 15 → "FizzBuzz"
• Clean code, docstrings, pylint >= 8.0
• Importable as module

Acceptance Criteria (from Strategist):
[AC-1] 100 lines of output
[AC-2] Line 15 = "FizzBuzz"
[AC-3] Line 9 = "Fizz"
[AC-4] Line 10 = "Buzz"
[AC-5] Line 1 = "1"
[AC-6] Pylint >= 8.0
[AC-7] Importable without side effects
[AC-8] Line 30 = "FizzBuzz"
```

**Polecat implements (no user involvement):**

```python
# fizzbuzz.py
"""FizzBuzz implementation - prints Fizz, Buzz, or FizzBuzz for 1-100."""


def fizzbuzz(n: int) -> str:
    """
    Return FizzBuzz result for a given number.

    Args:
        n: The number to evaluate

    Returns:
        "FizzBuzz" if divisible by 15,
        "Fizz" if divisible by 3,
        "Buzz" if divisible by 5,
        otherwise the number as string.
    """
    if n % 15 == 0:
        return "FizzBuzz"
    if n % 3 == 0:
        return "Fizz"
    if n % 5 == 0:
        return "Buzz"
    return str(n)


def main():
    """Print FizzBuzz results for numbers 1-100."""
    for i in range(1, 101):
        print(fizzbuzz(i))


if __name__ == "__main__":
    main()
```

**Polecat self-verifies against criteria, then submits:**

```
$ git add fizzbuzz.py
$ git commit -m "feat: Add FizzBuzz implementation (gt-fb001)"
$ gt done

✓ Work gt-fb001 submitted to merge queue
✓ Triggering verification gate...
```

---

## Phase 6: Verification (Inspector Ecosystem - Autonomous)

The entire verification runs automatically. User can observe in Inspector pane.

### Step 6a: Verifier Runs Tests

```
┌─────────────────────────────────────────────────────────────────────────┐
│ VERIFIER: Running acceptance criteria...                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ [AC-1] Output line count                                               │
│        $ python fizzbuzz.py | wc -l | grep -q '^100$'                  │
│        Result: PASS ✓                                                   │
│                                                                         │
│ [AC-2] FizzBuzz on line 15                                             │
│        $ python fizzbuzz.py | sed -n '15p' | grep -qx 'FizzBuzz'       │
│        Result: PASS ✓                                                   │
│                                                                         │
│ [AC-3] Fizz on line 9                                                  │
│        $ python fizzbuzz.py | sed -n '9p' | grep -qx 'Fizz'            │
│        Result: PASS ✓                                                   │
│                                                                         │
│ [AC-4] Buzz on line 10                                                 │
│        $ python fizzbuzz.py | sed -n '10p' | grep -qx 'Buzz'           │
│        Result: PASS ✓                                                   │
│                                                                         │
│ [AC-5] Plain number on line 1                                          │
│        $ python fizzbuzz.py | head -1 | grep -qx '1'                   │
│        Result: PASS ✓                                                   │
│                                                                         │
│ [AC-6] Code quality                                                    │
│        $ pylint fizzbuzz.py --fail-under=8.0                           │
│        Result: PASS ✓ (score: 10.00/10)                                │
│                                                                         │
│ [AC-7] Importable module                                               │
│        $ python -c "import fizzbuzz"                                   │
│        Result: PASS ✓ (no output)                                      │
│                                                                         │
│ [AC-8] FizzBuzz on line 30                                             │
│        $ python fizzbuzz.py | sed -n '30p' | grep -qx 'FizzBuzz'       │
│        Result: PASS ✓                                                   │
│                                                                         │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│ SUMMARY: 8/8 PASSED                                                     │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
└─────────────────────────────────────────────────────────────────────────┘
```

### Step 6b: Auditor Checks Compliance

```
┌─────────────────────────────────────────────────────────────────────────┐
│ AUDITOR: Compliance check                                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ Requirements vs Implementation:                                         │
│                                                                         │
│ ✓ Print numbers 1-100               → range(1, 101) confirmed          │
│ ✓ Multiples of 3 → "Fizz"           → n % 3 == 0 logic present         │
│ ✓ Multiples of 5 → "Buzz"           → n % 5 == 0 logic present         │
│ ✓ Multiples of 15 → "FizzBuzz"      → n % 15 == 0 checked FIRST ✓      │
│ ✓ Clean code with docstrings        → All functions documented         │
│ ✓ Pylint >= 8.0                     → Score: 10.00                     │
│ ✓ Importable as module              → __name__ guard present           │
│                                                                         │
│ COMPLIANCE: FULL                                                        │
└─────────────────────────────────────────────────────────────────────────┘
```

### Step 6c: Advocate Builds Defense

```
┌─────────────────────────────────────────────────────────────────────────┐
│ ADVOCATE: Defense Brief                                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ ARGUMENT FOR MERGING:                                                   │
│                                                                         │
│ 1. CORRECTNESS                                                          │
│    • All 8 acceptance criteria pass                                     │
│    • Logic handles the "15 before 3 or 5" edge case correctly          │
│    • Boundary conditions (1, 100) verified                              │
│                                                                         │
│ 2. CODE QUALITY                                                         │
│    • Perfect pylint score (10.00/10)                                   │
│    • Comprehensive docstrings with Args/Returns                        │
│    • Type hints on function signature                                   │
│    • Clean separation of concerns (fizzbuzz func + main)               │
│                                                                         │
│ 3. DESIGN                                                               │
│    • Reusable: fizzbuzz() can be called independently                  │
│    • Importable: no side effects on import                             │
│    • Testable: pure function with predictable output                   │
│                                                                         │
│ RECOMMENDATION: PASS                                                    │
└─────────────────────────────────────────────────────────────────────────┘
```

### Step 6d: Critic Builds Prosecution

```
┌─────────────────────────────────────────────────────────────────────────┐
│ CRITIC: Prosecution Brief                                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ POTENTIAL CONCERNS:                                                     │
│                                                                         │
│ 1. NO UNIT TESTS                                                        │
│    • Implementation has no test file                                    │
│    • However: Acceptance criteria provide functional coverage          │
│    • Verdict: Minor - not in requirements                              │
│                                                                         │
│ 2. MAGIC NUMBERS                                                        │
│    • 3, 5, 15 are hardcoded in function                                │
│    • Could be constants for readability                                │
│    • Verdict: Minor - acceptable for simple script                     │
│                                                                         │
│ 3. NO INPUT VALIDATION                                                  │
│    • fizzbuzz(0) returns "FizzBuzz" (is that correct?)                 │
│    • fizzbuzz(-15) returns "FizzBuzz" (negative handling?)             │
│    • Verdict: Out of scope - requirements specify 1-100 only          │
│                                                                         │
│ ASSESSMENT: All concerns are minor or out of scope                     │
│ RECOMMENDATION: PASS (with minor suggestions noted)                    │
└─────────────────────────────────────────────────────────────────────────┘
```

### Step 6e: Judge Delivers Verdict

```
┌─────────────────────────────────────────────────────────────────────────┐
│ JUDGE: Final Verdict                                                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ EVIDENCE REVIEWED:                                                      │
│ • Verifier: 8/8 criteria PASSED                                        │
│ • Auditor: FULL compliance                                             │
│ • Advocate: Recommends PASS                                            │
│ • Critic: Recommends PASS (minor suggestions only)                     │
│                                                                         │
│ ANALYSIS:                                                               │
│ The implementation satisfies all stated requirements. Code quality     │
│ is excellent (10/10 pylint). The Critic's concerns about edge cases    │
│ (0, negative numbers) are valid observations but explicitly out of     │
│ scope per the requirements ("numbers 1 through 100").                  │
│                                                                         │
│ ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ │
│ ┃                                                                     ┃ │
│ ┃   VERDICT: PASS                                                     ┃ │
│ ┃   CONFIDENCE: 0.95                                                  ┃ │
│ ┃                                                                     ┃ │
│ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ │
│                                                                         │
│ REASONING: All acceptance criteria met, code quality excellent,        │
│ implementation matches requirements precisely.                         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Phase 7: Merge (Refinery - Autonomous)

```
Refinery: Received PASS verdict. Merging...

$ git checkout main
$ git merge polecat/fizzbuzz-gt-fb001 --no-ff
$ git push origin main
$ git branch -d polecat/fizzbuzz-gt-fb001

✓ Merged gt-fb001 to main
✓ Cleaned up branch
✓ Closing bead
```

**Notifications:**

```
[MAIL to Mayor]
Subject: ✅ MERGED: gt-fb001 Create FizzBuzz app
Body: Verification passed (8/8 criteria). Code merged to main.

[MAIL to User - if subscribed]
Subject: ✅ Your FizzBuzz app is ready
Body: Implementation complete and merged.
```

---

## Complete User Interaction Summary

Here's everything the user actually typed:

```
┌─────────────── MAYOR PANE ───────────────┬─────────────── INSPECTOR PANE ───────────────┐
│                                          │                                               │
│ User: Create a fizzbuzz app in Python    │                                               │
│                                          │                                               │
│ [Designer proposes requirements]         │                                               │
│                                          │                                               │
│ User: Looks good. Also make sure it      │                                               │
│       can be imported as a module.       │                                               │
│                                          │                                               │
│                                          │ [Strategist proposes 7 criteria]              │
│                                          │                                               │
│                                          │ User: Good coverage. Also verify line 30      │
│                                          │       is FizzBuzz.                            │
│                                          │                                               │
│                                          │ User: Approve all.                            │
│                                          │                                               │
│ User: Yes, go ahead.                     │                                               │
│ (to confirm sling)                       │                                               │
│                                          │                                               │
│ [Everything else is automatic]           │ [Verification runs automatically]             │
│                                          │                                               │
│ [MAIL] ✅ MERGED: gt-fb001               │                                               │
│                                          │                                               │
└──────────────────────────────────────────┴───────────────────────────────────────────────┘
```

**Total user input: 5 short messages**
- 1 initial request
- 1 modification to requirements
- 1 addition to criteria
- 1 approval
- 1 confirmation to sling

Everything else was proposed by AI and executed automatically.

---

## Alternative Flow: Verification Fails

If Polecat had a bug (e.g., forgot to check 15 first):

### Verifier Catches It:

```
[AC-2] FizzBuzz on line 15
       $ python fizzbuzz.py | sed -n '15p' | grep -qx 'FizzBuzz'
       Expected: "FizzBuzz"
       Actual: "Fizz"
       Result: FAIL ✗

SUMMARY: 7/8 PASSED, 1 FAILED
```

### Judge Delivers FAIL Verdict:

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃   VERDICT: FAIL                                                    ┃
┃   CONFIDENCE: 0.99                                                 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛

REQUIRED FIX:
AC-2 failed - Line 15 should output "FizzBuzz" but got "Fizz".
The code must check for multiples of 15 BEFORE checking for 3 or 5.
```

### Mayor Notified:

```
[MAIL to Mayor]
Subject: ❌ REJECTED: gt-fb001 Create FizzBuzz app

Verification failed: 7/8 criteria passed

FAILED:
[AC-2] Line 15 should output "FizzBuzz" but got "Fizz"

REQUIRED FIX:
Check for multiples of 15 before checking for 3 or 5.

Action: Re-assigning to polecat for fix.
```

### Automatic Rework:

```
Mayor: Work rejected. Re-slinging with failure context...

$ gt sling gt-fb001 polecat-1 --rework

✓ Re-assigned gt-fb001 to polecat-1
✓ Attached: previous failure (AC-2), required fix
```

Polecat fixes, resubmits, passes, merges. **User doesn't need to do anything.**

---

## Summary: What User Does vs What AI Does

| Phase | AI Does | User Does |
|-------|---------|-----------|
| Request | - | "Create fizzbuzz" |
| Requirements | Designer proposes detailed spec | "Looks good" / "Also add X" |
| Criteria | Strategist proposes verification | "Approve" / "Add check for Y" |
| Dispatch | Mayor asks for confirmation | "Yes, go ahead" |
| Implementation | Polecat codes autonomously | Nothing |
| Verification | Inspector ecosystem runs tests | Nothing (can observe) |
| Merge/Reject | Refinery acts on verdict | Nothing |
| Rework | Mayor re-slings with context | Nothing |

**User effort: Minimal approval/adjustment**
**AI effort: All the detailed work**

---

## Appendix: Agent Roles

| Agent | Side | Interaction | Purpose |
|-------|------|-------------|---------|
| **Mayor** | Input | Direct | Coordinates work |
| **Designer** | Input | Proposes → User approves | Elaborates requirements |
| **Inspector** | Output | Direct | Manages verification |
| **Strategist** | Output | Proposes → User approves | Designs test criteria |
| **Polecat** | Middle | Autonomous | Implements code |
| **Verifier** | Output | Autonomous | Runs tests |
| **Auditor** | Output | Autonomous | Checks compliance |
| **Advocate** | Output | Autonomous | Argues FOR merging |
| **Critic** | Output | Autonomous | Argues AGAINST merging |
| **Judge** | Output | Autonomous | Delivers verdict |
| **Refinery** | Middle | Autonomous | Merges or rejects |
